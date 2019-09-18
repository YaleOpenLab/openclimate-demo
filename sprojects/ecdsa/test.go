package main

import (
	"crypto/rand"
	"crypto/rsa"
	"log"
	"math/big"

	btcutils "github.com/bithyve/research/utils"
	"github.com/btcsuite/btcd/btcec"
)

var (
	Curve *btcec.KoblitzCurve = btcec.S256() // take only the curve, can't use other stuff
	zero                      = new(big.Int).SetInt64(0)
	one                       = new(big.Int).SetInt64(1)
	two                       = new(big.Int).SetInt64(2)
)

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

func testPaillier() {
	k1, err := btcutils.NewPrivateKey()
	if err != nil {
		log.Fatal(err)
	}

	k2, err := btcutils.NewPrivateKey()
	if err != nil {
		log.Fatal(err)
	}

	K1 := btcutils.PointFromPrivkey(k1)
	K2 := btcutils.PointFromPrivkey(k2)
	K := btcutils.ScalarMult(K1, k2.Bytes())

	if !K.Cmp(btcutils.ScalarMult(K2, k1.Bytes())) {
		log.Fatal("k vals don't match")
	}

	k := K.X // k will be given along with s (the signature)

	primes, err := GetPrimes(2, 128)
	if err != nil {
		log.Fatal(err)
	}

	p := primes[0]
	q := primes[1]

	n := new(big.Int).Mul(p, q)
	nsq := new(big.Int).Exp(n, two, zero) // n^2

	g := new(big.Int).Add(n, one)
	// Public key: (n, g)

	pminus1 := new(big.Int).Sub(p, one) // p-1
	qminus1 := new(big.Int).Sub(q, one) // q-1

	gcd := new(big.Int).GCD(nil, nil, pminus1, qminus1)                 // lcm = (p*q) / gcd(p,q)
	lambda := new(big.Int).Div(new(big.Int).Mul(pminus1, qminus1), gcd) // lambda = lcm
	glambdamodn2 := new(big.Int).Exp(g, lambda, nsq)
	mu := new(big.Int).ModInverse(L(glambdamodn2, n), n)
	// Private Key: (lambda, mu)

	if (new(big.Int).Exp(k, lambda, n)).Cmp(one) != 0 {
		log.Fatal("mod exp wrong")
	}

	if (new(big.Int).Exp(k, new(big.Int).Mul(lambda, n), nsq)).Cmp(one) != 0 {
		log.Fatal("mod exp wrong")
	}

	testMsg := new(big.Int).SetBytes([]byte("Hello World"))
	testCipherText := Encrypt(testMsg, k, n, nsq)

	check1 := Decrypt(lambda, testCipherText, n, nsq, mu)
	if check1.Cmp(testMsg) != 0 {
		log.Fatal("test decryption of privkey not working")
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

	if !P1.Cmp(P2) {
		log.Fatal("ECDH points don't match")
	}

	// after exchanging DH keys, send the encrypted private key over
	ex1 := Encrypt(x1, k, n, nsq)                                // set encrypted text to private key of Party 1

	// signing message
	z := btcutils.Sha256([]byte("hello"))

	// Setup Random Points to get k
	k1, err := btcutils.NewPrivateKey()
	if err != nil {
		log.Fatal(err)
	}

	k2, err := btcutils.NewPrivateKey()
	if err != nil {
		log.Fatal(err)
	}

	K1 := btcutils.PointFromPrivkey(k1)
	K2 := btcutils.PointFromPrivkey(k2)
	K := btcutils.ScalarMult(K1, k2.Bytes())

	if !K.Cmp(btcutils.ScalarMult(K2, k1.Bytes())) {
		log.Fatal("k vals don't match")
	}

	k := K.X

	// Get primes p, q for Paillier signing
	primes, err := GetPrimes(2, 128)
	if err != nil {
		log.Fatal(err)
	}

	p := primes[0]
	q := primes[1]

	n := new(big.Int).Mul(p, q)
	nsq := new(big.Int).Exp(n, two, zero) // n^2

	g := new(big.Int).Add(n, one)
	// Public key: (n, g)

	pminus1 := new(big.Int).Sub(p, one) // p-1
	qminus1 := new(big.Int).Sub(q, one) // q-1

	gcd := new(big.Int).GCD(nil, nil, pminus1, qminus1)                 // lcm = (p*q) / gcd(p,q)
	lambda := new(big.Int).Div(new(big.Int).Mul(pminus1, qminus1), gcd) // lambda = lcm
	glambdamodn2 := new(big.Int).Exp(g, lambda, nsq)
	mu := new(big.Int).ModInverse(L(glambdamodn2, n), n)
	// Private Key: (lambda, mu)

	if (new(big.Int).Exp(k, lambda, n)).Cmp(one) != 0 {
		log.Fatal("mod exp wrong")
	}

	if (new(big.Int).Exp(k, new(big.Int).Mul(lambda, n), nsq)).Cmp(one) != 0 {
		log.Fatal("mod exp wrong")
	}

	/*
		Simple Method:
		gx := new(big.Int).Exp(g, x, zero)
		kn := new(big.Int).Exp(k, n, zero)
		gxkn := new(big.Int).Mul(gx, kn)
		ex := new(big.Int).Mod(gxkn, nsq)

		Better Method:
		ab mod(n) = (amod(n) * bmod(n)) mod n
		g^x * r^n mod(n^2) = (g^x mod(n^2)) * (r^n mod(n^2)) mod(n^2)
		g^x mod(n^2) = (1+n)^x mod(n^2)
		(1+n)^x = 1 + nx + n(n-1)/2*x^2 + ... = ((1+nx)*mod(n^2) + n^2(K1 + K2n + ...))
		(1+n)^x*mod(n^2) = (1+nx)mod(n^2)
		(g^x mod(n^2)) = (1+nx)mod(n^2)
	*/

	kex1 := new(big.Int).Exp(ex1, k, nsq) // (ex1^k)mod(n^2) = e(x1*k)

	kex1x2 := new(big.Int).Exp(kex1, x2, nsq) // (kex1^x2) mod(n^2) = e(x1*k*x2)

	gzmodn2 := new(big.Int).Exp(g, new(big.Int).SetBytes(z), nsq)       // g^z mod(n^2)
	mulpart := new(big.Int).Mod(new(big.Int).Mul(kex1x2, gzmodn2), nsq) // g^z mod(n^2)*e(x1*k*x2) = e(z+x1*k*x2)

	k2inv := new(big.Int).ModInverse(x2, nsq) // k2^-1 mod(n^2)

	sprime := new(big.Int).Exp(mulpart, k2inv, nsq) // since we need to multiply by k2inv, we exponentiate it to that
	// s' = e((z+x1*k*x2)/k2)
	log.Println("Partial sig: ", len(sprime.Bytes()))
	sig := Decrypt(lambda, sprime, n, nsq, mu)
	log.Println("SIG=", k, sig)
}
