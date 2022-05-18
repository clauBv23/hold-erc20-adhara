module.exports = async ({ getNamedAccounts, deployments }) => {
  const { deploy } = deployments;
  const { deployer } = await getNamedAccounts();

  // this contract is upgradeable through uups (EIP-1822)
  await deploy('HolderERC20', {
    from: deployer,
    proxy: {
      proxyContract: 'UUPSProxy',
      execute: {
        init: {
          methodName: 'initialize',
          args: ['holderToken', 'HTK'],
        },
      },
    },
    log: true,
    args: [],
  });
};

module.exports.tags = ['holder_erc20'];
