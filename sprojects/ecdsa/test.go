package main

import (
	"crypto/rand"
	"crypto/rsa"
	"log"
	"math/big"

	btcutils "github.com/bithyve/research/utils"
	"github.com/btcsuite/btcd/btcec"
)

var Curve *btcec.KoblitzCurve = btcec.S256() // take only the curve, can't use other stuff

func testDHExchange() {
	x1, err := btcutils.NewPrivateKey() // lets assume this to be the same as x
	if err != nil {
		log.Fatal(err)
	}

	x2, err := btcutils.NewPrivateKey() // lets assume this to be the same as x
	if err != nil {
		log.Fatal(err)
	}

	X1 := btcutils.PointFromPrivkey(x1)
	X2 := btcutils.PointFromPrivkey(x2)

	P1 := btcutils.ScalarMult(X1, x2.Bytes())
	P2 := btcutils.ScalarMult(X2, x1.Bytes())

	if !P1.Cmp(P2) {
		log.Println("dh key exchange failed")
	} else {
		log.Println("dh key exchange succeeded")
	}
}

func main() {
	x1, err := btcutils.NewPrivateKey() // lets assume this to be the same as x
	if err != nil {
		log.Fatal(err)
	}

	x2, err := btcutils.NewPrivateKey() // lets assume this to be the same as x
	if err != nil {
		log.Fatal(err)
	}

	X1 := btcutils.PointFromPrivkey(x1)
	X2 := btcutils.PointFromPrivkey(x2)

	P1 := btcutils.ScalarMult(X1, x2.Bytes())
	P2 := btcutils.ScalarMult(X2, x1.Bytes())

	log.Println("ECDH EXCHANGE: ", P1.Cmp(P2)) // dh key exchange complete
	P := P1

	z := btcutils.Sha256([]byte("hello"))

	k1, err := btcutils.NewPrivateKey() // lets assume this to be the same as x
	if err != nil {
		log.Fatal(err)
	}

	k2, err := btcutils.NewPrivateKey() // lets assume this to be the same as x
	if err != nil {
		log.Fatal(err)
	}

	K1 := btcutils.PointFromPrivkey(k1)
	K2 := btcutils.PointFromPrivkey(k2)

	// parties have generated random points

	// party 2 should receive K1, party 1 should receive K2

	K := btcutils.ScalarMult(K1, k2.Bytes())
	// point K will be given as part of the signature

	k := K.X.Bytes()

	// need to calculate Paillier encryption of x1 here

	primes, err := rsa.GenerateKey(rand.Reader, 2048) // 128 byte keys
	if err != nil {
		log.Fatal(err)
	}

	if len(primes.Primes) != 2 {
		log.Fatal("prime length not 2, quitting")
	}

	p := primes.Primes[0]
	q := primes.Primes[1]

	if len(p.Bytes()) != len(q.Bytes()) {
		log.Fatal("equal length primes needed")
	}

	zero := new(big.Int).SetInt64(0)
	one := new(big.Int).SetInt64(1)
	two := new(big.Int).SetInt64(2)

	n := new(big.Int).Mul(p, q)
	g := new(big.Int).Add(n, one)

	pminus1 := new(big.Int).Sub(p, one)
	qminus1 := new(big.Int).Sub(q, one)

	lambda := new(big.Int).Div(new(big.Int).Mul(pminus1, qminus1), two)

	if (new(big.Int).Exp(new(big.Int).SetBytes(k), lambda, n)).Cmp(one) != 0 {
		log.Fatal("mod exp wrong")
	}

	nsq := new(big.Int).Exp(n, two, zero)

	if (new(big.Int).Exp(new(big.Int).SetBytes(k), new(big.Int).Mul(lambda, n), nsq)).Cmp(one) != 0 {
		log.Fatal("mod exp wrong")
	}

	/*
		NAIVE METHOD
		gx := new(big.Int).Exp(g, x, zero)
		kn := new(big.Int).Exp(new(big.Int).SetBytes(k), n, zero)
		gxkn := new(big.Int).Mul(gx, kn)

		ex := new(big.Int).Mod(gxkn, nsq)
	*/

	x := x1

	ex1 := Encrypt(x, k, n)
	check1 := Decrypt(lambda, ex1, n, g)
	if check1.Cmp(x) != 0 {
		log.Fatal("test decryption of privkey not working")
	}

	kex1 := new(big.Int).Exp(ex1, new(big.Int).SetBytes(k), nsq) // e(x1 * k)

	kex1x2 := new(big.Int).Exp(kex1, x2, nsq) // e(x1 * k * x2)

	gzmodn2 := new(big.Int).Exp(g, new(big.Int).SetBytes(z), nsq)       // g^z
	mulpart := new(big.Int).Mod(new(big.Int).Mul(kex1x2, gzmodn2), nsq) // e(z + x1 * k * x2)

	k2inv := new(big.Int).ModInverse(x2, nsq)
	sprime := new(big.Int).Mod(new(big.Int).Mul(mulpart, k2inv), nsq)

	log.Println("SPAM: ", P, K2)
	log.Println("S' = ", sprime)
	sig := Decrypt(lambda, sprime, n, g)
	log.Println("SIG=", sig)
}

func Decrypt(lambda, ex, n, g *big.Int) *big.Int {
	zero := new(big.Int).SetInt64(0)
	two := new(big.Int).SetInt64(2)
	nsq := new(big.Int).Exp(n, two, zero)

	glambdamodn2 := new(big.Int).Exp(g, lambda, nsq)
	mu := new(big.Int).ModInverse(L(glambdamodn2, n), n)

	clambdamodnsq := new(big.Int).Exp(ex, lambda, nsq)
	Lc := L(clambdamodnsq, n)
	Lcmu := new(big.Int).Mul(Lc, mu)
	Lcmumodn := new(big.Int).Mod(Lcmu, n)

	return Lcmumodn
}

func L(x, n *big.Int) *big.Int {
	one := new(big.Int).SetInt64(1)

	xminusone := new(big.Int).Sub(x, one)
	divn := new(big.Int).Div(xminusone, n)

	return divn
}

func TestDecrypt(lambda, ex, x, n, g *big.Int) {
	zero := new(big.Int).SetInt64(0)
	one := new(big.Int).SetInt64(1)
	two := new(big.Int).SetInt64(2)

	nsq := new(big.Int).Exp(n, two, zero)

	glambdamodn2 := new(big.Int).Exp(g, lambda, nsq)
	mu := new(big.Int).ModInverse(L(glambdamodn2, n), n)

	exlambdamodnsq := new(big.Int).Exp(ex, lambda, nsq)

	// 1 + lambda * x * n == g ^ (lambda *x) mod (n^2) == g ^ (lambda *x) * r^(n*lambda) mod n^2 == e(x)^lambda mod n^2
	if (new(big.Int).Add(new(big.Int).Mul(new(big.Int).Mul(lambda, x), n), one)).Cmp(exlambdamodnsq) != 0 {
		log.Fatal("test decryption failed at 1+lambda*n*x")
	}

	exlambdamodnsqminone := new(big.Int).Sub(exlambdamodnsq, one)
	exlambdamodnsqminonedivn := new(big.Int).Div(exlambdamodnsqminone, n)

	if new(big.Int).Mul(lambda, x).Cmp(exlambdamodnsqminonedivn) != 0 {
		log.Fatal("test decryption failed")
	}

	decrypt := new(big.Int).Mod(new(big.Int).Mul(exlambdamodnsqminonedivn, mu), n)
	decrypt2 := new(big.Int).Div(exlambdamodnsqminonedivn, lambda)

	if decrypt.Cmp(decrypt2) != 0 {
		log.Fatal("decryption failed")
	} else {
		log.Println("test decryption works")
	}
}

func Encrypt(x *big.Int, k []byte, n *big.Int) *big.Int {
	one := new(big.Int).SetInt64(1)
	two := new(big.Int).SetInt64(2)
	nsq := new(big.Int).Exp(n, two, nil)

	xn := new(big.Int).Mul(x, n)

	xnplusonemodnsq := new(big.Int).Mod(new(big.Int).Add(xn, one), nsq)

	knmodnsq := new(big.Int).Exp(new(big.Int).SetBytes(k), n, nsq)

	return new(big.Int).Mod(new(big.Int).Mul(xnplusonemodnsq, knmodnsq), nsq)
}
