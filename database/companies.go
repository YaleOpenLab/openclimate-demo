package database

import (
	"encoding/json"
	edb "github.com/Varunram/essentials/database"
	globals "github.com/YaleOpenLab/openclimate/globals"
	"github.com/pkg/errors"
	"log"
)

type Company struct {
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
}

// RetrieveUser retrieves a particular User indexed by key from the database
func RetrieveCompany(key int) (Company, error) {
	var company Company
	companyBytes, err := RetrieveKey(CompanyBucket, key)
	if err != nil {
		return company, errors.Wrap(err, "could not marshal json, quitting")
	}
	err = json.Unmarshal(companyBytes, &company)
	if err != nil {
		return company, errors.Wrap(err, "could not unmarshal json, quitting")
	}
	return company, nil
}

// ValidateUser validates a particular user
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
		log.Println(err)
		return companies, errors.Wrap(err, "could not retrieve all user keys")
	}
	for _, val := range keys {
		companyBytes, err := json.Marshal(val)
		if err != nil {
			break
		}
		var x Company
		err = json.Unmarshal(companyBytes, &x)
		if err != nil {
			break
		}
		companies = append(companies, x)
	}

	return companies, nil
}

func NewCompany(name string, country string) (Company, error) {
	var company Company

	companies, err := RetrieveAllCompanies()
	if err != nil {
		return company, errors.Wrap(err, "could not retrieve all companies, quitting")
	}

	if len(companies) == 0 {
		company.Index = 1
	} else {
		company.Index = len(companies) + 1
	}

	company.Name = name
	company.Country = country
	return company, company.Save()
}

func (a *Company) Save() error {
	return edb.Save(globals.DbPath, CompanyBucket, a, a.Index)
}
