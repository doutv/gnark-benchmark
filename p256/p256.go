package p256

import (
	cryptoecdsa "crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"gnark-benchmark/utils"
	"log"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/solidity"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/std/math/emulated"
	"github.com/consensys/gnark/std/math/uints"
	"golang.org/x/crypto/cryptobyte"
	"golang.org/x/crypto/cryptobyte/asn1"
	"golang.org/x/crypto/sha3"
)

const NumSignatures = 1

var circuitName string

func init() {
	circuitName = "p256-" + strconv.Itoa(NumSignatures)
}

func compileCircuit(newBuilder frontend.NewBuilder) (constraint.ConstraintSystem, error) {
	circuit := EcdsaCircuit[emulated.P256Fp, emulated.P256Fr]{}
	r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), newBuilder, &circuit)
	if err != nil {
		return nil, err
	}
	return r1cs, nil
}

func generateWitness() (witness.Witness, error) {
	witness := EcdsaCircuit[emulated.P256Fp, emulated.P256Fr]{}
	perSignatureHashSize := 2*emulated.P256Fp{}.NbLimbs() + emulated.P256Fr{}.NbLimbs()
	hashIn := make([]byte, 0, NumSignatures*perSignatureHashSize)
	for i := 0; i < NumSignatures; i++ {
		// Keygen
		privKey, _ := cryptoecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		publicKey := privKey.PublicKey

		// Sign
		msg, err := genRandomBytes(i + 20)
		if err != nil {
			panic(err)
		}
		msgHash := keccak256(msg)
		sigBin, _ := privKey.Sign(rand.Reader, msgHash[:], nil)

		// Try verify
		var (
			r, s  = &big.Int{}, &big.Int{}
			inner cryptobyte.String
		)
		input := cryptobyte.String(sigBin)
		if !input.ReadASN1(&inner, asn1.SEQUENCE) ||
			!input.Empty() ||
			!inner.ReadASN1Integer(r) ||
			!inner.ReadASN1Integer(s) ||
			!inner.Empty() {
			panic("invalid sig")
		}
		flag := cryptoecdsa.Verify(&publicKey, msgHash[:], r, s)
		if !flag {
			println("can't verify signature")
		}

		// hashIn += Pub[i].X + Pub[i].Y + Msg[i]
		pubX := publicKey.X.Bytes()
		pubY := publicKey.Y.Bytes()
		println("pubX: ", hex.EncodeToString(pubX))
		println("pubY: ", hex.EncodeToString(pubY))
		println("msgHash: ", hex.EncodeToString(msgHash[:]))
		hashIn = append(hashIn, pubX[:]...)
		hashIn = append(hashIn, pubY[:]...)
		hashIn = append(hashIn, msgHash[:]...)
		// Assign to circuit witness
		witness.Sig[i] = Signature[emulated.P256Fr]{
			R: emulated.ValueOf[emulated.P256Fr](r),
			S: emulated.ValueOf[emulated.P256Fr](s),
		}
		witness.Msg[i] = emulated.ValueOf[emulated.P256Fr](msgHash[:])
		witness.Pub[i] = PublicKey[emulated.P256Fp, emulated.P256Fr]{
			X: emulated.ValueOf[emulated.P256Fp](publicKey.X),
			Y: emulated.ValueOf[emulated.P256Fp](publicKey.Y),
		}
	}
	hashOut := keccak256(hashIn)
	println("hashOut: ", hex.EncodeToString(hashOut[:]))
	copy(witness.Commitment[:], uints.NewU8Array(hashOut[:]))

	witnessData, err := frontend.NewWitness(&witness, ecc.BN254.ScalarField())
	if err != nil {
		panic(err)
	}

	return witnessData, nil
}

func Groth16Setup(fileDir string) {
	r1cs, err := compileCircuit(r1cs.NewBuilder)
	if err != nil {
		panic(err)
	}
	pk, vk, err := groth16.Setup(r1cs)
	if err != nil {
		panic(err)
	}
	// Write to file
	utils.WriteToFile(pk, fileDir+circuitName+".zkey")
	utils.WriteToFile(r1cs, fileDir+circuitName+".r1cs")
	utils.WriteToFile(vk, fileDir+circuitName+".vkey")
}

func Groth16Prove(fileDir string) {
	// proveStart := time.Now()
	// Witness generation
	start := time.Now()
	witnessData, err := generateWitness()
	if err != nil {
		panic(err)
	}
	elapsed := time.Since(start)
	log.Printf("Witness Generation: %d ms", elapsed.Milliseconds())

	// Read files
	start = time.Now()
	r1cs := groth16.NewCS(ecc.BN254)
	utils.ReadFromFile(r1cs, fileDir+circuitName+".r1cs")
	elapsed = time.Since(start)
	log.Printf("Read r1cs: %d ms", elapsed.Milliseconds())

	start = time.Now()
	pk := groth16.NewProvingKey(ecc.BN254)

	utils.UnsafeReadFromFile(pk, fileDir+circuitName+".zkey")
	elapsed = time.Since(start)
	log.Printf("Read zkey: %d ms", elapsed.Milliseconds())

	// Proof generation
	start = time.Now()
	proverOption := solidity.WithProverTargetSolidityVerifier(backend.GROTH16)
	proof, err := groth16.Prove(r1cs, pk, witnessData, proverOption)
	if err != nil {
		panic(err)
	}
	elapsed = time.Since(start)
	log.Printf("Prove: %d ms", elapsed.Milliseconds())

	// proveElapsed := time.Since(proveStart)
	// log.Printf("Prove: %d ms", proveElapsed.Milliseconds())

	utils.WriteToFile(proof, fileDir+circuitName+".proof")
	// Proof verification
	publicWitness, err := witnessData.Public()
	if err != nil {
		panic(err)
	}
	vk := groth16.NewVerifyingKey(ecc.BN254)
	utils.ReadFromFile(vk, fileDir+circuitName+".vkey")

	// proof to hex
	_proof, ok := proof.(interface{ MarshalSolidity() []byte })
	if !ok {
		panic("proof does not implement MarshalSolidity()")
	}
	proofBytes := _proof.MarshalSolidity()
	println("len(proof) =", len(proofBytes))
	printUint256(proofBytes)

	publicInput, err := publicWitness.MarshalBinary()
	if err != nil {
		panic(err)
	}
	// https://github.com/Consensys/gnark/blob/dc04a1d3b221dbe7571b5a8394b55d02c2872700/test/assert_solidity.go#L78-L83
	// that's quite dirty...
	// first 4 bytes -> nbPublic
	// next 4 bytes -> nbSecret
	// next 4 bytes -> nb elements in the vector (== nbPublic + nbSecret)
	publicInput = publicInput[12:]
	println("len(publicInput) =", len(publicInput))
	printUint256(publicInput)
	err = groth16.Verify(proof, vk, publicWitness, solidity.WithVerifierTargetSolidityVerifier(backend.GROTH16))
	if err != nil {
		panic(err)
	}
	// Export Solidity verifier
	f, _ := os.Create(fileDir + circuitName + "Verifier.sol")
	err = vk.ExportSolidity(f)
	if err != nil {
		panic(err)
	}
}

func genRandomBytes(size int) ([]byte, error) {
	blk := make([]byte, size)
	_, err := rand.Read(blk)
	if err != nil {
		return nil, err
	}
	return blk, nil
}

func keccak256(data []byte) (digest [32]byte) {
	h := sha3.NewLegacyKeccak256()
	h.Write(data)
	h.Sum(digest[:0])
	return
}

func printUint256(data []byte) {
	// println(hex.EncodeToString(data))
	for i := 0; i < len(data); i += 32 {
		println(hex.EncodeToString(data[i : i+32]))
	}
}