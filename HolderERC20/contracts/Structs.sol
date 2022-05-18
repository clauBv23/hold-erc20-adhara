// SPDX-License-Identifier: MIT
pragma solidity ^0.8.6;

/**
 * @dev The hold representation
 * @param operator The creator of the hold
 * @param holder The holder account
 * @param amount The hold amount
 */
struct Hold {
    address operator;
    address holder;
    uint256 amount;
}
