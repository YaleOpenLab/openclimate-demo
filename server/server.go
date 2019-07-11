package server

import (
	"log"
	"net/http"

	erpc "github.com/Varunram/essentials/rpc"
	"github.com/YaleOpenLab/openclimate/database"
)

func StartServer(port string, insecure bool) {

	database.Populate()
	erpc.SetupBasicHandlers()
	setupDBHandlers()
	setupSwytchApis()
	dataHandler()
	log.Println("Starting RPC Server on Port: ", port)
	if insecure {
		log.Fatal(http.ListenAndServe(":"+port, nil))
	} else {
		log.Fatal(http.ListenAndServeTLS(":"+port, "server.crt", "server.key", nil))
	}
}
