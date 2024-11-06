package mimc

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
)

type mimcCircuit struct {
	In       []frontend.Variable
	Expected frontend.Variable `gnark:",public"`
}

func (circuit *mimcCircuit) Define(api frontend.API) error {
	mimc, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}
	mimc.Write(circuit.In...)
	result := mimc.Sum()
	api.AssertIsEqual(result, circuit.Expected)
	return nil
} 