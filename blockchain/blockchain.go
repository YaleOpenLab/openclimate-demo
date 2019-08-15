package blockchain

import (
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"log"
	"bufio"
	"os"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereumclassic/go-ethereum/common"
)

type RootStorageContract struct {
	abi     abi.ABI
	address common.Address
}


func CommitToChain(data interface{}) error {
	keystore := keystore.NewKeyStore("/Users/pavelkrolevets/Documents/YaleOpenClimate/wallet", keystore.StandardScryptN, keystore.StandardScryptP)
	accounts := keystore.Accounts()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter passprase: ")
	passphrase, _ := reader.ReadString('\n')
	err := keystore.Unlock(accounts[0], passphrase)
	if err != nil {
		log.Fatal(err)
	}
	println("unlocked")
	return nil
}