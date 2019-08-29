# Draft ideas of smart contract applications to climate change markets

### The Problem

Existing climate change mitigation projects are aplenty and for a person who wants to offset his carbon emissions by a specific amount, there is no single place to do so. Definitions of offset also sometimes vary between this project, making it difficult for someone who's new to try and act on climate change.

### The Construction

There are two spinoff ideas which could be really useful:

1. The creation of climate hedge markets
2. The creation of secondary markets to trade climate assets

#### Climate Hedge Markets

The idea is to hedge a proposition X something against the future predicting that X will be affected by climate change affected by parameters Y1, Y2, ... Yn. More concretely, lets assume proposition X as:
```
In 2050, Kiribati will be submerged if climate change continues at 2019 rates
```
One side of the market would bet for the proposition and the other side would oppose this. To incentivise people to try and contribute to preventing climate change, we can incentivise one side by promising z% returns if Kiribati does exist in 2050. This would help the country for eg, to have mitigation mechanisms in place if climate change does adversely impact it. If the country survives, it pays back people who invested in the same. If the country doesn't, it can still use the money to migrate people to a different country.

This shouldn't be on something like Augur, where people vote on outcomes. Instead, the outcome must be ratifies by a central oracle like the UN to prevent ambiguity. This introduces a degree of central control to the problem but the bigger issue involves governments, so this doesn't add much to centralization.

#### Secondary Climate Markets

Primary climate markets are involved with purchasing Carbon Emission Mitigation equivalents (like RECs). Secondary climate markets can be used by people to trade arbitrary (including primary) climate assets. For eg, say someone wants to plant trees in their farm and wants to issue an asset for it. Traditional schemes won't work here since its expensive, cost and time wise, to try and certify the project. Instead, they can issue their own asset and let the market invest in their asset.

The owner creates an asset and fixes a price `t` for it on the secondary market. People who want to offset their Carbon Emissions can choose to bid on this asset `t` with their own price `t'` and depending on the number of bids, the owner can reduce / increase the price of the asset. This could also be a smart contract in the form of Dutch auctions, where price adjusts according to asset demand.

One should note here that secondary assets are not exchangeable for primary ones whereas the reverse is possible. Some secondary assets have no value on the primary market but may have value on the secondary market. Some secondary assets may transition to become primary assets in the future.

#### Monetary Reward

The idea is to link climate change with monetary rewards since in the past, movements that have had no monetary reward have failed to sustain themselves in the short/long run. By having a scheme for people to link climate change with monetary rewards, one can raise the importance of mitigating climate change without losing momentum due to lack of capital
