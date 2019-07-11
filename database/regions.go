package database

import (
	"encoding/json"
	edb "github.com/Varunram/essentials/database"
	globals "github.com/YaleOpenLab/openclimate/globals"
	"github.com/pkg/errors"
	"log"
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

/*
	ISSUE: edb.Save() asks for an key argument of type INT,
	but currently we are passing in a key argument of type string.
	This issue needs to be resolved. Could maybe just use a hash.

	RESOLVED: currently using solution previously implemented in OpenX;
	incrementing index for each new region, so the key is of type int.
*/

func NewRegion(name string, country string) (Region, error) {

	var new Region
	var err error

	// naive implementation of assigning keys to bucket items (simple indexing)
	regions, err := RetrieveAllRegions()
	if err != nil {
		log.Println(err)
		return new, errors.Wrap(err, "could not retreive all regions")
	}
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

	// // simply initializing these fields to nil for now
	// new.Area = 0
	// new.Iso = ""
	// new.Population = 0
	// new.Latitude = 0
	// new.Longitude = 0
	// new.Revenue = 0
	// new.CompanySize = 0
	// new.HQ = ""

	err = new.Save()
	return new, err

}

func (region *Region) Save() error {
	return edb.Save(globals.DbDir+"/openclimate.db", RegionBucket, region, region.Index)
}

func RetrieveRegion(key int) (Region, error) {
	var region Region
	temp, err := edb.Retrieve(globals.DbDir+"/openclimate.db", RegionBucket, key)

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

	for _, val := range allRegions {
		if val.Name == name && val.Country == country {
			region = val
			return region, nil
		}
	}

	return region, errors.New("could not find regions")
}

func RetrieveAllRegions() ([]Region, error) {
	var regions []Region
	keys, err := edb.RetrieveAllKeys(globals.DbDir+"/openclimate.db", RegionBucket)
	if err != nil {
		return regions, errors.Wrap(err, "error while retrieving all keys")
	}

	for _, value := range keys {
		var region Region
		regionBytes, err := json.Marshal(value)
		if err != nil {
			break
		}
		err = json.Unmarshal(regionBytes, region)
		if err != nil {
			return regions, err
		}
		regions = append(regions, region)
	}

	return regions, nil
}
