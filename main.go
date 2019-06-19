package main

import (
	"log"

	"github.com/YaleOpenLab/openclimate/database"
)

func main() {
	var user database.User
	user.Id = 2
	user.Name=  "GeorgeCool"
	user.Email = "george@example.com"
	user.Pwhash = "nicetry"
	err := user.Save()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(user)
}
