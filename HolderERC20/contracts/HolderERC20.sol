//SPDX-License-Identifier: Unlicense
pragma solidity ^0.8.6;

import "@openzeppelin/contracts-upgradeable/token/ERC20/ERC20Upgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts/utils/Counters.sol";

contract HolderERC20 is UUPSUpgradeable, ERC20Upgradeable, OwnableUpgradeable {
    using Counters for Counters.Counter;
    bytes32 public constant MINTER_ROLE = keccak256("MINTER_ROLE");

    constructor() initializer {}

    function initialize() external initializer {
        __Ownable_init();
        __UUPSUpgradeable_init();
    }

    /**
     * @dev Upgrades needed method
     */
    function _authorizeUpgrade(address) internal view override onlyOwner {}
}
