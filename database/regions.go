package database

import (
	"encoding/json"
	utils "github.com/Varunram/essentials/utils"
	"github.com/boltdb/bolt"
	"github.com/pkg/errors"
)

// includes states, regions, provinces, prefectures, etc.
type Region struct {
	
	Name			string
	Area			float64
	Iso				string
	EntityType		string
	Population		int
	Latitude		float64
	Longitude		float64
	Revenue			float64
	CompanySize		int
	HQ				string

}

func NewRegion(name string, area float64, population int) (Region, error) {

	var new Region
	var err error

	new.Name = name
	new.Area = area
	new.Population = population

	new.Iso = nil


}