package main

import (
	"gnark-benchmark/p256"
	"gnark-benchmark/sha2"
	"gnark-benchmark/sha3"
)

func main() {
	p256.Groth16Setup("build/p256/")
	p256.Groth16Prove("build/p256/")

	sha3.Groth16Setup("build/sha3/")
	sha3.Groth16Prove("build/sha3/")
	sha2.Groth16Setup("build/sha2/")
	sha2.Groth16Prove("build/sha2/")

	// attributes, err := json.Marshal(eddsa.Attributes{Attributes: []int{1, 2, 3}})
	// if err != nil {
	// 	panic(err)
	// }
	// credential, err := json.Marshal(ecdsa.KycCredential{
	// 	Credential: 12,
	// 	Age:        18,
	// 	Gender:     1,
	// 	Nation:     0b10,
	// 	ExpireTime: 123,
	// })
	// if err != nil {
	// 	panic(err)
	// }
	// if len(os.Args) > 1 && os.Args[1] == "plonk" {
	// 	println("------Plonk ECDSA Secp256k1------")
	// 	ecdsa.PlonkSetup("")
	// 	ecdsa.PlonkProve("", credential)
	// 	println("------Plonk EdDSA BN254------")
	// 	eddsa.PlonkSetup("")
	// 	eddsa.PlonkProve("", attributes)
	// } else {
	// 	println("------Groth16 Dummy Circuit 1200k------")
	// 	dummy1200k.Groth16Setup("")
	// 	dummy1200k.Groth16Prove("")
	// 	println("------Groth16 ECDSA Secp256k1------")
	// 	ecdsa.Groth16Setup("")
	// 	ecdsa.Groth16Prove("", credential)
	// 	println("------Groth16 EdDSA BN254------")
	// 	eddsa.Groth16Setup("")
	// 	eddsa.Groth16Prove("", attributes)
	// }
}
