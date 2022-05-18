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
     * @notice Emmited when a new hold is created
     * @param holdId The created hold Id
     * @param holdAmount The mount holded
     * @param holder The holder's address
     * @param operator The creator of the hold
     */
    event HoldCreated(uint256 indexed holdId, uint256 holdAmount, address holder, address operator);

    /**
     * @notice Emmited when a new hold is executed
     * @param holdId The created hold Id
     */
    event HoldExecuted(uint256 indexed holdId);

    /**
     * @notice Emmited when a new hold is removed by the contract owner
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
    function mint(address to_, uint256 amount_) public onlyOwner {
        _mint(to_, amount_);
    }

    /**
     * @dev Hold the specified amount
     * @param amount_ Amount of tokens to hold
     */
    function hold(uint256 amount_) public {
        _hold(msg.sender, msg.sender, amount_);
    }

    /**
     * @dev Hold the specified amount from another account
     * @param holder_ Owner of tokens to hold
     * @param amount_ Amount of tokens to hold
     */
    function holdFrom(address holder_, uint256 amount_) public {
        uint256 allowance_ = allowance(holder_, msg.sender);

        // check if the holder allow the caller to hold
        require(allowance_ >= amount_, "Not allowed");

        _hold(msg.sender, holder_, amount_);
    }

    /// @dev Helper function to create the hold
    function _hold(
        address operator_,
        address holder_,
        uint256 amount_
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
        _holds[holdId] = Hold({operator: operator_, holder: holder_, amount: amount_});

        emit HoldCreated(holdId, amount_, holder_, operator_);
    }

    /**
     * @dev Returns the holded tokens to the holder
     * @param holdId_ Id of the Hold
     */
    function executeHold(uint256 holdId_) public {
        Hold memory currentHold = _holds[holdId_];

        require(currentHold.operator == msg.sender, "The caller must be the hold creator");

        _revertHold(holdId_, currentHold.holder, currentHold.amount);

        emit HoldExecuted(holdId_);
    }

    /**
     * @notice Only owner allowed
     * @dev Returns the holded tokens to the holder
     * @param holdId_ Id of the Hold
     */
    function removeHold(uint256 holdId_) public onlyOwner {
        Hold memory currentHold = _holds[holdId_];

        require(currentHold.holder != address(0), "Undefined hold");

        _revertHold(holdId_, currentHold.holder, currentHold.amount);

        emit HoldRemoved(holdId_);
    }

    /// @dev Helper function to send back to the holder the holded tokens
    function _revertHold(
        uint256 holdId_,
        address holder_,
        uint256 amount_
    ) private {
        // delete the hold
        delete _holds[holdId_];

        // transfer the holded tokens to the holder
        ERC20Upgradeable(address(this)).transfer(holder_, amount_);
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
        holdAmount = _holds[holdId_].amount;
        holder = _holds[holdId_].holder;
        operator = _holds[holdId_].operator;

        require(holder != address(0), "Undefined hold");
    }

    /**
     * @dev Upgrades needed method
     */
    function _authorizeUpgrade(address) internal view override onlyOwner {}
}
