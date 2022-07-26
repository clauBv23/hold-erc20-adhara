package outputadapter

import (
	"cleanGo/api/outputinfra"
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"strings"
)

func NewHoldERC20Adapter(networkURL string, networkWS string, contractAddr string) adapter {

	client, err := ethclient.Dial(networkURL)
	if err != nil {
		log.Fatalf("Failed to the Ethereum network: %v", err)
	}

	wsClient, err := ethclient.Dial(networkWS)
	if err != nil {
		log.Fatalf("Failed to the Ethereum network: %v", err)
	}

	cAddress := common.HexToAddress(contractAddr)
	instance, err := outputinfra.NewMain(cAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	return adapter{instance: *instance, client: *client, wsClient: *wsClient, cAddr: cAddress}
}

func (a adapter) CreateHold(holderAddr string, amount int64) (int64, error) {
	privateKey, err := crypto.HexToECDSA(myPk)
	if err != nil {
		return 0, err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := a.client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		//log.Fatal(err)
		return 0, err
	}

	gasPrice, err := a.client.SuggestGasPrice(context.Background())
	if err != nil {
		//log.Fatal(err)
		return 0, err
	}

	auth := bind.NewKeyedTransactor(privateKey)

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	key := [32]byte{}
	value := [32]byte{}
	copy(key[:], []byte("foo"))
	copy(value[:], []byte("bar"))

	holder := common.HexToAddress(holderAddr)
	beneficiary := common.HexToAddress(myAddr)

	tx, err := a.instance.HoldFrom(auth, holder, big.NewInt(amount), beneficiary)
	if err != nil {
		//log.Fatal(err)
		return 0, err
	}
	ctx := context.Background()
	receipt, err := bind.WaitMined(ctx, &a.client, tx)
	if receipt.Status != types.ReceiptStatusSuccessful || err != nil {
		return 0, err
	}

	newHoldId := a.processTransaction(receipt)
	return newHoldId, nil
}

func (a adapter) processTransaction(receipt *types.Receipt) int64 {
	holdCreatedHash := common.HexToHash("0x1f04d8ba13156fb73e621b6df1a4a7aebc25167f7efbd455c45dfc4a3bbea61c")

	query := ethereum.FilterQuery{
		BlockHash: &receipt.BlockHash,
		Addresses: []common.Address{a.cAddr},
		Topics:    [][]common.Hash{{holdCreatedHash}},
	}

	logs, err := a.wsClient.FilterLogs(context.Background(), query)
	if err != nil {
		return 0
	}

	contractAbi, err := abi.JSON(strings.NewReader(string(outputinfra.MainABI)))

	hold, err := contractAbi.Unpack("HoldCreated", logs[0].Data)
	if err != nil {
		return 0

	}

	id := hold[0]
	idBint := id.(*big.Int)

	return idBint.Int64()
}

func (a adapter) TransferTo(toAddr string, amount int64) error {
	privateKey, err := crypto.HexToECDSA(myPk)
	if err != nil {
		return err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := a.client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		//log.Fatal(err)
		return err
	}

	gasPrice, err := a.client.SuggestGasPrice(context.Background())
	if err != nil {
		//log.Fatal(err)
		return err
	}

	auth := bind.NewKeyedTransactor(privateKey)

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	key := [32]byte{}
	value := [32]byte{}
	copy(key[:], []byte("foo"))
	copy(value[:], []byte("bar"))

	to := common.HexToAddress(toAddr)

	tx, err := a.instance.Transfer(auth, to, big.NewInt(amount))
	if err != nil {
		//log.Fatal(err)
		return err
	}
	ctx := context.Background()
	receipt, err := bind.WaitMined(ctx, &a.client, tx)
	if receipt.Status != types.ReceiptStatusSuccessful || err != nil {
		return err
	}

	return nil
}

func (a adapter) ExecuteHold(holdId int64) error {
	privateKey, err := crypto.HexToECDSA(myPk)
	if err != nil {
		return err
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	nonce, err := a.client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		//log.Fatal(err)
		return err
	}

	gasPrice, err := a.client.SuggestGasPrice(context.Background())
	if err != nil {
		//log.Fatal(err)
		return err
	}

	auth := bind.NewKeyedTransactor(privateKey)

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	key := [32]byte{}
	value := [32]byte{}
	copy(key[:], []byte("foo"))
	copy(value[:], []byte("bar"))

	tx, err := a.instance.ExecuteHold(auth, big.NewInt(holdId))
	if err != nil {
		//log.Fatal(err)
		return err
	}
	ctx := context.Background()
	receipt, err := bind.WaitMined(ctx, &a.client, tx)
	if receipt.Status != types.ReceiptStatusSuccessful || err != nil {
		return err
	}

	return nil
}
