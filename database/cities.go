package database

import (
	"encoding/json"
	edb "github.com/Varunram/essentials/database"
	utils "github.com/Varunram/essentials/utils"
	globals "github.com/YaleOpenLab/openclimate/globals"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

// includes cities, municipalities, towns, shires, villages, communes, etc.
type City struct {
	Index       int
	Name        string
	Region      string
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

func NewCity(name string, region string, country string) (City, error) {

	var new City
	var err error

	// naive implementation of assigning keys to bucket items (simple indexing)
	cities, err := RetrieveAllCities()
	lenCities := len(cities)

	if err != nil {
		return new, errors.Wrap(err, "Error while retrieving all cities from db")
	}

	if lenCities == 0 {
		new.Index = 1
	} else {
		new.Index = lenCities + 1
	}

	new.Name = name
	new.Region = region
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

/*
	ISSUE: edb.Save() asks for an key argument of type INT,
	but currently we are passing in a key argument of type string.
	This issue needs to be resolved. Could maybe just use a hash.

	RESOLVED: currently using solution previously implemented in OpenX;
	incrementing index for each new city, so the key is of type int.
*/

func (city *City) Save() error {
	return edb.Save(globals.DbDir, CityBucket, city, city.Index)
}

func RetrieveCity(key int) (City, error) {
	var city City
	temp, err := edb.Retrieve(globals.DbDir, CityBucket, key)

	if err != nil {
		return city, errors.Wrap(err, "Error while retrieving key from bucket")
	}

	city = temp.(City)
	return city, city.Save()
}

func RetrieveCityByName(name string, region string) (City, error) {
	var city City
	allCities, err := RetrieveAllCities()
	if err != nil {
		return city, errors.Wrap(err, "Error while retrieving all cities from database")
	}

	db, err := OpenDB()
	if err != nil {
		return city, errors.Wrap(err, "Could not open database, quitting")
	}

	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(CityBucket)

		limit := len(allCities) + 1
		for i := 1; i < limit; i++ {
			var tempCity City
			tempKey := bucket.Get(utils.ItoB(i))

			err := json.Unmarshal(tempKey, &tempCity)
			if err != nil {
				return errors.Wrap(err, "Could not unmarshal json, quitting")
			}

			if tempCity.Name == name && tempCity.Region == region {
				city = tempCity
				return nil
			}
		}
		return errors.New("City not found.")
	})
	return city, err
}

func RetrieveAllCities() ([]City, error) {
	var cities []City
	keys, err := edb.RetrieveAllKeys(globals.DbDir, CityBucket)
	if err != nil {
		return cities, errors.Wrap(err, "error while retrieving all keys")
	}

	for _, value := range keys {
		var city City
		err := json.Unmarshal(value, city)
		if err != nil {
			return cities, err
		}
		cities = append(cities, city)
	}

	return cities, nil
}
