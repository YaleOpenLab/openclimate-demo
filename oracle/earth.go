package oracle

import (
	"encoding/json"
	// erpc "github.com/Varunram/essentials/rpc"
	"github.com/pkg/errors"
	// "github.com/Varunram/essentials/utils"
	"net/http"
	// "log"
	"github.com/YaleOpenLab/openclimate/globals"
	"github.com/YaleOpenLab/openclimate/ipfs"

	"github.com/jlaffaye/ftp"

	"log"

	"io/ioutil"
	"time"
)

/*
	Holds all the functionality that verifies data concerning the
	state of the earth.
*/

func VerifyEarth(data interface{}) (ipfs.Earth, error) {
	var verifiedData ipfs.Earth
	return verifiedData, nil
}

func VerifyAtmosCO2() {

}

func VerifyGlobalTemp() {

}

func VerifySeaLevelRise() {

}

func VerifyLandUse() {

}

// func VerifyTropOzone() {

// }

// func VerifyStratOzone() {

// }

// func VerifyArcticIceMin() {

// }

// func VerifyIceSheets() {

// }


func RetrNoaaGlobalTrendDailyCO2() (int, error) {

	// log.Println("hit")

	c, err := ftp.Dial("aftp.cmdl.noaa.gov:21", ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return 0, errors.Wrap(err, "getNoaaGlobalDailyTrend() failed")
	}

	err = c.Login("anonymous", "anonymous")
	if err != nil {
		return 0, errors.Wrap(err, "getNoaaGlobalDailyTrend() failed")
	}
	
	resp, err := c.Retr("products/trends/co2/co2_trend_gl.txt")
	if err != nil {
		return 0, errors.Wrap(err, "getNoaaGlobalDailyTrend() failed")
	}
	// defer resp.Close()

	buf, err := ioutil.ReadAll(resp)
	if err != nil {
		return 0, errors.Wrap(err, "getNoaaGlobalDailyTrend() failed")
	}

	log.Println(string(buf))

	err = resp.Close()
	if err != nil {
		return 0, errors.Wrap(err, "getNoaaGlobalDailyTrend() failed")
	}

	if err := c.Quit(); err != nil {
		return 0, errors.Wrap(err, "getNoaaGlobalDailyTrend() failed")
	}	

	return 0, nil
}


func QueryNoaaSummary(datasetid string, startdate string, enddate string) (interface{}, error) {
	baseUrl := "https://www.ncdc.noaa.gov/cdo-web/webservices/v2/data"
	dataset := "datasetid=" + datasetid
	startdate = "startdate=" + startdate
	enddate = "enddate=" + enddate

	url := baseUrl + "?" + dataset + "&" + startdate + "&" + enddate

	var data interface{}
	body, err := GetRequest(url)
	if err != nil {
		return data, errors.Wrap(err, "NOAA query failed")
	}
	json.Unmarshal(body, &data)
	return data, nil
}


func GetRequest(url string) ([]byte, error) {
	var dummy []byte
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return dummy, errors.Wrap(err, "did not create new GET request")
	}
	req.Header.Add("Origin", "localhost")
	req.Header.Add("token", globals.NoaaToken)
	res, err := client.Do(req)
	if err != nil {
		return dummy, errors.Wrap(err, "did not make request")
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}
