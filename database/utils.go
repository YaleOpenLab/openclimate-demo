package database

import (
	"encoding/json"
	"github.com/pkg/errors"
	edb "github.com/Varunram/essentials/database"
	"github.com/Varunram/essentials/utils"
	"github.com/boltdb/bolt"
)


func Save(dir string, bucketName []byte, x interface{}) (int, error) {
	db, err := edb.OpenDB(dir)
	if err != nil {
		return 0, errors.Wrap(err, "could not open database")
	}
	defer db.Close()

	var id *int

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
		newID, _ := b.NextSequence()

		// Save the id to the int variable pointed to by *id
		*id = int(newID)

		iK, err := utils.ToByte(newID)
		if err != nil {
			return err
		}

		// Put bytes to bucket
		return b.Put(iK, encoded)
	})
	return *id, err
}