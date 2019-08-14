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

	// "log"

	"strconv"

	"strings"

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

func GetNoaaDailyCO2() (map[string][]float64, error) {

	final := make(map[string][]float64)

	globalTrendPath := "products/trends/co2/co2_trend_gl.txt"
	maunaLoaPath := "products/trends/co2/co2_annmean_mlo.txt"
	barrowPath := "data/trace_gases/co2/flask/surface/co2_brw_surface-flask_1_ccgg_month.txt"

	fs, err := RetrieveNoaaDailyCO2(globalTrendPath, maunaLoaPath, barrowPath)
	if err != nil {
		return final, err
	}

	globalData, err := ParseNoaaDailyCO2(fs[0], 5)
	final["global_daily"] = globalData

	maunaLoaData, err := ParseNoaaDailyCO2(fs[1], 3)
	final["mauna_loa_annual"] = maunaLoaData

	barrowData, err := ParseNoaaDailyCO2(fs[2], 3)
	final["barrow_monthly"] = barrowData

	return final, err
}


func ParseNoaaDailyCO2(filestring string, length int) ([]float64, error) {

	var err error

	substr := strings.Fields(filestring)
	temp := make([]float64, length)
	for i, elt := range substr[len(substr)-length:] {
		temp[i], err = strconv.ParseFloat(elt, 64)
		if err != nil {
			return temp, errors.Wrap(err, "ParseNoaaDailyCO2() failed")
		}
	}
	// idx := len(substr) - length
	return temp, nil
}


func RetrieveNoaaDailyCO2(filepaths ...string) ([]string, error) {

	var bufs []string

	c, err := ftp.Dial("aftp.cmdl.noaa.gov:21", ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return bufs, errors.Wrap(err, "getNoaaGlobalDailyTrend() failed")
	}

	err = c.Login("anonymous", "anonymous")
	if err != nil {
		return bufs, errors.Wrap(err, "getNoaaGlobalDailyTrend() failed")
	}
	
	for _, fp := range filepaths {
		resp, err := c.Retr(fp)
		if err != nil {
			return bufs, errors.Wrap(err, "getNoaaGlobalDailyTrend() failed")
		}

		buf, err := ioutil.ReadAll(resp)
		if err != nil {
			return bufs, errors.Wrap(err, "getNoaaGlobalDailyTrend() failed")
		}

		bufs = append(bufs, string(buf))

		err = resp.Close()
		if err != nil {
			return bufs, errors.Wrap(err, "getNoaaGlobalDailyTrend() failed")
		}
	}

	if err := c.Quit(); err != nil {
		return bufs, errors.Wrap(err, "getNoaaGlobalDailyTrend() failed")
	}

	return bufs, nil
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
