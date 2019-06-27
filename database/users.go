package database

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"math/big"

	aes "github.com/YaleOpenLab/openx/aes"
	utils "github.com/YaleOpenLab/openx/utils"
	"github.com/boltdb/bolt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	crypto "github.com/ethereum/go-ethereum/crypto"
)

type User struct {
	Index          int
	Name           string
	Email          string
	Pwhash         string
	EthereumWallet EthWallet
}

// EthWallet contains the structures needed for an ethereum wallet
type EthWallet struct {
	PrivateKey string
	PublicKey  string
	Address    string
}

func (a *User) GenKeys(seedpwd string) error {
	ecdsaPrivkey, err := crypto.GenerateKey()
	if err != nil {
		return errors.Wrap(err, "could not generate an ethereum keypair, quitting!")
	}

	privateKeyBytes := crypto.FromECDSA(ecdsaPrivkey)

	ek, err := aes.Encrypt([]byte(hexutil.Encode(privateKeyBytes)[2:]), seedpwd)
	if err != nil {
		return errors.Wrap(err, "error while encrypting seed")
	}

	a.EthereumWallet.PrivateKey = string(ek)
	a.EthereumWallet.Address = crypto.PubkeyToAddress(ecdsaPrivkey.PublicKey).Hex()

	publicKeyECDSA, ok := ecdsaPrivkey.Public().(*ecdsa.PublicKey)
	if !ok {
		return errors.Wrap(err, "error casting public key to ECDSA")
	}

	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	a.EthereumWallet.PublicKey = hexutil.Encode(publicKeyBytes)[4:] // an ethereum address is 65 bytes long and hte first byte is 0x04 for DER encoding, so we omit that

	if crypto.PubkeyToAddress(*publicKeyECDSA).Hex() != a.EthereumWallet.Address {
		return errors.Wrap(err, "addresses don't match, quitting!")
	}

	err = a.Save()
	return err
}

// NewUser creates a new user
func NewUser(name string, pwhash string, email string) (User, error) {
	var user User

	if len(pwhash) != 128 {
		return user, errors.New("pwhash not of length 128, quitting")
	}

	allUsers, err := RetrieveAllUsers()
	if err != nil {
		return user, errors.Wrap(err, "Error while retrieving all users from database")
	}

	// the ugly indexing thing again, need to think of something better here
	if len(allUsers) == 0 {
		user.Index = 1
	} else {
		user.Index = len(allUsers) + 1
	}

	user.Name = name
	user.Pwhash = pwhash
	user.Email = email

	return user, user.Save()
}

// Save inserts a passed User object into the database
func (a *User) Save() error {
	db, err := OpenDB()
	if err != nil {
		return errors.Wrap(err, "Error while opening database")
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(UserBucket)
		encoded, err := json.Marshal(a)
		if err != nil {
			return errors.Wrap(err, "Error while marshaling json")
		}
		return b.Put([]byte(utils.ItoB(a.Index)), encoded)
	})
	return err
}

// RetrieveAllUsers gets a list of all User in the database
func RetrieveAllUsers() ([]User, error) {
	var arr []User
	db, err := OpenDB()
	if err != nil {
		return arr, errors.Wrap(err, "Error while opening database")
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(UserBucket)
		for i := 1; ; i++ {
			var rUser User
			x := b.Get(utils.ItoB(i))
			if x == nil {
				return nil
			}
			err := json.Unmarshal(x, &rUser)
			if err != nil {
				return errors.Wrap(err, "Error while unmarshalling json")
			}
			arr = append(arr, rUser)
		}
	})
	return arr, err
}

// RetrieveUser retrieves a particular User indexed by key from the database
func RetrieveUser(key int) (User, error) {
	var inv User
	db, err := OpenDB()
	if err != nil {
		return inv, errors.Wrap(err, "error while opening database")
	}
	defer db.Close()
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(UserBucket)
		x := b.Get(utils.ItoB(key))
		if x == nil {
			return errors.New("retrieved user nil, quitting!")
		}
		return json.Unmarshal(x, &inv)
	})
	return inv, err
}

// ValidateUser validates a particular user
func ValidateUser(name string, pwhash string) (User, error) {
	var inv User
	temp, err := RetrieveAllUsers()
	if err != nil {
		return inv, errors.Wrap(err, "error while retrieving all users from database")
	}
	limit := len(temp) + 1
	db, err := OpenDB()
	if err != nil {
		return inv, errors.Wrap(err, "could not open db, quitting!")
	}
	defer db.Close()
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(UserBucket)
		for i := 1; i < limit; i++ {
			var rUser User
			x := b.Get(utils.ItoB(i))
			err := json.Unmarshal(x, &rUser)
			if err != nil {
				return errors.Wrap(err, "could not unmarshal json, quitting!")
			}
			// check names
			if rUser.Name == name && rUser.Pwhash == pwhash {
				inv = rUser
				return nil
			}
		}
		return errors.New("Not Found")
	})
	return inv, err
}

func (a *User) SendEthereumTx(address string, amount big.Int) error {
	chainId := big.NewInt(3) // Ropsten chain id
	senderPrivKey, err := crypto.HexToECDSA(a.EthereumWallet.PrivateKey)
	if err != nil {
		return errors.Wrap(err, "could not convert private key from hex to ecdsa")
	}
	recipientAddr := common.HexToAddress(address)

	nonce := uint64(7)
	gasLimit := uint64(100000)         // hardcode gas, max 100k exec limit
	gasPrice := big.NewInt(1000000000) // hardcode gas, 1 gwei price

	tx := types.NewTransaction(nonce, recipientAddr, &amount, gasLimit, gasPrice, nil)

	signer := types.NewEIP155Signer(chainId)
	signedTx, err := types.SignTx(tx, signer, senderPrivKey)
	if err != nil {
		return errors.Wrap(err, "could not sign transaction, quitting")
	}

	var buff bytes.Buffer
	signedTx.EncodeRLP(&buff)
	fmt.Printf("0x%x\n", buff.Bytes())

	// TODO: send this to infura or something once we decide to fit in a blockchain
	return nil
}
