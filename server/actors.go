package server

import (
	// "github.com/pkg/errors"
	"log"
	"net/http"

	// ipfs "github.com/Varunram/essentials/ipfs"
	erpc "github.com/Varunram/essentials/rpc"
	"github.com/YaleOpenLab/openclimate/database"
)

func setupActorsHandlers() {
	getAllCompanies()
	getCompany()
	getAllRegions()
	getRegion()
	getAllStates()
	getState()
	getStatesByCountry()
	getAllCities()
	getCity()
	getAllCountries()
	getCountry()
}

/*******************/
/* REGION HANDLERS */
/*******************/

func getAllRegions() {
	http.HandleFunc("/region/all", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		regions, err := database.RetrieveAllRegions()
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, regions)
	})
}

func getRegion() {
	http.HandleFunc("/region", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		if !checkReqdParams(w, r, "region_name", "region_country") {
			return
		}

		name := r.URL.Query()["region_name"][0]
		country := r.URL.Query()["region_country"][0]
		region, err := database.RetrieveRegionByName(name, country) //************ STOP ***********
		if err != nil {
			log.Println("Error while retrieving region, quitting")
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, region)
	})
}

/******************/
/* STATE HANDLERS */
/******************/

func getAllStates() {
	http.HandleFunc("/state/all", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		states, err := database.RetrieveAllStates()
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, states)
	})
}

func getStatesByCountry() {
	http.HandleFunc("/state/filter", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		if !checkReqdParams(w, r, "country") {
			return
		}

		states, err := database.FilterStatesByCountry(r.URL.Query()["country"][0])
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, states)
	})
}

func getState() {
	http.HandleFunc("/state", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		if !checkReqdParams(w, r, "state_name", "state_country") {
			return
		}

		name := r.URL.Query()["state_name"][0]
		country := r.URL.Query()["state_country"][0]
		state, err := database.RetrieveStateByName(name, country) //************ STOP ***********
		if err != nil {
			log.Println("Error while retrieving state, quitting")
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, state)
	})
}

/*****************/
/* CITY HANDLERS */
/*****************/

func getAllCities() {
	http.HandleFunc("/city/all", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		cities, err := database.RetrieveAllCities()
		if err != nil {
			log.Println("Error while retrieving all cities, quitting")
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, cities)
	})
}

func getCity() {
	http.HandleFunc("/city", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		if !checkReqdParams(w, r, "city_name", "city_region") {
			return
		}

		name := r.URL.Query()["city_name"][0]
		region := r.URL.Query()["city_region"][0]
		city, err := database.RetrieveCityByName(name, region)
		if err != nil {
			log.Println("Error while retrieving all cities, quitting")
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, city)
	})
}

/********************/
/* COMPANY HANDLERS */
/********************/

func getAllCompanies() {
	http.HandleFunc("/company/all", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		companies, err := database.RetrieveAllCompanies()
		if err != nil {
			log.Println("error while retrieving all companies, quitting")
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, companies)
	})
}

func getCompany() {
	http.HandleFunc("/company", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		if !checkReqdParams(w, r, "company_name", "company_country") {
			return
		}

		name := r.URL.Query()["company_name"][0]
		country := r.URL.Query()["company_country"][0]
		company, err := database.RetrieveCompanyByName(name, country)
		if err != nil {
			log.Println("error while retrieving all companies, quitting")
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, company)
	})
}

/**********************/
/* COUNTRIES HANDLERS */
/**********************/

func getAllCountries() {
	http.HandleFunc("/country/all", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		countries, err := database.RetrieveAllCountries()
		if err != nil {
			log.Println("error while retrieving all countries, quitting")
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, countries)
	})
}

func getCountry() {
	http.HandleFunc("/country", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		if !checkReqdParams(w, r, "country_name") {
			return
		}

		name := r.URL.Query()["country_name"][0]
		country, err := database.RetrieveCountryByName(name)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, country)
	})
}
