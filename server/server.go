package server

import (
	"log"
	"net/http"

	erpc "github.com/Varunram/essentials/rpc"
	ocdb "github.com/YaleOpenLab/openclimate/database"
)

func StartServer(port string, insecure bool) {

	ocdb.Populate()
	// log.Println("populated db")
	erpc.SetupBasicHandlers()
	setupDBHandlers()
	setupReportHandlers()
	setupSwytchApis()
	dataHandler()
	log.Println("Starting RPC Server on Port: ", port)
	if insecure {
		log.Fatal(http.ListenAndServe(":"+port, nil))
	} else {
		log.Fatal(http.ListenAndServeTLS(":"+port, "server.crt", "server.key", nil))
	}
}
