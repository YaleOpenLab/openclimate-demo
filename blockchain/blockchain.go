package blockchain

import (
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"log"
	"bufio"
	"os"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/YaleOpenLab/openclimate/blockchain/contracts/blockchain_storage"
	"strings"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"math/big"
	"context"
)

const (
	ipfsRootContractAddress = "0x1d0a334994a361111a193b98e6548bf0e8395879"
)

type RootStorage struct {
	abi     abi.ABI
	address common.Address
	opts *bind.TransactOpts
	timeStamp uint32
	rootHash []byte

	root *blockchain.IpfsRoot
}


func NewRoot(address common.Address, client *ethclient.Client, data  interface{}) (*RootStorage, error) {
	parsed, err := abi.JSON(strings.NewReader(blockchain.IpfsRootABI))
	if err != nil {
		return nil, err
	}
	Root, err := blockchain.NewIpfsRoot(address, client)
	if err != nil {
		return nil, err
	}
	return &RootStorage{
		abi:     parsed,
		address: address,
		root: Root,
	}, nil
}


func (root *RootStorage) commitRoot(keystore keystore.KeyStore, passphrase string) *types.Transaction {
	input, err := root.abi.Pack("insertRoot", root.timeStamp, root.rootHash)
	if err != nil {
		return nil
	}
	rawTx := types.NewTransaction(root.opts.Nonce.Uint64(), root.address, root.opts.Value, root.opts.GasLimit, root.opts.GasPrice, input)
	accounts := keystore.Accounts()
	signedTx, err := keystore.SignTxWithPassphrase(accounts[0], passphrase, rawTx, big.NewInt(42))
	return signedTx
}


func CommitToChain(data interface{}) error {
	// Connect to the chain
	//New Ethereum Client
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		log.Fatal(err)
	}

	// Unlock the account from JSON file using passphrase
	wallet := keystore.NewKeyStore("/Users/pavelkrolevets/Documents/YaleOpenClimate/wallet", keystore.StandardScryptN, keystore.StandardScryptP)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter passprase: ")
	passphrase, _ := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	println("unlocked")

	// Send transaction
	contractAddress := common.HexToAddress(ipfsRootContractAddress)
	newRoot, err := NewRoot(contractAddress, client, data)
	client.SendTransaction(context.Background(), newRoot.commitRoot(*wallet, passphrase))

	return nil
}
