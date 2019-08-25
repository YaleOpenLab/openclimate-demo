package database

import (
	"encoding/json"
	"github.com/pkg/errors"
	"sort"
	//"log"

	edb "github.com/Varunram/essentials/database"
	globals "github.com/YaleOpenLab/openclimate/globals"
)

// Our definition of "State" includes states,
// provinces, prefectures, etc. The
// following struct defines the relevant fields.
type State struct {

	// Identifying info
	Index   int
	Name    string
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
	// EntityType		string

	MRV string

	Pledges []int

	//	For states: children = companies (divided by state)
	Children []string

	// Data that is reported (through self-reporting, databases, IoT, etc.)
	// as opposed to data that is aggregated from its parts/children. Data
	// is stored on IPFS, so Reports holds the IPFS hashes.
	// Reports []RepData

	Emissions  map[string]string // accept whatever emissions the frontend passes
	Mitigation map[string]string
	Adaptation map[string]string
}

// Function that creates a new state object given its name and country
// and saves the object in the states bucket.
func NewState(name string, country string) (State, error) {
	var new State
	new.Name = name
	new.Country = country
	return new, new.Save()
}

// Given a key of type int, retrieves the corresponding state object
// from the database states bucket.
func RetrieveState(key int) (State, error) {
	var state State
	stateBytes, err := edb.Retrieve(globals.DbPath, StateBucket, key)
	if err != nil {
		return state, errors.Wrap(err, "error while retrieving key from bucket")
	}
	err = json.Unmarshal(stateBytes, &state)
	return state, err
}

// Given the name and country of the state, retrieves the
// corresponding state object from the database states bucket.
func RetrieveStateByName(name string, country string) (State, error) {
	var state State
	allStates, err := RetrieveAllStates()
	if err != nil {
		return state, errors.Wrap(err, "Error while retrieving all states from database")
	}

	for _, val := range allStates {
		if val.Name == name && val.Country == country {
			state = val
			return state, nil
		}
	}

	return state, errors.New("could not find states")
}

// Retrieves all states from the states bucket.
func RetrieveAllStates() ([]State, error) {
	var states []State
	keys, err := edb.RetrieveAllKeys(globals.DbPath, StateBucket)
	if err != nil {
		return states, errors.Wrap(err, "error while retrieving all states")
	}

	for _, val := range keys {
		var state State
		err = json.Unmarshal(val, &state)
		if err != nil {
			return states, errors.Wrap(err, "could not unmarshal json")
		}
		states = append(states, state)
	}

	return states, nil
}

// Retrieves and filters state by country.
func FilterStatesByCountry(country string) ([]State, error) {
	var states []State
	keys, err := edb.RetrieveAllKeys(globals.DbPath, StateBucket)
	if err != nil {
		return states, errors.Wrap(err, "error while retrieving filtered states")
	}

	for _, val := range keys {
		var state State
		err = json.Unmarshal(val, &state)
		if state.Country == country {
			if err != nil {
				return states, errors.Wrap(err, "could not unmarshal json")
			}
			states = append(states, state)
		}
	}

	sort.Slice(states, func(i, j int) bool { return states[i].Name < states[j].Name })
	return states, nil
}

func (c *State) AddPledges(pledgeIDs ...int) error {
	c.Pledges = append(c.Pledges, pledgeIDs...)
	return c.Save()
}

func (c State) GetPledges() ([]Pledge, error) {
	var pledges []Pledge

	for _, id := range c.Pledges {
		p, err := RetrievePledge(id)
		if err != nil {
			return pledges, errors.Wrap(err, "The State method GetPledges() failed")
		}
		pledges = append(pledges, p)
	}
	return pledges, nil
}
