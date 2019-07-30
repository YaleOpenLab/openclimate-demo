package wrappers

import (
	"encoding/json"
	erpc "github.com/Varunram/essentials/rpc"
	// "github.com/Varunram/essentials/utils"
	"net/http"
	// "log"
)

// func queryNoaa() (map[string]string, error) {
// 	url := "https://www.ncdc.noaa.gov/cdo-web/webservices/v2/"

// 	resp, err := erpc.GetRequest(url)


// }

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