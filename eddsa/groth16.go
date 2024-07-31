package eddsa

import (
	"gnark-benchmark/utils"
	"time"

	"os"
	"path/filepath"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark/frontend/cs/r1cs"

	"log"

	"github.com/consensys/gnark/backend/groth16"
)

func groth16Setup(fileDir string) {
	r1cs, err := compileCircuit(r1cs.NewBuilder)
	if err != nil {
		panic(err)
	}
	pk, vk, err := groth16.Setup(r1cs)
	if err != nil {
		panic(err)
	}
	// Write to file
	utils.WriteToFile(pk, fileDir+"eddsa.zkey")
	utils.WriteToFile(r1cs, fileDir+"eddsa.r1cs")
	utils.WriteToFile(vk, fileDir+"eddsa.vkey")
}

func Groth16Prove(fileDir string, attribute int64, op int64, value int64) {
	proveStart := time.Now()
	// Witness generation
	start := time.Now()
	witnessData, err := generateWitness(attribute, op, value)
	if err != nil {
		panic(err)
	}
	elapsed := time.Since(start)
	log.Printf("Witness Generation: %d ms", elapsed.Milliseconds())

	// Read files
	start = time.Now()
	r1cs := groth16.NewCS(ecc.BN254)
	utils.ReadFromFile(r1cs, filepath.Join(os.Getenv("HOME"), "Documents", "eddsa.r1cs"))
	elapsed = time.Since(start)
	log.Printf("Read r1cs: %d ms", elapsed.Milliseconds())

	start = time.Now()
	pk := groth16.NewProvingKey(ecc.BN254)
	utils.ReadFromFile(pk, filepath.Join(os.Getenv("HOME"), "Documents", "eddsa.zkey"))
	elapsed = time.Since(start)
	log.Printf("Read zkey: %d ms", elapsed.Milliseconds())

	// Proof generation
	start = time.Now()
	proof, err := groth16.Prove(r1cs, pk, witnessData)
	if err != nil {
		panic(err)
	}
	elapsed = time.Since(start)
	log.Printf("Prove: %d ms", elapsed.Milliseconds())

	proveElapsed := time.Since(proveStart)
	log.Printf("Total Prove time: %d ms", proveElapsed.Milliseconds())

	utils.WriteToFile(proof, filepath.Join(os.Getenv("HOME"), "Documents", "eddsa.proof"))

	// Proof verification
	// publicWitness, err := witnessData.Public()
	// if err != nil {
	// 	panic(err)
	// }
	// vk := groth16.NewVerifyingKey(ecc.BN254)
	// utils.ReadFromFile(vk, fileDir+"eddsa.vkey")
	// err = groth16.Verify(proof, vk, publicWitness)
	// if err != nil {
	// 	panic(err)
	// }

}

// func readFile(filename string) ([]byte, error) {
// 	path := filepath.Join(os.Getenv("HOME"), "Documents", filename)
// 	data, err := os.ReadFile(path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return data, nil
// }
