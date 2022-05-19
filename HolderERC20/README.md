# Holder ERC20

The contract is a regular ERC20 using the [EIP20](https://eips.ethereum.org/EIPS/eip-20) and the [OpenZeppelin ERC20](https://docs.openzeppelin.com/contracts/2.x/api/token/erc20) standard. The Hold functionality was added to allow the owners of the tokens to lock them in the contract to later send them to a predefined beneficiary (once the hold is executed). The Contract owner can remove the holds; in that case, the held tokens will go back to the holder.

## New Functions

- <ins>hold</ins>: The hold function receive a hold amount and a beneficiary. Once the tokens are held, the holder can't use them. The beneficiary won't be able to use them either until the hold is executed.
- <ins>holdFrom</ins>: This function is like the hold function, but an operator can make the hold on behalf of the token owner. The operator must have enough allowance.
- <ins>executeHold</ins>: The function can only be called by the hold creator. Once this function is executed, the held funds will go to the beneficiary, and he will be able to use those tokens.
- <ins>removeHold</ins>: The remove hold function can only be called by the contract owner. Once this function is called, the held tokens will be transferred back to the holder (not the beneficiary).

## Specifications

The HolderErc20 contract is `Upgradeable` using `UUPS Proxy` and [`Ownble`](https://docs.openzeppelin.com/contracts/2.x/api/ownership#Ownable). The project was developed using `HardHat` and `OpenZeppelin` libraries. The tests are on `Javascript`, using `mocha` and `chai` with `waffle`. We also used the OpenZeppelin test helpers to control the events and the reverted transactions, alongside the solidity coverage plugging to check that 100% coverage where raised.

### Quick Start

To use the contract, follow these steps:
- Run ``` yarn install``` to download the `node_modules` folder.
- Create a `.env` file with the specifications on the `env.example` file.
- Run ``` hardhat deploy --network <the defined network> ``` to deploy the contract on the desired network.

### Utils

- To check the tests run `hardhat test` on the console.
- To check the contract code coverage, run `hardhat coverage` on the console.
- The local network (a local `Ganache` network) and the `Rinkeby` network are configured on the project. You can go to the `hardhat.config.js` file and add the desired network configuration to configure a new one.

## Improvements

A different hold function could be added to allow the token owners to save tokens. It will be like a vault where each owner can hold their token to be used in the future and receive some profit depending on the held amount. 

Create a Minter role, allowing different accounts to mint tokens, not only the contract owner.
