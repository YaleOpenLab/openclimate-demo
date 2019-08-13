package server

import (
	"net/http"
	"log"
	erpc "github.com/Varunram/essentials/rpc"
	db "github.com/YaleOpenLab/openclimate/database"
	"github.com/YaleOpenLab/openclimate/ipfs"
)

func setupView() {
	viewCompanyPledges()
	ViewCompanyEarth()
	viewCompanyNational()
	viewCompanySubNationalByNational()
	viewCompanyAssetsBySubNational()
}

var viewUrl string = "/view"

func viewCompanyPledges() {
	http.HandleFunc(viewUrl + "/pledges", func(w http.ResponseWriter, r *http.Request) {
		user, err := CheckGetAuth(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		actor, err := user.RetrieveUserEntity()
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		pledges, err := actor.GetPledges()
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, pledges)
	})
}

func ViewCompanyEarth() {
	http.HandleFunc(viewUrl + "/earth", func(w http.ResponseWriter, r *http.Request) {

		_, err := CheckGetAuth(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		earthData, err := ipfs.GetFromIpfsEarthData()
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, earthData)
	})
}

func viewCompanyNational() {
	http.HandleFunc(viewUrl + "/national", func(w http.ResponseWriter, r *http.Request) {
		user, err := CheckGetAuth(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		company, err := db.RetrieveCompany(user.EntityID)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		countries, err := company.GetCountries()
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		countryPledges := make(map[string][]db.Pledge)
		for _, country := range countries {
			p, err := country.GetPledges()
			if err != nil {
				log.Println(err)
				erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			}
			countryPledges[country.Name] = p
		}

		final := make(map[string]interface{})
		final["countries"] = countries
		final["country_pledges"] = countryPledges

		erpc.MarshalSend(w, final)
	})
}

func viewCompanySubNationalByNational() {
	http.HandleFunc(viewUrl + "/subnational", func(w http.ResponseWriter, r *http.Request) {
		user, err := CheckGetAuth(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		company, err := db.RetrieveCompany(user.EntityID)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		states, err := company.GetStates()
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		countries, err := company.GetCountries()
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		statesByCountry := make(map[string][]db.State)

		for _, country := range countries {
			for _, state := range states {
				if country.Name == state.Country {
					statesByCountry[country.Name] = append(statesByCountry[country.Name], state)
				}
			}
		}

		statePledges := make(map[string][]db.Pledge)
		for _, state := range states {
			p, err := state.GetPledges()
			if err != nil {
				log.Println(err)
				erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			}
			statePledges[state.Name] = p
		}

		final := make(map[string]interface{})
		final["states_by_country"] = statesByCountry
		final["state_pledges"] = statePledges

		erpc.MarshalSend(w, final)
	})
}

func viewCompanyAssetsBySubNational() {
	http.HandleFunc(viewUrl + "/assets", func(w http.ResponseWriter, r *http.Request) {

		user, err := CheckGetAuth(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		if user.EntityType != "company" {
			log.Println("User entity type is not a company.")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		company, err := db.RetrieveCompany(user.EntityID)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		assetsByState := make(map[string][]db.Asset)

		for _, stateID := range company.States {
			s, err := db.RetrieveState(stateID)
			if err != nil {
				log.Println(err)
				erpc.ResponseHandler(w, erpc.StatusInternalServerError)
				return
			}
			assets, err := company.GetAssetsByState(s.Name)
			if err != nil {
				log.Println(err)
				erpc.ResponseHandler(w, erpc.StatusInternalServerError)
				return
			}
			assetsByState[s.Name] = assets
		}

		erpc.MarshalSend(w, assetsByState)
	})
}

