package main

import (
	"log"

	"github.com/YaleOpenLab/openclimate/database"
)

func main() {
	user, err := database.AuthUser("Jerry", "nicetry")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(user)
}
