package server

import (
	"encoding/json"
	// ipfs "github.com/Varunram/essentials/ipfs"
	erpc "github.com/Varunram/essentials/rpc"

	"github.com/YaleOpenLab/openclimate/oracle"

	"io/ioutil"
	//"log"
	"net/http"
)

func setupReportHandlers() {
	SelfReportData()
}

/**************************/
/* PLEDGE DATA STRUCTS */
/**************************/

type Pledges struct {

	// Meta-data
	UserID     int
	EntityType string

	// Info on specific pledges
	Pledges []PledgeData
}

type PledgeData struct {
	// * emissions reductions
	// * mitigation actions (energy efficiency, renewables, etc.)
	// * adaptation actions
	PledgeType string
	BaseYear   int
	TargetYear int
	Goal       int
	// is this goal determined by a regulator, or voluntarily
	// adopted by the climate actor?
	Regulatory bool
}

/**************************/
/* EMISSIONS DATA STRUCTS */
/**************************/

type Emissions struct {
	// Meta-data
	UserID     int
	EntityType string
	Year       int
	// Emissions data (by asset)
	// Country children: regions
	// Region children: companies & cities
	// Company children: assets
	ByChild []ChildEmissionsData
}

type ChildEmissionsData struct {
	ChildID      int
	ChildName    string
	ScopeICO2e   float64
	ScopeIICO2e  float64
	ScopeIIICO2e float64

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

/***************************/
/* MITIGATION DATA STRUCTS */
/***************************/

type Mitigation struct {
	// Meta-data
	UserID     int
	EntityType string
	Year       int

	// Emissions data (by asset)
	// Country children: regions
	// Region children: companies & cities
	// Company children: assets
	ByChild []ChildMitigationData
}

type ChildMitigationData struct {
	ChildID      int
	ChildName    string
	CarbonOffset float64
	EnergySaved  float64
	EnergyGen    float64

	// Where is the report and its data from?
	// (options: internally conducted report, consulting group, etc.)
	Source string

	// what methodology was used in the reporting and
	// verification of the mitigation data?
	Methodology string
}

/***************************/
/* ADAPTATION DATA STRUCTS */
/***************************/

type Adaptation struct {
}

type ChildAdaptationData struct {
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

		data := make(map[string]string)
		err = json.Unmarshal(bytes, &data)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		switch reportType := r.URL.Query()["Type"][0]; reportType {
		case "Emissions":
			oracle.IngestEmissions(data)
		case "Pledges":
			oracle.IngestPledges(data)
		case "Mitigation":
			oracle.IngestMitigation(data)
		case "Adaptation":
			oracle.IngestAdaptation(data)
		}

		erpc.MarshalSend(w, data)
	})
}
