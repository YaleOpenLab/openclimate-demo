package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	// "github.com/pkg/errors"

	erpc "github.com/Varunram/essentials/rpc"
	db "github.com/YaleOpenLab/openclimate/database"
	"github.com/YaleOpenLab/openclimate/oracle"
)

func setupReportHandlers() {
	SelfReportData()
	ConnectDatabase()
	AddPledge()
}

func AddPledge() {
	http.HandleFunc("/user/pledges/add", func(w http.ResponseWriter, r *http.Request) {
		user, err := CheckPostAuth(w, r)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		var entity interface{}

		switch user.EntityType {
		case "company":
			entity, err = db.RetrieveCompany(user.EntityID)
		case "city":
			entity, err = db.RetrieveCity(user.EntityID)
		case "state":
			entity, err = db.RetrieveState(user.EntityID)
		case "region":
			entity, err = db.RetrieveRegion(user.EntityID)
		case "country":
			entity, err = db.RetrieveCountry(user.EntityID)
		default:
			log.Println("Entity type of user is not valid.")
			erpc.ResponseHandler(w, erpc.StatusUnauthorized)
		}
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		erpc.MarshalSend(w, entity)
	})
}



func SelfReportData() {
	http.HandleFunc("/user/self-report", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckPost(w, r)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		if r.URL.Query()["report_type"] == nil {
			log.Println("report type not passed, quitting")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		reportType := r.URL.Query()["report_type"][0]
		bytes, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		var data interface{}
		err = json.Unmarshal(bytes, &data)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		ipfsHash, err := oracle.Verify(data, reportType)
		erpc.MarshalSend(w, ipfsHash)
	})
}

// Submit a request to connect with an external database that contains
// emissions/mitigation/adaptation data that users would like to report.
func ConnectDatabase() {
	http.HandleFunc("/user/request-database", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckPost(w, r)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		b, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		var request db.ConnectRequest
		err = json.Unmarshal(b, &request)
		if err != nil {
			log.Println("Error: failed to unmarshal bytes into Request struct")
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		db.NewRequest(request) // store request into request bucket, to be reviewed later
		erpc.MarshalSend(w, request)

		// log.Println("BYTES: ", b)

		// entityType := r.URL.Query()["entity_type"][0]
		// username := r.URL.Query()["username"][0]
		// user, err := database.RetrieveUserbyUsername(username)
		// if err != nil {
		// 	log.Println("failed to find user")
		// 	return
		// }

		// for _, db := range r.URL.Query()["database"] {

		// }

	})
}
