package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"log"
	"math/big"

	// btcutils "github.com/bithyve/research/utils"
	bech32 "github.com/bithyve/research/bech32"
	bitcoinrpc "github.com/bithyve/research/rpc"
	"github.com/btcsuite/btcd/btcec"
)

var Curve *btcec.KoblitzCurve = btcec.S256() // take only the curve, can't use other stuff

func Sha256(inputs ...[]byte) []byte {
	shaNew := sha256.New()
	for _, input := range inputs {
		shaNew.Write(input)
	}
	return shaNew.Sum(nil)
}

func NewPrivateKey() (*big.Int, error) {
	b := make([]byte, Curve.Params().BitSize/8+8)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		log.Fatal(err)
	}

	var one = new(big.Int).SetInt64(1)
	x := new(big.Int).SetBytes(b)
	n := new(big.Int).Sub(Curve.Params().N, one)
	x.Mod(x, n)
	x.Add(x, one)

	return x, nil
}

func PubkeyPointsFromPrivkey(privkey *big.Int) (*big.Int, *big.Int) {
	x, y := Curve.ScalarBaseMult(privkey.Bytes())
	return x, y
}

func testAddHomomorphic(P1x, P1y, P2x, P2y, x1, x2 *big.Int) error {
	Sumx, Sumy := Curve.Add(P1x, P1y, P2x, P2y)

	x1plusx2 := new(big.Int).Add(x1, x2)
	RSx, RSy := Curve.ScalarBaseMult(x1plusx2.Bytes())

	if Sumx.Cmp(RSx) != 0 && Sumy.Cmp(RSy) != 0 {
		log.Println("Additive homomorphiuc test failed")
		return errors.New("Additive homomorphiuc test failed")
	}
	return nil
}

func genCommitment() {
	x1, err := NewPrivateKey() // lets assume this to be the same as x
	if err != nil {
		log.Fatal(err)
	}

	x2, err := NewPrivateKey()
	if err != nil {
		log.Fatal(err)
	}

	P1x, P1y := PubkeyPointsFromPrivkey(x1)
	P2x, P2y := PubkeyPointsFromPrivkey(x2)

	err = testAddHomomorphic(P1x, P1y, P2x, P2y, x1, x2)
	if err != nil {
		log.Fatal(err)
	}

	shaBytes := Sha256(Curve.Params().Gx.Bytes(), Curve.Params().Gy.Bytes()) // SHA256(G)
	Hx, Hy := Curve.ScalarBaseMult(shaBytes)                                 // Point(SHA256(G))

	var a []byte
	a = []byte{1}

	aHx, aHy := Curve.ScalarMult(Hx, Hy, a)

	commitmentx, commitmenty := Curve.Add(P1x, P1y, aHx, aHy)
	log.Println("commitment: ", commitmentx, commitmenty)
}

func signCommitment() {
	x, err := NewPrivateKey() // lets assume this to be the same as x
	if err != nil {
		log.Fatal(err)
	}

	Px, Py := PubkeyPointsFromPrivkey(x) // P = x*G

	shaBytes := Sha256(Curve.Params().Gx.Bytes(), Curve.Params().Gy.Bytes()) // SHA256(G)
	Hx, Hy := Curve.ScalarBaseMult(shaBytes)                                 // Point(SHA256(G))

	// var a []byte
	// a = []byte{1}

	// aHx, aHy := Curve.ScalarMult(Hx, Hy, a)
	// Cx, Cy := Curve.Add(Px, Py, aHx, aHy)

	oneHx, oneHy := Curve.ScalarMult(Hx, Hy, []byte{1})

	Cprx, Cpry := Curve.Add(Px, Py, new(big.Int).Neg(oneHx), new(big.Int).Neg(oneHy))

	CprHash := Sha256(Cprx.Bytes(), Cpry.Bytes())
	CprHashHex := hex.EncodeToString(CprHash)

	privkey, err := bech32.PrivKeyToWIF("testnet", true, x.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	bitcoinrpc.SetBitcoindURL("http://localhost:18443/", "user", "password")
	sigBytes, err := bitcoinrpc.SignMessageWithPrivkey(privkey, CprHashHex)
	if err != nil {
		log.Println(err)
	}

	var sig bitcoinrpc.SignMessageWithPrivkeyReturn
	err = json.Unmarshal(sigBytes, &sig)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("SIG: ", sig.Result)
}

func main() {
	x, err := NewPrivateKey() // lets assume this to be the same as x
	if err != nil {
		log.Fatal(err)
	}

	k, err := NewPrivateKey() // lets assume this to be the same as x
	if err != nil {
		log.Fatal(err)
	}

	Px, Py := PubkeyPointsFromPrivkey(x) // P = x*G
	Kx, Ky := PubkeyPointsFromPrivkey(x) // K = x*G

	shaBytes := Sha256(Curve.Params().Gx.Bytes(), Curve.Params().Gy.Bytes()) // SHA256(G)
	Hx, Hy := Curve.ScalarBaseMult(shaBytes)                                 // Point(SHA256(G))

	oneHx, oneHy := Curve.ScalarMult(Hx, Hy, []byte{1})

	var a []byte
	a = []byte{1}

	aHx, aHy := Curve.ScalarMult(Hx, Hy, a)

	Cx, Cy := Curve.Add(Px, Py, aHx, aHy)                                           // P + a*H = commitment message
	C1x, C1y := Curve.Add(Px, Py, new(big.Int).Neg(oneHx), new(big.Int).Neg(oneHy)) // C1 = x*G - 1*H

	var m []byte
	m := append(m, Cx, Cy, C1x, C1y)
	E := []byte{2} // the second node in the ring

	eE := Sha256(Kx.Bytes(), Ky.Bytes(), m, E)
	se, err := NewPrivateKey()
	if err != nil {
		log.Fatal(err)
	}

	seGx, seGy := Curve.ScalarBaseMult(se.Bytes())
	eEPEx, eEPEy := Curve.ScalarMult(C1x, C1y, eE)

	KEx, KEy := Curve.Add(seGx, seGy, eEPEx, eEPEy) // C1's ring sig

	D := []byte{1}
	eD := Sha256(KEx.Bytes(), KEy.Bytes(), m, D)
	eDInt := new(big.Int).SetBytes(eD)
	sd := new(big.Int).Add(new(big.Int).Mul(eDInt, x), k)

	log.Println("P: ", Px, Py, "C1: ", C1x, C1y)
	log.Println("eE: ", new(big.Int).SetBytes(eE), "se: ", se, "ed: ", eDInt, "sd: ", sd)
}
