package server

import (
	"log"
	// "encoding/json"
	// "io/ioutil"
	erpc "github.com/Varunram/essentials/rpc"
	"github.com/YaleOpenLab/openclimate/blockchain"
	"github.com/YaleOpenLab/openclimate/database"
	"github.com/pkg/errors"
	"net/http"
	"strconv"
	"strings"
)

func frontendFns() {
	getNationStates()
	getMultiNationals()
	getNationStateId()
	getMultiNationalId()
	getActorId()
	getEarthStatus()
	getActors()
	postFiles()
	postRegister()
	postLogin()
}

func getId(w http.ResponseWriter, r *http.Request) (string, error) {
	var id string
	err := erpc.CheckGet(w, r)
	if err != nil {
		log.Println(err)
		return id, errors.New("request not get")
	}

	urlParams := strings.Split(r.URL.String(), "/")

	if len(urlParams) < 3 {
		return id, errors.New("no id provided, quitting")
	}

	id = urlParams[2]
	return id, nil
}

func getNationStates() {
	http.HandleFunc("/nation-states", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
		}

		nationStates, err := database.RetrieveAllCountries()
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		erpc.MarshalSend(w, nationStates)
	})
}

func getMultiNationals() {
	http.HandleFunc("/multinationals", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
		}

		multinationals, err := database.RetrieveAllMultiNationals()
		erpc.MarshalSend(w, multinationals)
	})
}

func getNationStateId() {
	http.HandleFunc("/nation-states/", func(w http.ResponseWriter, r *http.Request) {
		strID, err := getId(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(strID)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		nationState, err := database.RetrieveCountry(id)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		pledges, err := nationState.GetPledges()
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		results := make(map[string]interface{})
		results["name"] = nationState.Name
		results["full_name"] = nationState.Name
		results["description"] = nationState.Description
		results["pledges"] = pledges
		results["accountability"] = nationState.Accountability

		erpc.MarshalSend(w, results)
	})
}

func getMultiNationalId() {
	http.HandleFunc("/multinationals/", func(w http.ResponseWriter, r *http.Request) {
		strID, err := getId(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(strID)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		multinational, err := database.RetrieveCompany(id)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		pledges, err := multinational.GetPledges()
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		results := make(map[string]interface{})
		results["name"] = multinational.Name
		results["full_name"] = multinational.Name
		results["description"] = multinational.Description
		results["pledges"] = pledges
		results["accountability"] = multinational.Accountability
		results["locations"] = multinational.Locations

		erpc.MarshalSend(w, results)
	})
}

type NationState struct {
	Name        string
	Pledges     []database.Pledge
	Subnational []Subnational
}

type Subnational struct {
	Name    string
	Pledges []database.Pledge
	Assets  []database.Asset
}

func getActorId() {
	http.HandleFunc("/actors/", func(w http.ResponseWriter, r *http.Request) {
		strID, err := getId(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(strID)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		company, err := database.RetrieveCompany(id)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		urlParams := strings.Split(r.URL.String(), "/")
		if len(urlParams) < 4 {
			log.Println("insufficient amount of params")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		choice := urlParams[3]

		switch choice {
		case "dashboard":
			pledges, err := company.GetPledges()
			if err != nil {
				log.Println(err)
				erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			}
			results := make(map[string]interface{})
			results["full_name"] = company.Name
			results["description"] = company.Description
			results["locations"] = company.Locations
			results["accountability"] = company.Accountability
			results["pledges"] = pledges

		case "nation-states":
			nationStates, err := getActorIdNationStates(company, w, r)
			if err != nil {
				erpc.ResponseHandler(w, erpc.StatusInternalServerError)
				return
			}
			erpc.MarshalSend(w, nationStates)

		case "review":
			results := make(map[string]interface{})
			results["certificates"] = company.Certificates
			results["climate_reports"] = company.ClimateReports

			var err error
			results["emissions"], err = blockchain.RetrieveActorEmissions(id)
			if err != nil {
				erpc.MarshalSend(w, erpc.StatusInternalServerError)
				return
			}
			results["reductions"], err = blockchain.RetrieveActorEmissions(id)
			if err != nil {
				erpc.MarshalSend(w, erpc.StatusInternalServerError)
				return
			}
			erpc.MarshalSend(w, results)

		// case "manage":
		// 	w.Write([]byte("manage: " + strconv.Itoa(id)))

		case "climate-action-asset":
			if len(urlParams) < 5 {
				log.Println("insufficient amount of params")
				erpc.ResponseHandler(w, erpc.StatusBadRequest)
				return
			}
			id2 := urlParams[4]
			w.Write([]byte("climate-action-assets ids: " + strconv.Itoa(id) + id2))
		default:
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}
	})
}

func getActorIdNationStates(company database.Company, w http.ResponseWriter, r *http.Request) ([]NationState, error) {

	var nationStates []NationState

	countries, err := company.GetCountries()
	if err != nil {
		return nationStates, errors.Wrap(err, "getActorIdNationStates() failed")
	}

	for _, country := range countries {
		var nationState NationState
		states, err := company.GetStates()
		if err != nil {
			return nationStates, errors.Wrap(err, "getActorIdNationStates() failed")
		}

		pledges, err := country.GetPledges()
		if err != nil {
			return nationStates, errors.Wrap(err, "getActorIdNationStates() failed")
		}

		var subnationals []Subnational

		for _, s := range states {
			var subnational Subnational
			pledges, err := s.GetPledges()
			if err != nil {
				return nationStates, errors.Wrap(err, "getActorIdNationStates() failed")
			}
			assets, err := company.GetAssetsByState(s.Name)
			if err != nil {
				return nationStates, errors.Wrap(err, "getActorIdNationStates() failed")
			}

			subnational.Name = s.Name
			subnational.Pledges = pledges
			subnational.Assets = assets
			subnationals = append(subnationals, subnational)
		}

		nationState.Name = country.Name
		nationState.Pledges = pledges
		nationState.Subnational = subnationals
		nationStates = append(nationStates, nationState)
	}

	return nationStates, nil
}

type EarthStatusReturn struct {
	Warminginc               string `json:"warming_in_c"`
	Gtco2left                string `json:"gt_co2_left"`
	Atmosphericco2ppm        string `json:"atmospheric_co2_ppm"`
	Annualglobalemission     string `json:"annual_global_emission"`
	Estimatedbudgetdepletion string `json:"estimated_budget_depletion"`
}

func getEarthStatus() {
	http.HandleFunc("/earth-status", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
		}

		var x EarthStatusReturn
		x.Warminginc = "sample"
		x.Gtco2left = "sample"
		x.Atmosphericco2ppm = "sample"
		x.Annualglobalemission = "sample"
		x.Estimatedbudgetdepletion = "sample"

		erpc.MarshalSend(w, x)
	})
}

func getActors() {
	http.HandleFunc("/actors", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
		}

		w.Write([]byte("get actors"))
	})
}

func postFiles() {
	http.HandleFunc("/files", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckPost(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
		}

		w.Write([]byte("post files"))
	})
}

func postRegister() {
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckPost(w, r)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			log.Fatal(err)
		}

		err = r.ParseForm()
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		actor_id := r.FormValue("actor_id")
		actor_name := r.FormValue("actor_name")
		identification_file_id := r.FormValue("identification_file_id")
		employment_file_id := r.FormValue("employment_file_id")
		first_name := r.FormValue("first_name")
		last_name := r.FormValue("last_name")
		title := r.FormValue("title")
		email := r.FormValue("email")
		phone := r.FormValue("phone")
		account_type_id := r.FormValue("account_type_id")
		account_type := r.FormValue("account_type")

		switch account_type {
		case "country":
			log.Println("creating country")
		case "state":
			log.Println("creating state")
		case "region":
			log.Println("creating region")
		}

		// actorID := registerInfo["actor_id"].(int)
		// actorType := registerInfo["actor_type"].(string)

		// actor, err := RetrieveActor(actorType, actorID)
		// if err != nil {
		// 	erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		// 	log.Fatal(err)
		// }

		// // if RetrieveActor() returns nil for actor, that means the actor was not found
		// if actor == nil {
			
		// }

		// log.Println(registerInfo)

		log.Println(actor_id, actor_name, identification_file_id, employment_file_id,
			first_name, last_name, title, email, phone, account_type_id)
		w.Write([]byte("registered"))
	})
}

func postLogin() {
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckPost(w, r)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		err = r.ParseForm()
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		username := r.FormValue("username")
		pwhash := r.FormValue("pwhash")

		user, err := database.ValidateUser(username, pwhash)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		token, err := user.GenAccessToken()
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, token)
	})
}
