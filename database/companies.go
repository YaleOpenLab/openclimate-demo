package database

import (
	"encoding/json"
	utils "github.com/Varunram/essentials/utils"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
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
	db, err := OpenDB()
	if err != nil {
		return company, errors.Wrap(err, "error while opening database")
	}
	defer db.Close()
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(UserBucket)
		x := b.Get(utils.ItoB(key))
		if x == nil {
			return errors.New("retrieved user nil, quitting!")
		}
		return json.Unmarshal(x, &company)
	})
	return company, err
}

// ValidateUser validates a particular user
func RetrieveCompanyByName(name string, country string) (Company, error) {
	var company Company
	temp, err := RetrieveAllUsers()
	if err != nil {
		return company, errors.Wrap(err, "error while retrieving all users from database")
	}
	limit := len(temp) + 1
	db, err := OpenDB()
	if err != nil {
		return company, errors.Wrap(err, "could not open db, quitting!")
	}
	defer db.Close()
	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(UserBucket)
		for i := 1; i < limit; i++ {
			var rCompany Company
			x := b.Get(utils.ItoB(i))
			err := json.Unmarshal(x, &rCompany)
			if err != nil {
				return errors.Wrap(err, "could not unmarshal json, quitting!")
			}
			// check name
			if rCompany.Name == name && rCompany.Country == country {
				company = rCompany
				return nil
			}
		}
		return errors.New("Not Found")
	})
	return company, err
}

// RetrieveAllUsers gets a list of all User in the database
func RetrieveAllCompanies() ([]Company, error) {
	var arr []Company
	db, err := OpenDB()
	if err != nil {
		return arr, errors.Wrap(err, "Error while opening database")
	}
	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(UserBucket)
		for i := 1; ; i++ {
			var rCompany Company
			x := b.Get(utils.ItoB(i))
			if x == nil {
				return nil
			}
			err := json.Unmarshal(x, &rCompany)
			if err != nil {
				return errors.Wrap(err, "Error while unmarshalling json")
			}
			arr = append(arr, rCompany)
		}
	})
	return arr, err
}

// Save inserts a passed User object into the database
func (a *Company) Save() error {
	db, err := OpenDB()
	if err != nil {
		return errors.Wrap(err, "Error while opening database")
	}
	defer db.Close()
	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(UserBucket)
		encoded, err := json.Marshal(a)
		if err != nil {
			return errors.Wrap(err, "Error while marshaling json")
		}
		return b.Put([]byte(utils.ItoB(a.Index)), encoded)
	})
	return err
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
