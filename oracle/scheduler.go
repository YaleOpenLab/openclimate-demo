package oracle

import (
	"github.com/pkg/errors"
	"github.com/robfig/cron"
	"log"
	// "github.com/YaleOpenLab/openclimate/blockchain"
)

// Valid reportType(s):
// - "Atmospheric CO2"
// - "Global Temperature"
// - TBD: Emissions related reports
// - TBD: Mitigation related reports

const (
	// Earth data parameters
	earthCO2        = "Atmospheric CO2"
	earthTemp       = "Global Temperature"
	earthEntityID   = 1
	earthEntityType = "Earth"
)

type getExternalData func() (interface{}, error)

// GetVerifyCommit first retrieves data calling the API wrapper, then
// sends the data to the oracle for verification and storage on to the
// Ethereum blockchain. Returns a func() so that it can be passed as a
// parameter to c.AddFunc().
func getVerifyCommit(reportType string, entityType string, entityID int, get getExternalData) func() {
	return func() {
		data, err := get()
		if err != nil {
			log.Fatal(errors.Wrap(err, "GetVerifyCommit() failed"))
		}

		err = VerifyAndCommit(reportType, entityType, entityID, data)
		if err != nil {
			log.Fatal(errors.Wrap(err, "GetVerifyCommit() failed"))
		}
	}
}

// Schedule() schedules regular calls to APIs and FTP servers for data,
// sends the data to the oracle for verification and storage on IPFS
// and Ethereum. To add new scheduled API wrappers, simply add another
// c.AddFunc() call with your desired schedule, the report type of the
// data received from the external API, the type of climate actor that the
// data is associated with, the ID of the actor, and the function. The function
// must not receive any parameters and must output (interface{}, error); must
// conform to the getExternalData function type.
func Schedule() {

	c := cron.New()
	c.AddFunc("@daily", getVerifyCommit(earthCO2, earthEntityType, earthEntityID, getNoaaDailyCO2))
	c.AddFunc("@monthly", getVerifyCommit(earthCO2, earthEntityType, earthEntityID, getNoaaMonthlyCO2))

	c.Start()
}
