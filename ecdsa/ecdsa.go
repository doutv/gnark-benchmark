package ecdsa

import (
	"crypto/rand"
	"hash"
	"io"
	"os"
	"time"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	secp_mimc "github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	secp_ecdsa "github.com/consensys/gnark-crypto/ecc/secp256k1/ecdsa"
	"github.com/consensys/gnark/backend/plonk"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/frontend/cs/scs"
	"github.com/consensys/gnark/test/unsafekzg"

	"log"
	"math/big"

	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/algebra/emulated/sw_emulated"
	"github.com/consensys/gnark/std/hash/mimc"
	"github.com/consensys/gnark/std/math/emulated"
	"github.com/consensys/gnark/std/signature/ecdsa"
)

type kycCircuit struct {
	Age    frontend.Variable `gnark:"age,secret"`
	Gender frontend.Variable `gnark:"gender,secret"`
	Nation frontend.Variable `gnark:"nation,secret"`

	//credentail
	Credential frontend.Variable `gnark:",public"`
	ExpireTime frontend.Variable `gnark:",public"`
	//condition
	MinAge         frontend.Variable `gnark:",public"`
	MaxAge         frontend.Variable `gnark:",public"`
	ContainNations frontend.Variable `gnark:",public"`
	//isMale frontend.Variable `gnark:",public"`

	//signature
	PublicKey ecdsa.PublicKey[emulated.Secp256k1Fp, emulated.Secp256k1Fr] `gnark:",public"`
	Signature ecdsa.Signature[emulated.Secp256k1Fr]                       `gnark:",public"`
}

func (circuit *kycCircuit) Define(api frontend.API) error {
	// check signature
	//hash function for kyc credential
	hFunc, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}
	hFunc.Reset()
	// the signature is on h(nonce ∥ amount ∥ senderpubKey (x&y) ∥ receiverPubkey(x&y))
	hFunc.Write(circuit.Credential, circuit.Age, circuit.Gender, circuit.Nation, circuit.ExpireTime)
	message_hash := hFunc.Sum()

	scalarApi, err := emulated.NewField[emulated.Secp256k1Fr](api)
	if err != nil {
		return err
	}
	messageHashBits := api.ToBinary(message_hash, 256)
	gotMessageHash := scalarApi.FromBits(messageHashBits...)

	// signature verify
	circuit.PublicKey.Verify(api, sw_emulated.GetCurveParams[emulated.Secp256k1Fp](), gotMessageHash, &circuit.Signature)

	//// check age
	// minage < age
	// api.AssertIsEqual(frontend.Variable(-1), api.Cmp(circuit.MinAge, circuit.Age))
	// // maxage > age
	// api.AssertIsEqual(frontend.Variable(1), api.Cmp(circuit.MaxAge, circuit.Age))
	// ////check gender
	// //api.AssertIsEqual(circuit.isMale, circuit.Gender)
	// // contains nation
	// constrain_nations := api.ToBinary(circuit.ContainNations, 253)
	// nation := api.ToBinary(circuit.Nation, 253)
	// result := make([]frontend.Variable, 0)
	// for i := range constrain_nations {
	// 	result = append(result, api.And(constrain_nations[i], nation[i]))
	// }
	//api.AssertIsDifferent(api.FromBinary(result...), frontend.Variable(0))
	return nil
}

type kycCredential struct {
	Credential uint64
	Age        uint64
	Gender     uint64
	Nation     uint64
	Expirtime  uint64
}

// Sign signs a transaction
func (t *kycCredential) Sign(priv *secp_ecdsa.PrivateKey, h hash.Hash) (secp_ecdsa.Signature, []byte, error) {

	h.Reset()

	////var frNonce, msg fr.Element
	var credential fr.Element
	credential.SetUint64(t.Credential)
	b := credential.Bytes()
	_, _ = h.Write(b[:])

	var age fr.Element
	age.SetUint64(t.Age)
	b = age.Bytes()
	_, _ = h.Write(b[:])

	var gender fr.Element
	gender.SetUint64(t.Gender)
	b = gender.Bytes()
	_, _ = h.Write(b[:])

	var nation fr.Element
	nation.SetUint64(t.Nation)
	b = nation.Bytes()
	_, _ = h.Write(b[:])

	var expirTime fr.Element
	expirTime.SetUint64(t.Expirtime)
	b = expirTime.Bytes()
	_, _ = h.Write(b[:])

	msg := h.Sum(nil)

	sigBin, err := priv.Sign(msg, nil)
	if err != nil {
		return secp_ecdsa.Signature{}, nil, err
	}
	var sig secp_ecdsa.Signature
	if _, err := sig.SetBytes(sigBin); err != nil {
		return secp_ecdsa.Signature{}, nil, err
	}
	return sig, msg, nil
}

var hFunc = secp_mimc.NewMiMC()

func compileCircuit(newBuilder frontend.NewBuilder) (constraint.ConstraintSystem, error) {
	circuit := kycCircuit{}
	r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), newBuilder, &circuit)
	if err != nil {
		return nil, err
	}
	return r1cs, nil
}

func generateWitness(hFunc hash.Hash) (witness.Witness, error) {
	// generate parameters
	privKey, _ := secp_ecdsa.GenerateKey(rand.Reader)

	// sign
	credential := kycCredential{Credential: 12, Age: 18, Gender: 1, Nation: 0b10, Expirtime: 123}
	sigBin, _, err := credential.Sign(privKey, hFunc)
	if err != nil {
		panic(err)
	}

	// unmarshal signature
	r, s := new(big.Int), new(big.Int)
	r.SetBytes(sigBin.R[:32])
	s.SetBytes(sigBin.S[:32])

	witnessCircuit := kycCircuit{
		Signature: ecdsa.Signature[emulated.Secp256k1Fr]{
			R: emulated.ValueOf[emulated.Secp256k1Fr](r),
			S: emulated.ValueOf[emulated.Secp256k1Fr](s),
		},
		Credential: credential.Credential, Age: credential.Age, Gender: credential.Gender, Nation: credential.Nation, ExpireTime: credential.Expirtime,
		MinAge: 1,
		MaxAge: 60,
		//isMale: 1,
		ContainNations: 0b10011,
		PublicKey: ecdsa.PublicKey[emulated.Secp256k1Fp, emulated.Secp256k1Fr]{
			X: emulated.ValueOf[emulated.Secp256k1Fp](privKey.PublicKey.A.X),
			Y: emulated.ValueOf[emulated.Secp256k1Fp](privKey.PublicKey.A.Y),
		},
	}

	witnessData, err := frontend.NewWitness(&witnessCircuit, ecc.BN254.ScalarField())
	if err != nil {
		panic(err)
	}

	return witnessData, nil
}

func writeToFile(data io.WriterTo, fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = data.WriteTo(file)
	if err != nil {
		panic(err)
	}
}

func readFromFile(data io.ReaderFrom, fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// Use the ReadFrom method to read the file's content into data.
	if _, err := data.ReadFrom(file); err != nil {
		panic(err)
	}
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
	writeToFile(pk, fileDir+"ecdsa.zkey")
	writeToFile(r1cs, fileDir+"ecdsa.r1cs")
	writeToFile(vk, fileDir+"ecdsa.vkey")
}

func Groth16ProveAndVerify(fileDir string) {
	proveStart := time.Now()
	// Witness generation
	start := time.Now()
	witnessData, err := generateWitness(hFunc)
	if err != nil {
		panic(err)
	}
	elapsed := time.Since(start)
	log.Printf("Witness Generation: %d ms", elapsed.Milliseconds())

	// Read files
	start = time.Now()
	r1cs := groth16.NewCS(ecc.BN254)
	readFromFile(r1cs, fileDir+"ecdsa.r1cs")
	elapsed = time.Since(start)
	log.Printf("Read r1cs: %d ms", elapsed.Milliseconds())

	start = time.Now()
	pk := groth16.NewProvingKey(ecc.BN254)
	file, err := os.Open(fileDir+"ecdsa.zkey")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// UnsafeReadFrom is faster than ReadFrom
	if _, err := pk.UnsafeReadFrom(file); err != nil {
		panic(err)
	}
	// readFromFile(pk, "ecdsa.zkey")
	elapsed = time.Since(start)
	log.Printf("Read zkey: %d ms", elapsed.Milliseconds())

	// Proof generation
	start = time.Now()
	proof, err := groth16.Prove(r1cs, pk, witnessData)
	if err != nil {
		panic(err)
	}
	elapsed = time.Since(start)
	log.Printf("Prove: %d ms", elapsed.Milliseconds())

	proveElapsed := time.Since(proveStart)
	log.Printf("Total Prove time: %d ms", proveElapsed.Milliseconds())
	
	writeToFile(proof, "ecdsa.proof")
	// Proof verification
	publicWitness, err := witnessData.Public()
	if err != nil {
		panic(err)
	}
	vk := groth16.NewVerifyingKey(ecc.BN254)
	readFromFile(vk, fileDir+"ecdsa.vkey")
	err = groth16.Verify(proof, vk, publicWitness)
	if err != nil {
		panic(err)
	}
}

func PlonkSetup(fileDir string) {
	circuit := kycCircuit{}
	r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), scs.NewBuilder, &circuit)
	if err != nil {
		panic(err)
	}
	srs, srsLagrange, err := unsafekzg.NewSRS(r1cs)
	if err != nil {
		panic(err)
	}
	pk, vk, err := plonk.Setup(r1cs, srs, srsLagrange)
	if err != nil {
		panic(err)
	}
	// Write to file
	writeToFile(pk, fileDir+"ecdsa.plonk.zkey")
	writeToFile(r1cs, fileDir+"ecdsa.plonk.r1cs")
	writeToFile(vk, fileDir+"ecdsa.plonk.vkey")
}

func PlonkProveAndVerify(fileDir string) {
	proveStart := time.Now()
	witnessData, err := generateWitness(hFunc)
	if err != nil {
		panic(err)
	}
	// Read files
	start := time.Now()
	r1cs := plonk.NewCS(ecc.BN254)
	readFromFile(r1cs, fileDir+"ecdsa.plonk.r1cs")
	elapsed := time.Since(start)
	log.Printf("Read r1cs: %d ms", elapsed.Milliseconds())

	start = time.Now()
	pk := plonk.NewProvingKey(ecc.BN254)
	file, err := os.Open(fileDir+"ecdsa.plonk.zkey")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// UnsafeReadFrom is faster than ReadFrom
	if _, err := pk.UnsafeReadFrom(file); err != nil {
		panic(err)
	}
	elapsed = time.Since(start)
	log.Printf("Read zkey: %d ms", elapsed.Milliseconds())

	// Proof generation
	start = time.Now()
	proof, err := plonk.Prove(r1cs, pk, witnessData)
	if err != nil {
		panic(err)
	}
	elapsed = time.Since(start)
	log.Printf("Prove: %d ms", elapsed.Milliseconds())

	proveElapsed := time.Since(proveStart)
	log.Printf("Total Prove time: %d ms", proveElapsed.Milliseconds())
	writeToFile(proof, "ecdsa.plonk.proof")

	log.Println("start verify")
	publicWitness, err := witnessData.Public()
	if err != nil {
		panic(err)
	}
	vk := plonk.NewVerifyingKey(ecc.BN254)
	readFromFile(vk, fileDir+"ecdsa.plonk.vkey")
	err = plonk.Verify(proof, vk, publicWitness)
	if err != nil {
		panic(err)
	}
	log.Println("end verify")
}