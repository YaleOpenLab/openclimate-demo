package main

import (
	"log"
	"net/http"
)

func setupRelayHandlers() {
	setupRelayPingHandler()
}

// setupPingHandler is a ping route for remote callers to check if the platform is up
func setupRelayPingHandler() {
	http.HandleFunc("/relay/ping", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)
		responseHandler(w, StatusOK)
	})
}

func StartServer(port string, insecure bool) {
	setupRelayHandlers()
	log.Println("Starting Relay Server on Port: ", port)
	if insecure {
		log.Fatal(http.ListenAndServe(":"+port, nil))
	} else {
		log.Fatal(http.ListenAndServeTLS(":"+port, "server.crt", "server.key", nil))
	}
}

func main() {
	StartServer("8001", true)
}
