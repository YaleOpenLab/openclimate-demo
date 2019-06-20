package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
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

func RunPythonScript() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	cmd := exec.Command(dir + "/csv.py")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func main() {
	var err error

	err = RunPythonScript()
	if err != nil {
		log.Fatal(err)
	}

	StartServer("8001", true)
}
