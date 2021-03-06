// Rename to "staticdata.go"

package server

import (
	"encoding/json"
	erpc "github.com/Varunram/essentials/rpc"
	"github.com/Varunram/essentials/utils"
	"github.com/YaleOpenLab/openclimate/database"
	"github.com/YaleOpenLab/openclimate/globals"
	"io/ioutil"
	"log"
	"net/http"
	// "time"
)

func setupDataHandlers() {
	getUSStates()
	getUSCounties()
	getParisAgreement()
	getOceanData()
	getCountryId()
	getCarbonData()
	getCountriesEmissionsData()
}

/*****************************************/
/* US STATES & COUNTIES DATA API HANDLER */
/*****************************************/

type USStatesReturn struct {
	States []string
}

func getUSStates() {
	http.HandleFunc("/us/states", func(w http.ResponseWriter, r *http.Request) {
		_, err := CheckGetAuth(w, r)
		if err != nil {
			return
		}

		database.InitUSStates()
		var x USStatesReturn
		x.States = database.USStates
		erpc.MarshalSend(w, x)
	})
}

type USStateCountiesReturn struct {
	Counties map[string][]string
}

func getUSCounties() {
	http.HandleFunc("/us/counties", func(w http.ResponseWriter, r *http.Request) {
		_, err := CheckGetAuth(w, r)
		if err != nil {
			return
		}

		database.InitUSStates()
		var x USStateCountiesReturn
		x.Counties = database.USStateCities
		erpc.MarshalSend(w, x)
	})
}

/*******************************/
/* PARIS AGREEMENT API HANDLER */
/*******************************/

type ParisAgreementReturn struct {
	Code         string  `json:"Code"`
	Name         string  `json:"Name"`
	Signature    string  `json:"Signature"`
	Ratification string  `json:"Ratification"`
	Kind         string  `json:"Kind"`
	DateOfEffect string  `json:"Date-Of-Effect"`
	Emissions    float64 `json:"Emissions"`
	Percentage   float64 `json:"Percentage"`
	Year         float64 `json:"Year"`
}

type ParisAgreementReturnFinal struct {
	Name         string  `json:"Name"`
	Signature    string  `json:"Signature"`
	Ratification string  `json:"Ratification"`
	Kind         string  `json:"Kind"`
	DateOfEffect string  `json:"Date-Of-Effect"`
	Emissions    float64 `json:"Emissions"`
	Percentage   float64 `json:"Percentage"`
	Year         float64 `json:"Year"`
}

func getParisAgreement() {
	http.HandleFunc("/paris/data", func(w http.ResponseWriter, r *http.Request) {
		_, err := CheckGetAuth(w, r)
		if err != nil {
			return
		}

		data, err := ioutil.ReadFile(globals.StDataDir + "/paris_agreement_entry_into_force.json")
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		var x map[string]ParisAgreementReturn
		err = json.Unmarshal(data, &x)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		y := make(map[string]ParisAgreementReturnFinal)
		for _, values := range x {
			var temp ParisAgreementReturnFinal
			temp.Name = values.Name
			temp.Signature = values.Signature
			temp.Ratification = values.Ratification
			temp.Kind = values.Kind
			temp.DateOfEffect = values.DateOfEffect
			temp.Emissions = values.Emissions
			temp.Percentage = values.Percentage
			temp.Year = values.Year
			y[values.Code] = temp
		}
		erpc.MarshalSend(w, y)
	})
}

/**************************/
/* OCEAN DATA API HANDLER */
/**************************/

type OceanDataPrelim struct {
	Year              int     `json:"year"`
	OceanSink         float64 `json:"Ocean-Sink"`
	CCSMBEC           float64 `json:"CCSM-BEC"`
	MITgcmREcoM2      float64 `json:"MITgcm-REcoM2"`
	MPIOMHAMOCC       float64 `json:"MPIOM-HAMOCC"`
	NEMO36PISCESv2gas float64 `json:"NEMO3.6-PISCESv2-gas (CNRM)"`
	NEMOPISCESIPSL    float64 `json:"NEMO-PISCES (IPSL)"`
	NEMOPlankTOM5     float64 `json:"NEMO-PlankTOM5"`
	NorESMOC          float64 `json:"NorESM-OC"`
	Landschutzer      float64 `json:"Landschützer"`
	Rodenbeck         float64 `json:"Rödenbeck"`
}

type OceanDataFinal struct {
	OceanSink         float64 `json:"Ocean-Sink"`
	CCSMBEC           float64 `json:"CCSM-BEC"`
	MITgcmREcoM2      float64 `json:"MITgcm-REcoM2"`
	MPIOMHAMOCC       float64 `json:"MPIOM-HAMOCC"`
	NEMO36PISCESv2gas float64 `json:"NEMO3.6-PISCESv2-gas (CNRM)"`
	NEMOPISCESIPSL    float64 `json:"NEMO-PISCES (IPSL)"`
	NEMOPlankTOM5     float64 `json:"NEMO-PlankTOM5"`
	NorESMOC          float64 `json:"NorESM-OC"`
	Landschutzer      float64 `json:"Landschützer"`
	Rodenbeck         float64 `json:"Rödenbeck"`
}

func getOceanData() {
	http.HandleFunc("/ocean/data", func(w http.ResponseWriter, r *http.Request) {
		_, err := CheckGetAuth(w, r)
		if err != nil {
			return
		}

		data, err := ioutil.ReadFile(globals.StDataDir + "/ocean_sink.json")
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		var x map[string]OceanDataPrelim
		err = json.Unmarshal(data, &x)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		y := make(map[string]OceanDataFinal)
		for _, values := range x {
			var temp OceanDataFinal
			temp.OceanSink = values.OceanSink
			temp.CCSMBEC = values.CCSMBEC
			temp.MITgcmREcoM2 = values.MITgcmREcoM2
			temp.MPIOMHAMOCC = values.MPIOMHAMOCC
			temp.NEMO36PISCESv2gas = values.NEMO36PISCESv2gas
			temp.NEMOPISCESIPSL = values.NEMOPISCESIPSL
			temp.NEMOPlankTOM5 = values.NEMOPlankTOM5
			temp.NorESMOC = values.NorESMOC
			temp.Landschutzer = values.Landschutzer
			temp.Rodenbeck = values.Rodenbeck
			yearString, err := utils.ToString(values.Year)
			if err != nil {
				log.Println(err)
				erpc.ResponseHandler(w, erpc.StatusInternalServerError)
				return
			}
			y[yearString] = temp
		}
		erpc.MarshalSend(w, y)
	})
}

/*****************************************/
/* GLOBAL CARBON BUDGET DATA API HANDLER */
/*****************************************/

type CarbonDataPrelim struct {
	Year                   int     `json:"Year"`
	FossilFuelAndIndustry  float64 `json:"Fossil-Fuel-And-Industry"`
	LandUseChangeEmissions float64 `json:"Land-Use-Change-Emissions"`
	AtmosphericGrowth      float64 `json:"Atmospheric-Growth"`
	OceanSink              float64 `json:"Ocean-Sink"`
	LandSink               float64 `json:"Land-Sink"`
	BudgetImbalance        float64 `json:"Budget-Imbalance"`
}

type CarbonDataFinal struct {
	FossilFuelAndIndustry  float64 `json:"Fossil-Fuel-And-Industry"`
	LandUseChangeEmissions float64 `json:"Land-Use-Change-Emissions"`
	AtmosphericGrowth      float64 `json:"Atmospheric-Growth"`
	OceanSink              float64 `json:"Ocean-Sink"`
	LandSink               float64 `json:"Land-Sink"`
	BudgetImbalance        float64 `json:"Budget-Imbalance"`
}

func getCarbonData() {
	http.HandleFunc("/carbon/budget", func(w http.ResponseWriter, r *http.Request) {
		_, err := CheckGetAuth(w, r)
		if err != nil {
			return
		}

		data, err := ioutil.ReadFile(globals.StDataDir + "/global_carbon_budget.json")
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		var x map[string]CarbonDataPrelim
		err = json.Unmarshal(data, &x)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		y := make(map[string]CarbonDataFinal)
		for _, values := range x {
			var temp CarbonDataFinal
			temp.FossilFuelAndIndustry = values.FossilFuelAndIndustry
			temp.LandUseChangeEmissions = values.LandUseChangeEmissions
			temp.AtmosphericGrowth = values.AtmosphericGrowth
			temp.OceanSink = values.OceanSink
			temp.LandSink = values.LandSink
			temp.BudgetImbalance = values.BudgetImbalance
			yearString, err := utils.ToString(values.Year)
			if err != nil {
				log.Println(err)
				erpc.ResponseHandler(w, erpc.StatusInternalServerError)
				return
			}
			y[yearString] = temp
		}
		erpc.MarshalSend(w, y)
	})
}

/*******************************/
/* COUNTRIES EMISSIONS HANDLER */
/*******************************/

type CountriesEmissionsPrelim struct {
	Nation     string
	Year       int
	Total      int
	SolidFuel  float64
	LiquidFuel float64
	GasFuel    float64
	Cement     int
	GasFlaring float64
	PerCapita  float64
	Bunkers    int
}

type CountriesEmissionsFinal struct {
	Year       int
	Total      int
	SolidFuel  float64
	LiquidFuel float64
	GasFuel    float64
	Cement     int
	GasFlaring float64
	PerCapita  float64
	Bunkers    int
}

func getCountriesEmissionsData() {

	http.HandleFunc("/countries/emissions", func(w http.ResponseWriter, r *http.Request) {
		_, err := CheckGetAuth(w, r)
		if err != nil {
			return
		}

		data, err := ioutil.ReadFile(globals.StDataDir + "/countries_emissions_2014.json")
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		var dataMapPrelim map[string]CountriesEmissionsPrelim
		err = json.Unmarshal(data, &dataMapPrelim)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		dataMapFinal := make(map[string]CountriesEmissionsFinal)
		for _, values := range dataMapPrelim {

			var temp CountriesEmissionsFinal
			temp.Year = values.Year
			temp.Total = values.Total
			temp.SolidFuel = values.SolidFuel
			temp.LiquidFuel = values.LiquidFuel
			temp.GasFuel = values.GasFuel
			temp.Cement = values.Cement
			temp.GasFlaring = values.GasFlaring
			temp.PerCapita = values.PerCapita
			temp.Bunkers = values.Bunkers

			dataMapFinal[values.Nation] = temp

		}
		erpc.MarshalSend(w, dataMapFinal)
	})
}

/**************************/
/* COUNTRY ID API HANDLER */
/**************************/

type CountryIdResponse struct {
	CountryIds map[int]string
}

func getCountryId() {
	http.HandleFunc("/countries/id", func(w http.ResponseWriter, r *http.Request) {
		_, err := CheckGetAuth(w, r)
		if err != nil {
			return
		}

		countryIds := database.InitUSStates()
		var x CountryIdResponse
		x.CountryIds = countryIds
		erpc.MarshalSend(w, x)
	})
}
