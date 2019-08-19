package main

import (
	"encoding/hex"
	"encoding/json"
	"github.com/pkg/errors"
	"log"
	"math/big"

	// btcutils "github.com/bithyve/research/utils"
	utils "github.com/Varunram/essentials/utils"
	bech32 "github.com/bithyve/research/bech32"
	bitcoinrpc "github.com/bithyve/research/rpc"
	btcutils "github.com/bithyve/research/utils"
	"github.com/btcsuite/btcd/btcec"
)

var Curve *btcec.KoblitzCurve = btcec.S256() // take only the curve, can't use other stuff

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
	x1, err := btcutils.NewPrivateKey() // lets assume this to be the same as x
	if err != nil {
		log.Fatal(err)
	}

	x2, err := btcutils.NewPrivateKey()
	if err != nil {
		log.Fatal(err)
	}

	P1x, P1y := btcutils.PubkeyPointsFromPrivkey(x1)
	P2x, P2y := btcutils.PubkeyPointsFromPrivkey(x2)

	err = testAddHomomorphic(P1x, P1y, P2x, P2y, x1, x2)
	if err != nil {
		log.Fatal(err)
	}

	shaBytes := btcutils.Sha256(Curve.Params().Gx.Bytes(), Curve.Params().Gy.Bytes()) // btcutils.Sha256(G)
	Hx, Hy := Curve.ScalarBaseMult(shaBytes)                                          // btcutils.Point(btcutils.Sha256(G))

	var a []byte
	a = []byte{1}

	aHx, aHy := Curve.ScalarMult(Hx, Hy, a)

	commitmentx, commitmenty := Curve.Add(P1x, P1y, aHx, aHy)
	log.Println("commitment: ", commitmentx, commitmenty)
}

func signCommitment() {
	x, err := btcutils.NewPrivateKey() // lets assume this to be the same as x
	if err != nil {
		log.Fatal(err)
	}

	Px, Py := btcutils.PubkeyPointsFromPrivkey(x) // P = x*G

	shaBytes := btcutils.Sha256(Curve.Params().Gx.Bytes(), Curve.Params().Gy.Bytes()) // btcutils.Sha256(G)
	Hx, Hy := Curve.ScalarBaseMult(shaBytes)                                          // btcutils.Point(btcutils.Sha256(G))

	// var a []byte
	// a = []byte{1}

	// aHx, aHy := Curve.ScalarMult(Hx, Hy, a)
	// Cx, Cy := Curve.Add(Px, Py, aHx, aHy)

	oneHx, oneHy := Curve.ScalarMult(Hx, Hy, []byte{1})

	Cprx, Cpry := Curve.Add(Px, Py, new(big.Int).Neg(oneHx), new(big.Int).Neg(oneHy))

	CprHash := btcutils.Sha256(Cprx.Bytes(), Cpry.Bytes())
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

func Create21AOSSig() (*big.Int, *big.Int, *big.Int, btcutils.Point, btcutils.Point, []byte) {
	// Lets assume two parties - Brian and Dom.
	// Lets assume Brian and Dom's P are Pb and Pd. Lets assume their private keys
	// are b and d.
	// Then lets assume Dom wants to create a ring signature over C1 and C2 where
	// C1 = xG + 1H and C2 = xG where x is the blinding factor

	b, err := btcutils.NewPrivateKey()
	if err != nil {
		log.Fatal(err)
	}

	d, err := btcutils.NewPrivateKey()
	if err != nil {
		log.Fatal(err)
	}

	Pb := btcutils.PointFromPrivkey(b)
	Pd := btcutils.PointFromPrivkey(d)

	x, err := btcutils.NewPrivateKey()
	if err != nil {
		log.Fatal(err)
	}

	xG := btcutils.PointFromPrivkey(x)

	shaBytes := btcutils.Sha256(Curve.Params().Gx.Bytes(), Curve.Params().Gy.Bytes()) // btcutils.Sha256(G)

	var H btcutils.Point
	H.Set(Curve.ScalarBaseMult(shaBytes)) // H = btcutils.Point(btcutils.Sha256(G))

	one := []byte{1}
	oneH := btcutils.ScalarMult(H, one)

	C1 := btcutils.Add(xG, oneH) // xG + 1H

	var C2 btcutils.Point
	C2.Set(Curve.ScalarBaseMult(x.Bytes())) // xG

	var m []byte
	m = append(m, append(C1.Bytes(), C2.Bytes()...)...)

	kd, err := btcutils.NewPrivateKey() // random nonce for ring sig
	if err != nil {
		log.Fatal(err)
	}

	var Kd btcutils.Point
	Kd.Set(Curve.ScalarBaseMult(kd.Bytes())) // K = kd*G

	BrianNodeNumber := []byte{2} // assume brian has node number 2

	eb := btcutils.Sha256(Kd.Bytes(), m, BrianNodeNumber)

	sb, err := btcutils.NewPrivateKey() // choose a signature sb at random fro brian
	if err != nil {
		log.Fatal(err)
	}

	var sbG btcutils.Point
	sbG.Set(Curve.ScalarBaseMult(sb.Bytes())) // sb*G

	ebPb := btcutils.ScalarMult(Pb, eb) // eb * Pb

	Kb := btcutils.Sub(sbG, ebPb) // Kb = sb*G - eb*Pb

	DomNodeNumber := []byte{1}

	ed := btcutils.Sha256(Kb.Bytes(), m, DomNodeNumber) // ed = H(Kb || m || D)

	edd := new(big.Int).Mul(new(big.Int).SetBytes(ed), d) // ed * d

	sd := new(big.Int).Add(edd, kd) // ed*d + kd

	return new(big.Int).SetBytes(eb), sb, sd, Pb, Pd, m
}

func Verify21AOSSig() {
	eb, sb, sd, Pb, Pd, m := Create21AOSSig()

	BrianNodeNumber := []byte{2} // assume brian has node number 2
	DomNodeNumber := []byte{1}

	var sbG btcutils.Point
	sbG.Set(Curve.ScalarBaseMult(sb.Bytes())) // sb*G

	ebPb := btcutils.ScalarMult(Pb, eb.Bytes()) // eb*Pb

	Kb := btcutils.Sub(sbG, ebPb) // Kb = sb*G - eb*Pb

	ed := btcutils.Sha256(Kb.Bytes(), m, DomNodeNumber)
	// log.Println("ed: ", ed)

	var sdG btcutils.Point
	sdG.Set(Curve.ScalarBaseMult(sd.Bytes()))

	edPd := btcutils.ScalarMult(Pd, ed)

	Kd := btcutils.Sub(sdG, edPd)

	ebCheck := btcutils.Sha256(Kd.Bytes(), m, BrianNodeNumber)
	ebCheckInt := new(big.Int).SetBytes(ebCheck)

	if ebCheckInt.Cmp(eb) != 0 {
		log.Fatal("Signatures don't match")
	} else {
		log.Println("Ring signatures validated")
	}
}

func SubtractOnCurve(e []byte, Px, Py *big.Int) ([]byte, []byte, *big.Int) {
	s, err := btcutils.NewPrivateKey()
	if err != nil {
		log.Fatal(err)
	}

	sGx, sGy := Curve.ScalarBaseMult(s.Bytes())

	ePx, ePy := Curve.ScalarMult(Px, Py, e)

	minusedPy := new(big.Int).Neg(ePy)

	xX, xY := Curve.Add(sGx, sGy, ePx, new(big.Int).Mod(minusedPy, Curve.P))
	return xX.Bytes(), xY.Bytes(), s
}

func SubtractOnCurveS(e []byte, Px *big.Int, Py *big.Int, s *big.Int) ([]byte, []byte) {
	sGx, sGy := Curve.ScalarBaseMult(s.Bytes())

	ePx, ePy := Curve.ScalarMult(Px, Py, e)

	minusedPy := new(big.Int).Neg(ePy)
	xX, xY := Curve.Add(sGx, sGy, ePx, new(big.Int).Mod(minusedPy, Curve.P))
	return xX.Bytes(), xY.Bytes()
}

// ring sigs will not work since there are multiple ambiguities in the paper. Here,
// the math is implemented correctly and a person implementing it should consult with
// paper authors before making a decision

func testBorroeman() {
	P := make(map[int]map[int][]*big.Int)
	x := make(map[int]*big.Int)

	P[0] = make(map[int][]*big.Int, 3)
	P[1] = make(map[int][]*big.Int, 3)

	for i := 0; i < 3; i++ {
		P[0][i] = make([]*big.Int, 2)

		key, err := btcutils.NewPrivateKey()
		if err != nil {
			log.Fatal(err)
		}

		x[i] = key

		P[0][i][0], P[0][i][1] = btcutils.PubkeyPointsFromPrivkey(key)
	}

	for i := 0; i < 3; i++ {
		P[1][i] = make([]*big.Int, 2)

		key, err := btcutils.NewPrivateKey()
		if err != nil {
			log.Fatal(err)
		}

		x[3+i] = key

		P[1][i][0], P[1][i][1] = btcutils.PubkeyPointsFromPrivkey(key)
	}

	jistar := []int{1, 2, 3, 4, 5, 6} // indices of signer in each ring

	M := btcutils.Sha256([]byte("cool"))

	e := make(map[int]map[int][]byte)
	s := make(map[int]map[int]*big.Int)
	k := make(map[int]*big.Int)

	for i := 0; i < 6; i++ {
		ktemp, err := btcutils.NewPrivateKey()
		if err != nil {
			log.Fatal(err)
		}
		k[i] = ktemp
	}

	// start signing
	for i, loop := range P {
		iByte, err := utils.ToByte(i)
		if err != nil {
			log.Fatal(err)
		}

		e[i] = make(map[int][]byte, len(loop))
		s[i] = make(map[int]*big.Int, len(loop))

		kiGx, kiGy := btcutils.PubkeyPointsFromPrivkey(k[i])
		kiG := append(kiGx.Bytes(), kiGy.Bytes()...)

		jstari := jistar[i]
		jstariByte, err := utils.ToByte(jstari)
		if err != nil {
			log.Fatal(err)
		}

		e[i][jstari+1] = btcutils.Sha256(M, kiG, iByte, jstariByte)

		mi := len(loop)
		for j := jstari + 1; j < mi-1; j++ {
			jByte, err := utils.ToByte(j)
			if err != nil {
				log.Fatal(err)
			}

			var tempx, tempy []byte
			tempx, tempy, s[i][j] = SubtractOnCurve(e[i][j], P[i][j][0], P[i][j][1])
			e[i][j+1] = btcutils.Sha256(M, tempx, tempy, iByte, jByte)
			log.Println(e[i][j+1])
		}
	}

	toBeHashed := []byte("")
	for i := 0; i <= len(P)-1; i++ {
		var tempx, tempy []byte
		miMinusOne := 2
		tempx, tempy, s[i][miMinusOne] = SubtractOnCurve(e[i][miMinusOne], P[i][miMinusOne][0], P[i][miMinusOne][1])
		temp := append(tempx, tempy...)
		toBeHashed = append(toBeHashed, temp...)
	}

	e0 := btcutils.Sha256(toBeHashed)

	for i := 0; i <= 1; i++ {
		iByte, err := utils.ToByte(i)
		if err != nil {
			log.Fatal(err)
		}

		e[i][0] = e0

		for j := 0; j < jistar[i]; j++ {
			var tempx, tempy []byte
			jByte, err := utils.ToByte(j)
			if err != nil {
				log.Fatal(err)
			}

			tempx, tempy, s[i][j] = SubtractOnCurve(e[i][j], P[i][j][0], P[i][j][1])

			e[i][j+1] = btcutils.Sha256(M, tempx, tempy, iByte, jByte)

			eijstari := new(big.Int).SetBytes(e[i][jistar[i]])

			xieijstari := new(big.Int).Mul(x[i], eijstari)

			s[i][jistar[i]] = new(big.Int).Add(k[i], xieijstari)
		}
	}

	// log.Println("e0: ", e0)
	// log.Println("sigs: ", s)
	log.Println("e: ", e)

	ex := make(map[int]map[int][]byte)
	r := make(map[int]map[int][]byte)

	for i := 0; i <= 1; i++ {
		ex[i] = make(map[int][]byte)
		ex[i][0] = e0
		iByte, err := utils.ToByte(i)
		if err != nil {
			log.Fatal(err)
		}
		r[i] = make(map[int][]byte)

		for j := 0; j <= 2; j++ {
			jplusoneByte, err := utils.ToByte(j + 1)
			if err != nil {
				log.Fatal(err)
			}

			tempx, tempy := SubtractOnCurveS(ex[i][j], P[i][j][0], P[i][j][1], s[i][j])

			r[i][j+1] = append(tempx, tempy...)
			e[i][j+1] = btcutils.Sha256(M, tempx, tempy, iByte, jplusoneByte)
			log.Println(e[i][j+1])
		}
	}

	e0prime := btcutils.Sha256(r[0][2], r[1][2], M)
	log.Println(e0prime, e0)
}

func testElgamal() {
	// an elgamal commitment is a small upgrade from Pedersen commitments
	// xG + rH - Pedersen
	// (xG + rH, rG) - Elgamal

	// there is a small nuanace to how this work. lets assume we have xG + rH in Pedersen
	// if a person can enumerate over the entire input space, they can find r and as a result
	// binding is broken (the commitment binds to r but we now have a fake r). If binding
	// is broken, in some applications like CT, one can print infinite money. So pedersen
	// commitments have perfect hiding (no one can predict what value exists) but computational
	// binding (a person with finitely infinite resources can find another r)

	// In Elgamal, we commit to anohter point - rG. We can compute this since we know r (used anyway for Pedersen)
	// but where's the difference?
	// Lets assume an attacker finds r like in the above case. since a Pedersen commitment is xG + rH, they can
	// commit to another value r' where rH = r'H. In Elgamal, the mapping from r to rG is one-one, so we
	// get perfect binding (no one can commit to another value even if they have resources to find r). But
	// an attacker with resources can find r and then compute C - rH to find xG (the hidden commitment value).
	// Hence, Pedersen offers perfect binding and computational hiding

	// note that we're only implementing the commitment scheme here not the signature scheme. Signature scheme could be
	// AOS or something similar

	x, err := btcutils.NewPrivateKey()
	if err != nil {
		log.Fatal(err)
	}

	r, err := btcutils.NewPrivateKey()
	if err != nil {
		log.Fatal(err)
	}

	xG := btcutils.PointFromPrivkey(x)

	shaBytes := btcutils.Sha256(Curve.Params().Gx.Bytes(), Curve.Params().Gy.Bytes()) // btcutils.Sha256(G)

	var H btcutils.Point
	H.Set(Curve.ScalarBaseMult(shaBytes)) // H = btcutils.Point(btcutils.Sha256(G))

	var rH btcutils.Point
	rH.Set(Curve.ScalarMult(H.X, H.Y, r.Bytes()))

	var pedersen btcutils.Point
	pedersen.Add(xG, rH)

	var rG btcutils.Point
	rG = btcutils.PointFromPrivkey(r)

	var elgamal btcutils.Elgamal
	elgamal.Set(pedersen, rG)

	log.Println("elgamal: ", elgamal)
}

func main() {
	x, err := btcutils.NewPrivateKey()
	if err != nil {
		log.Fatal(err)
	}

	r, err := btcutils.NewPrivateKey()
	if err != nil {
		log.Fatal(err)
	}

	xG := btcutils.PointFromPrivkey(x)

	shaBytes := btcutils.Sha256(Curve.Params().Gx.Bytes(), Curve.Params().Gy.Bytes()) // btcutils.Sha256(G)

	var H btcutils.Point
	H.Set(Curve.ScalarBaseMult(shaBytes)) // H = btcutils.Point(btcutils.Sha256(G))

	var rH btcutils.Point
	rH.Set(Curve.ScalarMult(H.X, H.Y, r.Bytes()))

	var pedersen btcutils.Point
	pedersen.Add(xG, rH)

	var rG btcutils.Point
	rG = btcutils.PointFromPrivkey(r)

	// xG+(r+H(xG+rH||rG))H is the switch commitment
	// let the ugly term be p ie the commitment is xG = pH
	insideHash := btcutils.Sha256(pedersen.Bytes(), rG.Bytes())
	insideHashNumber := new(big.Int).SetBytes(insideHash)
	p := new(big.Int).Add(r, insideHashNumber)

	var pH btcutils.Point
	pH.Set(Curve.ScalarMult(H.X, H.Y, p.Bytes()))

	var switchCmt btcutils.Point
	switchCmt.Add(xG, pH)

	log.Println("switch commitment: ", switchCmt)
}
