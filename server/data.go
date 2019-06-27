package server

import (
	"github.com/YaleOpenLab/openclimate/database"
	"log"
	"net/http"
)

func dataHandler() {
	getUSStates()
	getUSCounties()
}

type USStatesReturn struct {
	States []string
}

func getUSStates() {
	http.HandleFunc("/us/states", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)

		_, err := authorizeUser(r)
		if err != nil {
			log.Println("could not retrieve user from the database, quitting")
			responseHandler(w, StatusBadRequest)
			return
		}

		database.InitUSStates()
		var x USStatesReturn
		x.States = database.USStates
		MarshalSend(w, x)
	})
}

type USStateCountiesReturn struct {
	Counties map[string][]string
}

func getUSCounties() {
	http.HandleFunc("/us/counties", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)

		_, err := authorizeUser(r)
		if err != nil {
			log.Println("could not retrieve user from the database, quitting")
			responseHandler(w, StatusBadRequest)
			return
		}

		database.InitUSStates()
		var x USStateCountiesReturn
		x.Counties = database.USStateCities
		MarshalSend(w, x)
	})
}
