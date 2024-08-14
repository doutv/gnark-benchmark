package main

import (
	"gnark-benchmark/dummy1200k"
	"gnark-benchmark/ecdsa"
	"gnark-benchmark/eddsa"
	"os"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "plonk" {
		println("------Plonk ECDSA Secp256k1------")
		ecdsa.PlonkSetup("")
		ecdsa.PlonkProve("")
		println("------Plonk EdDSA BN254------")
		eddsa.PlonkSetup("")
		eddsa.PlonkProve("")
	} else {
		println("------Groth16 Dummy Circuit 1200k------")
		dummy1200k.Groth16Setup("")
		dummy1200k.Groth16Prove("")
		println("------Groth16 ECDSA Secp256k1------")
		ecdsa.Groth16Setup("")
		ecdsa.Groth16Prove("")
		println("------Groth16 EdDSA BN254------")
		eddsa.Groth16Setup("")
		eddsa.Groth16Prove("")
	}
}
