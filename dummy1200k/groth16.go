package dummy1200k

import (
	"gnark-benchmark/utils"
	"time"

	"github.com/consensys/gnark-crypto/ecc"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark/frontend/cs/r1cs"

	"log"

	"github.com/consensys/gnark/backend/groth16"
	"github.com/consensys/gnark/frontend"
)

const circuitName = "dummy1200k"

func Groth16Setup(fileDir string) {
	circuit := dummyCircuit{}
	cs, err := frontend.Compile(ecc.BN254.ScalarField(), r1cs.NewBuilder, &circuit)
	if err != nil {
		panic(err)
	}
	pk1, vk, err := groth16.Setup(cs)
	if err != nil {
		panic(err)
	}
	// Write to file
	utils.WriteToFile(pk1, fileDir+circuitName+".zkey")
	utils.WriteToFile(cs, fileDir+circuitName+".r1cs")
	utils.WriteToFile(vk, fileDir+circuitName+".vkey")
}

func Groth16Prove(fileDir string) {
	proveStart := time.Now()

	// Witness generation
	start := time.Now()
	gw1200k := dummyCircuit{A: 3, C: generateMimcHash(3, 3636)}
	witnessData, err := frontend.NewWitness(&gw1200k, ecc.BN254.ScalarField())
	if err != nil {
		panic(err)
	}
	elapsed := time.Since(start)
	log.Printf("Witness Generation: %d ms", elapsed.Milliseconds())

	// Read files
	start = time.Now()
	r1cs := groth16.NewCS(ecc.BN254)
	utils.ReadFromFile(r1cs, fileDir+circuitName+".r1cs")
	elapsed = time.Since(start)
	log.Printf("Read r1cs: %d ms", elapsed.Milliseconds())

	start = time.Now()
	pk := groth16.NewProvingKey(ecc.BN254)
	utils.UnsafeReadFromFile(pk, fileDir+circuitName+".zkey")
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

	utils.WriteToFile(proof, fileDir+circuitName+".proof")
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
