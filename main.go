package main

import (
	"log"

	"github.com/YaleOpenLab/openclimate/database"
)

func main() {
	user, err := database.RetrieveUser("Jerry")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(user)
	user, err = database.PutUser(user)
	if err != nil {
		log.Fatal(err)
	}
}
