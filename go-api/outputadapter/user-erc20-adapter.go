package outputadapter

import (
	"cleanGo/api/outputinfra"
	"context"
	"crypto/ecdsa"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

type adapter struct {
	instance outputinfra.Main
	client   ethclient.Client
	wsClient ethclient.Client
	cAddr    common.Address
}

const myPk = "431fb01d4f91d7afb25d2214e08e19edfda753e613018975ea8c017dc26a13a9"
const myAddr = "0x01D1dBd8D5796881A59a3822F6def9e5FF77B9e4"

func NewUserERC20Adapter(networkURL string, contractAddr string) adapter {

	client, err := ethclient.Dial(networkURL)
	if err != nil {
		log.Fatalf("Failed to the Ethereum network: %v", err)
	}

	cAddress := common.HexToAddress(contractAddr)

	instance, err := outputinfra.NewMain(cAddress, client)
	if err != nil {
		log.Fatal(err)
	}
	return adapter{instance: *instance, client: *client}
}

func (a adapter) GetUserBalance(address string) (*int, error) {

	response, err := a.instance.BalanceOf(nil, common.HexToAddress(address))
	if err != nil {
		return nil, err
	}
	b := int(response.Int64())

	return &b, nil
}

func (a adapter) MintTokensToUser(userAddr string, amount int64) error {
	// mint tokens to the registered user

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

	toAddr := common.HexToAddress(userAddr)
	_, err = a.instance.Mint(auth, toAddr, big.NewInt(amount))

	if err != nil {
		//log.Fatal(err)
		return err
	}
	return nil
}
