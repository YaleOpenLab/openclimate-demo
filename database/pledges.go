package database

import (
	// "encoding/json"
	// edb "github.com/Varunram/essentials/database"
	globals "github.com/YaleOpenLab/openclimate/globals"
	// "github.com/pkg/errors"
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
	// log.Println("INSIDE SAVE()")
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


func (p *Pledge) SetID(id int) {
	p.ID = id
}












