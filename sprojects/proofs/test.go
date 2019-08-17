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
	utils "github.com/Varunram/essentials/utils"
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
	// Lets assume Brian and Dom's P are Pb and Pd. Lets assume their private keys
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

func Verify21AOSSig() {
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

func SubtractOnCurve(e, PxByte, PyByte []byte) ([]byte, []byte) {
	s, err := NewPrivateKey()
	if err != nil {
		log.Fatal(err)
	}

	sGx, sGy := Curve.ScalarBaseMult(s.Bytes())

	Px, Py := new(big.Int).SetBytes(PxByte), new(big.Int).SetBytes(PyByte)

	ePx, ePy := Curve.ScalarMult(Px, Py, e)

	minusedPy := new(big.Int).Neg(ePy)

	xX, xY := Curve.Add(sGx, sGy, ePx, new(big.Int).Mod(minusedPy, Curve.P))
	return append(xX.Bytes(), xY.Bytes()...), s.Bytes()
}

func SubtractOnCurveS(e []byte, Px *big.Int, Py *big.Int, s []byte) []byte {
	sGx, sGy := Curve.ScalarBaseMult(s)

	ePx, ePy := Curve.ScalarMult(Px, Py, e)

	minusedPy := new(big.Int).Neg(ePy)
	xX, xY := Curve.Add(sGx, sGy, ePx, new(big.Int).Mod(minusedPy, Curve.P))
	return append(xX.Bytes(), xY.Bytes()...)
}

func main() {
	P := make(map[int]map[int][]byte)
	x := make(map[int][]byte)

	P[0] = make(map[int][]byte, 3)
	P[1] = make(map[int][]byte, 3)

	for i := 0; i < 3; i++ {
		key, err := NewPrivateKey()
		if err != nil {
			log.Fatal(err)
		}

		x[i] = key.Bytes()

		tempx, tempy := PubkeyPointsFromPrivkey(key)
		P[0][i] = append(tempx.Bytes(), tempy.Bytes()...)
	}

	for i := 0; i < 3; i++ {
		key, err := NewPrivateKey()
		if err != nil {
			log.Fatal(err)
		}

		x[3+i] = key.Bytes()

		tempx, tempy := PubkeyPointsFromPrivkey(key)
		P[1][i] = append(tempx.Bytes(), tempy.Bytes()...)
	}

	jistar := []int{1, 2, 3, 4, 5, 6} // indices of signer in each ring

	M := Sha256([]byte("cool"))

	e := make(map[int]map[int][]byte)
	s := make(map[int]map[int][]byte)

	k := make(map[int][]byte)
	for i := 0; i < 6; i++ {
		ktemp, err := NewPrivateKey()
		if err != nil {
			log.Fatal(err)
		}
		k[i] = append(k[i], ktemp.Bytes()...)
	}

	// start signing
	for i, loop := range P {
		iByte, err := utils.ToByte(i)
		if err != nil {
			log.Fatal(err)
		}

		e[i] = make(map[int][]byte, len(loop))
		s[i] = make(map[int][]byte, len(loop))

		kiGx, kiGy := PubkeyPointsFromPrivkey(new(big.Int).SetBytes(k[i]))
		kiG := append(kiGx.Bytes(), kiGy.Bytes()...)

		jstari := jistar[i]
		jstariByte, err := utils.ToByte(jstari)
		if err != nil {
			log.Fatal(err)
		}

		e[i][jstari+1] = Sha256(M, kiG, iByte, jstariByte)

		mi := len(loop)
		for j := jstari + 1; j < mi-1; j++ {
			jByte, err := utils.ToByte(j)
			if err != nil {
				log.Fatal(err)
			}

			var temp []byte
			temp, s[i][j] = SubtractOnCurve(e[i][j], P[i][j][0:31], P[i][j][32:63])
			e[i][j+1] = Sha256(M, temp, iByte, jByte)
		}
	}

	toBeHashed := []byte("")
	for i := 0; i <= len(P)-1; i++ {
		var temp []byte
		miMinusOne := 2 // len(loop) - 1
		temp, s[i][miMinusOne] = SubtractOnCurve(e[i][miMinusOne], P[i][miMinusOne][0:31], P[i][miMinusOne][32:63])
		toBeHashed = append(toBeHashed, temp...)
	}

	e0 := Sha256(toBeHashed)

	for i := 0; i <= 1; i++ {
		iByte, err := utils.ToByte(i)
		if err != nil {
			log.Fatal(err)
		}

		e[i][0] = e0

		for j := 0; j < jistar[i]; j++ {
			var temp []byte
			jByte, err := utils.ToByte(j)
			if err != nil {
				log.Fatal(err)
			}

			temp, s[i][j] = SubtractOnCurve(e[i][j], P[i][j][0:31], P[i][j][32:63])
			e[i][j+1] = Sha256(M, temp, iByte, jByte)

			ki := new(big.Int).SetBytes(k[i])
			xi := new(big.Int).SetBytes(x[i])
			eijstari := new(big.Int).SetBytes(e[i][jistar[i]])

			xieijstari := new(big.Int).Mul(xi, eijstari)

			s[i][jistar[i]] = new(big.Int).Add(ki, xieijstari).Bytes()
		}
	}

	r := make(map[int]map[int][]byte)
	for i := 0; i <= 1; i++ {
		iByte, err := utils.ToByte(i)
		if err != nil {
			log.Fatal(err)
		}
		r[i] = make(map[int][]byte)

		for j := 0; j <= 2; j++ {
			Pijx, Pijy := new(big.Int).SetBytes(P[i][j][0:31]), new(big.Int).SetBytes(P[i][j][32:63])

			jplusoneByte, err := utils.ToByte(j + 1)
			if err != nil {
				log.Fatal(err)
			}

			r[i][j+1] = SubtractOnCurveS(e[i][j], Pijx, Pijy, s[i][j])

			e[i][j+1] = Sha256(M, r[i][j+1], iByte, jplusoneByte)
		}
	}

	e0prime := Sha256(r[0][2], r[1][2])
	log.Println(e0prime, e0)
}
