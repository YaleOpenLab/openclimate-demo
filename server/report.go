package server

import (
	"io/ioutil"
	"log"
	"net/http"
	"encoding/json"

	erpc "github.com/Varunram/essentials/rpc"
	"github.com/YaleOpenLab/openclimate/oracle"
)

func setupReportHandlers() {
	SelfReportData()
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
		}

		reportType := r.URL.Query()["report_type"][0]
		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		data := make(map[string]string)
		err = json.Unmarshal(bytes, &data)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		ipfsHash, err := oracle.Verify(data, reportType)
		erpc.MarshalSend(w, ipfsHash)
	})
}

func ConnectDatabase() {
	http.HandleFunc("/user/existing-database", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckPost(w, r)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

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



