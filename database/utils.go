package database

import (
	"encoding/json"
	edb "github.com/Varunram/essentials/database"
	utils "github.com/Varunram/essentials/utils"
	"github.com/YaleOpenLab/openclimate/globals"
	"github.com/pkg/errors"
	// "github.com/Varunram/essentials/utils"
	"github.com/boltdb/bolt"
	// "log"
)

func RetrieveAllUsers() ([]User, error) {
	var arr []User
	x, err := edb.RetrieveAllKeys(globals.DbPath, UserBucket)
	if err != nil {
		return arr, errors.Wrap(err, "error while retrieving all users")
	}
	for _, value := range x {
		var temp User
		err := json.Unmarshal(value, &temp)
		if err != nil {
			return arr, errors.New("error while unmarshalling json, quitting")
		}
		arr = append(arr, temp)
	}

	return arr, nil
}

func RetrieveAllPledges() ([]Pledge, error) {
	var arr []Pledge
	x, err := edb.RetrieveAllKeys(globals.DbPath, PledgeBucket)
	if err != nil {
		return arr, errors.Wrap(err, "error while retrieving all users")
	}
	for _, value := range x {
		var temp Pledge
		err := json.Unmarshal(value, &temp)
		if err != nil {
			return arr, errors.New("error while unmarshalling json, quitting")
		}
		arr = append(arr, temp)
	}
	return arr, nil
}

func Save(dir string, bucketName []byte, x BucketItem) error {
	db, err := edb.OpenDB(dir)
	if err != nil {
		return errors.Wrap(err, "could not open database")
	}
	defer db.Close()

	err = db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		if b == nil {
			return errors.New("Bucket missing")
		}

		// Generate and set ID for the user.
		// This returns an error only if the Tx is closed or not writeable.
		// That can't happen in an Update() call so I ignore the error check.
		keyBytes, err := utils.ToByte(x.GetID())
		if err != nil {
			return err
		}

		elemExists := b.Get(keyBytes) // try to fetch the element from the db
		if elemExists == nil {        // if the element does not exist, assign a new index and create it
			id, _ := b.NextSequence()
			x.SetID(int(id))
		}

		encoded, err := json.Marshal(x)
		if err != nil {
			return errors.Wrap(err, "error while marshaling json struct")
		}

		// Put bytes to bucket
		keyBytes, err = utils.ToByte(x.GetID())
		if err != nil {
			return err
		}
		return b.Put(keyBytes, encoded)
	})
	return err
}
