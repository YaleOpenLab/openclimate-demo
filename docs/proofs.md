# Range Proofs using Pederson Commitments and Ring Signatures

### The Problem

We can imagine two scenarios:

1. A specific nation state entity does not want others nations to know the specific compositions of its emissions (say its building a secret plant underground or something)
2. A nation uses a black box scheme to which only specific queries are possible (ie the nation provides an oracle for others to query information from and claims that it is following the NDCs it promised as part of the Paris agreement)

Both these problems require blinding of specific emission amounts while still proving that they somehow are valid and add up to the NDC emission amount specified

### The Construction

We can use Pedersen Commitments and Ring Signatures in order to achieve this. The idea of this is simple:

Assume that you have numbers 1, -3, -5 and 7 and you need to prove that the sum is zero. One would do:

```
(1+7) + (=3+-5) = 0
```

The key is doing this without revealing the numbers themselves, which we explore.

Say you have some data D. We define a commitment as the hash H = H(D). But since this is prone to a grinding attack (one can guess what the data is if its simple) we add a blinding factor `x` and define

```
H = H(x || D
```

As seen in our simple number trick, the commitments must be additively homomorphic ie

```
C(1) + C(2) = C(3)
```

(since `1 + 2 = 3`). We use blockstream's elements sidechain for reference, which uses ECC to do this.

In ECC, `P = x*G` and if  
`P1 = x1*G`, `P2 = x2*G`,  
`P1 + P2 = x1*G + x2*G = (x1+x2)*G`  
(operations mod n, where n is the group order wherever applicable).

Elements uses a different generator point H (which is a point on the curve) for Pedersen commitments, which can be derived from generator G
```
H = H(G)
```
and define our commitment as
```
C = x*G + a*H
```
where x is the blinding factor and a is the amount we want to commit to. It is easy to verify that this commitment has the same properties we desire. But there's a problem if we're using this for non negative numbers (which the NDCs have): one could create negative emissions like so: `1 + 3 + (-5) + 7` and report their net emissions as 6 instead of 16 since there's nothing preventing them from doing so.

For this, we need to provide a range proof (a proof that proves that a particular value is between two pre-defined values). In the context of climate data, we can assume that this range is `0-600*(10**9)` ie `0-2**40`. So to construct a proof that one knows
```
C = xG + 1H
```
we construct
```
C' = C - 1H
```
and ask the prover to provide a signature with C'. This would mean that the person doesn't reveal his blinding factor x and still proves the commitment towards 1. But assume that we want to not tell someone that the amount is 1 and instead just want to reveal that emissions are in a range. For this, we need a ring signature. To prove that the prover's commitment amount is 0 / 1, we provide a signature over `{C, C'}` in the above case. If we want to prove that  C is in the range `[0, 32)`, we provide the following proofs:

```
C1 is 0 or 1 C2 is 0 or 2 C3 is 0 or 4 C4 is 0 or 8 C5 is 0 or 16.
```

and pick the blinding factors such that `C1+C2+C3+C4+C5 = C`

### Use in NDC / Climate Change Stuff

In the context of the problems described earlier, we need to construct something such that we know only the sum of their emissions (if needed). We should construct an appropriate ring signature scheme and the oracle can provide these proofs while still not revealing the commitment amounts to the verifier.

## References (not complete, will add once final)
1. https://blockstream.com/bitcoin17-final41.pdf
2. https://elementsproject.org/features/confidential-transactions/investigation
