package main

import (
	"gnark-benchmark/ecdsa"
	"gnark-benchmark/eddsa"
	"os"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "plonk" {
		println("------Plonk ECDSA Secp256k1------")
		ecdsa.PlonkSetup("")
		ecdsa.PlonkProveAndVerify("")
		println("------Plonk EdDSA BN254------")
		eddsa.PlonkSetup("")
		eddsa.PlonkProveAndVerify("")
	} else {
		println("------Groth16 ECDSA Secp256k1------")
		ecdsa.Groth16Setup("")
		ecdsa.Groth16ProveAndVerify("")
		println("------Groth16 EdDSA BN254------")
		eddsa.Groth16Setup("")
		eddsa.Groth16ProveAndVerify("")
	}
}