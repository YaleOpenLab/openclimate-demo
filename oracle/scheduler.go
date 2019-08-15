package oracle

import (

	"github.com/robfig/cron"
	"log"
	"github.com/pkg/errors"

)

func GetAndCommitDaily() {

	dailyNoaaData, err := GetNoaaDailyCO2()
	if err != nil {
		log.Fatal(errors.Wrap(err, "GetAndCommitDaily() failed"))
	}

	err = Verify("AtmosCO2", "Earth", 1, dailyNoaaData)
	if err != nil {
		log.Fatal(errors.Wrap(err, "GetAndCommitDaily() failed"))
	}
}

func GetAndCommitMonthly() {

	monthlyNoaaData, err := GetNoaaMonthlyCO2()
	if err != nil {
		log.Fatal(errors.Wrap(err, "GetAndCommitMonthly() failed"))
	}

	err = Verify("AtmosCO2", "Earth", 1, monthlyNoaaData)
	if err != nil {
		log.Fatal(errors.Wrap(err, "GetAndCommitMonthly() failed"))
	}
}

func ScheduleNoaaCO2() {

	c := cron.New()
	c.AddFunc("@daily", GetAndCommitDaily)
	c.AddFunc("@monthly", GetAndCommitMonthly)

	c.Start()
}