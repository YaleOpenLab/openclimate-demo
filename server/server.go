package server

import (
	"log"
	"net/http"

	erpc "github.com/Varunram/essentials/rpc"
	ocdb "github.com/YaleOpenLab/openclimate/database"
)

func StartServer(port string, insecure bool) {

	ocdb.Populate()
	erpc.SetupBasicHandlers()

	setupView()
	setupManage()
	setupUser()
	setupReport()

	setupActorsHandlers()
	setupIpfsHandlers()

	setupSwytchApis()
	setupDataHandlers()
	frontendFns()

	log.Println("Starting RPC Server on Port: ", port)
	if insecure {
		log.Fatal(http.ListenAndServe(":"+port, nil))
	} else {
		log.Fatal(http.ListenAndServeTLS(":"+port, "certs/server.crt", "certs/server.key", nil))
	}
}
