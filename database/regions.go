package database

import (
	"encoding/json"
	edb "github.com/Varunram/essentials/database"
	globals "github.com/YaleOpenLab/openclimate/globals"
	"github.com/pkg/errors"
	"log"
)

// Our definition of "City" includes states, 
// regions, provinces, prefectures, etc. The
// following struct defines the relevant fields.
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

// Function that creates a new region object given its name and country
// and saves the object in the regions bucket.
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

	err = new.Save()
	return new, err

}

// Saves region object in regions bucket. Called by NewRegion
func (region *Region) Save() error {
	return edb.Save(globals.DbPath, RegionBucket, region, region.Index)
}

// Given a key of type int, retrieves the corresponding region object
// from the database regions bucket.
func RetrieveRegion(key int) (Region, error) {
	var region Region
	regionBytes, err := edb.Retrieve(globals.DbPath, RegionBucket, key)
	if err != nil {
		return region, errors.Wrap(err, "error while retrieving key from bucket")
	}
	err = json.Unmarshal(regionBytes, &region)
	if err != nil {
		return region, errors.Wrap(err, "could not unmarshal json, quitting")
	}
	return region, nil
}

// Given the name and country of the region, retrieves the 
// corresponding region object from the database regions bucket.
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

// Retrieves all regions from the regions bucket.
func RetrieveAllRegions() ([]Region, error) {
	var regions []Region
	keys, err := edb.RetrieveAllKeys(globals.DbPath, RegionBucket)
	if err != nil {
		return regions, errors.Wrap(err, "error while retrieving all keys")
	}

	for _, val := range keys {
		var region Region
		err = json.Unmarshal(val, &region)
		if err != nil {
			return regions, err
		}
		regions = append(regions, region)
	}

	return regions, nil
}
