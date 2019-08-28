pragma solidity ^0.5.0;

import "./GRC20Mintable.sol";

/**
 * @dev Extension of `GRC20Mintable` that adds a cap to the supply of tokens.
 */
contract GRC20Capped is GRC20Mintable {
    uint256 private _cap;

    /**
     * @dev Sets the value of the `cap`. This value is immutable, it can only be
     * set once during construction.
     */
    constructor (uint256 cap) public {
        require(cap > 0, "GRC20Capped: cap is 0");
        _cap = cap;
    }

    /**
     * @dev Returns the cap on the token's total supply.
     */
    function cap() public view returns (uint256) {
        return _cap;
    }

    /**
     * @dev See `GRC20Mintable.mint`.
     *
     * Requirements:
     *
     * - `value` must not cause the total supply to go over the cap.
     */
    function _mint(address account, uint256 value) internal {
        require(totalSupply().add(value) <= _cap, "GRC20Capped: cap exceeded");
        super._mint(account, value);
    }
}
