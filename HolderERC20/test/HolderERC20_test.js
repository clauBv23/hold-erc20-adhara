const HolderERC20 = artifacts.require('HolderERC20');

contract('HolderERC20', function () {
  beforeEach(async () => {
    await deployments.fixture(['holder_erc20']);
    let deployment = await deployments.get('HolderERC20');

    this.holderErc20 = await HolderERC20.at(deployment.address);
  });

  it('should be deployed', async () => {
    assert.isOk(this.holderErc20.address);
  });
});
