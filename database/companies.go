package database

import (
	"encoding/json"
	edb "github.com/Varunram/essentials/database"
	globals "github.com/YaleOpenLab/openclimate/globals"
	"github.com/pkg/errors"
	// "log"
)

// Our definition of "Company" includes ....
// The following struct defines the relevant fields.
type Company struct {

	// Identifying info
	Index   int
	Name    string
	Description string

	Locations []Location
	Accountability []DistributionRecord

	Country string
	Address string

	UserIDs []int

	// Contextual data
	Area        float64
	Iso         string
	Population  int
	// Latitude    float64
	// Longitude   float64
	Revenue     float64
	CompanySize int
	HQ          string

	MultiNational []string // an array of all the countries a company is in; if not an MNC, leave empty
	ForProfit     bool
	Industry      bool

	MRV string // the company's chosen MRV reporting methodology

	Pledges []int

	// Nested scopes
	States    []int
	Regions   []int
	Countries []int

	// The entity IDs of all the company's physical assets
	Assets []int

	// IDs of all the company's financial/regulatory assets (e.g. RECs, climate bonds, etc.)
	Credits []int

	// // Data that is reported (through self-reporting, databases, IoT, etc.)
	// // as opposed to data that is aggregated from its parts/children. Data
	// // is stored on IPFS, so Reports holds the IPFS hashes.
	// Reports []RepData

	// Emissions  map[string]string // accept whatever emissions the frontend passes
	// Mitigation map[string]string
	// Adaptation map[string]string
}

func (c *Company) AddPledges(pledgeIDs ...int) error {
	c.Pledges = append(c.Pledges, pledgeIDs...)
	return c.Save()
}

func (c Company) GetPledges() ([]Pledge, error) {
	var pledges []Pledge

	for _, id := range c.Pledges {
		p, err := RetrievePledge(id)
		if err != nil {
			return pledges, errors.Wrap(err, "The Company method GetPledges() failed.")
		}
		pledges = append(pledges, p)
	}
	return pledges, nil
}

func (c *Company) AddAssets(assetIDs ...int) error {
	c.Assets = append(c.Assets, assetIDs...)
	return c.Save()
}

func (c *Company) GetAssetsByState(state string) ([]Asset, error) {
	var assets []Asset
	for _, id := range c.Assets {
		a, err := RetrieveAsset(id)
		if err != nil {
			return assets, errors.Wrap(err, "The Company method GetAssetsByState() failed.")
		}
		if a.State == state {
			assets = append(assets, a)
		}
	}
	return assets, nil
}

func (c *Company) AddStates(stateIDs ...int) error {
	c.States = append(c.States, stateIDs...)
	return c.Save()
}

func (c *Company) GetStates() ([]State, error) {
	var states []State
	for _, id := range c.States {
		s, err := RetrieveState(id)
		if err != nil {
			return states, errors.Wrap(err, "The Company method GetStates() failed.")
		}
		states = append(states, s)
	}
	return states, nil
}

func (c *Company) AddRegions(regionIDs ...int) error {
	c.Regions = append(c.Regions, regionIDs...)
	return c.Save()
}

func (c *Company) GetRegions() ([]Region, error) {
	var regions []Region
	for _, id := range c.Regions {
		r, err := RetrieveRegion(id)
		if err != nil {
			return regions, errors.Wrap(err, "The Company method GetRegions() failed.")
		}
		regions = append(regions, r)
	}
	return regions, nil
}

func (c *Company) AddCountries(countryIDs ...int) error {
	c.Countries = append(c.Countries, countryIDs...)
	return c.Save()
}

func (c *Company) GetCountries() ([]Country, error) {
	var countries []Country
	for _, id := range c.Countries {
		country, err := RetrieveCountry(id)
		if err != nil {
			return countries, errors.Wrap(err, "The Company method GetCountries() failed.")
		}
		countries = append(countries, country)
	}
	return countries, nil
}

// Function that creates a new company object given its name
// and country and saves the object in the countries bucket.
func NewCompany(name string, country string) (Company, error) {
	var company Company
	company.Name = name
	company.Country = country
	return company, company.Save()
}

// Given a key of type int, retrieves the corresponding company object
// from the database companies bucket.
func RetrieveCompany(key int) (Company, error) {
	var company Company
	companyBytes, err := edb.Retrieve(globals.DbPath, CompanyBucket, key)
	if err != nil {
		return company, errors.Wrap(err, "error while retrieving key from bucket")
	}
	err = json.Unmarshal(companyBytes, &company)
	return company, err
}

// Given a name and country, retrieves the corresponding company object
// from the database companies bucket.
func RetrieveCompanyByName(name string, country string) (Company, error) {
	var company Company
	temp, err := RetrieveAllCompanies()
	if err != nil {
		return company, errors.Wrap(err, "error while retrieving all users from database")
	}

	for _, company := range temp {
		if company.Name == name && company.Country == country {
			return company, nil
		}
	}

	return company, errors.New("company not found, quitting")
}

// RetrieveAllCompanies gets a list of all companies in the database
func RetrieveAllCompanies() ([]Company, error) {
	var companies []Company
	keys, err := edb.RetrieveAllKeys(globals.DbPath, CompanyBucket)
	if err != nil {
		return companies, errors.Wrap(err, "could not retrieve all user keys")
	}
	for _, val := range keys {
		var x Company
		err = json.Unmarshal(val, &x)
		if err != nil {
			break
		}
		companies = append(companies, x)
	}

	return companies, nil
}


func RetrieveAllMultiNationals() ([]Company, error) {

	var multinationals []Company

	companies, err := RetrieveAllCompanies()
	if err != nil {
		return multinationals, errors.Wrap(err, "RetrieveAllMultiNationals() failed")
	}

	for _, c := range companies {
		if c.MultiNational != nil {
			multinationals = append(multinationals, c)
		}
	}

	return multinationals, nil
}

