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
	SelfReportEmissionsData()
}

/*

	RPC handlers to allow users of the platform to report data.
	----------------------------------------------------------------

	* Self-report: emissions, mitigation, adaptation, pledge data
	* Set up automatic reporting (connect to database, etc.)

	* TODO: Add logic (probably a function) that checks if
		the methodology used is acceptable and thus verified.

	* TODO: hook up IoT devices?

*/

/**************************/
/* PLEDGE DATA STRUCTS */
/**************************/

type Pledges struct {

	// Meta-data
	UserID 			int
	EntityType 		string

	// Info on specific pledges
	Pledges			[]PledgeData
}

type PledgeData struct {

	// * emissions reductions
	// * mitigation actions (energy efficiency, renewables, etc.)
	// * adaptation actions
	PledgeType 		string
	BaseYear		int
	TargetYear 		int
	Goal 			int
	// is this goal determined by a regulator, or voluntarily
	// adopted by the climate actor?
	Regulatory 		bool 

}

/**************************/
/* EMISSIONS DATA STRUCTS */
/**************************/

type Emissions struct {

	// Meta-data
	UserID 			int
	EntityType		string
	Year    		int

	// Emissions data (by asset)
	// Country children: regions
	// Region children: companies & cities
	// Company children: assets
	ByChild 		[]ChildEmissionsData
}

type ChildEmissionsData struct {

	ChildID     	int
	ChildName 		string
	ScopeICO2e 		float64
	ScopeIICO2e 	float64
	ScopeIIICO2e 	float64

	// Where is the report and its data from?
	// (options: internally conducted report, consulting group, etc.)
	Source string

	// what methodology was used in the reporting and
	// verification of the emissions data?
	Methodology string

	// // "verified" represents if the data is sufficiently reviewed
	// // and confirmed/corroborated (from oracle, third-party auditor, etc)
	// Verified string
}

/******************************/
/* OTHER ACTIONS DATA STRUCTS */
/******************************/

type Mitigation struct {

}

type Adaptation struct {

}


func SelfReportEmissionsData() {
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
		var data interface{}

		// switch reportType := r.URL.Query()["Type"]; reportType {
		// case "Emissions":	
		// 	var data Emissions
		// 	err = json.Unmarshal(bytes, &data)
		// 	if err != nil {
		// 		erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		// 		return
		// 	}
		// case "Pledges":
		// 	var data Pledges
		// 	err = json.Unmarshal(bytes, &data)
		// 	if err != nil {
		// 		erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		// 		return
		// 	}
		// case "Mitigation":
		// 	var data Mitigation
		// 	err = json.Unmarshal(bytes, &data)
		// 	if err != nil {
		// 		erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		// 		return
		// 	}
		// case "Adaptation":
		// 	var data Adaptation
		// 	err = json.Unmarshal(bytes, &data)
		// 	if err != nil {
		// 		erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		// 		return
		// 	}
		// }

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


func SelfReportPledge() {
}
