package database

import (
	"encoding/json"
	edb "github.com/Varunram/essentials/database"
	globals "github.com/YaleOpenLab/openclimate/globals"
	"github.com/pkg/errors"
	// "log"
)

type Pledge struct {
	ID        int
	ActorType string
	ActorID   int
	Coop      bool

	/*
		Pledges can be:
		emissions reductions,
		mitigation actions (energy efficiency, renewables, etc.),
		or adaptation actions
	*/
	PledgeType string

	BaseYear float64

	TargetYear float64

	Goal float64

	// is this goal determined by a regulator, or voluntarily
	// adopted by the climate actor?
	Regulatory bool
}

func NewPledge(pledgeType string, baseYear float64, targetYear float64, goal float64,
	regulatory bool, actorType string, actorID int) (Pledge, error) {
	var p Pledge
	p.PledgeType = pledgeType
	p.BaseYear = baseYear
	p.TargetYear = targetYear
	p.Goal = goal
	p.Regulatory = regulatory
	p.ActorType = actorType
	p.ActorID = actorID

	err := p.Save()
	if err != nil {
		return p, errors.Wrap(err, "NewPledge() failed")
	}

	actor, err := RetrieveActor(actorType, actorID)
	if err != nil {
		return p, errors.Wrap(err, "NewPledge() failed")
	}

	err = actor.AddPledges(p.ID)
	if err != nil {
		return p, errors.Wrap(err, "NewPledge() failed")
	}

	return p, nil
}

func UpdatePledge(key int, updated Pledge) error {
	pledge, err := RetrievePledge(key)
	if err != nil {
		return errors.Wrap(err, "UpdatePledge() failed (likely because pledge doesn't exist)")
	}

	// ActorID and PledgeType are not updated because
	// these attributes should not change.

	pledge.BaseYear = updated.BaseYear
	pledge.TargetYear = updated.TargetYear
	pledge.Goal = updated.Goal
	pledge.Regulatory = updated.Regulatory
	return pledge.Save()
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
