package ecdsa

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/algebra/emulated/sw_emulated"
	"github.com/consensys/gnark/std/hash/mimc"
	"github.com/consensys/gnark/std/math/emulated"
	"github.com/consensys/gnark/std/signature/ecdsa"
)

type KycCircuit struct {
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

func (circuit *KycCircuit) Define(api frontend.API) error {
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
