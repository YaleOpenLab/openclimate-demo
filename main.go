package main

import (
	"github.com/pkg/errors"
	"log"

	//"github.com/YaleOpenLab/openclimate/database"
	"github.com/YaleOpenLab/openclimate/blockchain"
	"github.com/YaleOpenLab/openclimate/globals"
	"github.com/YaleOpenLab/openclimate/server"
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

	log.Println(blockchain.CommitToChain("ethereum", "blah", "cool"))
	server.StartServer("8001", true)
}
