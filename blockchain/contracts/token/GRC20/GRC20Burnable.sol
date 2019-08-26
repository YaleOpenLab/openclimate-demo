pragma solidity ^0.5.0;

import "./GRC20.sol";

/**
 * @dev Extension of `GRC20` that allows token holders to destroy both their own
 * tokens and those that they have an allowance for, in a way that can be
 * recognized off-chain (via event analysis).
 */
contract GRC20Burnable is GRC20 {
    /**
     * @dev Destoys `amount` tokens from the caller.
     *
     * See `GRC20._burn`.
     */
    function burn(uint256 amount) public {
        _burn(msg.sender, amount);
    }

    /**
     * @dev See `GRC20._burnFrom`.
     */
    function burnFrom(address account, uint256 amount) public {
        _burnFrom(account, amount);
    }
}
