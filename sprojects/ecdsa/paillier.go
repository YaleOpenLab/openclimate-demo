package main

import (
	"math/big"
)

func GetPrimes(len int, bitlen int) ([]*big.Int, error) {
	var primes []*big.Int

	for i := 0; i < (len+1)/2; i++ {
		rsakeys, err := rsa.GenerateKey(rand.Reader, bitlen*16) // 2x128 byte keys give 2048 bit rsa keys
		if err != nil {
			return primes, err
		}

		p := rsakeys.Primes[0]
		q := rsakeys.Primes[1]

		if p.BitLen() != q.BitLen() {
			log.Fatal("equal length primes needed")
		}

		primes = append(primes, p, q)
	}

	return primes[0:len], nil
}

func Encrypt(x *big.Int, k, n, nsq *big.Int) *big.Int {
	xn := new(big.Int).Mul(x, n)

	xnplusonemodnsq := new(big.Int).Mod(new(big.Int).Add(xn, one), nsq)

	knmodnsq := new(big.Int).Exp(k, n, nsq)

	return new(big.Int).Mod(new(big.Int).Mul(xnplusonemodnsq, knmodnsq), nsq)
}

func Decrypt(lambda, ex, n, nsq, mu *big.Int) *big.Int {
	clambdamodnsq := new(big.Int).Exp(ex, lambda, nsq)
	Lc := L(clambdamodnsq, n)
	Lcmu := new(big.Int).Mul(Lc, mu)
	Lcmumodn := new(big.Int).Mod(Lcmu, n)

	return Lcmumodn
}

func L(x, n *big.Int) *big.Int {
	xminusone := new(big.Int).Sub(x, one)
	divn := new(big.Int).Div(xminusone, n)

	return divn
}
