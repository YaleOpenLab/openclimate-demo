package main

import (
	// "github.com/YaleOpenLab/openclimate/blockchain"
	"github.com/YaleOpenLab/openclimate/database"
	// "github.com/YaleOpenLab/openclimate/oracle"
	"github.com/YaleOpenLab/openclimate/server"
	"log"
	// "math/big"
)

func main() {
	// oracle.Schedule()
	// blockchain.CheckTokenBalance()
	// blockchain.CommitToChain(big.NewInt(1565752648), "0x4920636172652061626f757420636c696d617465")
	database.FlushDB()
	database.CreateHomeDir()

	var a database.User
	a.Username = "cool"
	a.Pwhash = "cool"
	a.Index = 1
	err := a.Save()
	if err != nil {
		log.Fatal(err)
	}
	log.Println("flushed and created new db")
	server.StartServer("8001", true)
}
