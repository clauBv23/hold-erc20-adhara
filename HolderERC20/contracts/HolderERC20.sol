//SPDX-License-Identifier: Unlicense
pragma solidity ^0.8.6;

import "@openzeppelin/contracts-upgradeable/token/ERC20/ERC20Upgradeable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";
import "@openzeppelin/contracts/utils/Counters.sol";
import "./Structs.sol";

contract HolderERC20 is UUPSUpgradeable, ERC20Upgradeable, OwnableUpgradeable {
    using Counters for Counters.Counter;

    /// @dev The incremental counter for the holds ids
    Counters.Counter private _holdIdCounter;

    /// @dev Mapping for the registered holds
    mapping(uint256 => Hold) private _holds;

    /**
     * @dev Emmited when a new hold is created
     * @param holdId The created hold Id
     * @param holdAmount The amount held
     * @param holder The holder's address
     * @param operator The creator of the hold
     */
    event HoldCreated(uint256 holdId, uint256 holdAmount, address holder, address operator);

    /**
     * @dev Emmited when a new hold is executed
     * @param holdId The created hold Id
     */
    event HoldExecuted(uint256 indexed holdId);

    /**
     * @dev Emmited when a new hold is removed by the contract owner
     * @param holdId The created hold Id
     */
    event HoldRemoved(uint256 indexed holdId);

    constructor() initializer {}

    function initialize(string calldata name_, string calldata symbol_) external initializer {
        __ERC20_init(name_, symbol_);
        __Ownable_init();
        __UUPSUpgradeable_init();
    }

    /**
     * @notice Only owner allowed to mint
     * @dev Mint current tokens and assign them to a specific account
     * @param to_ Address to owns the minted tokens
     * @param amount_ Amount of tokens to be minted
     */
    function mint(address to_, uint256 amount_) external onlyOwner {
        _mint(to_, amount_);
    }

    /**
     * @notice Hold the specified amount
     * @param amount_ Amount of tokens to hold
     * @param beneficiary_ Account to transfer the held tokens
     */
    function hold(uint256 amount_, address beneficiary_) external {
        _hold(msg.sender, msg.sender, amount_, beneficiary_);
    }

    /**
     * @notice Hold the specified amount from another account
     * @param holder_ Owner of tokens to hold
     * @param amount_ Amount of tokens to hold
     * @param beneficiary_ Account to transfer the held tokens to
     */
    function holdFrom(
        address holder_,
        uint256 amount_,
        address beneficiary_
    ) external {
        uint256 allowance_ = allowance(holder_, msg.sender);

        // check if the holder allow the caller to hold
        require(allowance_ >= amount_, "Not allowed");

        _hold(msg.sender, holder_, amount_, beneficiary_);
    }

    /// @dev Helper function to create the hold
    function _hold(
        address operator_,
        address holder_,
        uint256 amount_,
        address beneficiary_
    ) private returns (uint256 holdId) {
        // start the token id in 1 instead of 0
        _holdIdCounter.increment();
        holdId = _holdIdCounter.current();

        // send the funds to the current contract
        if (operator_ == holder_) {
            transfer(address(this), amount_);
        } else {
            transferFrom(holder_, address(this), amount_);
        }

        // store the hold
        _holds[holdId] = Hold({
            operator: operator_,
            holder: holder_,
            amount: amount_,
            beneficiary: beneficiary_
        });

        emit HoldCreated(holdId, amount_, holder_, operator_);
    }

    /**
     * @notice Returns the held tokens to the holder
     * @param holdId_ Id of the Hold
     */
    function executeHold(uint256 holdId_) external {
        Hold memory currentHold = _holds[holdId_];

        require(currentHold.operator == msg.sender, "The caller must be the hold creator");

        // execute the hold and send the holed tokens to the beneficiary
        _sendHeldTokens(holdId_, currentHold.beneficiary, currentHold.amount);

        emit HoldExecuted(holdId_);
    }

    /**
     * @notice Returns the held tokens to the holder
     * @dev Only owner allowed
     * @param holdId_ Id of the Hold
     */
    function removeHold(uint256 holdId_) external onlyOwner {
        Hold memory currentHold = _holds[holdId_];

        require(currentHold.holder != address(0), "Undefined hold");

        // revert the hold and send the held tokens back to the holder
        _sendHeldTokens(holdId_, currentHold.holder, currentHold.amount);

        emit HoldRemoved(holdId_);
    }

    /// @dev Helper function to send the held tokens
    function _sendHeldTokens(
        uint256 holdId_,
        address to_,
        uint256 amount_
    ) private {
        // delete the hold
        delete _holds[holdId_];

        // transfer the held tokens to the holder
        ERC20Upgradeable(address(this)).transfer(to_, amount_);
    }

    /**
     * @dev Getter for the holds
     * @param holdId_ The id of the Hold
     */
    function getHold(uint256 holdId_)
        public
        view
        returns (
            uint256 holdAmount,
            address holder,
            address operator
        )
    {
        Hold storage currentHold = _holds[holdId_];
        holdAmount = currentHold.amount;
        holder = currentHold.holder;
        operator = currentHold.operator;

        require(holder != address(0), "Undefined hold");
    }

    /**
     * @dev Upgrades needed method
     */
    function _authorizeUpgrade(address) internal view override onlyOwner {}
}
