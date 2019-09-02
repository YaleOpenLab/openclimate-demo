package database

import (
	"encoding/json"
	edb "github.com/Varunram/essentials/database"
	globals "github.com/YaleOpenLab/openclimate/globals"
	"github.com/pkg/errors"
)

// Nation-states, countries
type Country struct {

	// Identifying info
	Index       int
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`

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

	// For countries: children = regions
	Children       []string
	Credits        []int
	Pledges        []int                `json:"pledges"`
	Logo           string               `json:"logo"`
	Accountability []DistributionRecord `json:"accountability"`

	// Data that is reported (through self-reporting, databases, IoT, etc.)
	// as opposed to data that is aggregated from its parts/children. Data
	// is stored on IPFS, so Reports holds the IPFS hashes.
	// Reports []RepData

	Emissions  map[string]string `json:"emissions"` // accept whatever emissions data the frontend passes
	Mitigation map[string]string
	Adaptation map[string]string

	LastUpdated string
	Files       []string // an a rray of all the necessary documents to validate this specific country
}

// Function that creates a new country object given its name and saves
// the object in the countries bucket.
func NewCountry(name string) (Country, error) {
	var new Country
	new.Name = name
	return new, new.Save()
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
	return country, err
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
			return countries, errors.Wrap(err, "could not unmarshal json")
		}
		countries = append(countries, country)
	}

	return countries, nil
}

func (c *Country) AddPledges(pledgeIDs ...int) error {
	c.Pledges = append(c.Pledges, pledgeIDs...)
	return c.Save()
}

func (c Country) GetPledges() ([]Pledge, error) {
	var pledges []Pledge

	for _, id := range c.Pledges {
		p, err := RetrievePledge(id)
		if err != nil {
			return pledges, errors.Wrap(err, "The Country method GetPledges() failed.")
		}
		pledges = append(pledges, p)
	}
	return pledges, nil
}

// func (c Country) GetStates() ([]State, error) {
// 	var states []State
// 	states, err := RetrieveAllStates()
// 	if err != nil {
// 		return states, errors.Wrap(err, "country.GetStates() failed")
// 	}

// 	for _, state := range states {
// 		if state.Country == c.Name {
// 			states := append(states, state)
// 		}
// 	}
// 	return states, nil
// }
