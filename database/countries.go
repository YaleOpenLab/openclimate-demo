package database

import (
	"encoding/json"
	edb "github.com/Varunram/essentials/database"
	globals "github.com/YaleOpenLab/openclimate/globals"
	"github.com/pkg/errors"
	"log"
)

// Nation-states, countries
type Country struct {

	// Identifying info
	Index int
	Name  string

	// Contextual data
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

// Function that creates a new country object given its name and saves
// the object in the countries bucket.
func NewCountry(name string) (Country, error) {

	var new Country
	var err error

	// naive implementation of assigning keys to bucket items (simple indexing)
	countries, err := RetrieveAllCountries()
	if err != nil {
		log.Println(err)
		return new, errors.Wrap(err, "could not retreive all countries")
	}
	lenCountries := len(countries)
	if err != nil {
		return new, errors.Wrap(err, "Error while retrieving all countries from db")
	}

	if lenCountries == 0 {
		new.Index = 1
	} else {
		new.Index = lenCountries + 1
	}

	new.Name = name

	err = new.Save()
	return new, err

}

// Saves country object in countries bucket. Called by NewCountry
func (country *Country) Save() error {
	return edb.Save(globals.DbPath, CountryBucket, country, country.Index)
}

// Given a key of type int, retrieves the corresponding country object
// from the database countries bucket.
func RetrieveCountry(key int) (Country, error) {
	var country Country
	countryBytes, err := edb.Retrieve(globals.DbPath, CountryBucket, key)
	if err != nil {
		return country, errors.Wrap(err, "error while retrieving key from bucket")
	}
	err = json.Unmarshal(countryBytes, &country)
	if err != nil {
		return country, errors.Wrap(err, "could not unmarshal json, quitting")
	}
	return country, nil
}

// Given the name of the country, retrieves the corresponding country object
// from the database countries bucket.
func RetrieveCountryByName(name string) (Country, error) {
	var country Country
	allCountries, err := RetrieveAllCountries()
	if err != nil {
		return country, errors.Wrap(err, "Error while retrieving all countries from database")
	}

	for _, val := range allCountries {
		if val.Name == name {
			country = val
			return country, nil
		}
	}

	return country, errors.New("could not find countries")
}

// Retrieves all countries from the countries bucket.
func RetrieveAllCountries() ([]Country, error) {
	var countries []Country
	keys, err := edb.RetrieveAllKeys(globals.DbPath, CountryBucket)
	if err != nil {
		return countries, errors.Wrap(err, "error while retrieving all keys")
	}

	for _, val := range keys {
		var country Country
		err = json.Unmarshal(val, &country)
		if err != nil {
			return countries, err
		}
		countries = append(countries, country)
	}

	return countries, nil
}
