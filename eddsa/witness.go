package eddsa

import (
	"fmt"
	"math/big"
	"math/rand"
	"time"

	"github.com/consensys/gnark-crypto/ecc/bn254/fr/mimc"
	"github.com/consensys/gnark/std/algebra/native/twistededwards"

	"github.com/consensys/gnark-crypto/ecc"
	tedwards "github.com/consensys/gnark-crypto/ecc/twistededwards"
	"github.com/consensys/gnark-crypto/hash"
	"github.com/consensys/gnark-crypto/signature/eddsa"
	"github.com/consensys/gnark/backend/witness"

	"github.com/consensys/gnark/frontend"
)

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}
func generateWitness(attribute int64, op int64, value int64) (witness.Witness, error) {
	//age,gender,natioinality
	//gender:0 female,1 male
	//nationality,0 America,1 China
	//user info先写死: 18,1,0
	attributes := []int{18, 1, 0}
	fields := make([][]byte, len(attributes))

	curve := tedwards.BN254
	hashType := hash.MIMC_BN254
	snarkField, err := twistededwards.GetSnarkField(curve)
	panicIfErr(err)

	userId := *big.NewInt(1)
	var hFunc = mimc.NewMiMC()

	userIdBytes := convertToBytes(userId, snarkField)
	hFunc.Write(userIdBytes)
	for i := 0; i < len(attributes); i++ {
		msg := *big.NewInt(int64(attributes[i]))
		msgData := convertToBytes(msg, snarkField)
		fields[i] = msgData
		hFunc.Write(msgData)
	}

	hashMsg := hFunc.Sum(nil)

	var fieldsInVar []frontend.Variable
	for _, d := range fields {
		fieldsInVar = append(fieldsInVar, d)
	}

	// generate parameters for the signatures
	seed := time.Now().Unix()
	randomness := rand.New(rand.NewSource(seed)) //#nosec G404 -- This is a false positive
	privKey, err := eddsa.New(tedwards.BN254, randomness)
	panicIfErr(err)

	// generate signature
	signature, err := privKey.Sign(hashMsg, hashType.New())
	panicIfErr(err)

	// check if there is no problem in the signature
	pubKey := privKey.Public()
	checkSig, err := pubKey.Verify(signature, hashMsg, hashType.New())
	panicIfErr(err)
	// assert.True(checkSig, "signature verification failed")
	fmt.Println(checkSig)

	// calculate claimHash
	claimAttributeBytes := convertToBytes(*big.NewInt(attribute), snarkField)
	claimOperatorBytes := convertToBytes(*big.NewInt(op), snarkField)
	claimValueBytes := convertToBytes(*big.NewInt(value), snarkField)

	hFunc.Reset()
	hFunc.Write(claimAttributeBytes)
	hFunc.Write(claimOperatorBytes)
	hFunc.Write(claimValueBytes)

	claimHash := hFunc.Sum(nil)

	// calculate pubkeyHash
	hFunc.Reset()
	hFunc.Write(pubKey.Bytes())

	pubkeyHash := hFunc.Sum(nil)

	witnessCircuit := kycCircuit{
		Attributes:     fieldsInVar,
		Expire:         userIdBytes,
		ClaimAttribute: claimAttributeBytes,
		ClaimOperator:  claimOperatorBytes,
		ClaimValue:     claimValueBytes,
		Timestamp:      userIdBytes,
		Address:        userIdBytes,
		UserId:         userIdBytes,
		ClaimHash:      claimHash,
		PublicKeyHash:  pubkeyHash,
	}
	witnessCircuit.PublicKey.Assign(curve, pubKey.Bytes())
	witnessCircuit.Signature.Assign(curve, signature)

	witnessData, err := frontend.NewWitness(&witnessCircuit, ecc.BN254.ScalarField())
	panicIfErr(err)

	return witnessData, nil
}

func convertToBytes(msg big.Int, snarkField *big.Int) []byte {
	msgDataUnpadded := msg.Bytes()
	msgData := make([]byte, len(snarkField.Bytes()))
	copy(msgData[len(msgData)-len(msgDataUnpadded):], msgDataUnpadded)
	return msgData
}
