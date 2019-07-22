package database

import (
	"encoding/json"
	"github.com/pkg/errors"
	globals "github.com/YaleOpenLab/openclimate/globals"
	edb "github.com/Varunram/essentials/database"
)

type ConnectRequest struct {

	Index 			int
	DBName 			string
	OrgName			string

	DBActorTypes 	[]string 	// what type of actors does the DB cover?
	DBActionTypes 	[]string 	// what type of actions does the DB track?

	API 			string
	Links 			[]string 	// links with more info
}


func NewRequest(request ConnectRequest) error {
	allRequests, err := RetrieveAllRequests()
	if err != nil {
		return errors.Wrap(err, "could not retrieve all requests, quitting")
	}

	if len(allRequests) == 0 {
		request.Index = 1
	} else {
		request.Index = len(allRequests) + 1
	}

	return request.Save()
}


func (a *ConnectRequest) Save() error {
	return edb.Save(globals.DbPath, RequestBucket, a, a.Index)
}


func RetrieveAllRequests() ([]ConnectRequest, error) {
	var requests []ConnectRequest
	keys, err := edb.RetrieveAllKeys(globals.DbPath, RequestBucket)
	if err != nil {
		return requests, errors.Wrap(err, "error while retrieving all requests")
	}

	for _, val := range keys {
		var request ConnectRequest
		err = json.Unmarshal(val, &request)
		if err != nil {
			return requests, errors.Wrap(err, "could not unmarshal json")
		}
		requests = append(requests, request)
	}

	return requests, nil
}