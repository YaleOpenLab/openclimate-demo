package main

import (
	"github.com/pkg/errors"
	"log"

	"github.com/YaleOpenLab/openclimate/blockchain"
	"github.com/YaleOpenLab/openclimate/database"
	"github.com/YaleOpenLab/openclimate/globals"
	"github.com/YaleOpenLab/openclimate/server"
	//"github.com/YaleOpenLab/openclimate/notif"
	"github.com/spf13/viper"
)

func loadGlobals() error {
	var err error
	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err = viper.ReadInConfig()
	if err != nil {
		return errors.New("Error while reading platform email from config file")
	}

	globals.PrivateKey = viper.Get("privkey").(string)
	globals.PrivateKeyPassword = viper.Get("privkeypassword").(string)
	return nil
}

func main() {
	var err error
	err = loadGlobals()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("PRIVATEKEY: %s \nPRIVATE KEY PASSWORD: %s", globals.PrivateKey, globals.PrivateKeyPassword)

	_, err = database.RetrieveAllUsers()
	if err != nil {
		database.CreateHomeDir()
	}

	user, err := database.NewUser("name", "9a768ace36ff3d1771d5c145a544de3d68343b2e76093cb7b2a8ea89ac7f1a20c852e6fc1d71275b43abffefac381c5b906f55c3bcff4225353d02f1d3498758", "email")
	log.Println(user)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(blockchain.CommitToChain("ethereum", "blah", "cool"))
	server.StartServer("8001", true)
}
