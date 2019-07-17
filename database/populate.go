 package database

import (
	"log"
)

func Populate() {
	PopulateUSStates()
	PopulateCountries()
	PopulateRegionsTest()
	PopulateAdminUsers()
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
	if err != nil {
		log.Println(err, "Failed populate regions test")
		return
	}
	_, err = NewRegion("Osaka", "Japan")
	if err != nil {
		log.Println(err, "Failed populate regions test")
		return
	}
	_, err = NewRegion("Cancun", "Mexico")
	if err != nil {
		log.Println(err, "Failed populate regions test")
		return
	}
	_, err = NewRegion("Addis Ababa", "Ethiopia")
	if err != nil {
		log.Println(err, "Failed populate regions test")
		return
	}
}


func PopulateAdminUsers() {
	_, err := NewUser("amanda", "9a768ace36ff3d1771d5c145a544de3d68343b2e76093cb7b2a8ea89ac7f1a20c852e6fc1d71275b43abffefac381c5b906f55c3bcff4225353d02f1d3498758", "amanda@test.com", "individual")
	if err != nil {
		log.Println(err, "Failed populate users test")
		return
	}
	_, err = NewUser("brian", "e9a75486736a550af4fea861e2378305c4a555a05094dee1dca2f68afea49cc3a50e8de6ea131ea521311f4d6fb054a146e8282f8e35ff2e6368c1a62e909716", "brian@test.com", "individual")
	if err != nil {
		log.Println(err, "Failed populate users test")
		return
	}
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
