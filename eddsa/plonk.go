package eddsa

import (
	"bytes"

	"github.com/consensys/gnark/backend/plonk"
	"github.com/consensys/gnark/frontend/cs/scs"
	"github.com/consensys/gnark/test/unsafekzg"

	"log"
)

func PlonkTest() {
	cs, witnessData, err := generateWitness(scs.NewBuilder)
	if err != nil {
		panic(err)
	}
	// 1. One time setup
	srs, srsLagrange, err := unsafekzg.NewSRS(cs)
	if err != nil {
		panic(err)
	}

	pk, vk, err := plonk.Setup(cs, srs, srsLagrange)

	if err != nil {
		panic(err)
	}

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
	proof, err := plonk.Prove(cs, pk, witnessData)
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
