const HolderERC20 = artifacts.require('HolderERC20');

const { BN, expectRevert, expectEvent } = require('@openzeppelin/test-helpers');

contract('HolderERC20', function () {
  const mint = async (to, amount) => {
    await this.holderErc20.mint(to, amount);
  };

  beforeEach(async () => {
    await deployments.fixture(['holder_erc20']);
    let deployment = await deployments.get('HolderERC20');

    this.holderErc20 = await HolderERC20.at(deployment.address);
  });

  it('should be deployed', async () => {
    assert.isOk(this.holderErc20.address);
  });

  describe('check HolderERC20 contract', async () => {
    beforeEach(async () => {
      const { deployer, bob, alice } = await getNamedAccounts();
      // mint tokens to the accounts
      await mint(deployer, 10);
      await mint(bob, 10);
      await mint(alice, 10);
    });

    it('check mint method', async () => {
      const { deployer, bob, alice } = await getNamedAccounts();
      let aliceBalance = await this.holderErc20.balanceOf(alice);
      let deployerBalance = await this.holderErc20.balanceOf(deployer);
      let bobBalance = await this.holderErc20.balanceOf(bob);

      assert.strictEqual(aliceBalance.toString(), '10');
      assert.strictEqual(deployerBalance.toString(), '10');
      assert.strictEqual(bobBalance.toString(), '10');
    });

    describe('check hold method', async () => {
      beforeEach(async () => {
        const { alice } = await getNamedAccounts();

        // hold 3 tokens
        this.holdTX = await this.holderErc20.hold(3, { from: alice });
      });

      it('check the HoldCreated event where emited', async () => {
        const { alice } = await getNamedAccounts();

        expectEvent(this.holdTX, 'HoldCreated', { holdAmount: '3', holder: alice, operator: alice });
      });

      it('check hold where stored', async () => {
        const { alice } = await getNamedAccounts();

        let aliceHold = await this.holderErc20.getHold(1);

        assert.strictEqual(aliceHold[0].toString(), '3');
        assert.strictEqual(aliceHold[1], alice);
        assert.strictEqual(aliceHold[2], alice);
      });

      it('check the tokens were holded', async () => {
        const { alice } = await getNamedAccounts();

        let aliceBalance = await this.holderErc20.balanceOf(alice);
        let contractBalance = await this.holderErc20.balanceOf(this.holderErc20.address);

        assert.strictEqual(aliceBalance.toString(), '7');
        assert.strictEqual(contractBalance.toString(), '3');
      });

      describe('check execute hold method', async () => {
        beforeEach(async () => {
          const { alice } = await getNamedAccounts();

          // hold 5 tokens
          this.holdTX = await this.holderErc20.hold(3, { from: alice });

          // execute the second hold created
          this.holdTX = await this.holderErc20.executeHold(2, { from: alice });
        });

        it('check the HoldExecuted event where emited', async () => {
          expectEvent(this.holdTX, 'HoldExecuted', { holdId: '2' });
        });

        it('check hold where deleted', async () => {
          await expectRevert(this.holderErc20.getHold(2), 'Undefined hold');
        });

        it('check the tokens were reverted', async () => {
          const { alice } = await getNamedAccounts();

          let aliceBalance = await this.holderErc20.balanceOf(alice);
          let contractBalance = await this.holderErc20.balanceOf(this.holderErc20.address);

          assert.strictEqual(aliceBalance.toString(), '7');
          assert.strictEqual(contractBalance.toString(), '3');
        });
      });
    });

    describe('check hold from method', async () => {
      beforeEach(async () => {
        const { alice, bob } = await getNamedAccounts();

        // hold 5 tokens from alice with bob
        this.holderErc20.approve(bob, 5, { from: alice });
        this.holdTX = await this.holderErc20.holdFrom(alice, 5, { from: bob });
      });

      it('check the HoldCreated event where emited', async () => {
        const { alice, bob } = await getNamedAccounts();

        expectEvent(this.holdTX, 'HoldCreated', { holdAmount: '5', holder: alice, operator: bob });
      });

      it('check hold where stored', async () => {
        const { alice, bob } = await getNamedAccounts();

        let aliceHold = await this.holderErc20.getHold(1);

        assert.strictEqual(aliceHold[0].toString(), '5');
        assert.strictEqual(aliceHold[1], alice);
        assert.strictEqual(aliceHold[2], bob);
      });

      it('check the tokens were transferred', async () => {
        const { alice, bob } = await getNamedAccounts();

        let aliceBalance = await this.holderErc20.balanceOf(alice);
        let bobBalance = await this.holderErc20.balanceOf(bob);

        let contractBalance = await this.holderErc20.balanceOf(this.holderErc20.address);

        assert.strictEqual(aliceBalance.toString(), '5');
        assert.strictEqual(bobBalance.toString(), '10');
        assert.strictEqual(contractBalance.toString(), '5');
      });

      describe('check execute hold method when the operator is not the holder', async () => {
        beforeEach(async () => {
          const { alice, bob } = await getNamedAccounts();

          // hold 1 token from alice with bob
          this.holderErc20.approve(bob, 1, { from: alice });
          this.holdTX = await this.holderErc20.holdFrom(alice, 1, { from: bob });

          // chech that alice can't execute the hold

          await expectRevert(this.holderErc20.executeHold(2, { from: alice }), 'The caller must be the hold creator');

          // execute the second hold created
          this.holdTX = await this.holderErc20.executeHold(2, { from: bob });
        });

        it('check the HoldExecuted event where emited', async () => {
          expectEvent(this.holdTX, 'HoldExecuted', { holdId: '2' });
        });

        it('check hold where deleted', async () => {
          await expectRevert(this.holderErc20.getHold(2), 'Undefined hold');
        });

        it('check the tokens were reverted', async () => {
          const { alice } = await getNamedAccounts();

          let aliceBalance = await this.holderErc20.balanceOf(alice);
          let contractBalance = await this.holderErc20.balanceOf(this.holderErc20.address);

          assert.strictEqual(aliceBalance.toString(), '5');
          assert.strictEqual(contractBalance.toString(), '5');
        });
      });
    });

    describe('check remove hold method', async () => {
      beforeEach(async () => {
        const { deployer, alice } = await getNamedAccounts();

        // hold 5 tokens with alice
        await this.holderErc20.hold(2, { from: alice });

        // remove the hold
        this.removeTX = await this.holderErc20.removeHold(1, { from: deployer });
      });

      it('check the HoldRemoved event where emited', async () => {
        expectEvent(this.removeTX, 'HoldRemoved', { holdId: '1' });
      });

      it('check hold where deleted', async () => {
        await expectRevert(this.holderErc20.getHold(1), 'Undefined hold');
      });

      it('check the tokens where reverted', async () => {
        const { alice } = await getNamedAccounts();

        let aliceBalance = await this.holderErc20.balanceOf(alice);
        let contractBalance = await this.holderErc20.balanceOf(this.holderErc20.address);

        assert.strictEqual(aliceBalance.toString(), '10');
        assert.strictEqual(contractBalance.toString(), '0');
      });
    });

    // -----------------requires------------------
    it('check hold from when is not allowed', async () => {
      const { alice, bob } = await getNamedAccounts();

      // hold 5 tokens from alice with bob with no allowance
      await expectRevert(this.holderErc20.holdFrom(alice, 5, { from: bob }), 'Not allowed');
    });

    it('check the removeHold can only be called by the owner', async () => {
      const { bob } = await getNamedAccounts();

      await expectRevert(this.holderErc20.removeHold(1, { from: bob }), 'Ownable: caller is not the owner');
    });

    it('check the removeHold method reverts if the hold is not defined', async () => {
      const { deployer } = await getNamedAccounts();

      await expectRevert(this.holderErc20.removeHold(1, { from: deployer }), 'Undefined hold');
    });
  });
});
