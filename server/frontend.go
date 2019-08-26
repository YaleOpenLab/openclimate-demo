package server

import (
	"log"
	// "encoding/json"
	// "io/ioutil"
	erpc "github.com/Varunram/essentials/rpc"
	"github.com/pkg/errors"
	"net/http"
	"strings"
	"strconv"
	"encoding/json"
	"io/ioutil"
	"github.com/YaleOpenLab/openclimate/database"
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
	http.HandleFunc("/multinationals/", func(w http.ResponseWriter, r *http.Request) {
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
		results["nation_state"] = nationState
		results["pledges"] = pledges

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
	Name string
	Pledges []database.Pledge
	Subnational []Subnational
}

type Subnational struct {
	Name string
	Pledges []database.Pledge
	Assets []database.Asset
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
				log.Fatal(err)
			}
			erpc.MarshalSend(w, nationStates)

		case "review":
			results := make(map[string] interface{})
			results["certificates"] = company.Certificates
			results["climate_reports"] = company.ClimateReports
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

func getEarthStatus() {
	http.HandleFunc("/earth-status", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
		}

		w.Write([]byte("earth status"))
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
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
		}

		bytes, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			log.Fatal(err)
		}

		var registerInfo map[string]interface{}
		err = json.Unmarshal(bytes, &registerInfo)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			log.Fatal(err)
		}

		log.Println(registerInfo)
	})
}

func postLogin() {
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckPost(w, r)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			log.Fatal(err)
		}

		bytes, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			log.Fatal(err)
		}

		var credentials map[string]string
		err = json.Unmarshal(bytes, &credentials)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			log.Fatal(err)
		}

		username := credentials["username"]
		pwhash := credentials["pwhash"]

		_, err = database.ValidateUser(username, pwhash)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			log.Fatal(err)
		}

		accessToken := "placeholder"
		erpc.MarshalSend(w, accessToken)
	})
}

