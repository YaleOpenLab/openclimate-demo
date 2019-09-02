package main

import (
	// "github.com/YaleOpenLab/openclimate/blockchain"
	"github.com/YaleOpenLab/openclimate/database"
	// "github.com/YaleOpenLab/openclimate/oracle"
	"github.com/YaleOpenLab/openclimate/globals"
	"github.com/YaleOpenLab/openclimate/server"
	flags "github.com/jessevdk/go-flags"
	"log"
	"os"
	// "math/big"
)

var opts struct {
	Insecure bool `short:"i" description:"Start the API using http. Not recommended"`
	Port     int  `short:"p" description:"The port on which the server runs on. Default: HTTPS/8080"`
}

// ParseConfig parses CLI parameters passed
func ParseConfig(args []string) (bool, int, error) {
	_, err := flags.ParseArgs(&opts, args)
	if err != nil {
		return false, -1, err
	}
	port := globals.DefaultRpcPort
	if opts.Port != 0 {
		port = opts.Port
	}
	return opts.Insecure, port, nil
}

func main() {
	// oracle.Schedule()
	// blockchain.CheckTokenBalance()
	// blockchain.CommitToChain(big.NewInt(1565752648), "0x4920636172652061626f757420636c696d617465")
	database.FlushDB()
	database.CreateHomeDir()
	database.Populate()

	insecure, port, err := ParseConfig(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("flushed and created new db")
	server.StartServer(port, insecure)
}
