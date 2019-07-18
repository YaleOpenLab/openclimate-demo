package database

import (
	"github.com/pkg/errors"
	"log"
	"os"

	edb "github.com/Varunram/essentials/database"
	"github.com/YaleOpenLab/openclimate/globals"
	"github.com/boltdb/bolt"
)

var UserBucket = []byte("Users")
var CompanyBucket = []byte("Companies")
var RegionBucket = []byte("Regions")
var CityBucket = []byte("Cities")
var CountryBucket = []byte("Countries")
var AssetBucket = []byte("Assets")

// CreateHomeDir creates a home directory
func CreateHomeDir() error {
	if _, err := os.Stat(globals.HomeDir); os.IsNotExist(err) {
		// directory does not exist, create one
		err = os.MkdirAll(globals.HomeDir, os.ModePerm)
		if err != nil {
			return errors.Wrap(err, "could not create directory")
		}
	}

	if _, err := os.Stat(globals.DbDir); os.IsNotExist(err) {
		err = os.MkdirAll(globals.DbDir, os.ModePerm)
		if err != nil {
			return errors.Wrap(err, "could not create directory")
		}
		_, err = os.Create(globals.DbPath)
		db, err := edb.CreateDB(globals.DbPath, UserBucket, CompanyBucket, RegionBucket, CityBucket, CountryBucket)
		if err != nil {
			return errors.Wrap(err, "could not create database")
		}
		db.Close()
	}

	log.Println("created new db and home directory")
	return nil
}

func FlushDB() error {
	if _, err := os.Stat(globals.HomeDir); os.IsNotExist(err) {
	} else {
		// directory exists, flush db
		log.Println("deleting database")
		return os.RemoveAll(globals.HomeDir)
	}
	return nil
}

// don't lock since boltdb can only process one operation at a time. As the application
// grows bigger, this would be a major reason to search for a new db system

// OpenDB opens the db
func OpenDB() (*bolt.DB, error) {
	return edb.CreateDB(globals.DbPath, UserBucket, CompanyBucket, RegionBucket, CityBucket, CountryBucket)
}

// DeleteKeyFromBucket deletes a given key from the bucket bucketName but doesn
// not shift indices of elements succeeding the deleted element's index
func DeleteKeyFromBucket(key int, bucketName []byte) error {
	return edb.DeleteKeyFromBucket(globals.DbPath, key, bucketName)
}
