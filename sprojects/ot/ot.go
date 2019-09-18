package main

import (
	"errors"
	"log"
	"math/big"

	btcutils "github.com/bithyve/research/utils"
)

var Curve = btcutils.Curve

func AliceSign(T btcutils.Point, a *big.Int, A, H btcutils.Point, m0, m1 []byte) (*big.Int, btcutils.Point, *big.Int, btcutils.Point) {

	c := new(big.Int).SetInt64(1)
	cH := btcutils.ScalarMult(H, c.Bytes())

	T0 := T
	T1 := btcutils.Sub(T0, cH)

	s0, R0 := aliceObliviousSig(T0, a, A, m0)
	s1, R1 := aliceObliviousSig(T1, a, A, m1)

	return s0, R0, s1, R1
}

func BobVerify(s0, s1 *big.Int, P, R0, R1, T, H, A btcutils.Point, m0, m1 []byte) error {
	c := new(big.Int).SetInt64(1)
	cH := btcutils.ScalarMult(H, c.Bytes())

	T0 := T
	T1 := btcutils.Sub(T0, cH)

	if !verifyAdaptorSig(s0, A, R0, T0, m0) {
		return errors.New("adaptor sig verification failed")
	}
	if !verifyAdaptorSig(s1, A, R1, T1, m1) {
		return errors.New("adaptor sig verification failed")
	}

	return nil
}

func aliceObliviousSig(T btcutils.Point, a *big.Int, A btcutils.Point, m []byte) (*big.Int, btcutils.Point) {
	r, err := btcutils.NewPrivateKey()
	if err != nil {
		log.Fatal(err)
	}

	R := btcutils.PointFromPrivkey(r)
	RplusT := btcutils.Add(R, T) // R+T

	HPRTm := btcutils.Sha256(A.Bytes(), RplusT.Bytes(), m)    // H(P||R+T||m)
	HPRTmx := new(big.Int).Mul(btcutils.BytesToNum(HPRTm), a) // H(P||R+T||m) * x
	s := new(big.Int).Add(r, HPRTmx)                          // s = r + H(P||R+T||m) * x
	return s, R
}

func verifyAdaptorSig(s *big.Int, P, R, T btcutils.Point, m []byte) bool {
	sG := btcutils.ScalarBaseMult(s.Bytes())

	RplusT := btcutils.Add(R, T)

	HPRTm := btcutils.Sha256(P.Bytes(), RplusT.Bytes(), m) // H(P||R+T||m)

	HPRTmP := btcutils.ScalarMult(P, HPRTm) // H(P||R+T||m) * P

	RplusHPRTmP := btcutils.Add(R, HPRTmP) // R + H(P||R+T||m) * P
	// s*G == R + H(P||R+T||m) * P
	return sG.Cmp(RplusHPRTmP)
}

func main() {
	/*
		The idea behind Oblivious Transfers is quite simple. Assume two parties Alice and Bob.
		An oblivious transfer is where Alice transmits 1/n messages to Bob of Bob's choosing but
		Alice doesn't learn which one they transmitted (Bob knows)

		This idea is somewhat remotely related to AOS sigs where a party that sees n signatures doesn't
		know which one is the one but knows that one among those is the real signature (aka ring sigs)

		The way this is going to work is illustrated in the talk at https://www.youtube.com/watch?time_continue=10377&v=-gdfxNalDIc
		but these are the steps we're going to do:

		- Bob has a choice for c - 0 or 1
		- Bob chooses a random y
		- Bob transmits a Pedersen commitment T = b*G + c*H to Alice
		- Alice signs two transactions - T, Th^-1 and sends them to Bob
		- Bob verifies the adaptor signatures, completes one of them and transmits that to the blockchain

		Everything above is done with the help of adaptor signatures
	*/

	b, err := btcutils.NewPrivateKey()
	if err != nil {
		log.Fatal(err)
	}

	var H btcutils.Point
	shaBytes := btcutils.Sha256(Curve.Params().Gx.Bytes(), Curve.Params().Gy.Bytes()) // btcutils.Sha256(G)
	H.Set(Curve.ScalarBaseMult(shaBytes))                                             // H = btcutils.Point(btcutils.Sha256(G))

	bG := btcutils.ScalarBaseMult(b.Bytes())

	T := bG // ie we choose c as 0 in this case

	m0 := []byte("Alice wins")
	m1 := []byte("Alice loses")

	// Generate Alice's private key
	a, err := btcutils.NewPrivateKey()
	if err != nil {
		log.Fatal(err)
	}

	A := btcutils.PointFromPrivkey(a)

	s0, R0, s1, R1 := AliceSign(T, a, A, H, m0, m1)

	BobVerify(s0, s1, A, R0, R1, T, H, A, m0, m1)

	log.Println("Oblivious Signing test successful")
	// now Bob should verify these adaptor signatures
}
