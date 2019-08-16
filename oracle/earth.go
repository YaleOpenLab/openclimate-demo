package oracle

import (
	// "encoding/json"
	"github.com/pkg/errors"
	// "net/http"
	// "log"
	// "github.com/YaleOpenLab/openclimate/globals"
	"github.com/jlaffaye/ftp"
	"io/ioutil"
	"strconv"
	"strings"
	"time"
)

const (
	// NOAA FTP addresses
	noaaFtpAddress  = "aftp.cmdl.noaa.gov:21"
	globalTrendPath = "products/trends/co2/co2_trend_gl.txt"
	barrowPath      = "data/trace_gases/co2/flask/surface/co2_brw_surface-flask_1_ccgg_month.txt"
	maunaLoaPath    = "data/trace_gases/co2/flask/surface/co2_mlo_surface-flask_1_ccgg_month.txt"
	southPolePath   = "data/trace_gases/co2/flask/surface/co2_spo_surface-flask_1_ccgg_month.txt"
	amSamoaPath     = "data/trace_gases/co2/flask/surface/co2_smo_surface-flask_1_ccgg_month.txt"

	// NOAA HTTP URLs
	noaaBaseUrl = "https://www.ncdc.noaa.gov/cdo-web/webservices/v2/data"
)

type GlobalTemp struct {
	Source    string
	Location  string
	Frequency string // annually, monthly, daily, etc.
	Year      int
	Month     int
	Day       int
	Cycle     float64
	Trend     float64
}

// The GlobalCO2 struct defines the data and meta-data that will be attached
// To the atmospheric CO2 measurements observed at various sites and stored
// on the NOAA ESRL FTP server.
type GlobalCO2 struct {
	Source    string
	Location  string
	Frequency string // annually, monthly, daily, etc.
	Year      int
	Month     int
	Day       int
	Cycle     float64
	Trend     float64
}

// GetNoaaDailyCO2() retrieves data from the file that stores all atmospheric
// CO2 global estimates on the NOAA ESRL FTP server. The function is called
// by the scheduler daily (see oracle/scheduler.go), and the data from this
// function is sent to the oracle for processing & verification.
func GetNoaaDailyCO2() (interface{}, error) {

	var dataArr []GlobalCO2

	filestrings, err := RetrieveNoaaCO2(globalTrendPath)
	if err != nil {
		return dataArr, err
	}

	for _, fs := range filestrings {

		raw, err := ParseNoaaCO2(fs.FileStr, 5)
		if err != nil {
			return dataArr, errors.Wrap(err, "GetNoaaMonthlyCO2() failed")
		}

		var data GlobalCO2
		data.Source = "NOAA"
		data.Location = fs.Name
		data.Frequency = "Daily"
		data.Year = int(raw[0])
		data.Month = int(raw[1])
		data.Day = int(raw[2])
		data.Cycle = raw[3]
		data.Trend = raw[4]

		dataArr = append(dataArr, data)
	}

	return dataArr, nil
}

// GetNoaaDailyCO2() retrieves data from the files that store atmospheric
// CO2 measurements observed from various sites on the NOAA ESRL FTP server.
// The function is called by the scheduler monthly (see oracle/scheduler.go),
// and the data from this function is sent to the oracle for processing &
// verification.
func GetNoaaMonthlyCO2() (interface{}, error) {

	var dataArr []GlobalCO2

	filestrings, err := RetrieveNoaaCO2(barrowPath, maunaLoaPath, southPolePath, amSamoaPath)
	if err != nil {
		return dataArr, errors.Wrap(err, "GetNoaaMonthlyCO2() failed")
	}

	for _, fs := range filestrings {

		raw, err := ParseNoaaCO2(fs.FileStr, 3)
		if err != nil {
			return dataArr, errors.Wrap(err, "GetNoaaMonthlyCO2() failed")
		}

		var data GlobalCO2
		data.Source = "NOAA"
		data.Location = fs.Name
		data.Frequency = "Monthly"
		data.Year = int(raw[0])
		data.Month = int(raw[1])
		// data.Day is left blank
		data.Cycle = raw[2]
		// data.Trend is left blank

		dataArr = append(dataArr, data)
	}

	return dataArr, nil
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

type FileString struct {
	Name    string
	FileStr string
}

func RetrieveNoaaCO2(filepaths ...string) ([]FileString, error) {
	var fstrings []FileString

	c, err := ftp.Dial(noaaFtpAddress, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return fstrings, errors.Wrap(err, "getNoaaGlobalDailyTrend() failed")
	}

	err = c.Login("anonymous", "anonymous")
	if err != nil {
		return fstrings, errors.Wrap(err, "getNoaaGlobalDailyTrend() failed")
	}

	for _, fp := range filepaths {
		resp, err := c.Retr(fp)
		if err != nil {
			return fstrings, errors.Wrap(err, "getNoaaGlobalDailyTrend() failed")
		}

		buf, err := ioutil.ReadAll(resp)
		if err != nil {
			return fstrings, errors.Wrap(err, "getNoaaGlobalDailyTrend() failed")
		}

		var fs FileString
		fpslice := strings.Split(fp, "/")
		filename := fpslice[len(fpslice)-1]
		fs.Name = strings.Split(filename, ".")[0]
		fs.FileStr = string(buf)

		fstrings = append(fstrings, fs)

		err = resp.Close()
		if err != nil {
			return fstrings, errors.Wrap(err, "getNoaaGlobalDailyTrend() failed")
		}
	}

	if err := c.Quit(); err != nil {
		return fstrings, errors.Wrap(err, "getNoaaGlobalDailyTrend() failed")
	}

	return fstrings, nil
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
