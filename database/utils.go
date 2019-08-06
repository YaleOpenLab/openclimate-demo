package database

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/YaleOpenLab/openclimate/globals"
	edb "github.com/Varunram/essentials/database"
	// "github.com/Varunram/essentials/utils"
	"github.com/boltdb/bolt"
	"encoding/binary"
	"log"
)


func RetrieveAllUsers() ([]User, error) {

	var arr []User
	db, err := edb.OpenDB(globals.DbPath)
	if err != nil {
		return arr, errors.Wrap(err, "RetrieveAllKeys() failed.")
	}

	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(UserBucket)
		if b == nil {
			return errors.New("Bucket is missing")
		}

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var actor User
			err = json.Unmarshal(v, &actor)
			if err != nil {
				return errors.Wrap(err, "RetrieveAllActors() failed")
			}
			arr = append(arr, actor)
		}
		return nil
	})
	return arr, nil
}


func RetrieveAllPledges() ([]Pledge, error) {

	var arr []Pledge
	db, err := edb.OpenDB(globals.DbPath)
	if err != nil {
		return arr, errors.Wrap(err, "RetrieveAllPledges() failed.")
	}

	defer db.Close()

	err = db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(PledgeBucket)
		if b == nil {
			return errors.New("Bucket is missing")
		}

		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			var pledge Pledge
			err = json.Unmarshal(v, &pledge)
			if err != nil {
				return errors.Wrap(err, "RetrieveAllPledges() failed")
			}
			arr = append(arr, pledge)
		}
		return nil
	})
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

		// Generate ID for the user.
        // This returns an error only if the Tx is closed or not writeable.
        // That can't happen in an Update() call so I ignore the error check.
		id, _ := b.NextSequence()
		log.Println(id)

		// Set the id of the user
		x.SetID(int(id))

		encoded, err := json.Marshal(x)
		if err != nil {
			return errors.Wrap(err, "error while marshaling json struct")
		}

		// Put bytes to bucket
		return b.Put(itob(int(id)), encoded)
	})
	return err
}


func itob(v int) []byte {
    b := make([]byte, 8)
    binary.BigEndian.PutUint64(b, uint64(v))
    return b
}
