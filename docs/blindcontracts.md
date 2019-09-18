# Blind Contracts

The idea behind Random Contracts is simple - an organization is presented with a set of choices (or options) and can choose one but it doesn't know which one it chooses until the transaction hits the blockchain. This idea is built upon the concept of Oblivious Transfers.

### Oblivious Transfers

Lets assume that there are two parties in a system - The Receiver R and the Signer S. Lets assume that the receiver has a set of messages m0, m1, ... mn that it wants the signer to sign but it doesn't want the signer to know which one it signs until it (R) broadcasts the same (on a blockchain in this example). For the sake of simplicity, assume there are two messages (and related outcomes O0 and O1) m0 and m1 of which the receiver would like one to be committed towards. We use a combination of Pedersen commitments and Scriptless Scripts on Bitcoin in order to achieve the following protocol.

1. Receiver chooses c from {0,1}
2. Receiver generates a random number b, generates a pedersen commitment (T) bG + cH
3. Receiver passes along Pedersen commitment to Signer
4. Signer signs against T and T -1H. Since c belongs to {0,1},
  - if the receiver knows the discrete log of bG (ie c=0), it doesn't know the discrete log of bG - 1H
  - if the receiver knows the discrete log of bG + 1H (ie c=1), it doesn't know the discrete log of bG
5. Signer sends (s0, R0) and (s1, R1) to Receiver
6. Receiver validates adaptor signatures, signs one of them and broadcasts it to the blockchain

All signatures above are Schnorr Adaptor Signatures.

Constructed the right way, we can tweak this construction for lotteries, signing choices, etc which would prove to be interesting applications.

### Applications

1. Blind Commitments - if a company / nation state wants to commit to one out of a set of pledges but can not choose which one, we can use this protocol to ensure that the country knows it is signing to commit to one of the pledges but doesn't know which one until the transaction is broadcast the the blockchain. This can be useful in conferences where a specific NGO can propose a list of alternatives that are acceptable to all nations and nations are encouraged to commit to one of the specific outcomes.

2. Lotteries and games - This construction could be used to have games where people win 50% of the time, 60% of the time, etc depending on the number of messages and parties. This, along with the idea of financial markets could be useful in creating financial incentives for parties to take action on climate change.

3. Future pledges - We could have a smart contract which stores actors' pledges towards a set of objectives (eg. Country A will reduce car production by 50% if temperature increase at the global level is 1 degree by the end of 2030) but cannot broadcast them until the date (2030 in the example) has arrived. This would ensure that countries can still remember the pledges that they made in 2019 about climate change and act on it. This contract could also setup a collective which collects funds and donate that towards country A's environmental initiatives if the event (1 degree before 2030) triggers

### Resources

1. https://www.youtube.com/watch?time_continue=10375&v=-gdfxNalDIc
2. http://www.cs.nccu.edu.tw/~raylin/UndergraduateCourse/ComtenporaryCryptography/Spring2009/TSOINSPET2007.pdf
