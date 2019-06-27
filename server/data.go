package server

import (
	"encoding/json"
	"github.com/YaleOpenLab/openclimate/database"
	"io/ioutil"
	"log"
	"net/http"
)

func dataHandler() {
	getUSStates()
	getUSCounties()
	getParisAgreement()
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

type ParisAgreementReturn struct {
	Code         string  `json:"Code"`
	Name         string  `json:"Name"`
	Signature    string  `json:"Signature"`
	Ratification string  `json:"Ratification"`
	Kind         string  `json:"Kind"`
	DateOfEffect string  `json:"Date-Of-Effect"`
	Emissions    float64 `json:"Emissions"`
	Percentage   float64 `json:"Percentage"`
	Year         float64 `json:"Year"`
}

type ParisAgreementReturnFinal struct {
	Name         string  `json:"Name"`
	Signature    string  `json:"Signature"`
	Ratification string  `json:"Ratification"`
	Kind         string  `json:"Kind"`
	DateOfEffect string  `json:"Date-Of-Effect"`
	Emissions    float64 `json:"Emissions"`
	Percentage   float64 `json:"Percentage"`
	Year         float64 `json:"Year"`
}

func getParisAgreement() {
	http.HandleFunc("/paris/data", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)

		_, err := authorizeUser(r)
		if err != nil {
			log.Println("could not retrieve user from the database, quitting")
			responseHandler(w, StatusBadRequest)
			return
		}

		data, err := ioutil.ReadFile("load-data/json_data/paris_agreement_entry_into_force.json")
		if err != nil {
			log.Println(err)
			responseHandler(w, StatusInternalServerError)
			return
		}

		var x map[string]ParisAgreementReturn
		err = json.Unmarshal(data, &x)
		if err != nil {
			log.Println(err)
			responseHandler(w, StatusInternalServerError)
			return
		}

		y := make(map[string]ParisAgreementReturnFinal)
		for _, values := range x {
			var temp ParisAgreementReturnFinal
			temp.Name = values.Name
			temp.Signature = values.Signature
			temp.Ratification = values.Ratification
			temp.Kind = values.Kind
			temp.DateOfEffect = values.DateOfEffect
			temp.Emissions = values.Emissions
			temp.Percentage = values.Percentage
			temp.Year = values.Year
			y[values.Code] = temp
		}
		MarshalSend(w, y)
	})
}
