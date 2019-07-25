package database

import (
	"encoding/json"
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

	osOrgs, err := RetrieveAllOsOrgs()
	if err != nil {
		return osOrg, errors.Wrap(err, "could not retrieve all oversight organizations, quitting")
	}

	if len(osOrgs) == 0 {
		osOrg.Index = 1
	} else {
		osOrg.Index = len(osOrgs) + 1
	}

	osOrg.Name = name
	return osOrg, osOrg.Save()
}

func (a *Oversight) Save() error {
	return edb.Save(globals.DbPath, OversightBucket, a, a.Index)
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
