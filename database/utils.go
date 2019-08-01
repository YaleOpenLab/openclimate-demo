package database

import (
	"encoding/json"
	"github.com/pkg/errors"
	edb "github.com/Varunram/essentials/database"
	"github.com/Varunram/essentials/utils"
	"github.com/boltdb/bolt"
)


func Save(dir string, bucketName []byte, x Actor) error {
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

		encoded, err := json.Marshal(x)
		if err != nil {
			return errors.Wrap(err, "error while marshaling json struct")
		}

		// Generate ID for the user.
        // This returns an error only if the Tx is closed or not writeable.
        // That can't happen in an Update() call so I ignore the error check.
		id, _ := b.NextSequence()

		// Set the id of the user
		x.SetID(int(id))

		iK, err := utils.ToByte(id)
		if err != nil {
			return err
		}

		// Put bytes to bucket
		return b.Put(iK, encoded)
	})
	return err
}