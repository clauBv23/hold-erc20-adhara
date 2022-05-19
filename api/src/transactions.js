require('dotenv').config();
const MY_ADDRESS = process.env.MY_ADDRESS;
const PRIVATE_KEY = process.env.PRIVATE_KEY;
const PROVIDER_URL = process.env.PROVIDER_URL;
const CONTRACT_ADDRESS = process.env.CONTRACT_ADDRESS;

const contractABI = require('./contract-abi.json');
const contractAddress = CONTRACT_ADDRESS;

console.log('contract===================', contractAddress);
const Web3 = require('web3');
const ethTx = require('ethereumjs-tx');
const readline = require('readline');

// provider url
var provider = PROVIDER_URL;
console.log('******************************************');
console.log('Using provider : ' + provider);
console.log('******************************************');

var web3 = new Web3(new Web3.providers.HttpProvider(provider));
web3.transactionConfirmationBlocks = 1;

//my address and private key
const myAddress = MY_ADDRESS;
const privKey = Buffer.from(PRIVATE_KEY, 'hex');

// contract instance
var contract = new web3.eth.Contract(contractABI, contractAddress);

async function callTransaction(data) {
  // calculate the gas limit
  let gasLimit = await web3.eth.estimateGas({
    to: contractAddress,
    data: data,
  });

  // Get the address transaction count in order to specify the correct nonce
  let txnCount = await web3.eth.getTransactionCount(myAddress, 'pending');

  // Create the transaction object
  var txObject = {
    nonce: web3.utils.numberToHex(txnCount),
    gasPrice: web3.utils.numberToHex(1000),
    gasLimit,
    to: contractAddress,
    value: '0x00',
    data,
  };

  // Sign the transaction with the private key
  var tx = new ethTx(txObject);
  tx.sign(privKey);

  //Convert to raw transaction string
  var serializedTx = tx.serialize();
  var rawTxHex = '0x' + serializedTx.toString('hex');

  // send the signed transaction
  let txHash = await web3.eth.sendSignedTransaction(rawTxHex).catch((error) => {
    console.log('Error: ', error.message);
  });

  return txHash;
}

const mint = async (to) => {
  console.log('---MINT---');

  // call data
  var data = contract.methods.mint(to, web3.utils.numberToHex(100)).encodeABI();

  return await callTransaction(data);
};

const holdFrom = async (holder) => {
  console.log('---HOLD-FROM---');

  // call data
  var data = contract.methods.holdFrom(holder, web3.utils.numberToHex(5), myAddress).encodeABI();

  // send the signed transaction
  let txHash = await callTransaction(data);

  let receipt = await web3.eth.getTransactionReceipt(txHash.transactionHash);

  // get the logs the Event HoldCreated hash = 0x1f04d8ba13156fb73e621b6df1a4a7aebc25167f7efbd455c45dfc4a3bbea61c
  let log = receipt.logs.filter((element) => {
    return element.topics[0] == '0x1f04d8ba13156fb73e621b6df1a4a7aebc25167f7efbd455c45dfc4a3bbea61c';
  })[0];

  // parse the log to get the event values
  let hold = web3.eth.abi.decodeLog(
    [
      {
        indexed: true,
        name: 'holdId',
        type: 'uint256',
      },
      {
        name: 'holdAmount',
        type: 'uint256',
      },
      {
        name: 'holder',
        type: 'address',
      },
      {
        name: 'operator',
        type: 'address',
      },
    ],
    log.data,
    [log.topics[1]]
  );

  // return the transaction and the created hold id
  return { tx: txHash, id: hold.holdId };
};

const executeHold = async (holdId) => {
  console.log('---EXECUTE-HOLD---');

  // call data
  var data = contract.methods.executeHold(holdId).encodeABI();

  // call to send the signed transaction
  return await callTransaction(data);
};

const transfer = async (to) => {
  console.log('---TRANSFER---');

  // call data
  var data = contract.methods.transfer(to, web3.utils.numberToHex(20)).encodeABI();

  return await callTransaction(data);
};

const getBalance = async (address) => {
  console.log('---GET-BALANCE---');

  // call data
  const rawBalance = await contract.methods.balanceOf(address).call();

  return rawBalance;
};

module.exports = { mint, holdFrom, executeHold, transfer, getBalance };
