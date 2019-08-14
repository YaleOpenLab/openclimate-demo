pragma solidity 0.5.3;


/**
 After oracle verification and commit to the IPFS we need to save Root Merkle Trie hashes on chain to verify the data.
*/

contract IPFSMerklelRoot {

    // General struct to save Merkle Trie Root on chain.
    struct Root {
        bytes32 rootHash;
        uint index;
    }

    // we will store IPFS hashes by timestamps in UNIX format (integer).
    mapping(uint => Root) private Roots;
    uint[] private dataIndex;

    // emit after inserting the root
    event LogNewRoot (uint indexed timeStamp, uint index, bytes32 rootHash);

    // Util safety check
    function checkTimeStamp(uint timeStamp)
    public
    view
    returns(bool isIndeed)
    {
        if(dataIndex.length == 0) return false;
        return (dataIndex[Roots[timeStamp].index] == timeStamp);
    }

    // insert new Root with the timeStamp
    function insertRoot(
        uint timeStamp,
        bytes32 rootHash)
    public
    returns(uint index)
    {
        if(checkTimeStamp(timeStamp)) revert();
        Roots[timeStamp].rootHash         = rootHash;
        Roots[timeStamp].index             = dataIndex.push(timeStamp)-1;
        emit LogNewRoot(
            timeStamp,
            Roots[timeStamp].index,
            rootHash);
        return dataIndex.length-1;
    }
    // get Root by a timeStamp
    function getRoot(uint timeStamp)
    public
    view
    returns(bytes32 rootHash, uint index)
    {
        if(!checkTimeStamp(timeStamp)) revert();
        return(
        Roots[timeStamp].rootHash,
        Roots[timeStamp].index);
    }

    function getRootCount()
    public
    view
    returns(uint count)
    {
        return dataIndex.length;
    }

    function getRootAtIndex(uint index)
    public
    view
    returns(uint timeStamp)
    {
        return dataIndex[index];
    }
}
