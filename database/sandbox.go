package database

import (
	"github.com/pkg/errors"
	"github.com/Varunram/essentials/utils"
)

func Populate() {
	PopulateUSStates()
	PopulateCountries()
	PopulateRegionsTest()
	PopulateTestUsers()
	PopulateAdminUsers()
}

// Test function populating the regions bucket with the US states
func PopulateUSStates() error {
	USStates = []string{"Alabama", "Alaska", "American Samoa", "Arizona", "Arkansas", "California", "Colorado", "Connecticut", "Delaware", "District of Columbia", "Federated States of Micronesia", "Florida", "Georgia", "Guam", "Hawaii", "Idaho", "Illinois", "Indiana", "Iowa", "Kansas", "Kentucky", "Louisiana", "Maine", "Marshall Islands", "Maryland", "Massachusetts", "Michigan", "Minnesota", "Mississippi", "Missouri", "Montana", "Nebraska", "Nevada", "New Hampshire", "New Jersey", "New Mexico", "New York", "North Carolina", "North Dakota", "Northern Mariana Islands", "Ohio", "Oklahoma", "Oregon", "Palau", "Pennsylvania", "Puerto Rico", "Rhode Island", "South Carolina", "South Dakota", "Tennessee", "Texas", "Utah", "Vermont", "Virgin Island", "Virginia", "Washington", "West Virginia", "Wisconsin", "Wyoming"}
	for _, state := range USStates {
		_, err := NewRegion(state, "USA")
		if err != nil {
			return errors.Wrap(err, "could not populate US States")
		}
	}
	return nil
}

func PopulateRegionsTest() error {
	_, err := NewRegion("Shanghai", "China")
	if err != nil {
		return errors.Wrap(err, "Failed populate regions test")
	}
	_, err = NewRegion("Osaka", "Japan")
	if err != nil {
		return errors.Wrap(err, "Failed populate regions test")
	}
	_, err = NewRegion("Cancun", "Mexico")
	if err != nil {
		return errors.Wrap(err, "Failed populate regions test")
	}
	_, err = NewRegion("Addis Ababa", "Ethiopia")
	if err != nil {
		return errors.Wrap(err, "Failed populate regions test")
	}
	return nil
}

func PopulateTestUsers() error {
	pwhash := utils.SHA3hash("a")
	user, err := NewUser("testuser", pwhash, "user@test.com", "country", "USA", "")
	if err != nil {
		return errors.Wrap(err, "failed to create test user in country: USA")
	}
	user.Verified = true
	user.Save()
	return nil
}

func PopulateAdminUsers() error {
	pwhash := utils.SHA3hash("p")

	_, err := NewUser("amanda", pwhash, "amanda@test.com", "individual", "", "")
	if err != nil {
		return errors.Wrap(err, "failed to populate user amanda")
	}
	_, err = NewUser("brian", pwhash, "brian@test.com", "individual", "", "")
	if err != nil {
		return errors.Wrap(err, "failed to populate user brian")
	}

	return nil
}

// Test function populating the countries bucket with dummy values
// to test the rpc endpoint for countries
func PopulateCountries() error {
	countries := []string{"USA", "China", "Japan", "Mexico", "Ethiopia"}
	for _, country := range countries {
		_, err := NewCountry(country)
		if err != nil {
			return errors.Wrap(err, "Failed to add countries")
		}
	}
	return nil
}
