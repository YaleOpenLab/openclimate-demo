package database

import (
	"encoding/json"
	edb "github.com/Varunram/essentials/database"
	globals "github.com/YaleOpenLab/openclimate/globals"
	"github.com/pkg/errors"
	// "log"
)

type Pledge struct {

	ID 			int
	ActorID 	int

	/* 
		Pledges can be:
		emissions reductions,
		mitigation actions (energy efficiency, renewables, etc.),
		or adaptation actions 
	*/
	PledgeType	string

	BaseYear  	int

	TargetYear 	int

	Goal 		float64

	// is this goal determined by a regulator, or voluntarily
	// adopted by the climate actor?
	Regulatory	bool
}


func (p *Pledge) Save() error {
	return Save(globals.DbPath, PledgeBucket, p)
}


func NewPledge(pledgeType string, baseYear int, targetYear int, goal float64, regulatory bool, actorID int) (Pledge, error) {

	var p Pledge
	p.PledgeType = pledgeType
	p.BaseYear = baseYear
	p.TargetYear = targetYear
	p.Goal = goal
	p.Regulatory = regulatory
	p.ActorID = actorID

	return p, p.Save()
}


func RetrievePledge(key int) (Pledge, error) {
	var pledge Pledge
	pledgeBytes, err := edb.Retrieve(globals.DbPath, PledgeBucket, key)
	if err != nil {
		return pledge, errors.Wrap(err, "error while retrieving key from bucket")
	}
	err = json.Unmarshal(pledgeBytes, &pledge)
	return pledge, err
}


// func UpdatePledge(pledgeID int, pledgeType string, baseYear int, targetYear int, goal float64, regulatory bool, actorID int) (Pledge, error) {

// 	p := RetrievePledge(pledgeID)
// 	if pledgeType != "" {

// 	}
// }


func (p *Pledge) SetID(id int) {
	p.ID = id
}












