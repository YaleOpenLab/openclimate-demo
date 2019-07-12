package database

import (
	"encoding/json"
	edb "github.com/Varunram/essentials/database"
	globals "github.com/YaleOpenLab/openclimate/globals"
	"github.com/pkg/errors"
)

// Our definition of "City" includes cities, municipalities, 
// towns, shires, villages, communes, etc. The following struct
// defines the relevant fields.
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

// Function that creates a new city object given its name, region,
// and country and saves the object in the countries bucket.
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

// Saves city object in cities bucket. Called by NewCity
func (city *City) Save() error {
	return edb.Save(globals.DbPath, CityBucket, city, city.Index)
}

// Given a key of type int, retrieves the corresponding city object
// from the database cities bucket.
func RetrieveCity(key int) (City, error) {
	var city City
	cityBytes, err := edb.Retrieve(globals.DbPath, CityBucket, key)
	if err != nil {
		return city, errors.Wrap(err, "error while retrieving key from bucket")
	}

	err = json.Unmarshal(cityBytes, &city)
	if err != nil {
		return city, errors.Wrap(err, "could not unmarshal json, quitting")
	}
	return city, nil
}

// Given a name and region, retrieves the corresponding city object
// from the database cities bucket.
func RetrieveCityByName(name string, region string) (City, error) {
	var city City
	allCities, err := RetrieveAllCities()
	if err != nil {
		return city, errors.Wrap(err, "Error while retrieving all cities from database")
	}

	for _, val := range allCities {
		if val.Name == name && val.Region == region {
			city = val
			return city, nil
		}
	}

	return city, errors.New("city not found")
}

// Retrieves all countries from the countries bucket.
func RetrieveAllCities() ([]City, error) {
	var cities []City
	keys, err := edb.RetrieveAllKeys(globals.DbPath, CityBucket)
	if err != nil {
		return cities, errors.Wrap(err, "error while retrieving all keys")
	}

	for _, val := range keys {
		var city City
		err = json.Unmarshal(val, city)
		if err != nil {
			return cities, err
		}
		cities = append(cities, city)
	}

	return cities, nil
}
