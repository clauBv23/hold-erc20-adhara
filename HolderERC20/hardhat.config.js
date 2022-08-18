require('@nomiclabs/hardhat-truffle5');
require('@nomiclabs/hardhat-waffle');
require('@nomiclabs/hardhat-ethers');
require('@nomiclabs/hardhat-web3');
require('hardhat-contract-sizer');
require('hardhat-gas-reporter');
require('solidity-coverage');
require('hardhat-deploy');
require('dotenv').config();
require('hardhat-abi-exporter');


const PRIVATE_KEY = process.env.PRIVATE_KEY;
const RINKEBY_RPC_URL = process.env.RINKEBY_RPC_URL;

module.exports = {
  networks: {
    hardhat: {
      tags: ['local'],
      allowUnlimitedContractSize: true,
    },
    localhost: {
      url: 'http://127.0.0.1:7545',
      tags: ['local'],
    },
    // rinkeby: {
    //   url: RINKEBY_RPC_URL,
    //   accounts: [PRIVATE_KEY],
    //   tags: ['testnet'],
    // },
  },
  upgradeable: {
    uups: ['HolderERC20'],
  },
  namedAccounts: {
    deployer: 0,
    user: 1,
    bob: 2,
    alice: 3,
  },
  gasReporter: {
    enabled: false,
  },
  mocha: {
    timeout: 999999,
  },
  solidity: {
    compilers: [
      {
        version: '0.8.9',
        settings: {
          optimizer: {
            enabled: true,
            runs: 200,
          },
        },
      },
    ],
  },
};
