package blockchain

import (
	"bufio"
	"context"
	"fmt"
	"github.com/YaleOpenLab/openclimate/blockchain/contracts/blockchain_storage"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"os"
	"strings"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"runtime"
)

const (
	ipfsRootContractAddress = "0x1d0a334994a361111a193b98e6548bf0e8395879"
)

type RootStorage struct {
	client    *ethclient.Client
	abi       abi.ABI
	address   common.Address
	opts      *bind.TransactOpts
	timeStamp *big.Int
	rootHash  [32]byte

	root *blockchain.IpfsRoot
}

func NewRoot(address common.Address, client *ethclient.Client, timeStamp *big.Int, rootHash [32]byte) (*RootStorage, error) {
	parsed, err := abi.JSON(strings.NewReader(blockchain.IpfsRootABI))
	if err != nil {
		return nil, err
	}
	Root, err := blockchain.NewIpfsRoot(address, client)
	if err != nil {
		return nil, err
	}
	return &RootStorage{
		client: client,
		abi:     parsed,
		address: address,
		timeStamp: timeStamp,
		rootHash: rootHash,

		root:    Root,
	}, nil
}

func (root *RootStorage) commitRoot(keystore keystore.KeyStore, passphrase string) (*types.Transaction, error) {
	input, err := root.abi.Pack("insertRoot", root.timeStamp, root.rootHash)
	if err != nil {
		return nil, err
	}
	nonce, err := root.client.PendingNonceAt(context.Background(), keystore.Accounts()[0].Address)
	if err != nil {
		return nil, err
	}
	fmt.Println("Nonce", nonce)
	gasPrice, err := root.client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, err
	}
	fmt.Println("gasPrice")

	rawTx := types.NewTransaction(nonce, root.address, big.NewInt(0), 300000, gasPrice, input)
	fmt.Println("RawTx", rawTx)

	accounts := keystore.Accounts()
	signedTx, err := keystore.SignTxWithPassphrase(accounts[0], passphrase, rawTx, big.NewInt(42))
	if err != nil {
		return nil, err
	}
	return signedTx, err
}

func CommitToChain(timeStamp *big.Int, rootHash string) error {
	// Connect to the chain
	//New Ethereum Client
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		log.Fatal(err)
	}

	// Unlock the account from JSON file using passphrase
	wallet := keystore.NewKeyStore("./blockchain/wallet/", keystore.StandardScryptN, keystore.StandardScryptP)
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter passprase: ")
	passphrase, _ := reader.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}
	rootHashBytes, err := hexutil.Decode(rootHash)
	if err != nil {
		log.Fatal(err)
	}
	var rootHashBytes32 [32]byte
	copy(rootHashBytes32[:], rootHashBytes)
	// Send transaction
	contractAddress := common.HexToAddress(ipfsRootContractAddress)
	newRoot, err := NewRoot(contractAddress, client, timeStamp, rootHashBytes32)
	if err != nil {
		_, fn, line, _ := runtime.Caller(1)
		log.Printf("[error] %s:%d %v", fn, line, err)
	}

	newTx, err:= newRoot.commitRoot(*wallet, strings.TrimSpace(passphrase))
	if err != nil {
		_, fn, line, _ := runtime.Caller(1)
		log.Printf("[error] %s:%d %v", fn, line, err)
	}
	client.SendTransaction(context.Background(), newTx)
	fmt.Println("Successfuly commited new ipfs root.")
	return nil
}
