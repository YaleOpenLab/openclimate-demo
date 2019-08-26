pragma solidity ^0.5.0;

import "./GRC20.sol";
import "../../access/roles/MinterRole.sol";

/**
 * @dev Extension of `GRC20` that adds a set of accounts with the `MinterRole`,
 * which have permission to mint (create) new tokens as they see fit.
 *
 * At construction, the deployer of the contract is the only minter.
 */
contract GRC20Mintable is GRC20, MinterRole {
    /**
     * @dev See `GRC20._mint`.
     *
     * Requirements:
     *
     * - the caller must have the `MinterRole`.
     */
    function mint(address account, uint256 amount) public onlyMinter returns (bool) {
        _mint(account, amount);
        return true;
    }
}
