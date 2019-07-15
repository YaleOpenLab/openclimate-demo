package server

import (

	"net/http"

	// ipfs "github.com/Varunram/essentials/ipfs"
	// erpc "github.com/Varunram/essentials/rpc"

)

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
	UserIdx			int
	Year			int

	// Emissions data (by asset)
	Assets 			[]AssetData
}

type AssetData struct {

	ScopeICO2e		float64
	ScopeIICO2e		float64
	ScopeIIICO2e	float64

	// Where is the report and its data from? 
	// (options: internally conducted report, consulting group, etc.)
	Source			string

	// what methodology was used in the reporting and
	// verification of the emissions data?
	Methodology		string

	// "verified" represents if the data is sufficiently reviewed
	// and confirmed/corroborated (from oracle, third-party auditor, etc)
	Verified		bool 

}


func SelfReportData() {
	http.HandleFunc("/user/self-report", func(w http.ResponseWriter, r *http.Request) {
		user, err := CheckGetAuth(w, r)
		if err != nil {
			return
		}

		year := r.URL.Query()["year"][0]

		repAssetData := r.URL.Query()["asset_data"][0]
		
		var newData AssetData
		for _, asset := range repAssetData {

			newData.ScopeICO2e = asset["ScopeICO2e"]
			newData.ScopeIICO2e = asset["ScopeIICO2e"]
			newData.ScopeIIICO2e = asset["ScopeIIICO2e"]

			newData.Source = asset["Source"]
			newData.Methodology = asset["Methodology"]
			newData.Verified = false // all self-reported data is not verified yet

		}

		// Commit the data to ipfs

	})
}

