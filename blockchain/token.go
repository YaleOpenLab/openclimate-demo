package blockchain

import (
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"github.com/ethereum/go-ethereum/common"
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/YaleOpenLab/openclimate/blockchain/contracts/token"
	"strings"
)

//  YOCL (GRC20) token is deployed on Kovan testnet. Initial supply is 1000000000.
// https://kovan.etherscan.io/tx/0x998aadf4602f79b8f8e4e84def8b84191237d3b6abe73e715546303c072259b5
// To interact with the token we are using Web3go library. Connection through infura endpoint.

const (
	contractAddress = "0xb19ac159b87a4491b3b8bef4554b59da2bf42555"
	ownerAddress    = "0xfE1827f2F1C366c04d458b3c07B8Bd207D42eab4"
	rpcUrl          = "https://kovan.infura.io/v3/def7370cf49d49d791b9df949986b9a0"
)

type Token struct {
	client    *ethclient.Client
	abi       abi.ABI
	address   common.Address

	token *blockchain.YToken
}

func NewToken(address common.Address, client *ethclient.Client) (*Token, error) {
	parsed, err := abi.JSON(strings.NewReader(blockchain.YTokenABI))
	if err != nil {
		return nil, err
	}
	token, err := blockchain.NewYToken(address, client)
	if err != nil {
		return nil, err
	}
	return &Token{
		client:    client,
		abi:       parsed,
		address:   address,

		token: token,
	}, nil
}

func CheckTokenBalance() {
	//New Ethereum Client
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		log.Fatal(err)
	}

	address := common.HexToAddress(ownerAddress)
	if err != nil {
		log.Fatal(err)
	}
	// 3. get ether balance from the latest block
	balance, err := client.BalanceAt(context.Background(), address, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(balance)

	//Access the Token
	contractAddress := common.HexToAddress(contractAddress)
	YToken, err := NewToken(contractAddress, client)
	if err != nil {
		log.Fatal(err)
	}

	/*
		//Check if the address is a Token owner
		isMinter, _ := YToken.IsMinter(address)
		if isMinter {
			fmt.Println("True")
		}else{
			fmt.Println("False")
		}
	*/

	tokenBalance, err := YToken.token.BalanceOf(nil, address)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("YOCL Token Balance", tokenBalance.String())
}
