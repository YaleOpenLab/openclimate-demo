package oracle

import (
	"github.com/pkg/errors"
	"github.com/robfig/cron"
	"log"

	// "github.com/YaleOpenLab/openclimate/blockchain"
)

const (
	// Earth data parameters
	earthCO2 = "Atmospheric CO2"
	earthTemp = "Global Temperature"
	earthEntityID = 1
	earthEntityType = "Earth"
)

func GetVerifyCommit(reportType string, entityType string, entityID int, get func() (interface{}, error)) func() {
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

func Schedule() {
	
	c := cron.New()

	c.AddFunc("@daily", GetVerifyCommit(earthCO2, earthEntityType, earthEntityID, GetNoaaDailyCO2))
	c.AddFunc("@monthly", GetVerifyCommit(earthCO2, earthEntityType, earthEntityID, GetNoaaMonthlyCO2))

	c.Start()
}
