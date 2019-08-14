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

	// Identifying info
	Index   int
	Name    string
	Region  string
	Country string

	// Contextual data
	Area        float64
	Iso         string
	Population  int
	Latitude    float64
	Longitude   float64
	Revenue     float64
	CompanySize int
	HQ          string

	MRV string

	Children []string

	Pledges []int

	// Data that is reported (through self-reporting, databases, IoT, etc.)
	// as opposed to data that is aggregated from its parts/children. Data
	// is stored on IPFS, so Reports holds the IPFS hashes.
	Reports []RepData

	Emissions  map[string]string // accept whatever emissions the frontend passes
	Mitigation map[string]string
	Adaptation map[string]string
}

// Function that creates a new city object given its name, region,
// and country and saves the object in the countries bucket.
func NewCity(name string, region string, country string) (City, error) {
	var new City
	new.Name = name
	new.Country = country
	return new, new.Save()
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
			return cities, errors.Wrap(err, "could not unmarshal struct")
		}
		cities = append(cities, city)
	}

	return cities, nil
}

func (c *City) AddPledges(pledgeIDs ...int) error {
	c.Pledges = append(c.Pledges, pledgeIDs...)
	return c.Save()
}

func (c City) GetPledges() ([]Pledge, error) {
	var pledges []Pledge

	for _, id := range c.Pledges {
		p, err := RetrievePledge(id)
		if err != nil {
			return pledges, errors.Wrap(err, "The Region method GetPledges() failed.")
		}
		pledges = append(pledges, p)
	}
	return pledges, nil
}
