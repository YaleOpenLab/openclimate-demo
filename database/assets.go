package database

import (
	"encoding/json"
	edb "github.com/Varunram/essentials/database"
	globals "github.com/YaleOpenLab/openclimate/globals"
	"github.com/pkg/errors"
	"log"
)

type Asset struct {
	Index   	int
	Name    	string
	Company 	string
	Location 	string

	Type 		string
	Reports   []RepData
}

// Puts asset object in assets bucket. Called by NewAsset
func (a *Asset) Save() error {
	return edb.Save(globals.DbPath, AssetBucket, a, a.Index)
}

func NewAsset(name string, company string) (Asset, error) {
	var asset Asset

	assets, err := RetrieveAllAssets()
	if err != nil {
		return asset, errors.Wrap(err, "could not retrieve all assets, quitting")
	}

	if len(assets) == 0 {
		asset.Index = 1
	} else {
		asset.Index = len(assets) + 1
	}

	asset.Name = name
	asset.Company = company
	return asset, asset.Save()
}

func UpdateAsset(key int, info Asset) error {
	asset, err := RetrieveAsset(key)
	if err != nil {
		return errors.Wrap(err, "UpdateAsset() failed (likely because asset doesn't exist)")
	}
	asset.Name = info.Name
	asset.Location = info.Location
	asset.Type = info.Type
	asset.Save()

	return nil
}

// Given a key of type int, retrieves the corresponding asset object
// from the database assets bucket.
func RetrieveAsset(key int) (Asset, error) {
	var asset Asset
	assetBytes, err := edb.Retrieve(globals.DbPath, AssetBucket, key)
	if err != nil {
		return asset, errors.Wrap(err, "error while retrieving key from bucket")
	}
	err = json.Unmarshal(assetBytes, &asset)
	return asset, err
}

// Given a name and company, retrieves the corresponding asset object
// from the database assets bucket.
func RetrieveAssetByName(name string, company string) (Asset, error) {
	var asset Asset
	allAssets, err := RetrieveAllAssets()
	if err != nil {
		return asset, errors.Wrap(err, "error while retrieving all users from database")
	}

	for _, asset := range allAssets {
		if asset.Name == name && asset.Company == company {
			return asset, nil
		}
	}

	return asset, errors.New("asset not found, quitting")
}

// RetrieveAllAssets gets a list of all assets in the database
func RetrieveAllAssets() ([]Asset, error) {
	var assets []Asset
	keys, err := edb.RetrieveAllKeys(globals.DbPath, AssetBucket)
	if err != nil {
		log.Println(err)
		return assets, errors.Wrap(err, "could not retrieve all user keys")
	}
	for _, val := range keys {
		var x Asset
		err = json.Unmarshal(val, &x)
		if err != nil {
			break
		}
		assets = append(assets, x)
	}

	return assets, nil
}
