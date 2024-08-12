package eddsa

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
)

type DummyCircuit struct {
	A frontend.Variable `gnark:"a,secret"`
	C frontend.Variable `gnark:"c,public"`
}

func (circuit *DummyCircuit) Define(api frontend.API) error {
	f, _ := mimc.NewMiMC(api)
	h := circuit.A
	// 3636 \approx 1.2M constraints
	for i := 0; i < 1000; i++ {
		f.Reset()
		f.Write(h)
		h = f.Sum()
	}
	api.AssertIsEqual(h, circuit.C)

	return nil
}
