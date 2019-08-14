package main

import (
	// "encoding/json"
	"log"

	"github.com/YaleOpenLab/openclimate/blockchain"
	"github.com/YaleOpenLab/openclimate/database"
	"github.com/YaleOpenLab/openclimate/oracle"
	"github.com/YaleOpenLab/openclimate/server"
	//"github.com/Varunram/essentials/ipfs"
	//"github.com/YaleOpenLab/openclimate/notif"
)

func main() {
	// Interact with the blockchain and check token balance
	data, err := oracle.GetNoaaDailyCO2()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(data)


	blockchain.CheckTokenBalance()
	database.FlushDB()
	database.CreateHomeDir()
	log.Println("flushed and created new db")
	server.StartServer("8001", true)
}
