package database

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"github.com/pkg/errors"
	"log"
	"math/big"

	// keys "github.com/cosmos/cosmos-sdk/crypto/keys"
	aes "github.com/Varunram/essentials/aes"
	edb "github.com/Varunram/essentials/database"
	"github.com/YaleOpenLab/openclimate/globals"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	crypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type User struct {
	Index int

	FirstName string
	LastName  string
	Username  string
	Email     string
	Pwhash    string
	EIN       string

	EntityType string // choices are: individual, company, city, state, region, country, oversight
	EntityID   int    // index of the entity the user is associated with
	Verified   bool   // if the user is a verified member of the entity they purport to be a part of
	Admin      bool   // is the user an admin for its entity?

	EthereumWallet EthWallet
	//CosmosWallet   CosmWallet

}

// EthWallet contains the structures needed for an ethereum wallet
type EthWallet struct {
	EncryptedPrivateKey string
	PublicKey           string
	Address             string
}

/*
type CosmWallet struct {
	PrivateKey string
	PublicKey  string
}
*/

// NewUser creates a new user
func NewUser(username string, pwhash string, email string, entityType string, entityName string, entityParent string) (User, error) {
	var user User
	var err error

	if len(pwhash) != 128 {
		return user, errors.New("pwhash not of length 128, quitting")
	}

	user.Username = username
	user.Pwhash = pwhash
	user.Email = email

	if entityType == "" {
		return user, errors.New("Entity type not specified, quitting")
	}

	user.EntityType = entityType

	switch entityType {
	case "company":
		var entity Company
		entity, err = RetrieveCompanyByName(entityName, entityParent)
		user.EntityID = entity.Index
	case "city":
		var entity City
		entity, err = RetrieveCityByName(entityName, entityParent)
		user.EntityID = entity.Index
	case "state":
		var entity State
		entity, err = RetrieveStateByName(entityName, entityParent)
		user.EntityID = entity.Index
	case "region":
		var entity Region
		entity, err = RetrieveRegionByName(entityName, entityParent)
		user.EntityID = entity.Index
	case "country":
		var entity Country
		entity, err = RetrieveCountryByName(entityName)
		user.EntityID = entity.Index
	case "oversight":
		var entity Oversight
		entity, err = RetrieveOsOrgByName(entityName)
		user.EntityID = entity.Index
	default:
		err = errors.New("Actor type not found.")
	}

	if err != nil {
		return user, errors.Wrap(err, "NewUser() failed.")
	}

	return user, user.Save()
}

// // RetrieveAllUsers gets a list of all User in the database
// func RetrieveAllUsers() ([]User, error) {
// 	var users []User
// 	keys, err := edb.RetrieveAllKeys(globals.DbPath, UserBucket)
// 	if err != nil {
// 		log.Println(err)
// 		return users, errors.Wrap(err, "could not retrieve all user keys")
// 	}
// 	for _, val := range keys {
// 		var x User
// 		err = json.Unmarshal(val, &x)
// 		if err != nil {
// 			break
// 		}
// 		users = append(users, x)
// 	}

// 	return users, nil
// }

// RetrieveUser retrieves a particular User indexed by key from the database
func RetrieveUser(key int) (User, error) {
	var user User
	userBytes, err := edb.Retrieve(globals.DbPath, UserBucket, key)
	if err != nil {
		return user, errors.Wrap(err, "error while retrieving key from bucket")
	}
	err = json.Unmarshal(userBytes, &user)
	return user, err
}

func RetrieveUserByUsername(username string) (User, error) {
	var user User
	allUsers, err := RetrieveAllUsers()
	if err != nil {
		return user, errors.Wrap(err, "Could not retrieve all users from db")
	}

	for _, val := range allUsers {
		if val.Username == username {
			user = val
			return user, nil
		}
	}
	return user, errors.New("User not found")
}

// ValidateUser validates a particular user
func ValidateUser(username string, pwhash string) (User, error) {
	var user User
	users, err := RetrieveAllUsers()
	if err != nil {
		return user, errors.Wrap(err, "error while retrieving all users from database")
	}

	for _, user := range users {
		if user.Username == username && user.Pwhash == pwhash {
			return user, nil
		}
	}

	return user, errors.New("user not found")
}

// Empty function, simply allows User to match "Actor" interface methods
func (u *User) AddPledge(pledge Pledge) {
	return
}

func (user *User) GetUserActor() (Actor, error) {

	var entity Actor
	var err error

	switch user.EntityType {
	case "company":
		var x Company
		x, err = RetrieveCompany(user.EntityID)
		entity = &x
	case "city":
		var x City
		x, err = RetrieveCity(user.EntityID)
		entity = &x
	case "state":
		var x State
		x, err = RetrieveState(user.EntityID)
		entity = &x
	case "region":
		var x Region
		x, err = RetrieveRegion(user.EntityID)
		entity = &x
	case "country":
		var x Country
		x, err = RetrieveCountry(user.EntityID)
		entity = &x
	case "oversight":
		var x Oversight
		x, err = RetrieveOsOrg(user.EntityID)
		entity = &x
	default:
		return entity, errors.New("User's entity type is not valid.")
	}

	if err != nil {
		return entity, errors.Wrap(err, "User's linked actor was not found")
	}

	return entity, nil
}

func (a *User) SendEthereumTx(address string, amount big.Int) (string, error) {
	client, err := ethclient.Dial("https://ropsten.infura.io")
	if err != nil {
		return "", errors.Wrap(err, "could not contact infura")
	}

	privateKey, err := crypto.HexToECDSA(a.EthereumWallet.EncryptedPrivateKey)
	if err != nil {
		return "", errors.Wrap(err, "could not generate private key")
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return "", errors.Wrap(err, "could not derive publickey from private key")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return "", errors.Wrap(err, "could not derive nonce, quitting")
	}

	gasLimit := uint64(21000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return "", errors.Wrap(err, "could not get gas price from infura, quitting")
	}

	toAddress := common.HexToAddress(address)
	var data []byte
	tx := types.NewTransaction(nonce, toAddress, &amount, gasLimit, gasPrice, data)
	signedTx, err := types.SignTx(tx, types.HomesteadSigner{}, privateKey)
	if err != nil {
		return "", errors.Wrap(err, "could not sing transaction, quitting")
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return "", errors.Wrap(err, "could not send transaction to infura, quitting")
	}

	log.Printf("tx sent: %s", signedTx.Hash().Hex())
	return signedTx.Hash().Hex(), nil
}

func (a *User) GenEthKeys(seedpwd string) error {
	ecdsaPrivkey, err := crypto.GenerateKey()
	if err != nil {
		return errors.Wrap(err, "could not generate an ethereum keypair, quitting!")
	}

	privateKeyBytes := crypto.FromECDSA(ecdsaPrivkey)

	ek, err := aes.Encrypt([]byte(hexutil.Encode(privateKeyBytes)[2:]), seedpwd)
	if err != nil {
		return errors.Wrap(err, "error while encrypting seed")
	}

	a.EthereumWallet.EncryptedPrivateKey = string(ek)
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

	return a.Save()
}

/*
func (a *User) GenCosmosKeys() error {
	// Select the encryption and storage for your cryptostore
	cstore := keys.NewInMemory()

	sec := keys.Secp256k1

	// Add keys and see they return in alphabetical order
	bob, _, err := cstore.CreateMnemonic("Bob", keys.English, "friend", sec)
	if err != nil {
		// this should never happen
		log.Println(err)
	} else {
		// return info here just like in List
		log.Println(bob.GetName())
	}
	_, _, _ = cstore.CreateMnemonic("Alice", keys.English, "secret", sec)
	_, _, _ = cstore.CreateMnemonic("Carl", keys.English, "mitm", sec)
	info, _ := cstore.List()
	for _, i := range info {
		log.Println(i.GetName())
	}

	// We need to use passphrase to generate a signature
	tx := []byte("deadbeef")
	sig, pub, err := cstore.Sign("Bob", "friend", tx)
	if err != nil {
		log.Println("don't accept real passphrase")
	}

	// and we can validate the signature with publicly available info
	binfo, _ := cstore.Get("Bob")
	if !binfo.GetPubKey().Equals(bob.GetPubKey()) {
		log.Println("Get and Create return different keys")
	}

	if pub.Equals(binfo.GetPubKey()) {
		log.Println("signed by Bob")
	}
	if !pub.VerifyBytes(tx, sig) {
		log.Println("invalid signature")
	}
}
*/
