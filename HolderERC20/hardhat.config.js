require('@nomiclabs/hardhat-waffle');
require('@nomiclabs/hardhat-ethers');
require('@nomiclabs/hardhat-waffle');
require('@nomiclabs/hardhat-web3');
require('solidity-coverage');
require('hardhat-gas-reporter');
require('hardhat-contract-sizer');
require('hardhat-deploy');

require('dotenv').config();

const MNEMONIC = process.env.MNEMONIC;
const PRIVATE_KEY = process.env.PRIVATE_KEY;
const MUMBAI_RPC_URL = process.env.MUMBAI_RPC_URL;
const INFURA_API_KEY = process.env.INFURA_API_KEY;
const TENDERLY_PROJECT = process.env.TENDERLY_PROJECT;
const TENDERLY_USERNAME = process.env.TENDERLY_USERNAME;
const INFURA_POLYGON_RPC_URL = process.env.INFURA_POLYGON_RPC_URL;


module.exports = {
  networks: {
    hardhat: {
      tags: ['local'],
      allowUnlimitedContractSize: true,
    },
    localhost: {
      url: 'http://127.0.0.1:8545',
      tags: ['local'],
    },
    mumbai: {
      url: MUMBAI_RPC_URL,
      accounts: [PRIVATE_KEY],
      tags: ['testnet'],
    },
    polygon: {
      url: INFURA_POLYGON_RPC_URL,
      accounts: [PRIVATE_KEY],
      tags: ['mainnet'],
    },
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
