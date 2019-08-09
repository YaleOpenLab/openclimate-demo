package server

import (
	"net/http"
	"log"
	erpc "github.com/Varunram/essentials/rpc"
	db "github.com/YaleOpenLab/openclimate/database"
)

func setupView() {
	viewPledges()
	viewCompanyNational()
	viewCompanySubNationalByNational()
	viewCompanyAssetsBySubNational()
}


func viewPledges() {
	http.HandleFunc("/user/pledges/view", func(w http.ResponseWriter, r *http.Request) {
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


func viewCompanyNational() {
	http.HandleFunc("/company/national", func(w http.ResponseWriter, r *http.Request) {
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

		erpc.MarshalSend(w, countries)
	})
}


func viewCompanySubNationalByNational() {
	http.HandleFunc("/company/subnational", func(w http.ResponseWriter, r *http.Request) {
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

		erpc.MarshalSend(w, statesByCountry)
	})
}


func viewCompanyAssetsBySubNational() {
	http.HandleFunc("/company/assets/filter", func(w http.ResponseWriter, r *http.Request) {

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
