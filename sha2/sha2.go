package sha2

import (
	"crypto/sha256"
	"gnark-benchmark/utils"
	"log"
	"time"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/std/math/uints"
)
const circuitName = "sha2"
const inputLength = 128

func compileCircuit(newBuilder frontend.NewBuilder) (constraint.ConstraintSystem, error) {
	circuit := sha2Circuit{
		In: make([]uints.U8, inputLength),
	}
	r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), newBuilder, &circuit)
	if err != nil {
		return nil, err
	}
	return r1cs, nil
}

func generateWitness() (witness.Witness, error) {
	input := make([]byte, inputLength)
	dgst := sha256.Sum256(input)
	witness := sha2Circuit{
		In: uints.NewU8Array(input[:]),
	}
	copy(witness.Expected[:], uints.NewU8Array(dgst[:]))

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
	proof, err := groth16.Prove(r1cs, pk, witnessData)
	if err != nil {
		panic(err)
	}
	elapsed = time.Since(start)
	log.Printf("Prove: %d ms", elapsed.Milliseconds())

	// proveElapsed := time.Since(proveStart)
	// log.Printf("Prove: %d ms", proveElapsed.Milliseconds())

	utils.WriteToFile(proof, fileDir+circuitName+".proof")
	// Proof verification
	// publicWitness, err := witnessData.Public()
	// if err != nil {
	// 	panic(err)
	// }
	// vk := groth16.NewVerifyingKey(ecc.BN254)
	// utils.ReadFromFile(vk, fileDir+circuitName+".vkey")
	// err = groth16.Verify(proof, vk, publicWitness)
	// if err != nil {
	// 	panic(err)
	// }
}