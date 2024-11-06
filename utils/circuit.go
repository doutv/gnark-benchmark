package utils

import (
	"log"
	"os"
	"time"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/r1cs"
)

// CircuitCompiler represents a function that compiles a circuit
type CircuitCompiler func(newBuilder frontend.NewBuilder) (constraint.ConstraintSystem, error)

// WitnessGenerator represents a function that generates witness data
type WitnessGenerator func() (witness.Witness, error)

// Groth16Setup performs the common setup for Groth16 proving system
func Groth16Setup(fileDir, circuitName string, compiler CircuitCompiler) {
	r1cs, err := compiler(r1cs.NewBuilder)
	if err != nil {
		panic(err)
	}
	pk, vk, err := groth16.Setup(r1cs)
	if err != nil {
		panic(err)
	}
	// Write to file
	// Create directory if not exists
	if _, err := os.Stat(fileDir); os.IsNotExist(err) {
		err := os.MkdirAll(fileDir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	WriteToFile(pk, fileDir+circuitName+".zkey")
	WriteToFile(r1cs, fileDir+circuitName+".r1cs")
	WriteToFile(vk, fileDir+circuitName+".vkey")
}

// Groth16Prove performs the common proving process for Groth16
func Groth16Prove(curveId ecc.ID, fileDir, circuitName string, witnessGen WitnessGenerator) {
	// Witness generation
	start := time.Now()
	witnessData, err := witnessGen()
	if err != nil {
		panic(err)
	}
	elapsed := time.Since(start)
	log.Printf("Witness Generation: %d ms", elapsed.Milliseconds())

	// Read files
	start = time.Now()
	r1cs := groth16.NewCS(curveId)
	ReadFromFile(r1cs, fileDir+circuitName+".r1cs")
	elapsed = time.Since(start)
	log.Printf("Read r1cs: %d ms", elapsed.Milliseconds())

	start = time.Now()
	pk := groth16.NewProvingKey(curveId)
	UnsafeReadFromFile(pk, fileDir+circuitName+".zkey")
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

	WriteToFile(proof, fileDir+circuitName+".proof")
} 