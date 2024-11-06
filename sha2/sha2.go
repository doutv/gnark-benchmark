package sha2

import (
	"crypto/sha256"
	"gnark-benchmark/utils"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
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
	utils.Groth16Setup(fileDir, circuitName, compileCircuit)
}

func Groth16Prove(fileDir string) {
	utils.Groth16Prove(fileDir, circuitName, generateWitness)
}