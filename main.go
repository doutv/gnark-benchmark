package main

import (
	"gnark-benchmark/ecdsa"
	"os"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "plonk" {
		ecdsa.PlonkSetup("")
		ecdsa.PlonkProveAndVerify("")
	} else {
		ecdsa.Setup("")
		ecdsa.ProveAndVerify("")
	}
}