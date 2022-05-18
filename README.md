# Adhara Coding Assessment 

## Implement an improved ERC20 Contract supporting balances on hold

Implement a Holdable ERC20 contract. Provide the right tests to ensure the contract works as expected. Use ganache for the development and testing process.

### Guidelines

Use next ERC20 interface as starting point 
https://github.com/OpenZeppelin/openzeppelin-contracts/blob/master/contracts/token/ERC20/IERC20.sol

New operations need to be provided: 
- hold 
- holdFrom 

These methods are similar to transfer and transferFrom with an extra field (holdId that must be unique).

Hold operation is a special type of transfer in 2 steps. When the hold is created the balance is moved to on hold. The balance on hold can not be used for transfers but it is part of the total supply.

An extra method is needed to execute the hold: 
- executeHold(holdId) 

When the method is executed the money is transferred from the on hold balance to the wallet where the hold was instructed. It can only be executed by the same address that created the hold.
- removeHold(holdId)

Returns the money on hold and all the balance is available for transfers. This function can only be called by the owner of the contract.

## Implement a simple backend providing some endpoints to interact with the contract 

The backend will provide a simple betting game:
1. Offer a service to register a user. Once registered the user balance starts in 100 b. As a user you can play the game betting 5 tokens (fixed, always that value). When the bet is done, the money of the user wallet is placed on hold using the contract implemented before(to be transferred to a wallet owned by the game/owner). When there are 4 players the holds are executed and a winner is selected (random 1..4) and all the money is transferred to the winner. 
2. As a user I should be able to check my bets and my balance. 
Implement the service in your prefered language including tests to ensure it works as expected. 
Write a README.md file explaining the approach followed to design and build the smart contract and also the backend. Try to cover patterns, testability, extensibility, good practices, etc… Also include improvements that you consider that can be done that were not implemented.
