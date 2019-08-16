package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	// "strconv"
	// "github.com/pkg/errors"

	erpc "github.com/Varunram/essentials/rpc"
	// db "github.com/YaleOpenLab/openclimate/database"
	// "github.com/YaleOpenLab/openclimate/ipfs"
	"github.com/YaleOpenLab/openclimate/oracle"
)

func setupReport() {
	reportDirect()
	// connectDatabase()
}

/*
	Handler that allows actors to self-report their climate action data.
	The data in the body of the POST request must follow the format of
	either the Emissions, Mitigation, or Adaptation structs defined in
	ipfs/data.go.
*/
func reportDirect() {
	http.HandleFunc("/report/direct", func(w http.ResponseWriter, r *http.Request) {
		user, err := CheckPostAuth(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		// entity, err := user.RetrieveUserEntity()
		// if err != nil {
		// 	log.Println(err)
		// 	erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		// 	return
		// }

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

		var data map[string][]float64
		err = json.Unmarshal(bytes, &data)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		err = oracle.Verify(reportType, user.EntityType, user.EntityID, data)
		if err != nil {
			log.Fatal(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}
		// commit to blockchain

		// erpc.MarshalSend(w, ipfsHash)
	})
}

type ReportIpcc struct {
}

// Report data using the IPCC methodology.
func reportIpcc(data interface{}) (ReportIpcc, error) {
	var empty ReportIpcc
	return empty, nil
}
