package mimc

import (
	"gnark-benchmark/utils"
	"math/big"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/hash"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend"
)

const (
	circuitName = "mimc"
	inputLength = 4 // You can adjust this based on your needs
)

func compileCircuit(newBuilder frontend.NewBuilder) (constraint.ConstraintSystem, error) {
	circuit := mimcCircuit{
		In: make([]frontend.Variable, inputLength),
	}
	r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), newBuilder, &circuit)
	if err != nil {
		return nil, err
	}
	return r1cs, nil
}

func generateWitness() (witness.Witness, error) {
	// Generate input data
	modulus := ecc.BN254.ScalarField()
	var data [inputLength]big.Int
	data[0].Sub(modulus, big.NewInt(1))
	for i := 1; i < inputLength; i++ {
		data[i].Add(&data[i-1], &data[i-1]).Mod(&data[i], modulus)
	}

	// Calculate MiMC hash
	goMimc := hash.MIMC_BN254.New()
	for i := 0; i < inputLength; i++ {
		goMimc.Write(data[i].Bytes())
	}
	expectedHash := goMimc.Sum(nil)

	// Create witness
	witness := mimcCircuit{
		In: make([]frontend.Variable, inputLength),
	}
	for i := 0; i < inputLength; i++ {
		witness.In[i] = data[i].String()
	}
	witness.Expected = expectedHash

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