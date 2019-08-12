package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	// "strconv"
	// "github.com/pkg/errors"

	erpc "github.com/Varunram/essentials/rpc"
	db "github.com/YaleOpenLab/openclimate/database"
	// "github.com/YaleOpenLab/openclimate/ipfs"
	"github.com/YaleOpenLab/openclimate/oracle"
)

func setupReport() {
	report()
	connectDatabase()
}


/*
	Handler that allows actors to self-report their climate action data.
	The data in the body of the POST request must follow the format of
	either the Emissions, Mitigation, or Adaptation structs defined in
	ipfs/data.go.
*/
func report() {
	http.HandleFunc("/user/self-report", func(w http.ResponseWriter, r *http.Request) {
		user, err := CheckPostAuth(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		entity, err := user.RetrieveUserEntity()
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

type ReportIpcc struct {

}

// Report data using the IPCC methodology.
func reportIpcc(data interface{}) (ReportIpcc, error) {
	var empty ReportIpcc
	return empty, nil
}


// Submit a request to connect with an external database that contains
// emissions/mitigation/adaptation data that users would like to report.
func connectDatabase() {
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
	})
}
