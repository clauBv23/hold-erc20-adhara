# Holder ERC20

The contract is a regular ERC20 using the [EIP20](https://eips.ethereum.org/EIPS/eip-20) and the [OpenZeppelin ERC20](https://docs.openzeppelin.com/contracts/2.x/api/token/erc20) standard The Hold functionality was added to allow the owners of the tokens hold tokens and send them to a beneficiary when the hold is executed. The COntract owner can remove the holds, in that case, the held tokens will go back to the holder.

## New Functions

- <ins>hold</ins>: The hold function receive a hold amount and a beneficiary, once the tokens are held, the holder can't use them. The beneficiary won't be able to use the tokens held in his favor until the hold is executed.
- <ins>holdFrom</ins>: This function is like the hold function but an operator can make the hold. The operator must have an allowance to hold the holder tokens.
- <ins>executeHold</ins>: The function can only be called by the hold creator. Once this function is executed the held funds will go to the beneficiary and it will be able to use those tokens.
- <ins>removeHold</ins>: The remove hold function can only be called by the contract owner. Once this function is called the held tokens will be transferred back to the holder (not the beneficiary).

## Specifications

The HolderErc20 contract is Upgradeable using UUPS Proxy and [Ownble](https://docs.openzeppelin.com/contracts/2.x/api/ownership#Ownable).
The Project was developed using HardHat and OpenZeppelin libraries. The tests are on js, using mocha and chai with waffle, in the tests were also used the OpenZeppelin test helpers to control the events and the reverted transactions. The solidity coverage plugging was used to check that 100% coverage where raised.

### Quick Start
 To use the contract fallow this steps:
 
 - Run ``` yarn install``` to download the `node_modules` folder.
 - Create a `.env` file with the specifications on the `evn.example` file
 - Run ``` hardhat deploy --network <the defined network> ``` to deloy the contract on the desired network

### Utils 
- To check the tests run `hardhat test` on the console
- To check the contract code coverage run on the console ` hardhat coverage`
- The local network (a local ganache network) and the rinkeby network are configured on the project. To add a new one go to the `hardhat.config.js` file and the desired network configuration.


## Improvements

A different Hold function could be added to allow the token owners to save tokens. It will be like a vault where each owner can hold his token to be used in the future and receive some profit depending on the held amount.
Create a Minter Roll and give him permission to Mint, allowing him to mint tokens not only to the contract owner.
