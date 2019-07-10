package database

import (
	"log"
	"os"

	utils "github.com/Varunram/essentials/utils"
	"github.com/YaleOpenLab/openclimate/globals"
	"github.com/boltdb/bolt"
)

var UserBucket = []byte("Users")
var CompanyBucket = []byte("Companies")
var RegionBucket = []byte("Regions")
var CityBucket = []byte("Cities")

// CreateHomeDir creates a home directory
func CreateHomeDir() {
	if _, err := os.Stat(globals.HomeDir); os.IsNotExist(err) {
		// directory does not exist, create one
		log.Println("Creating home directory")
		os.MkdirAll(globals.HomeDir, os.ModePerm)
	}

	if _, err := os.Stat(globals.DbDir); os.IsNotExist(err) {
		os.MkdirAll(globals.DbDir, os.ModePerm)
	}
}

// don't lock since boltdb can only process one operation at a time. As the application
// grows bigger, this would be a major reason to search for a new db system

// OpenDB opens the db
func OpenDB() (*bolt.DB, error) {
	// we need to check and create this directory if it doesn't exist
	db, err := bolt.Open(globals.DbDir+"/openclimate.db", 0600, nil) // store this in its ownd database
	if err != nil {
		log.Println("Couldn't open database, exiting!")
		return db, err
	}
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists(UserBucket) // the projects bucket contains all our projects
		if err != nil {
			log.Println("Error while creating projects bucket", err)
			return err
		}
		_, err = tx.CreateBucketIfNotExists(CompanyBucket) // the projects bucket contains all our projects
		if err != nil {
			log.Println("Error while creating projects bucket", err)
			return err
		}
		return nil
	})
	return db, err
}

// DeleteKeyFromBucket deletes a given key from the bucket bucketName but doesn
// not shift indices of elements succeeding the deleted element's index
func DeleteKeyFromBucket(key int, bucketName []byte) error {
	// deleting project might be dangerous since that would mess with the other
	// functions, have it in here for now, don't do too much with it / fiox retrieve all
	// to handle this case
	db, err := OpenDB()
	if err != nil {
		return err
	}
	defer db.Close()
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucketName)
		b.Delete(utils.ItoB(key))
		return nil
	})
}
