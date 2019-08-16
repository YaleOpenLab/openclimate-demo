package oracle

import (
	"github.com/pkg/errors"
	"github.com/robfig/cron"
	"log"

	"github.com/YaleOpenLab/openclimate/blockchain"

)

func GetAndCommitDaily() {

	dailyNoaaData, err := GetNoaaDailyCO2()
	if err != nil {
		log.Fatal(errors.Wrap(err, "GetAndCommitDaily() failed"))
	}

	val, err = VerifyAtmosCO2(dailyNoaaData)
	if err != nil {
		log.Fatal(errors.Wrap(err, "GetAndCommitDaily() failed"))
	}

	err = blockchain.CommitToChain(val)
	if err != nil {
		log.Fatal(errors.Wrap(err, "GetAndCommitDaily() failed"))
	}
}

func GetAndCommitMonthly() {

	monthlyNoaaData, err := GetNoaaMonthlyCO2()
	if err != nil {
		log.Fatal(errors.Wrap(err, "GetAndCommitMonthly() failed"))
	}

	_, err = VerifyAtmosCO2(monthlyNoaaData)
	if err != nil {
		log.Fatal(errors.Wrap(err, "GetAndCommitMonthly() failed"))
	}
}

func Schedule() {

	c := cron.New()
	c.AddFunc("@daily", GetAndCommitDaily)
	c.AddFunc("@monthly", GetAndCommitMonthly)

	c.Start()
}
