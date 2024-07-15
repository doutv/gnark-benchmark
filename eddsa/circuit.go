package eddsa

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/hash/mimc"
	"github.com/consensys/gnark/std/selector"
	"github.com/consensys/gnark/std/signature/eddsa"

	tedwards "github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark/std/algebra/native/twistededwards"
)

type KycCircuit struct {
	Attributes     []frontend.Variable `gnark:"fields,secret"`
	Expire         frontend.Variable   `gnark:"expire,secret"`
	ClaimAttribute frontend.Variable   `gnark:",secret"`
	ClaimOperator  frontend.Variable   `gnark:",secret"`
	ClaimValue     frontend.Variable   `gnark:",secret"`

	// public witness fields
	Timestamp frontend.Variable `gnark:",public"`
	Address   frontend.Variable `gnark:",public"`

	// public inputs
	UserId        frontend.Variable `gnark:",public"`
	ClaimHash     frontend.Variable `gnark:",public"`
	PublicKeyHash frontend.Variable `gnark:",public"`

	//signature
	PublicKey eddsa.PublicKey `gnark:",public"`
	Signature eddsa.Signature `gnark:",public"`
}

func (circuit *KycCircuit) Define(api frontend.API) error {
	hFunc, err := mimc.NewMiMC(api)
	if err != nil {
		return err
	}

	// signature verify
	hFunc.Write(circuit.UserId)
	for i := 0; i < len(circuit.Attributes); i++ {
		hFunc.Write(circuit.Attributes[i])
	}
	message_hash := hFunc.Sum()

	curve, err := twistededwards.NewEdCurve(api, tedwards.BN254)
	if err != nil {
		return err
	}

	hFunc.Reset()
	eddsa.Verify(curve, circuit.Signature, message_hash, circuit.PublicKey, &hFunc)

	// check claim hash
	hFunc.Reset()
	hFunc.Write(circuit.ClaimAttribute, circuit.ClaimOperator, circuit.ClaimValue)
	api.AssertIsEqual(circuit.ClaimHash, hFunc.Sum())

	// check claim
	// Operators:
	// 0: equal
	// 1: not equal
	// 2: less than
	// 3: greater than
	selectedValue := selector.Mux(api, circuit.ClaimAttribute, circuit.Attributes...)

	compareResult := api.Cmp(selectedValue, circuit.ClaimValue)
	o0 := api.IsZero(compareResult)
	o1 := api.Sub(1, o0)
	o2 := api.IsZero(api.Add(compareResult, 1))
	o3 := api.IsZero(api.Sub(compareResult, 1))

	claimResult := selector.Mux(api, circuit.ClaimOperator, []frontend.Variable{o0, o1, o2, o3}...)

	api.AssertIsEqual(claimResult, 1)

	// check public key hash
	hFunc.Reset()
	hFunc.Write(circuit.PublicKey.A.X, circuit.PublicKey.A.Y)
	// pubkeyHash := hFunc.Sum()
	// api.Println("public key hash", pubkeyHash)
	// api.AssertIsEqual(circuit.PublicKeyHash, hFunc.Sum())

	return nil
}
