# Two Party ECDSA

ECDSA is a well knmown algorithm used in multiple applications. One of the main limitations pointed of ECDSA pointed out is its inability to aggregate signatures (whereas Schnorr signatures can aggregate signatures). A limitation that results is that constructing threshold signatures is hard - [until recently](https://eprint.iacr.org/2017/552.pdf), it was not known to be efficient.

## Construction

We perform operations using the [Paillier cryptosystem](https://en.wikipedia.org/wiki/Paillier_cryptosystem). This results in one extra assumption (ECDLP + Paillier) in the proof system.

###  Paillier

#### Key Setup

1. Choose two numbers p and q of equal length.
2. n = p*q and lambda = lcm(p-1, q-1) = (p-1) * (q-1) / gcd(p-1, q-1)
3. g = n + 1
4. Compute gl = g^lambda mod(n^2)
5. Define L(x, n) = (x-1)/n
6. mu = L(gl, n)^(-1) mod(n)
7. Public encryption key - (n, g)
8. Private decryption key - (lambda, mu)

#### Encryption

1. m is the message to be encrypted
2. choose random r
3. Compute ciphertext - g^m * r^n mod(n^2)

#### Decryption

1. let c be the ciphertext and m be the message
2. cl = c^lambda mod(n^2)
3. m = L(cl) * mu mod(n)

### 2 Party ECDSA

#### Setup

1. Assume two parties - Parties A and B
2. Let parties A and B setup a secure ECDH channel between them
3. Party A encrypts private key a under Paillier and gives e(a) to Party B
4. A and B agree on a common hashed message H(m)
5. A and B generate two random numbers ra and rb with public keys Ra and Rb
6. A and B exchange Ra and Rb similar to ECDH and arrive at a common point r (ra*Rb = rb*Ra)
7. B calculates the following signature: (z + r * e(a) * b) / r2)
  - calculate rea = r * e(a) as e(a) ^ r mod(n^2) under Paillier
  - calculate reab = rea * b as rea ^ b mod (n^2) under Paillier
  - calculate zreab = z + reab as rea * g^z mod(n^2) under Paillier
  - calculate sig = zreab / r2 as zreab ^ (r2 ^(-1)) mod (n^2)
8. B gives sig to A
9. A decrypts sig under Paillier, multiples by a^(-1) and publishes signature

## Applications

1. Untrusted signing - Two Parties A and B who don't trust each other to sign the message first can agree on the common message m and then use 2P ECDSA to coordinate signing. In this way, neither party is forced to trust the other and if one party aborts the scheme, the other stays unaffected.

2. Threshold signing - Threshold signing is an extension of 2 party ECDSA and can be used to coordinate m of n signature schemes. The advantage of ECDSA is that it would be compatible with existing blockchains like Bitcoin.

## Resources

1. https://eprint.iacr.org/2017/552.pdf
2. https://medium.com/cryptoadvance/ecdsa-is-not-that-bad-two-party-signing-without-schnorr-or-bls-1941806ec36f
3. https://en.wikipedia.org/wiki/Paillier_cryptosystem
