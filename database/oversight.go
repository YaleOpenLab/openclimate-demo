package database

import (
	"encoding/json"
	// xlm "github.com/Varunram/essentials/crypto/xlm"
	edb "github.com/Varunram/essentials/database"
	globals "github.com/YaleOpenLab/openclimate/globals"
	"github.com/pkg/errors"
	"log"
)

type Oversight struct {
	Index int
	Name  string

	UserIDs []int

	// options:
	// - "international organization"
	// - "UNFCCC Body"
	// - "Civil Society"
	// - "Financial Institution"
	// - "Research Organization"
	// - "Independent Reviewer"
	// - "Oracle System"
	OrgType string

	Scope  string // where does the actor operate?
	Weight int    // their weight, based on the organization's reputation

}

func NewOsOrg(name string) (Oversight, error) {
	var osOrg Oversight
	osOrg.Name = name
	return osOrg, osOrg.Save()
}

func RetrieveOsOrg(key int) (Oversight, error) {
	var osOrg Oversight
	bytes, err := edb.Retrieve(globals.DbPath, OversightBucket, key)
	if err != nil {
		return osOrg, errors.Wrap(err, "error while retrieving key from bucket")
	}
	err = json.Unmarshal(bytes, &osOrg)
	return osOrg, err
}

func RetrieveOsOrgByName(name string) (Oversight, error) {
	var osOrg Oversight
	temp, err := RetrieveAllOsOrgs()
	if err != nil {
		return osOrg, errors.Wrap(err, "error while retrieving all users from database")
	}

	for _, osOrg := range temp {
		if osOrg.Name == name {
			return osOrg, nil
		}
	}

	return osOrg, errors.New("osOrg not found, quitting")
}

func RetrieveAllOsOrgs() ([]Oversight, error) {
	var osOrgs []Oversight
	keys, err := edb.RetrieveAllKeys(globals.DbPath, OversightBucket)
	if err != nil {
		log.Println(err)
		return osOrgs, errors.Wrap(err, "could not retrieve all keys")
	}
	for _, val := range keys {
		var x Oversight
		err = json.Unmarshal(val, &x)
		if err != nil {
			break
		}
		osOrgs = append(osOrgs, x)
	}

	return osOrgs, nil
}

func (os *Oversight) AddPledges(pledgeIDs ...int) error {
	return nil
}

func (os Oversight) GetPledges() ([]Pledge, error) {
	var empty []Pledge
	return empty, nil
}

// func CommitToStellar(ipfsHash string, seed string, pubkey string) (string, string, error) {

// 	memo := "IPFSHASH: " + ipfsHash // add padding to the ipfs hash length

// 	firstHalf := memo[:28]
// 	secondHalf := memo[28:]

// 	_, tx1, err := xlm.SendXLM(pubkey, "1", seed, firstHalf)
// 	if err != nil {
// 		return "", "", err
// 	}
// 	_, tx2, err := xlm.SendXLM(pubkey, "1", seed, secondHalf)
// 	if err != nil {
// 		return "", "", err
// 	}

// 	log.Printf("tx hash: %s, tx2 hash: %s", tx1, tx2)
// 	return tx1, tx2, nil
// }
