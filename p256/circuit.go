package p256

import (
	"github.com/consensys/gnark/frontend"
	"github.com/consensys/gnark/std/algebra/emulated/sw_emulated"
	"github.com/consensys/gnark/std/math/emulated"
)

type EcdsaCircuit[T, S emulated.FieldParams] struct {
	// Commitment [32]uints.U8 `gnark:",public"` // Commit(Sig[0], Msg[0], Sig[1], Msg[1], ...)

	Sig [NumSignatures]Signature[S] `gnark:",secret"`
	Msg [NumSignatures]emulated.Element[S] `gnark:",secret"`
	Pub [NumSignatures]PublicKey[T, S] `gnark:",secret"`
}

func (c *EcdsaCircuit[T, S]) Define(api frontend.API) error {
	// Verify all ECDSA-P256 signatures
	for i := range c.Sig {
		c.Pub[i].Verify(api, sw_emulated.GetCurveParams[T](), &c.Msg[i], &c.Sig[i])
	}
	// SHA-3 (Keccak256) Commit to all signatures
	// h, err := sha3.New256(api)
	// if err != nil {
	// 	return err
	// }
	// uapi, err := uints.New[uints.U64](api)
	// if err != nil {
	// 	return err
	// }

	// for i := 0; i < NumSignatures; i++ {
	// 	h.Write(c.Sig[i].R)
	// 	h.Write(c.Sig[i].S)
	// 	h.Write(c.Msg[i])
	// }
	// res := h.Sum()

	// for i := range c.Commitment {
	// 	uapi.ByteAssertEq(c.Commitment[i], res[i])
	// }
	return nil
}

// Signature represents the signature for some message.
type Signature[Scalar emulated.FieldParams] struct {
	R, S emulated.Element[Scalar]
}

// PublicKey represents the public key to verify the signature for.
type PublicKey[Base, Scalar emulated.FieldParams] sw_emulated.AffinePoint[Base]

// Verify asserts that the signature sig verifies for the message msg and public
// key pk. The curve parameters params define the elliptic curve.
//
// We assume that the message msg is already hashed to the scalar field.
func (pk PublicKey[T, S]) Verify(api frontend.API, params sw_emulated.CurveParams, msg *emulated.Element[S], sig *Signature[S]) {
	cr, err := sw_emulated.New[T, S](api, params)
	if err != nil {
		panic(err)
	}
	scalarApi, err := emulated.NewField[S](api)
	if err != nil {
		panic(err)
	}
	baseApi, err := emulated.NewField[T](api)
	if err != nil {
		panic(err)
	}
	pkpt := sw_emulated.AffinePoint[T](pk)
	sInv := scalarApi.Inverse(&sig.S)
	msInv := scalarApi.MulMod(msg, sInv)
	rsInv := scalarApi.MulMod(&sig.R, sInv)

	// q = [rsInv]pkpt + [msInv]g
	q := cr.JointScalarMulBase(&pkpt, rsInv, msInv)
	qx := baseApi.Reduce(&q.X)
	qxBits := baseApi.ToBits(qx)
	rbits := scalarApi.ToBits(&sig.R)
	if len(rbits) != len(qxBits) {
		panic("non-equal lengths")
	}
	for i := range rbits {
		api.AssertIsEqual(rbits[i], qxBits[i])
	}
}