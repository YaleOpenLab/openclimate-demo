package main

import (
	"log"

	"github.com/YaleOpenLab/openclimate/database"
)

func main() {
	err := database.DeleteUser("Jerry", 3)
	if err != nil {
		log.Fatal(err)
	}
}
