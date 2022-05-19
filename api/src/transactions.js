const contractABI = require('./contract-abi.json');
const contractAddress = '0x366084Adae7F4FB8C8506B655eab1ecaD6D17aa5';

const Web3 = require('web3');
const ethTx = require('ethereumjs-tx');
const readline = require('readline');

const args = process.argv.slice(2);

// Ganache url
var provider = args[0] || 'http://localhost:7545';
console.log('******************************************');
console.log('Using provider : ' + provider);
console.log('******************************************');

var web3 = new Web3(new Web3.providers.HttpProvider(provider));
web3.transactionConfirmationBlocks = 1;

//my address and private key
const myAddress = '0x8Dd9d0e36183052402bE9796149312bE23d4DF06';
const privKey = Buffer.from('15878c3bb1b17cc2c3fcf2a8d0e78d470ecc7251d07ae2ca0519b78db017e21d', 'hex');

// contract instance
var contract = new web3.eth.Contract(contractABI, contractAddress);

async function callTransaction(data, gasLimit) {
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
  // call data
  var data = contract.methods.mint(to, web3.utils.numberToHex(100)).encodeABI();

  //gas limit
  let gasLimit = await web3.eth.estimateGas({
    to: contractAddress,
    data: data,
  });
  console.log('---MINT---');
  // return the transaction hash
  return await callTransaction(data, gasLimit);
};

module.exports = { mint };
