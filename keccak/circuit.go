package keccak

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/sha3"
	"github.com/consensys/gnark/std/math/uints"
)

type keccakCircuit struct {
	In       []uints.U8 `gnark:",secret"`
	Expected [32]uints.U8 `gnark:",public"`
}

func (circuit *keccakCircuit) Define(api frontend.API) error {
	keccak, err := sha3.NewLegacyKeccak256(api)
	if err != nil {
		return err
	}
	uapi, err := uints.New[uints.U64](api)
	if err != nil {
		return err
	}

	keccak.Write(circuit.In)
	result := keccak.Sum()

	// Verify each byte of the hash
	for i := range circuit.Expected {
		uapi.ByteAssertEq(circuit.Expected[i], result[i])
	}
	return nil
} 