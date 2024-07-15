package eddsa

import (
	"bytes"
	"fmt"
	"time"

	"github.com/consensys/gnark/frontend/cs/r1cs"

	"log"

	"github.com/consensys/gnark/backend/groth16"
)

func Groth16Test() {
	start := time.Now()

	cs, witnessData, err := generateWitness(r1cs.NewBuilder)
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

	log.Println("pk ", "nG1", pk.NbG1(), "nG2", pk.NbG2())
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

	fmt.Printf("setup %v\n", time.Since(start))
	start = time.Now()

	// 2. Proof creation
	proof, err := groth16.Prove(cs, pk, witnessData)
	if err != nil {
		panic(err)
	}

	fmt.Printf("prove %v\n", time.Since(start))
	// start = time.Now()

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
