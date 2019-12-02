pragma solidity ^0.5.10;

import "../parisagreement/paris_agreement.sol";

contract DAOInterface {
    
    struct Proposal {
        // A plain text description of the proposal
        string description;
        // A unix timestamp, denoting the end of the voting period
        uint votingDeadline;
        // True if the proposal's votes have yet to be counted, otherwise False
        bool open;
        // True if quorum has been reached, the votes have been counted, and
        // the majority said yes
        bool proposalPassed;
        // A hash to check validity of a proposal
        bytes32 proposalHash;
        // Number of GHG offset Tokens in favor of the proposal
        uint yea;
        // Number of GHG offset Tokens opposed to the proposal
        uint nay;
        // Simple mapping to check if a shareholder has voted for it
        mapping (address => bool) votedYes;
        // Simple mapping to check if a shareholder has voted against it
        mapping (address => bool) votedNo;
        // Address of the shareholder who created the proposal
        address creator;
    }
      // Proposals
    Proposal[] public proposals;
    
    function newProposal(
        string _description,
        bytes _transactionData,
        uint _debatingPeriod,
    ) returns (uint _proposalID);
    
    function vote(uint _proposalID, bool _supportsProposal);
    
    function halveMinQuorum() returns (bool _success);
     event ProposalAdded(
        uint indexed proposalID,
        string description
    );
    event Voted(uint indexed proposalID, bool position, address indexed voter);
    event ProposalTallied(uint indexed proposalID, bool result, uint quorum);
}

contract DAO is DAOInterface {
    
     // Import NDCs contract
    DetermedContributions contract_ndcs;
    function ImportNDCS(address _t) public {
        contract_ndcs = DetermedContributions(_t);
    }
    
     // Import PA contract
    ParisAgreement contract_pa;
    function ImportPA(address _t) public {
        contract_pa = ParisAgreement(_t);
    }
    
    function newProposal(
        string _description,
        bytes _transactionData,
        uint64 _debatingPeriod
    ) returns (uint _proposalID) {
        require(contract_ndcs.isCounrty(msg.sender), "Country address isnt registered");

        // to prevent curator from halving quorum before first proposal
        if (proposals.length == 1) // initial length is 1 (see constructor)
            lastTimeMinQuorumMet = now;

        _proposalID = proposals.length++;
        Proposal p = proposals[_proposalID];
        p.description = _description;
        p.proposalHash = sha3(_recipient, _amount, _transactionData);
        p.votingDeadline = now + _debatingPeriod;
        p.open = true;
        //p.proposalPassed = False; // that's default
        p.creator = msg.sender;
        p.proposalDeposit = msg.value;

        sumOfProposalDeposits += msg.value;

        ProposalAdded(
            _proposalID,
            _description
        );
    }
    
    function vote(uint _proposalID, bool _supportsProposal) {

        Proposal p = proposals[_proposalID];

        unVote(_proposalID);

        if (_supportsProposal) {
            p.yea += token.balanceOf(msg.sender);
            p.votedYes[msg.sender] = true;
        } else {
            p.nay += token.balanceOf(msg.sender);
            p.votedNo[msg.sender] = true;
        }

        if (blocked[msg.sender] == 0) {
            blocked[msg.sender] = _proposalID;
        } else if (p.votingDeadline > proposals[blocked[msg.sender]].votingDeadline) {
            // this proposal's voting deadline is further into the future than
            // the proposal that blocks the sender so make it the blocker
            blocked[msg.sender] = _proposalID;
        }

        votingRegister[msg.sender].push(_proposalID);
        Voted(_proposalID, _supportsProposal, msg.sender);
    }
    
     function unVote(uint _proposalID){
        Proposal p = proposals[_proposalID];

        if (now >= p.votingDeadline) {
            throw;
        }

        if (p.votedYes[msg.sender]) {
            p.yea -= token.balanceOf(msg.sender);
            p.votedYes[msg.sender] = false;
        }

        if (p.votedNo[msg.sender]) {
            p.nay -= token.balanceOf(msg.sender);
            p.votedNo[msg.sender] = false;
        }
    }
    
    function closeProposal(uint _proposalID) internal {
        Proposal p = proposals[_proposalID];
        if (p.open)
            sumOfProposalDeposits -= p.proposalDeposit;
        p.open = false;
    }
    
    function minQuorum(uint _value) internal constant returns (uint _minQuorum) {
        // minimum of 14.3% and maximum of 47.6%
        return token.totalSupply() / minQuorumDivisor +
            (_value * token.totalSupply()) / (3 * (actualBalance()));
    }
    
    function halveMinQuorum() returns (bool _success) {
        // this can only be called after `quorumHalvingPeriod` has passed or at anytime after
        // fueling by the curator with a delay of at least `minProposalDebatePeriod`
        // between the calls
        if ((lastTimeMinQuorumMet < (now - quorumHalvingPeriod) || msg.sender == curator)
            && lastTimeMinQuorumMet < (now - minProposalDebatePeriod)
            && proposals.length > 1) {
            lastTimeMinQuorumMet = now;
            minQuorumDivisor *= 2;
            return true;
        } else {
            return false;
        }
    }

    function numberOfProposals() constant returns (uint _numberOfProposals) {
        // Don't count index 0. It's used by getOrModifyBlocked() and exists from start
        return proposals.length - 1;
    }
  

}
