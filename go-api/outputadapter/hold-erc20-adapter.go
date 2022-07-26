package outputadapter

import (
	"cleanGo/api/outputinfra"
	"context"
	"crypto/ecdsa"
	"fmt"
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
	fmt.Println("hererere")
	privateKey, err := crypto.HexToECDSA("431fb01d4f91d7afb25d2214e08e19edfda753e613018975ea8c017dc26a13a9")
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

	fmt.Println("before")
	tx, err := a.instance.HoldFrom(auth, holder, big.NewInt(amount), beneficiary)
	if err != nil {
		//log.Fatal(err)
		return 0, err
	}
	fmt.Println("waiting")
	ctx := context.Background()
	receipt, err := bind.WaitMined(ctx, &a.client, tx)
	fmt.Println("done")
	if receipt.Status != types.ReceiptStatusSuccessful || err != nil {
		return 0, err
	}

	fmt.Println("calling method")
	newHoldId := a.processTransaction(receipt)
	return newHoldId, nil
}

func (a adapter) processTransaction(receipt *types.Receipt) int64 {
	fmt.Println("on method")
	//fmt.Println(receipt)

	holdCreatedHash := common.HexToHash("0x1f04d8ba13156fb73e621b6df1a4a7aebc25167f7efbd455c45dfc4a3bbea61c")

	query := ethereum.FilterQuery{
		BlockHash: &receipt.BlockHash,
		Addresses: []common.Address{a.cAddr},
		Topics:    [][]common.Hash{{holdCreatedHash}},
	}
	//fmt.Println("query")
	//fmt.Println(query)

	logs, err := a.wsClient.FilterLogs(context.Background(), query)
	if err != nil {
		return 0
	}

	contractAbi, err := abi.JSON(strings.NewReader(string(outputinfra.MainABI)))

	hold, err := contractAbi.Unpack("HoldCreated", logs[0].Data)
	if err != nil {
		return 0

	}
	fmt.Println(hold)

	id := hold[0]

	fmt.Println(id)
	fmt.Printf("%T\\n", id)

	//fmt.Println(logs[0].Data)
	//if err != nil {
	//	log.Fatal(err)
	//}
	idBint := id.(*big.Int)

	return idBint.Int64()
}
