package server

import (
	"log"
	"net/http"

	erpc "github.com/Varunram/essentials/rpc"
	utils "github.com/Varunram/essentials/utils"
)

func checkReqdParams(w http.ResponseWriter, r *http.Request, options ...string) bool {
	for _, option := range options {
		if r.URL.Query()[option] == nil {
			log.Println("reqd param: ", option, "not found")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return false
		}
	}
	return true
}

func checkReqdPostParams(w http.ResponseWriter, r *http.Request, options ...string) bool {
	for _, option := range options {
		if r.FormValue(option) == "" {
			log.Println("reqd param: ", option, "not found")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return false
		}
	}
	return true
}

func StartServer(portx int, insecure bool) {

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

	port, err := utils.ToString(portx)
	if err != nil {
		log.Fatal("Port not string")
	}

	log.Println("Starting RPC Server on Port: ", port)
	if insecure {
		log.Fatal(http.ListenAndServe(":"+port, nil))
	} else {
		log.Fatal(http.ListenAndServeTLS(":"+port, "certs/server.crt", "certs/server.key", nil))
	}
}
