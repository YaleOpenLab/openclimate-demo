# Climate Explorer

Climate Explorer is an idea to have a searchable interface where people can query transactions and pledges that were committed onto the blockchain. For a chain like Ethereum, we could have a smart contract that keeps track of this for us and we can query the smart contract directly but for a chain like bitcoin, we would need to go through the transaction history each time to figure out which transactions are related to open climate.

This would require that we run an indexing / explorer node in order to watch transactions and then publish those transactions that are related to climate change. These transactions can then be shown on a UI similar to other block explorers.

This would also mean that we need to tag the outputs of open climate in a way that is easily queryable by the explorer node. This can be done with the use of something like an OP_RETURN output in bitcoin which can be tracked or by burning funds to a specific address in the case of Ethereum. Other methods depending on their complexity and fees can also be adopted. An interesting thing would be for this to be made customizable so that parties who are willing to pay more for a custom search output can do so and then the explorer node should be able to add this condition in.

An additional thing that would be interesting is for oracles to certify that pledges related to certain transactions were fulfilled and these can also be tracked and displayed in a dashboard. These oracles would have to be validated by the main smart contract so might be something that's restricted to chains with Turing complete VM solutions.

### Applications

1. Searching for pledges
2. Searching for future commitments (similar to other ideas in `docs/`)
3. Comparing two different commitments using the explorer UI
4. Preparing a dashboard (commitments in the last month, last year, average number of commitments, etc)
5. Tracking oracle commitments and progress
