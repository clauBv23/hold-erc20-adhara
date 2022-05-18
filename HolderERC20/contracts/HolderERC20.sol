//SPDX-License-Identifier: Unlicense
pragma solidity ^0.8.6;

import "@openzeppelin/contracts-upgradeable/token/ERC20/ERC20Upgradeable.sol";
import "@openzeppelin/contracts/utils/Counters.sol";

contract HolderERC20 is ERC20Upgradeable {
    using Counters for Counters.Counter;

    constructor() {}
}
