package keccak

import (
	"gnark-benchmark/utils"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/math/uints"
	"golang.org/x/crypto/sha3"
)

const (
	circuitName = "keccak"
	inputLength = 128
)

func compileCircuit(newBuilder frontend.NewBuilder) (constraint.ConstraintSystem, error) {
	circuit := keccakCircuit{
		In: make([]uints.U8, inputLength),
	}
	r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), newBuilder, &circuit)
	if err != nil {
		return nil, err
	}
	return r1cs, nil
}

func generateWitness() (witness.Witness, error) {
	// Create input data
	input := make([]byte, inputLength)
	
	// Calculate Keccak256 hash
	hasher := sha3.NewLegacyKeccak256()
	hasher.Write(input)
	dgst := hasher.Sum(nil)

	// Create witness
	witness := keccakCircuit{
		In: uints.NewU8Array(input),
	}
	copy(witness.Expected[:], uints.NewU8Array(dgst))

	// Create witness data
	witnessData, err := frontend.NewWitness(&witness, ecc.BN254.ScalarField())
	if err != nil {
		return nil, err
	}
	return witnessData, nil
}

func Groth16Setup(fileDir string) {
	utils.Groth16Setup(fileDir, circuitName, compileCircuit)
}

func Groth16Prove(fileDir string) {
	utils.Groth16Prove(fileDir, circuitName, generateWitness)
} 