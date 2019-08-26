pragma solidity ^0.5.0;

import "../../ownership/Ownable.sol";
import "../../access/roles/MinterRole.sol";
import "../../access/roles/PauserRole.sol";


/**
 * @dev Extension of `ERC20` that adds a set of accounts with the `MinterRole`,
 * which have permission to mint (create) new tokens as they see fit.
 *
 * At construction, the deployer of the contract is the only minter.
 */
contract GRC20TransferOSabel is Ownable, MinterRole, PauserRole {
    
    function transferGRCOwnership(address newOwner) public onlyOwner {
        //1.转移合约所有权
        transferOwnership(newOwner);
        //2.转移增发权
        addMinter(newOwner);
        renounceMinter();
        //3.转移禁止流通/开始流通权
        addPauser(newOwner);
        renouncePauser();
    }

}
 