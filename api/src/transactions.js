const contractABI = require('./contract-abi.json');
const contractAddress = '0xE8e23A1A2c1B462013d883246248B2daFdebF322';

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

  return await callTransaction(data, gasLimit);
};

const holdFrom = async (holder) => {
  console.log('---HOLD-FROM---');

  // call data
  var data = contract.methods.holdFrom(holder, web3.utils.numberToHex(5), myAddress).encodeABI();

  // send the signed transaction
  let txHash = await callTransaction(data, gasLimit);

  let receipt = await web3.eth.getTransactionReceipt(txHash.transactionHash);

  // get the logs
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
  return await callTransaction(data, gasLimit);
};

const transfer = async (to) => {
  console.log('---TRANSFER---');

  // call data
  var data = contract.methods.transfer(to, web3.utils.numberToHex(20)).encodeABI();

  return await callTransaction(data, gasLimit);
};

const getBalance = async (address) => {
  console.log('---GET-BALANCE---');

  // call data
  const rawBalance = await contract.methods.balanceOf(address).call();

  return rawBalance;
};

module.exports = { mint, holdFrom, executeHold, transfer, getBalance };