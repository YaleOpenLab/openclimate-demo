package server

import (
	"encoding/json"
	"log"
	"net/http"
)

// WriteToHandler constructs a reply to the passed writer
func WriteToHandler(w http.ResponseWriter, jsonString []byte) {
	w.Header().Add("Access-Control-Allow-Headers", "Accept, Authorization, Cache-Control, Content-Type")
	w.Header().Add("Access-Control-Allow-Methods", "*")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonString)
}

// MarshalSend marshals and writes a json string into the writer
func MarshalSend(w http.ResponseWriter, x interface{}) {
	xJson, err := json.Marshal(x)
	if err != nil {
		log.Println("did not marshal json", err)
		errString := "Internal Server Error"
		WriteToHandler(w, []byte(errString))
		return
	}
	log.Println("JSON: ", string(xJson))
	WriteToHandler(w, xJson)
}

// checkOrigin checks the origin of the incoming request
func checkOrigin(w http.ResponseWriter, r *http.Request) {
	// re-enable this function for all private routes
	// if r.Header.Get("Origin") != "localhost" { // allow only our frontend UI to connect to our RPC instance
	// 	http.Error(w, "404 page not found", http.StatusNotFound)
	// }
}

// checkGet checks if the invoming request is a GET request
func checkGet(w http.ResponseWriter, r *http.Request) {
	checkOrigin(w, r)
	if r.Method != "GET" {
		responseHandler(w, StatusNotFound)
		return
	}
}

func setupBasicHandlers() {
	setupPingHandler()
}

// setupPingHandler is a ping route for remote callers to check if the platform is up
func setupPingHandler() {
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)
		responseHandler(w, StatusOK)
	})
}

func StartServer(port string, insecure bool) {
	setupBasicHandlers()
	setupDBHandlers()
	log.Println("Starting RPC Server on Port: ", port)
	if insecure {
		log.Fatal(http.ListenAndServe(":"+port, nil))
	} else {
		log.Fatal(http.ListenAndServeTLS(":"+port, "server.crt", "server.key", nil))
	}
}
