package wrappers

import (
	"encoding/json"
	erpc "github.com/Varunram/essentials/rpc"
	"github.com/pkg/errors"
	// "github.com/Varunram/essentials/utils"
	"net/http"
	// "log"
)

func queryNoaaGlobalSummary() (interface{}, error) {
	base_url := "https://www.ncdc.noaa.gov/cdo-web/webservices/v2/data"
	data_set := "datasetid=gov.noaa.ncdc:C00947"
	startdate := "startdate=2009-01-01"
	enddate := "enddate=2019-01-01"

	url := base_url + "?" + data_set + "&" + startdate + "&" + enddate

	var data interface{}
	body, err := erpc.GetRequest(url)
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
	req.Header.Add("token", main.noaaToken)
	res, err := client.Do(req)
	if err != nil {
		return dummy, errors.Wrap(err, "did not make request")
	}
	defer res.Body.Close()
	return ioutil.ReadAll(res.Body)
}