package p256

import (
	cryptoecdsa "crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"gnark-benchmark/utils"
	"math/big"
	"strconv"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/math/emulated"
	"golang.org/x/crypto/cryptobyte"
	"golang.org/x/crypto/cryptobyte/asn1"
)

const NumSignatures = 1

var circuitName string

const curveId = ecc.BW6_761

func init() {
	circuitName = "p256-" + strconv.Itoa(NumSignatures)
}

func compileCircuit(newBuilder frontend.NewBuilder) (constraint.ConstraintSystem, error) {
	circuit := EcdsaCircuit[emulated.P256Fp, emulated.P256Fr]{}
	r1cs, err := frontend.Compile(curveId.ScalarField(), newBuilder, &circuit)
	if err != nil {
		return nil, err
	}
	return r1cs, nil
}

func generateWitness() (witness.Witness, error) {
	witness := EcdsaCircuit[emulated.P256Fp, emulated.P256Fr]{}
	for i := 0; i < NumSignatures; i++ {
		// Keygen
		privKey, _ := cryptoecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		publicKey := privKey.PublicKey

		// Sign
		msg, err := genRandomBytes(i + 20)
		if err != nil {
			panic(err)
		}
		msgHash := sha256.Sum256(msg)
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

	witnessData, err := frontend.NewWitness(&witness, curveId.ScalarField())
	if err != nil {
		panic(err)
	}

	return witnessData, nil
}

func Groth16Setup(fileDir string) {
	utils.Groth16Setup(fileDir, circuitName, compileCircuit)
}

func Groth16Prove(fileDir string) {
	utils.Groth16Prove(curveId, fileDir, circuitName, generateWitness)
	// // proveStart := time.Now()
	// // Witness generation
	// start := time.Now()
	// witnessData, err := generateWitness()
	// if err != nil {
	// 	panic(err)
	// }
	// elapsed := time.Since(start)
	// log.Printf("Witness Generation: %d ms", elapsed.Milliseconds())

	// // Read files
	// start = time.Now()
	// r1cs := groth16.NewCS(curveId)
	// utils.ReadFromFile(r1cs, fileDir+circuitName+".r1cs")
	// elapsed = time.Since(start)
	// log.Printf("Read r1cs: %d ms", elapsed.Milliseconds())

	// start = time.Now()
	// pk := groth16.NewProvingKey(curveId)

	// utils.UnsafeReadFromFile(pk, fileDir+circuitName+".zkey")
	// elapsed = time.Since(start)
	// log.Printf("Read zkey: %d ms", elapsed.Milliseconds())

	// // Proof generation
	// start = time.Now()
	// proof, err := groth16.Prove(r1cs, pk, witnessData)
	// if err != nil {
	// 	panic(err)
	// }
	// elapsed = time.Since(start)
	// log.Printf("Prove: %d ms", elapsed.Milliseconds())

	// // proveElapsed := time.Since(proveStart)
	// // log.Printf("Prove: %d ms", proveElapsed.Milliseconds())

	// utils.WriteToFile(proof, fileDir+circuitName+".proof")
	// // Proof verification
	// publicWitness, err := witnessData.Public()
	// if err != nil {
	// 	panic(err)
	// }
	// vk := groth16.NewVerifyingKey(curveId)
	// utils.ReadFromFile(vk, fileDir+circuitName+".vkey")
	// err = groth16.Verify(proof, vk, publicWitness)
	// if err != nil {
	// 	panic(err)
	// }
	// // Export Solidity verifier
	// f, _ := os.Create(fileDir + circuitName + "Verifier.sol")
	// err = vk.ExportSolidity(f)
	// if err != nil {
	// 	panic(err)
	// }
}

func genRandomBytes(size int) ([]byte, error) {
	blk := make([]byte, size)
	_, err := rand.Read(blk)
	if err != nil {
		return nil, err
	}
	return blk, nil
}
