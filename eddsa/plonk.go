package eddsa

import (
	"gnark-benchmark/utils"
	"time"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/backend/plonk"
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/frontend/cs/scs"
	"github.com/consensys/gnark/test/unsafekzg"

	"log"
)

func PlonkSetup(fileDir string) {
	circuit := kycCircuit{
		Attributes: make([]frontend.Variable, 4),
	}
	r1cs, err := frontend.Compile(ecc.BN254.ScalarField(), scs.NewBuilder, &circuit)
	if err != nil {
		panic(err)
	}
	srs, srsLagrange, err := unsafekzg.NewSRS(r1cs)
	if err != nil {
		panic(err)
	}
	pk, vk, err := plonk.Setup(r1cs, srs, srsLagrange)
	if err != nil {
		panic(err)
	}
	// Write to file
	utils.WriteToFile(pk, fileDir+"ecdsa.plonk.zkey")
	utils.WriteToFile(r1cs, fileDir+"ecdsa.plonk.r1cs")
	utils.WriteToFile(vk, fileDir+"ecdsa.plonk.vkey")
}

func PlonkProveAndVerify(fileDir string) {
	proveStart := time.Now()
	witnessData, err := generateWitness()
	if err != nil {
		panic(err)
	}
	// Read files
	start := time.Now()
	r1cs := plonk.NewCS(ecc.BN254)
	utils.ReadFromFile(r1cs, fileDir+"ecdsa.plonk.r1cs")
	elapsed := time.Since(start)
	log.Printf("Read r1cs: %d ms", elapsed.Milliseconds())

	start = time.Now()
	pk := plonk.NewProvingKey(ecc.BN254)
	utils.UnsafeReadFromFile(pk, fileDir+"ecdsa.plonk.zkey")
	elapsed = time.Since(start)
	log.Printf("Read zkey: %d ms", elapsed.Milliseconds())

	// Proof generation
	start = time.Now()
	proof, err := plonk.Prove(r1cs, pk, witnessData)
	if err != nil {
		panic(err)
	}
	elapsed = time.Since(start)
	log.Printf("Prove: %d ms", elapsed.Milliseconds())

	proveElapsed := time.Since(proveStart)
	log.Printf("Total Prove time: %d ms", proveElapsed.Milliseconds())
	utils.WriteToFile(proof, fileDir+"ecdsa.plonk.proof")

	log.Println("start verify")
	publicWitness, err := witnessData.Public()
	if err != nil {
		panic(err)
	}
	vk := plonk.NewVerifyingKey(ecc.BN254)
	utils.ReadFromFile(vk, fileDir+"ecdsa.plonk.vkey")
	err = plonk.Verify(proof, vk, publicWitness)
	if err != nil {
		panic(err)
	}
	log.Println("end verify")
}