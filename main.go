package main

import (
	"github.com/pkg/errors"
	"log"

	//"github.com/YaleOpenLab/openclimate/blockchain"
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
	// Interact with the blockchain and check token balance
	blockchain.CheckTokenBalance()

	var err error
	err = loadGlobals()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("PRIVATEKEY: %s \nPRIVATE KEY PASSWORD: %s", globals.PrivateKey, globals.PrivateKeyPassword)

	_, err = database.RetrieveAllUsers()
	if err != nil {
		log.Println(err)
		database.CreateHomeDir()
	}


	// log.Println(blockchain.CommitToChain("ethereum", "top", "secret"))
	server.StartServer("8001", true)

}
