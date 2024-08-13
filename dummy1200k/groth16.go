package main

import (
	"fmt"
	"time"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark/constraint"
	"github.com/consensys/gnark/frontend/cs/r1cs"

	"log"

	"io"
	"os"

	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/backend/witness"
	"github.com/consensys/gnark/frontend"
)

func ReadFromFile(data io.ReaderFrom, fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	// Use the ReadFrom method to read the file's content into data.
	if _, err := data.ReadFrom(file); err != nil {
		panic(err)
	}
}
func WriteToFile(data io.WriterTo, fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = data.WriteTo(file)
	if err != nil {
		panic(err)
	}
}

func Groth16Test() {
	// setup

	circuit := DummyCircuit{}
	cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	pk1, vk, err := groth16.Setup(cs)
	if err != nil {
		panic(err)
	}
	// Write to file
	WriteToFile(pk1, "dummy.zkey")
	WriteToFile(cs, "dummy.r1cs")
	WriteToFile(vk, "dummy.vkey")

	//setup end

	r1cs := groth16.NewCS(ecc.BN254)
	ReadFromFile(r1cs, "dummy.r1cs")

	pk := groth16.NewProvingKey(ecc.BN254)
	ReadFromFile(pk, "dummy.zkey")

	gw1200k := DummyCircuit{A: 3, C: generateMimcHash(3, 3636)}

	witnessData, _ := frontend.NewWitness(&gw1200k, ecc.BN254.ScalarField())

	start := time.Now()

	// 2. Proof creation
	_, err = groth16.Prove(r1cs, pk, witnessData)
	if err != nil {
		panic(err)
	}

	fmt.Printf("prove %v\n", time.Since(start))
	// start = time.Now()

	log.Println("end proof")

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

	cs, err := frontend.Compile(ecc.BN254.ScalarField(), newBuilder, circuit)
	if err != nil {
		panic(err)
	}

	return cs, witnessData, nil
}
