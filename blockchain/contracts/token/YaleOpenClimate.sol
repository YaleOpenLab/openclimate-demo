pragma solidity ^0.5.0;

import './GRC20/GRC20Mintable.sol';
import './GRC20/GRC20Pausable.sol';
import './GRC20/GRC20Burnable.sol';
import './GRC20/GRC20TransferOSable.sol';
 
contract YaleOpenClimate is GRC20Mintable, GRC20Pausable, GRC20Burnable, GRC20TransferOSabel{ 
    string private _name;
    string private _symbol;
    uint8 private _decimals;

    /**
     * @dev Sets the values for `name`, `symbol`, and `decimals`. All three of 
     * these values are immutable: they can only be set once during
     * construction.
     */
    constructor (string memory name, string memory symbol, uint8 decimals, uint256 totalSupply, address owner) public {
        _name = name;
        _symbol = symbol;
        _decimals = decimals;
        _totalSupply = totalSupply;
        _balances[owner] = totalSupply;
        if (msg.sender != owner){
        transferGRCOwnership(owner); 
        }
       
    }

    /**
     * @dev Returns the name of the token.
     */
    function name() public view returns (string memory) {
        return _name;
    }

    /**
     * @dev Returns the symbol of the token, usually a shorter version of the
     * name.
     */
    function symbol() public view returns (string memory) {
        return _symbol;
    }

    /**
     * @dev Returns the number of decimals used to get its user representation.
     * For example, if `decimals` equals `2`, a balance of `505` tokens should
     * be displayed to a user as `5,05` (`505 / 10 ** 2`).
     *
     * Tokens usually opt for a value of 18, imitating the relationship between
     * Ether and Wei.
     *
     * > Note that this information is only used for _display_ purposes: it in
     * no way affects any of the arithmetic of the contract, including
     * `IGRC20.balanceOf` and `IGRC20.transfer`.
     */
    function decimals() public view returns (uint8) {
        return _decimals;
    }

}