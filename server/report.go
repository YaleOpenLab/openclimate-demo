package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
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
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		// Read the data in the body of the response into bytes,
		// which will be parsed into a map[string]string.
		bytes, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		// Unmarshal the json bytes into the map[string]string data type
		// to be parsed for the arguments to create a new Pledge item and
		// inserted into the pledge bucket.
		var pledge map[string]string
		err = json.Unmarshal(bytes, &pledge)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		// Convert string data into the correct data type to be passed
		// to the NewPledge() function, in order to create the new pledge.

		actorID := user.EntityID

		pledgeType := pledge["pledge_type"]

		baseYear, err := strconv.Atoi(pledge["base_year"])
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		targetYear, err := strconv.Atoi(pledge["target_year"])
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		goal, err := strconv.ParseFloat(pledge["goal"], 64)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		regulatory, err := strconv.ParseBool(pledge["regulatory"])
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		// Call NewPledge() with all the arguments, which have been typecasted
		// into the proper types required by the NewPledge function.
		new, err := db.NewPledge(pledgeType, baseYear, targetYear, goal, regulatory, actorID)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		// TODO: Convert pledge into smart contract condition
		erpc.MarshalSend(w, new)
	})
}

func UpdatePledge() {
	http.HandleFunc("/user/pledges/update", func(w http.ResponseWriter, r *http.Request) {

		user, err := CheckPostAdmin(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		pledgeID, err := strconv.Atoi(r.URL.Query()["pledge_id"][0])
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		bytes, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		var pledge db.Pledge
		err = json.Unmarshal(bytes, &pledge)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		pledge.ID = pledgeID
		pledge.ActorID = user.EntityID

		err = db.UpdatePledge(pledgeID, pledge)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		erpc.MarshalSend(w, pledge)
	})
}


func CommitPledge() {
	http.HandleFunc("user/pledges/commit", func (w http.ResponseWriter, r *http.Request) {
		_, err := CheckGetAuth(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		if r.URL.Query()["pledge_ID"] == nil {
			log.Println("pledge ID not passed, quitting")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		pledgeID, err := strconv.Atoi(r.URL.Query()["pledge_ID"][0])
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		pledge, err := db.RetrievePledge(pledgeID)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		ipfsHash, err := oracle.IpfsCommitData(pledge)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, ipfsHash)
	})
}


func SelfReportData() {
	http.HandleFunc("/user/self-report", func(w http.ResponseWriter, r *http.Request) {
		user, err := CheckPostAuth(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		entity, err := user.GetUserActor()
		if err != nil {
			log.Println(err)
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

		ipfsHash, err := oracle.Verify(reportType, entity, data)

		// commit to blockchain

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
