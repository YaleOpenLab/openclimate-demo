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

func Create21AOSSig() (*big.Int, *big.Int, *big.Int, *big.Int, *big.Int, *big.Int, *big.Int, []byte) {
	// Lets assume two parties - Brian and Dom.
	// Lets assume Brian and Dom's pubkeys are Pb and Pd. Lets assume their private keys
	// are b and d.
	// Then lets assume Dom wants to create a ring signature over C1 and C2 where
	// C1 = xG + 1H and C2 = xG where x is the blinding factor

	b, err := NewPrivateKey()
	if err != nil {
		log.Fatal(err)
	}

	d, err := NewPrivateKey()
	if err != nil {
		log.Fatal(err)
	}

	Pbx, Pby := PubkeyPointsFromPrivkey(b)
	Pdx, Pdy := PubkeyPointsFromPrivkey(d)

	x, err := NewPrivateKey()
	if err != nil {
		log.Fatal(err)
	}

	xGx, xGy := Curve.ScalarBaseMult(x.Bytes()) // xG

	shaBytes := Sha256(Curve.Params().Gx.Bytes(), Curve.Params().Gy.Bytes()) // SHA256(G)
	Hx, Hy := Curve.ScalarBaseMult(shaBytes)                                 // H = Point(SHA256(G))

	one := []byte{1}
	oneHx, oneHy := Curve.ScalarMult(Hx, Hy, one) // 1*H

	C1x, C1y := Curve.Add(xGx, xGy, oneHx, oneHy) // xG + 1H
	C2x, C2y := Curve.ScalarBaseMult(x.Bytes())   // xG

	var m []byte
	C1 := append(C1x.Bytes(), C1y.Bytes()...)
	C2 := append(C2x.Bytes(), C2y.Bytes()...)
	m = append(m, append(C1, C2...)...)

	kd, err := NewPrivateKey() // random nonce for ring sig
	if err != nil {
		log.Fatal(err)
	}

	Kdx, Kdy := Curve.ScalarBaseMult(kd.Bytes()) // K = kd*G

	BrianNodeNumber := []byte{2} // assume brian has node number 2

	eb := Sha256(Kdx.Bytes(), Kdy.Bytes(), m, BrianNodeNumber)

	sb, err := NewPrivateKey() // choose a signature sb at random fro brian
	if err != nil {
		log.Fatal(err)
	}

	sbGx, sbGy := Curve.ScalarBaseMult(sb.Bytes()) // sb*G
	ebPbx, ebPby := Curve.ScalarMult(Pbx, Pby, eb) // eb * Pb

	minusedPby := new(big.Int).Neg(ebPby)
	Kbx, Kby := Curve.Add(sbGx, sbGy, ebPbx, new(big.Int).Mod(minusedPby, Curve.P)) // Kb = sb*G - eb*Pb

	DomNodeNumber := []byte{1}

	ed := Sha256(Kbx.Bytes(), Kby.Bytes(), m, DomNodeNumber) // ed = H(Kb || m || D)
	// log.Println("ed: ", ed)

	edd := new(big.Int).Mul(new(big.Int).SetBytes(ed), d) // ed * d

	sd := new(big.Int).Add(edd, kd) // ed*d + kd

	return new(big.Int).SetBytes(eb), sb, sd, Pbx, Pby, Pdx, Pdy, m
}

func main() {
	eb, sb, sd, Pbx, Pby, Pdx, Pdy, m := Create21AOSSig()

	BrianNodeNumber := []byte{2} // assume brian has node number 2
	DomNodeNumber := []byte{1}

	sbGx, sbGy := Curve.ScalarBaseMult(sb.Bytes())         // sb*G
	ebPbx, ebPby := Curve.ScalarMult(Pbx, Pby, eb.Bytes()) // eb*Pb

	minusedPby := new(big.Int).Neg(ebPby)
	Kbx, Kby := Curve.Add(sbGx, sbGy, ebPbx, new(big.Int).Mod(minusedPby, Curve.P))

	ed := Sha256(Kbx.Bytes(), Kby.Bytes(), m, DomNodeNumber)
	// log.Println("ed: ", ed)

	sdGx, sdGy := Curve.ScalarBaseMult(sd.Bytes())
	edPdx, edPdy := Curve.ScalarMult(Pdx, Pdy, ed)

	minusedPdy := new(big.Int).Neg(edPdy)
	Kdx, Kdy := Curve.Add(sdGx, sdGy, edPdx, new(big.Int).Mod(minusedPdy, Curve.P))

	ebCheck := Sha256(Kdx.Bytes(), Kdy.Bytes(), m, BrianNodeNumber)
	ebCheckInt := new(big.Int).SetBytes(ebCheck)

	if ebCheckInt.Cmp(eb) != 0 {
		log.Fatal("Signatures don't match")
	} else {
		log.Println("Ring singatures validated")
	}
}
