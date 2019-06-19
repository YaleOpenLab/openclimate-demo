package database

import (
	"log"
	"testing"
)

func TestUsers(t *testing.T) {
	// delete the db here first and create it again. Use only for local testing at the moment
	users, err := RetrieveAllUsers()
	if err != nil {
		log.Fatal(err)
	}

	initLength := len(users)
	var user User
	user.Name = "Jerry"
	user.Email = "jerry@test.com"
	user.Pwhash = "nicetry"

	user, err = PutUser(user)
	if err != nil {
		log.Fatal(err)
	}

	//retrieve the user that we just inserted into the database
	user, err = RetrieveUser("Jerry")
	if err != nil {
		log.Fatal(err)
	}

	users, err = RetrieveAllUsers()
	if err != nil {
		log.Fatal(err)
	}

	if len(users)-initLength != 1 {
		t.Fatal("user length not 1, quitting")
	}

	user.Name = "GeorgeCool"
	user.Email = "george@example.com"
	user.Pwhash = "nicetry"
	err = user.Save()
	if err != nil {
		t.Fatal(err)
	}
}
