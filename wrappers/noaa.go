package wrappers

import (
	"encoding/json"
	erpc "github.com/Varunram/essentials/rpc"
	"github.com/pkg/errors"
	// "github.com/Varunram/essentials/utils"
	"net/http"
	// "log"
	"time"
	"github.com/YaleOpenLab/openclimate/globals"
	"io/ioutil"
)

func queryNoaaGlobalSummary(startdate string, enddate string) (interface{}, error) {
	baseUrl := "https://www.ncdc.noaa.gov/cdo-web/webservices/v2/data"
	dataset := "datasetid=gov.noaa.ncdc:C00947"
	startdate = "startdate=" + startdate
	enddate = "enddate=" + enddate

	url := baseUrl + "?" + dataset + "&" + startdate + "&" + enddate

	var data interface{}
	body, err := erpc.GetRequest(url)
	if err != nil {
		return data, errors.Wrap(err, "NOAA query failed")
	}
	json.Unmarshal(body, &data)
	return data, nil
}

// func queryNoaaMonth


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