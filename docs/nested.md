# Nested Scope Accounting

### The Problem
Assume three laters of reporting:

1. City
2. State
3. Nation

and assume that city, state and nation level companies are present. Assume that emissions of each of these layers are Ec, Es, En respectively. Assume a specific company C's emissions on the city, state and nation levels are Cc, Cs, Cn respectively. We also assume that on a global level, there are

1. Nations
2. Companies

and the aggregate emissions data for the whole planet is En' + Ec' where En' is the aggregate of emissions reported for nations and Ec' is the aggregate emissions reported for companies

The problem is that En is a superset of Es and En, Es is a superset of En on the reporting level and Cc that of Ec, Cs that of Es and Cn that of En. We need to find a way in which Cc's emissions are accounted for only once in the reporting of the Nation's emissions report.

### The simple view

We assume that when entities report certain data, they reference where the emissions came from. For eg, say country X is reporting data on its emissions, we assume that it references that city Y emitted Z amount of CO2 in the particular reporting period. We also assume that companies report that their plant in city Y emitted Z' CO2 equivalent in the observation period ie

`Z = Z' + Zt`, where `Zt` is the true emission when emissions are accounted for at a global level. We need to find: `sum(Zt) = sum(Z) - sum(Z')`. Even though this seems straightforward, the model does not take into account that multiple data points are missing at each level. As a result, we're forced to track emitting entities individually in order to make sure that their emissions are not accounted for.

### The Construction

Lets assume there are n companies in a city AAA, each emitting E1, E2, ... En. The emissions by humans in the city can be assumed to be Eh. Now,

Ec = Eh + E1 + E2 + ... + En
Cc = E1 + E2 + ... + En

Lets assign unique identifiers to each Ei and call them Ui. On a nation level, we would need to go through these unique identifiers and then sum them exactly once in order to solve the double accounting problem. Companies on the other hand, need to identify which of these assets belong to them and would need to sum over them in order to find their net emissions. So the problem is, in effect reduced to creating and tracking these individual assets and knowing where these assets belong to / whom these assets belong to.

We define a structure with the following keys used to identify the asset:
```
Asset Name (A)
Latitude (LAT)
Longitude (LNG)
City (C)
State (S)
Country (CO)
Company (CMP)
Emissions (E)
```

And take the hash `H = H(A || LAT || LNG || C || S || CO || CMP || E)`. We define a structure for reporting as follows:

```
Hash (H)
City (C)
State (S)
Country (CO)
Country (CO)
Company Name (CMP)
```

The company field is only if the specific entity is a company. When a nation / company wants to report its emissions, it references the hashes of all the assets contained in the country. Keeping keeping track of all the hashes however, is inefficient and can be improved by replacing this with a merkle tree.

```
             H(n)
            /    \
        H(West) H(East)
         /  \     /  \
       H(N) H(S) H(N) H(S)
```

H(N) refers to Hash(North) and H(S) refers to Hash(South). The keys would further expand as:

```
            H(N)
          /      \
        H(W)     H(E)
       /   \      |   \
    Wash.  Oreg.  Mon. Wyom.
```

and each of the states would spin into structures containing cities within them, etc. An entity which wants to report its emissions can just report the merkle hash at its level (eg. if Washington state wants to report its emissions, it reports H(Washington)). When the nation (in this example, the US) wants to report its emissions, it updates the hashes in its merkle tree and publishes them. Verification at each level is quite easy - assume we want to verify that Seattle's emissions are correct having the emissions of the US at hand. We would

1. Get H(US), traverse to the left branch
2. Get H(W) for the West States
3. Get H(W) again since Washington state is in the west among the West states
4. Get H(Wash.) since Seattle is in Washington state
5. Get H(Seattle) depending on how Washington state chooses to structure its merkle tree.

A construction for companies would be similar but the structure of the tree would vary. Taking the example of Amazon and Seattle,

```
          H(Amazon)
        /          \
      H(W)         H(E)
     /   \         /  \
  H(EU)  E(A)
        /    \
      H(NA)  H(SA)
     /   \
H(US||C) H(M||O)
   /   \
 H(US)  H(C)
```

followed by a tree similar to the country tree above, this time depending on where Amazon's locations in the US are. Again, each entity at a level can report its emissions by just reporting the hash of the merkle tree at the level (eg. Amazon EU can report its emissions by reporting H(EU))

One might note that this doesn't take care of emission scopes (Scope 1, 2 and 3). This data can be put into any level but its better to maintain separate trees for each scope to enable easy accounting. As a result, the overall structure for a level would look something like

```
Hash (H)
City (C)
State (S)
Country (CO)
Country (CO)
Company Name (CMP)
Scope1Hash
Scope2Hash
Scope3Hash
```

with each level filled by entities that report at these levels.

This does not take care of verification of emissions themselves, this is left to oracles and other entities which may be willing to take this up. This doesn't also prevent countries from omitting certain regions when reporting their emissions but such practices can be easily identified by traversing through the merkle tree. The structure of the merkle tree can vary between entity and entity. Assuming that data is stored on ipfs / on the blockchain, the root hash will be of constant length across reporting entities.

An interesting thing here is to explore unique identities with the help of Decentralized IDs, which could help give us hashes to track for reporting without having to construct a tree ourselves. This ID could also be used in other places. In essence, that is exactly what a merkle tree accomplishes (giving a hash based identifier to each emitting unit)

These hashes could either be stored in a database system (decentralized / centralized) or could be part of a bigger smart contract that is used for other purposes. In the context of a smart contract on Ethereum, the hashes at each level will be state variables, stored on the blockchain and given the structure of the merkle tree, the contract can verify whether the entity is omitting certain reports. At the same time, with the help of something like the Paris Agreement Smart Contract, it could verify whether the emissions reported are inline with the country's NDCs.
