# Betting Game API

This is an API made on `NodeJs` using `Express` and `Web3`. The goal is to simulate a Betting Game where several users can make bets and one of them will be randomly selected as a winner. To simulate this game the [HolderERC20](https://github.com/Adhara-Tech/claubv23-backend-code-assessment/tree/feature/api/HolderERC20) Smart Contract will be used. The API allows to make bets, check the registered users, list the active bets, and check a user balance and a user's active bets.

## Quick Start
To use this API have to configure the blockchain URL and the Smart Contract Address, and also have to define the contract owner address and private key, this will be the transactions signer.

After that, 
- run ` yarn install` to download the `node_modules` folder.
- run `node app.js ` to run the server.


## Endpoints

##### GET active holds

The Get active holds endpoint is `/holds/`, it will return the active holds.

##### GET registered users

The Get registered users endpoint is `/users/`, it will return the registered users (`usersId`, `userAddress`).


##### GET a user balance

The Get a user balance is `/balance/{id}`, it will return the balance of the defined user.

Sample endpoint with parameters:  `/balance?id=1/`


##### GET a user holds

The Get user holds endpoint is `/user/holds/{ids}`, it will return the active holds of the defined user.

Sample endpoint with parameters:  `/user/holds?id=1/`

##### POST register a user

The Post register a user endpoint is `/reg`, it will receive on the body the user address(`addr`) and will return the user Id. On this endpoint if is the first time the user is registered 100 `HoldERC20` tokens will be minted on the user address and the transaction hash will also be returned.

Sample endpoint :  `/reg/`

Sample body: ``` {
    "add": "0xeAf0dE96EBe1F9cA3dA8969905C0b6CD930BaeAE"
    }```

##### POST a user make a bet

The Post to a user make a bet endpoint is `/bet`, it will receive on the body a the user id to make a bet. This endpoint will send a transaction to the `HolderERC20` to hold the tokens, the user must allow the contract owner's address the amount to bet, or the transaction will fail.

Sample endpoint :  `/bet/`

Sample body: ``` {
    "userId": "1"
    }```

If the current bet is the 4th bet the contract will execute the hold (will send the held amount to the contract's owner address) and will randomly define a winner. The contract's owner will send the betted amount to the winner.


## Improvements

The first important needed feature is to connect a DB to store the information, actually, the values are stored statically and the information will not persist.

Also, a front-end application is needed to have a user-friendly app ad to integrate with `Metamask` to allow the user to sign transactions.

The current restriction could be removed, and allow bets with more than 4 users, or bets bigger than 5.

Other functionalities could be added like: 
- 1vs1 bets when two users bet as much money as they want and with 50% of provability the winner takes all the tokens.
- Show the historical performance of a player, all the bets, all the wins, the losses the total bet amount, and the won amount.
