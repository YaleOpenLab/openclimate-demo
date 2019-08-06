package database

import (
	"encoding/json"
	edb "github.com/Varunram/essentials/database"
	globals "github.com/YaleOpenLab/openclimate/globals"
	"github.com/pkg/errors"
)

type ConnectRequest struct {
	Index   int
	DBName  string
	OrgName string

	DBActorTypes  []string // what type of actors does the DB cover?
	DBActionTypes []string // what type of actions does the DB track?

	ContactInfo string
	Links       []string // links with more info
}

func NewRequest(request ConnectRequest) error {
	return request.Save()
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
