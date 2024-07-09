package ecdsa

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"hash"
	"time"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	secp_mimc "github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	secp_ecdsa "github.com/consensys/gnark-crypto/ecc/secp256k1/ecdsa"
	"github.com/consensys/gnark/backend/plonk"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend/cs/r1cs"
	"github.com/consensys/gnark/frontend/cs/scs"
	"github.com/consensys/gnark/test/unsafekzg"

	"log"
	"math/big"

	"github.com/consensys/gnark/backend/groth16"
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

type KycCredential struct {
	Credential uint64
	Age        uint64
	Gender     uint64
	Nation     uint64
	Expirtime  uint64
}

// Sign signs a transaction
func (t *KycCredential) Sign(priv *secp_ecdsa.PrivateKey, h hash.Hash) (secp_ecdsa.Signature, []byte, error) {

	h.Reset()

	////var frNonce, msg fr.Element
	var credential fr.Element
	credential.SetUint64(t.Credential)
	b := credential.Bytes()
	_, _ = h.Write(b[:])

	var age fr.Element
	age.SetUint64(t.Age)
	b = age.Bytes()
	_, _ = h.Write(b[:])

	var gender fr.Element
	gender.SetUint64(t.Gender)
	b = gender.Bytes()
	_, _ = h.Write(b[:])

	var nation fr.Element
	nation.SetUint64(t.Nation)
	b = nation.Bytes()
	_, _ = h.Write(b[:])

	var expirTime fr.Element
	expirTime.SetUint64(t.Expirtime)
	b = expirTime.Bytes()
	_, _ = h.Write(b[:])

	msg := h.Sum(nil)

	sigBin, err := priv.Sign(msg, nil)
	if err != nil {
		return secp_ecdsa.Signature{}, nil, err
	}
	var sig secp_ecdsa.Signature
	if _, err := sig.SetBytes(sigBin); err != nil {
		return secp_ecdsa.Signature{}, nil, err
	}
	return sig, msg, nil
}

var hFunc = secp_mimc.NewMiMC()

func generateWitness(newBuilder frontend.NewBuilder) (constraint.ConstraintSystem, witness.Witness, error) {
	// generate parameters
	privKey, _ := secp_ecdsa.GenerateKey(rand.Reader)

	// sign
	credential := KycCredential{Credential: 12, Age: 18, Gender: 1, Nation: 0b10, Expirtime: 123}
	sigBin, _, err := credential.Sign(privKey, hFunc)
	if err != nil {
		panic(err)
	}

	// unmarshal signature
	r, s := new(big.Int), new(big.Int)
	r.SetBytes(sigBin.R[:32])
	s.SetBytes(sigBin.S[:32])

	circuit := KycCircuit{}
	witnessCircuit := KycCircuit{
		Signature: ecdsa.Signature[emulated.Secp256k1Fr]{
			R: emulated.ValueOf[emulated.Secp256k1Fr](r),
			S: emulated.ValueOf[emulated.Secp256k1Fr](s),
		},
		Credential: credential.Credential, Age: credential.Age, Gender: credential.Gender, Nation: credential.Nation, ExpireTime: credential.Expirtime,
		MinAge: 1,
		MaxAge: 60,
		//isMale: 1,
		ContainNations: 0b10011,
		PublicKey: ecdsa.PublicKey[emulated.Secp256k1Fp, emulated.Secp256k1Fr]{
			X: emulated.ValueOf[emulated.Secp256k1Fp](privKey.PublicKey.A.X),
			Y: emulated.ValueOf[emulated.Secp256k1Fp](privKey.PublicKey.A.Y),
		},
	}

	witnessData, err := frontend.NewWitness(&witnessCircuit, ecc.BN254.ScalarField())
	if err != nil {
		panic(err)
	}

	r1cs_temp, err := frontend.Compile(ecc.BN254.ScalarField(), newBuilder, &circuit)

	schema, _ := frontend.NewSchema(&witnessCircuit)
	ret, _ := witnessData.ToJSON(schema)

	var b bytes.Buffer
	json.Indent(&b, ret, "", "\t")
	//log.Println("start proof: witness", b.String())
	return r1cs_temp, witnessData, nil
}

func PlonkTest() {
	r1cs_temp, witnessData, err := generateWitness(scs.NewBuilder)
	// 1. One time setup
	srs, srsLagrange, err := unsafekzg.NewSRS(r1cs_temp)

	pk, vk, err := plonk.Setup(r1cs_temp, srs, srsLagrange)

	if err != nil {
		panic(err)
	}

	var pkbuffer bytes.Buffer
	pkn, err := pk.WriteTo(&pkbuffer)
	if err != nil {
		panic(err)
	}
	var r1csbuffer bytes.Buffer
	r1csn, err := r1cs_temp.WriteTo(&r1csbuffer)
	if err != nil {
		panic(err)
	}

	log.Printf("end setup. size: %vmb, pk: %vmb constrain: %v mb", (float64(pkn+r1csn))/(1024.0*1024), (float64(pkn))/(1024.0*1024), (float64(r1csn))/(1024.0*1024))

	// 2. Proof creation
	proof, err := plonk.Prove(r1cs_temp, pk, witnessData)
	if err != nil {
		panic(err)
	}

	log.Println("end proof")

	log.Println("start verify")
	publicWitness, err := witnessData.Public()
	if err != nil {
		panic(err)
	}
	// 3. Proof verification
	err = plonk.Verify(proof, vk, publicWitness)
	if err != nil {
		panic(err)
	}
	log.Println("end verify")
}

func Groth16Test() {
	start := time.Now()
	r1cs_temp, witnessData, err := generateWitness(r1cs.NewBuilder)
	elapsed := time.Since(start)
	log.Printf("Witness Generation: %d ms", elapsed.Milliseconds())
	if err != nil {
		panic(err)
	}
	// 1. One time setup
	log.Println("start setup")
	pk, vk, err := groth16.Setup(r1cs_temp)
	if err != nil {
		panic(err)
	}
	log.Println("pk ", "nG1", pk.NbG1(), "nG2", pk.NbG2())
	var pkbuffer bytes.Buffer
	pkn, err := pk.WriteTo(&pkbuffer)
	if err != nil {
		panic(err)
	}
	var r1csbuffer bytes.Buffer
	r1csn, err := r1cs_temp.WriteTo(&r1csbuffer)
	if err != nil {
		panic(err)
	}

	log.Printf("end setup. size: %vmb, pk: %vmb constrain: %v mb", (float64(pkn+r1csn))/(1024.0*1024), (float64(pkn))/(1024.0*1024), (float64(r1csn))/(1024.0*1024))

	// 2. Proof creation
	proof, err := groth16.Prove(r1cs_temp, pk, witnessData)
	if err != nil {
		panic(err)
	}

	log.Println("end proof")

	log.Println("start verify")
	publicWitness, err := witnessData.Public()
	if err != nil {
		panic(err)
	}
	// 3. Proof verification
	err = groth16.Verify(proof, vk, publicWitness)
	if err != nil {
		panic(err)
	}
	log.Println("end verify")
}
