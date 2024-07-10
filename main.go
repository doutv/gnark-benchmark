package main

import (
	"gnark-benchmark/ecdsa"
	"os"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "full" {
		ecdsa.Setup("")
	}
	ecdsa.ProveAndVerify("")
}