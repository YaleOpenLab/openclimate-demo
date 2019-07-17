 package database

import (
	"log"
)

func Populate() {
	PopulateUSStates()
	PopulateCountries()
}

// Test function populating the regions bucket with the US states
func PopulateUSStates() {
	USStates = []string{"Alabama", "Alaska", "American Samoa", "Arizona", "Arkansas", "California", "Colorado", "Connecticut", "Delaware", "District of Columbia", "Federated States of Micronesia", "Florida", "Georgia", "Guam", "Hawaii", "Idaho", "Illinois", "Indiana", "Iowa", "Kansas", "Kentucky", "Louisiana", "Maine", "Marshall Islands", "Maryland", "Massachusetts", "Michigan", "Minnesota", "Mississippi", "Missouri", "Montana", "Nebraska", "Nevada", "New Hampshire", "New Jersey", "New Mexico", "New York", "North Carolina", "North Dakota", "Northern Mariana Islands", "Ohio", "Oklahoma", "Oregon", "Palau", "Pennsylvania", "Puerto Rico", "Rhode Island", "South Carolina", "South Dakota", "Tennessee", "Texas", "Utah", "Vermont", "Virgin Island", "Virginia", "Washington", "West Virginia", "Wisconsin", "Wyoming"}
	for _, state := range USStates {
		_, err := NewRegion(state, "USA")
		if err != nil {
			log.Println("INSIDE POPULATE: ", err)
			return
		}
	}
}

func PopulateRegionsTest() {
	_, err := NewRegion("Shanghai", "China")
	_, err = NewRegion("Osaka", "Japan")
	_, err = NewRegion("Cancun", "Mexico")
	_, err = NewRegion("Addis Ababa", "Ethiopia")
	log.Println(err)
}

// Test function populating the countries bucket with dummy values
// to test the rpc endpoint for countries
func PopulateCountries() {
	countries := []string{"USA", "China", "Japan", "Mexico", "Ethiopia"}
	for _, country := range countries {
		_, err := NewCountry(country)
		if err != nil {
			log.Println(err, "Failed to add countries")
			return
		}
	}
}
