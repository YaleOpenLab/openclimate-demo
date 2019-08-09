package server

import (
	"net/http"
	"log"
	erpc "github.com/Varunram/essentials/rpc"

)

func setupView() {
	ViewPledges()
}


func ViewPledges() {
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
		log.Println(actor)

		pledges, err := actor.GetPledges()
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, pledges)
	})
}


func getCompanyStates() {
	http.HandleFunc("/company/states", func(w http.ResponseWriter, r *http.Request) {
		user, err := CheckGetAuth(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		log.Println(user)

		// user.RetrieveUserEntity()

		// if r.URL.Query()["company_name"] == nil || r.URL.Query()["company_country"] == nil {
		// 	log.Println(err)
		// 	erpc.ResponseHandler(w, erpc.StatusBadRequest)
		// 	return
		// }

		// // Given its name and country, retrieve the company from the database

		// name := r.URL.Query()["company_name"][0]
		// country := r.URL.Query()["company_country"][0]
		// company, err := database.RetrieveCompanyByName(name, country)
		// if err != nil {
		// 	log.Println(err)
		// 	erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		// 	return
		// }

		// Get the information of the states that the company is in

		// states, err := company.GetStates()
		// if err != nil {
		// 	log.Println(err)
		// 	erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		// 	return
		// }

		// erpc.MarshalSend(w, states)
	})
}