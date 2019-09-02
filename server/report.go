package server

import (
	// "encoding/json"
	// "io/ioutil"
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

		err = r.ParseForm()
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		if !checkReqdParams(w, r, "report_type", "data") {
			return
		}

		reportType := r.FormValue("report_type")
		data := r.FormValue("data")
		err = oracle.VerifyAndCommit(reportType, user.EntityType, user.EntityID, data)
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
