package main

import (
	"gnark-benchmark/eddsa"
	"os"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "plonk" {
		println("------Plonk ECDSA Secp256k1------")
		// ecdsa.PlonkSetup("")
		// ecdsa.PlonkProveAndVerify("")
		// println("------Plonk EdDSA BN254------")
		// eddsa.PlonkSetup("")
		// eddsa.PlonkProveAndVerify("")
	} else {

		// println("------Groth16 ECDSA Secp256k1------")
		// ecdsa.groth16Setup("")
		// ecdsa.Groth16Prove("")
		println("------Groth16 EdDSA BN254------")
		// eddsa.groth16Setup("")
		eddsa.Groth16Prove("")
	}
}
