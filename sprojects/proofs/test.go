package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/pkg/errors"
	"io"
	"log"
	"math"
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

func getkG(e, Px, Py []byte) ([]byte, []byte) {
	s, err := NewPrivateKey()
	if err != nil {
		log.Fatal(err)
	}

	sGx, sGy := Curve.ScalarBaseMult(s.Bytes())
	ePx, ePy := Curve.ScalarMult(new(big.Int).SetBytes(Px), new(big.Int).SetBytes(Py), e)

	xX, xY := Curve.Add(sGx, sGy, ePx, ePy)
	return append(xX.Bytes(), xY.Bytes()...), s.Bytes()
}

/*
def get_kG(e, P, s=None):
    '''Use EC operation: kG = sG +eP.
    If s (signature) is not provided, it is generated
    randomly and returned.
    e - hash value, 32 bytes binary
    P - verification pubkey
    s - 32 bytes binary'''
    if not s:
	s = os.urandom(32)
    sG = btc.fast_multiply(btc.G,btc.decode(s,256))
    eP = btc.fast_multiply(P,btc.decode(e,256))
    return (btc.fast_add(sG, eP), s)
*/
func main() {
	pubkeys := make(map[int]map[int][]byte)
	privkeys := make(map[int][]byte)

	pubkeys[0] = make(map[int][]byte, 3)
	pubkeys[1] = make(map[int][]byte, 3)

	for i := 0; i < 3; i++ {
		x, err := NewPrivateKey()
		if err != nil {
			log.Fatal(err)
		}

		privkeys[i] = x.Bytes()

		tempx, tempy := PubkeyPointsFromPrivkey(x)
		pubkeys[0][i] = append(tempx.Bytes(), tempy.Bytes()...)
	}

	for i := 0; i < 3; i++ {
		x, err := NewPrivateKey()
		if err != nil {
			log.Fatal(err)
		}

		privkeys[3+i] = x.Bytes()

		tempx, tempy := PubkeyPointsFromPrivkey(x)
		pubkeys[1][i] = append(tempx.Bytes(), tempy.Bytes()...)
	}

	signingIndices := []int{1, 2, 3, 4, 5, 6} // replace this with a single element and test

	M := Sha256([]byte("cool"))

	e := make(map[int]map[int][]byte)
	s := make(map[int]map[int][]byte)

	startIndex := make([]int, 100)
	for i, _ := range startIndex {
		startIndex[i] = 0
	}

	k := make(map[int][]byte)
	for i := 0; i < 6; i++ {
		ktemp, err := NewPrivateKey()
		if err != nil {
			log.Fatal(err)
		}
		k[i] = append(k[i], ktemp.Bytes()...)
	}

	toBeHashed := []byte("")

	for i, loop := range pubkeys {
		iByte, err := utils.ToByte(i)
		if err != nil {
			log.Fatal(err)
		}
		e[i] = make(map[int][]byte, len(loop))
		s[i] = make(map[int][]byte, len(loop))

		kGx, kGy := PubkeyPointsFromPrivkey(new(big.Int).SetBytes(k[i]))
		kG := append(kGx.Bytes(), kGy.Bytes()...)

		startIndex[i] = int(math.Mod(float64(signingIndices[i]+1), float64(len(loop))))

		if startIndex[i] == 0 {
			toBeHashed = append(toBeHashed, kG...)
			continue
		}

		startIndexByte, err := utils.ToByte(startIndex[i])
		if err != nil {
			log.Fatal(err)
		}

		e[i][startIndex[i]] = Sha256(M, kG, iByte, startIndexByte)

		for x := startIndex[i]; x < len(loop); x++ {
			var y []byte
			xByte, err := utils.ToByte(x)
			if err != nil {
				log.Fatal(err)
			}

			y, s[i][x-1] = getkG(e[i][x-1], pubkeys[i][x-1][0:31], pubkeys[i][x-1][32:63])
			e[i][x] = Sha256(M, y, iByte, xByte)
		}

		var kGend []byte
		kGend, s[i][len(pubkeys[i])-1] = getkG(e[i][len(pubkeys[i])-1], pubkeys[i][len(pubkeys[i])-1][0:31], pubkeys[i][len(pubkeys[i])-1][32:63])
		toBeHashed = append(toBeHashed, kGend...)
	}

	toBeHashed = append(toBeHashed, M...)
	e0 := Sha256(toBeHashed)

	for i, _ := range pubkeys {
		iByte, err := utils.ToByte(i)
		if err != nil {
			log.Fatal(err)
		}
		zeroByteArray := []byte{0}
		e[i][0] = Sha256(M, e0, iByte, zeroByteArray)
	}

	for i, _ := range pubkeys {
		iByte, err := utils.ToByte(i)
		if err != nil {
			log.Fatal(err)
		}
		for x := 1; x < signingIndices[i]+1; x++ {
			var y []byte
			xByte, err := utils.ToByte(x)
			if err != nil {
				log.Fatal(err)
			}

			y, s[i][x-1] = getkG(e[i][x-1], pubkeys[i][x-1][0:31], pubkeys[i][x-1][32:63])
			e[i][x] = Sha256(M, y, iByte, xByte)

			kiNum := new(big.Int).SetBytes(k[i])
			privkeyNum := new(big.Int).SetBytes(privkeys[i])
			esigI := new(big.Int).SetBytes(e[i][signingIndices[i]])

			pknminusesigI := new(big.Int).Add(privkeyNum, new(big.Int).Neg(esigI))
			// aboveModP := new(big.Int).Mod(minusedPby, Curve.P)
			subNum := new(big.Int).Add(kiNum, pknminusesigI)
			subNumMod := new(big.Int).Mod(subNum, Curve.P)
			s[i][signingIndices[i]] = subNumMod.Bytes()
		}
	}

	log.Println("e0: ", e0, "sig: ", s, len(s))
}
