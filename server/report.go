package server

import (
	"encoding/json"
	ipfs "github.com/Varunram/essentials/ipfs"
	erpc "github.com/Varunram/essentials/rpc"
	"io/ioutil"
	//"log"
	"net/http"
)


func setupReportHandlers() {
	SelfReportData()
}

/*

	RPC handlers to allow users of the platform to self-report data.
	----------------------------------------------------------------
	*** Currently working on support for companies to report data. ***

	* TODO: support for regions
	* TODO: support for countries
	* TODO: support for cities

	* TODO: Add logic (probably a function) that checks if
		the methodology used is acceptable and thus verified.

*/

type CompanyData struct {

	// Meta-data
	UserID 		int
	EntityType	string
	Year    	int

	// Emissions data (by asset)
	Assets []AssetData
}

type AssetData struct {
	AssetID      int
	AssetName    string
	ScopeICO2e   float64
	ScopeIICO2e  float64
	ScopeIIICO2e float64

	// Where is the report and its data from?
	// (options: internally conducted report, consulting group, etc.)
	Source string

	// what methodology was used in the reporting and
	// verification of the emissions data?
	Methodology string

	// "verified" represents if the data is sufficiently reviewed
	// and confirmed/corroborated (from oracle, third-party auditor, etc)
	Verified string
}

func SelfReportData() {
	http.HandleFunc("/user/self-report", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckPost(w, r)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		bytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		// Check if the bytes is in a valid JSON format
		var data CompanyData
		err = json.Unmarshal(bytes, &data)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		// add data (in byte format) to IPFS
		hash, err := ipfs.IpfsAddBytes(bytes)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		// *** COMMIT HASH TO A BLOCKCHAIN, HASH LOOKUP USING SMART CONTRACT *** 

		erpc.MarshalSend(w, hash)
	})
}
