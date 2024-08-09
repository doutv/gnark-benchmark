package eddsa

import (
	"bytes"
	"fmt"
	"time"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend/cs/r1cs"

	"gnark-benchmark/utils"
	"log"

	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/frontend"
)

func Groth16Prove(fileDir string) {
	//setupstart

	//setup end

	gc1200k := DummyCircuit{}
	gw1200k := DummyCircuit{A: 3, C: generateMimcHash(3, 3636)}
	cs, witnessData, err := generateLoadWitness(r1cs.NewBuilder, &gc1200k, &gw1200k)

	r1cs := groth16.NewCS(ecc.BN254)
	utils.ReadFromFile(r1cs, fileDir+"eddsa.r1cs")

	pk := groth16.NewProvingKey(ecc.BN254)

	utils.UnsafeReadFromFile(pk, fileDir+"eddsa.zkey")

	var pkbuffer bytes.Buffer
	pkn, err := pk.WriteTo(&pkbuffer)
	if err != nil {
		panic(err)
	}
	var r1csbuffer bytes.Buffer
	r1csn, err := cs.WriteTo(&r1csbuffer)
	if err != nil {
		panic(err)
	}

	log.Printf("end setup. size: %vmb, pk: %vmb constrain: %v mb", (float64(pkn+r1csn))/(1024.0*1024), (float64(pkn))/(1024.0*1024), (float64(r1csn))/(1024.0*1024))

	// 2. Proof creation
	groth16.Prove(cs, pk, witnessData)
	if err != nil {
		panic(err)
	}

	// start = time.Now()

	log.Println("end proof")

	// log.Println("start verify")
	// publicWitness, err := witnessData.Public()
	// if err != nil {
	// 	panic(err)
	// }
	// // 3. Proof verification
	// err = groth16.Verify(proof, vk, publicWitness)
	// if err != nil {
	// 	panic(err)
	// }
	// log.Println("end verify")
}

func generateMimcHash(seed uint64, number int) []byte {
	var hFunc = mimc.NewMiMC()

	var t fr.Element
	t.SetUint64(seed)
	var b = t.Bytes()
	var bb = b[:]
	for i := 0; i < number; i++ {
		hFunc.Reset()
		hFunc.Write(bb)
		bb = hFunc.Sum(nil)
	}

	return bb
}
func generateLoadWitness(newBuilder frontend.NewBuilder, circuit frontend.Circuit, witness frontend.Circuit) (
	constraint.ConstraintSystem, witness.Witness, error) {

	witnessData, err := frontend.NewWitness(witness, ecc.BN254.ScalarField())
	if err != nil {
		panic(err)
	}

	r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), newBuilder, circuit)
	if err != nil {
		panic(err)
	}

	return r1cs, witnessData, nil
}
func Groth16Setup(fileDir string) {
	start := time.Now()

	// cs, witnessData, err := generateWitness(r1cs.NewBuilder)
	gc1200k := DummyCircuit{}
	gw1200k := DummyCircuit{A: 3, C: generateMimcHash(3, 3636)}
	cs, _, err := generateLoadWitness(r1cs.NewBuilder, &gc1200k, &gw1200k)

	if err != nil {
		panic(err)
	}
	// fmt.Printf("%+v\n", witnessData)
	fmt.Printf("generate witness %v\n", time.Since(start))
	start = time.Now()
	// 1. One time setup
	pk, vk, err := groth16.Setup(cs)
	if err != nil {
		panic(err)
	}
	// Write to file
	utils.WriteToFile(pk, fileDir+"eddsa.zkey")
	utils.WriteToFile(cs, fileDir+"eddsa.r1cs")
	utils.WriteToFile(vk, fileDir+"eddsa.vkey")
}
