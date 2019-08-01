package database

import (
	"encoding/json"
	"github.com/pkg/errors"
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

	Pledges []Pledge

	//	For states: children = companies (divided by state)
	Children []string

	// Data that is reported (through self-reporting, databases, IoT, etc.)
	// as opposed to data that is aggregated from its parts/children. Data
	// is stored on IPFS, so Reports holds the IPFS hashes.
	Reports []RepData

	Emissions  map[string]string // accept whatever emissions the frontend passes
	Mitigation map[string]string
	Adaptation map[string]string
}

// Function that creates a new state object given its name and country
// and saves the object in the states bucket.
func NewState(name string, country string) (State, error) {
	var new State
	var err error
	var lenStates int
	// naive implementation of assigning keys to bucket items (simple indexing)
	states, err := RetrieveAllStates()
	if err != nil {
		// states doesn't exist yet
		lenStates = 0
	} else {
		lenStates = len(states)
	}

	new.Index = lenStates + 1
	new.Name = name
	new.Country = country

	return new, new.Save()
}

// Saves state object in states bucket. Called by NewState
func (state *State) Save() error {
	return edb.Save(globals.DbPath, StateBucket, state, state.Index)
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

func (s *State) SetID(id int) {
	s.Index = id
}

func (s *State) AddPledge(pledge Pledge) {
	s.Pledges = append(s.Pledges, pledge)
}
