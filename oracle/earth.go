package oracle

import (
	// "encoding/json"
	"github.com/pkg/errors"
	// "net/http"
	// "log"
	// "github.com/YaleOpenLab/openclimate/globals"
	"github.com/jlaffaye/ftp"
	"strconv"
	"strings"
	"io/ioutil"
	"time"
)

const (
	noaaFtpAddress = "aftp.cmdl.noaa.gov:21"
	globalTrendPath = "products/trends/co2/co2_trend_gl.txt"
	barrowPath = "data/trace_gases/co2/flask/surface/co2_brw_surface-flask_1_ccgg_month.txt"
	maunaLoaPath = "data/trace_gases/co2/flask/surface/co2_mlo_surface-flask_1_ccgg_month.txt"
	southPolePath = "data/trace_gases/co2/flask/surface/co2_spo_surface-flask_1_ccgg_month.txt"
	amSamoaPath = "data/trace_gases/co2/flask/surface/co2_smo_surface-flask_1_ccgg_month.txt"

	noaaBaseUrl = "https://www.ncdc.noaa.gov/cdo-web/webservices/v2/data"

)


func GetNoaaDailyCO2() (map[string][]float64, error) {

	data := make(map[string][]float64)

	fs, err := RetrieveNoaaCO2(globalTrendPath)
	if err != nil {
		return data, err
	}
	
	globalLatest, err := ParseNoaaCO2(fs[0], 5)
	data["global_trend_daily"] = globalLatest

	return data, nil
}


func GetNoaaMonthlyCO2() (map[string][]float64, error) {

	data := make(map[string][]float64)

	fs, err := RetrieveNoaaCO2(barrowPath, maunaLoaPath, southPolePath, amSamoaPath)
	if err != nil {
		return data, err
	}

	barrowLatest, err := ParseNoaaCO2(fs[0], 3)
	data["barrow_monthly"] = barrowLatest

	maunaLoaLatest, err := ParseNoaaCO2(fs[1], 3)
	data["mauna_loa_monthly"] = maunaLoaLatest

	southPoleLatest, err := ParseNoaaCO2(fs[2], 3)
	data["south_pole_monthly"] = southPoleLatest

	amSamoaLatest, err := ParseNoaaCO2(fs[3], 3)
	data["am_samoa_monthly"] = amSamoaLatest

	return data, nil
}


func GetNoaaAnnualCO2() ([]float64, error) {
	maunaLoaPath := "products/trends/co2/co2_annmean_mlo.txt"

	var maunaLoaData []float64

	fs, err := RetrieveNoaaCO2(maunaLoaPath)
	if err != nil {
		return maunaLoaData, err
	}

	maunaLoaData, err = ParseNoaaCO2(fs[0], 3)
	return maunaLoaData, nil
}


func ParseNoaaCO2(filestring string, length int) ([]float64, error) {

	var err error

	substr := strings.Fields(filestring)
	temp := make([]float64, length)
	for i, elt := range substr[len(substr)-length:] {
		temp[i], err = strconv.ParseFloat(elt, 64)
		if err != nil {
			return temp, errors.Wrap(err, "ParseNoaaCO2() failed")
		}
	}
	return temp, nil
}


func RetrieveNoaaCO2(filepaths ...string) ([]string, error) {
	var bufs []string

	c, err := ftp.Dial(noaaFtpAddress, ftp.DialWithTimeout(5*time.Second))
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

// func QueryNoaaSummary(datasetid string, startdate string, enddate string) (interface{}, error) {

// 	baseUrl := noaaBaseUrl
// 	dataset := "datasetid=" + datasetid
// 	startdate = "startdate=" + startdate
// 	enddate = "enddate=" + enddate

// 	url := baseUrl + "?" + dataset + "&" + startdate + "&" + enddate

// 	var data interface{}

// 	body, err := getRequest(url)
// 	if err != nil {
// 		return data, errors.Wrap(err, "NOAA query failed")
// 	}

// 	json.Unmarshal(body, &data)
// 	return data, nil
// }

// func getRequest(url string) ([]byte, error) {

// 	var dummy []byte
// 	client := &http.Client{
// 		Timeout: 5 * time.Second,
// 	}

// 	req, err := http.NewRequest("GET", url, nil)
// 	if err != nil {
// 		return dummy, errors.Wrap(err, "did not create new GET request")
// 	}

// 	req.Header.Add("Origin", "localhost")
// 	req.Header.Add("token", globals.NoaaToken)

// 	res, err := client.Do(req)
// 	if err != nil {
// 		return dummy, errors.Wrap(err, "did not make request")
// 	}

// 	defer res.Body.Close()
// 	return ioutil.ReadAll(res.Body)
// }
