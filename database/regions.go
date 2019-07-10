package database

import (
	"encoding/json"
	edb "github.com/Varunram/essentials/database"
	utils "github.com/Varunram/essentials/utils"
	globals "github.com/YaleOpenLab/openclimate/globals"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

// includes states, regions, provinces, prefectures, etc.
type Region struct {
	Index       int
	Name        string
	Country     string
	Area        float64
	Iso         string
	Population  int
	Latitude    float64
	Longitude   float64
	Revenue     float64
	CompanySize int
	HQ          string
	// EntityType		string
}

func NewRegion(name string, country string) (Region, error) {

	var new Region
	var err error

	// naive implementation of assigning keys to bucket items (simple indexing)
	regions, err := RetrieveAllRegions()
	lenRegions := len(regions)

	if err != nil {
		return new, errors.Wrap(err, "Error while retrieving all regions from db")
	}

	if lenRegions == 0 {
		new.Index = 1
	} else {
		new.Index = lenRegions + 1
	}

	new.Name = name
	new.Country = country

	// simply initializing these fields to nil for now
	new.Area = 0
	new.Iso = ""
	new.Population = 0
	new.Latitude = 0
	new.Longitude = 0
	new.Revenue = 0
	new.CompanySize = 0
	new.HQ = ""

	err = new.Save()
	return new, err

}

/*
	ISSUE: edb.Save() asks for an key argument of type INT,
	but currently we are passing in a key argument of type string.
	This issue needs to be resolved. Could maybe just use a hash.

	RESOLVED: currently using solution previously implemented in OpenX;
	incrementing index for each new region, so the key is of type int.
*/

func (region *Region) Save() error {
	return edb.Save(globals.DbDir, RegionBucket, region, region.Index)
}

func RetrieveRegion(key int) (Region, error) {
	var region Region
	temp, err := edb.Retrieve(globals.DbDir, RegionBucket, key)

	if err != nil {
		return region, errors.Wrap(err, "Error while retrieving key from bucket")
	}

	region = temp.(Region)
	return region, region.Save()
}

func RetrieveRegionByName(name string, country string) (Region, error) {
	var region Region
	allRegions, err := RetrieveAllRegions()
	if err != nil {
		return region, errors.Wrap(err, "Error while retrieving all regions from database")
	}

	db, err := OpenDB()
	if err != nil {
		return region, errors.Wrap(err, "Could not open database, quitting")
	}

	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(RegionBucket)

		limit := len(allRegions) + 1
		for i := 1; i < limit; i++ {
			var tempRegion Region
			tempKey := bucket.Get(utils.ItoB(i))

			err := json.Unmarshal(tempKey, &tempRegion)
			if err != nil {
				return errors.Wrap(err, "Could not unmarshal json, quitting")
			}

			if tempRegion.Name == name && tempRegion.Country == country {
				region = tempRegion
				return nil
			}
		}
		return errors.New("Region not found.")
	})
	return region, err
}

func RetrieveAllRegions() ([]Region, error) {
	var regions []Region
	keys, err := edb.RetrieveAllKeys(globals.DbDir, RegionBucket)
	if err != nil {
		return regions, errors.Wrap(err, "error while retrieving all keys")
	}

	for _, value := range keys {
		regions = append(regions, value.(Region))
	}

	return regions, nil
}
