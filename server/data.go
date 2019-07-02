package server

import (
	"encoding/json"
	utils "github.com/Varunram/essentials/utils"
	"github.com/YaleOpenLab/openclimate/database"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func dataHandler() {
	getUSStates()
	getUSCounties()
	getParisAgreement()
	getOceanData()
	queryNazca()
	queryNazcaCountry()
	getCountryId()
	getCarbonData()
}

type USStatesReturn struct {
	States []string
}

func getUSStates() {
	http.HandleFunc("/us/states", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)

		_, err := authorizeUser(r)
		if err != nil {
			log.Println("could not retrieve user from the database, quitting")
			responseHandler(w, StatusBadRequest)
			return
		}

		database.InitUSStates()
		var x USStatesReturn
		x.States = database.USStates
		MarshalSend(w, x)
	})
}

type USStateCountiesReturn struct {
	Counties map[string][]string
}

func getUSCounties() {
	http.HandleFunc("/us/counties", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)

		_, err := authorizeUser(r)
		if err != nil {
			log.Println("could not retrieve user from the database, quitting")
			responseHandler(w, StatusBadRequest)
			return
		}

		database.InitUSStates()
		var x USStateCountiesReturn
		x.Counties = database.USStateCities
		MarshalSend(w, x)
	})
}

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
		checkGet(w, r)
		checkOrigin(w, r)

		_, err := authorizeUser(r)
		if err != nil {
			log.Println("could not retrieve user from the database, quitting")
			responseHandler(w, StatusBadRequest)
			return
		}

		data, err := ioutil.ReadFile("data/json_data/paris_agreement_entry_into_force.json")
		if err != nil {
			log.Println(err)
			responseHandler(w, StatusInternalServerError)
			return
		}

		var x map[string]ParisAgreementReturn
		err = json.Unmarshal(data, &x)
		if err != nil {
			log.Println(err)
			responseHandler(w, StatusInternalServerError)
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
		MarshalSend(w, y)
	})
}

type OceanDataPrelim struct {
	Year              int     `json:"year"`
	OceanSink         float64 `json:"Ocean-Sink"`
	CCSMBEC           float64 `json:"CCSM-BEC"`
	MITgcmREcoM2      float64 `json:"MITgcm-REcoM2"`
	MPIOMHAMOCC       float64 `json:"MPIOM-HAMOCC"`
	NEMO36PISCESv2gas float64 `json:"NEMO3.6-PISCESv2-gas (CNRM)""`
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
	NEMO36PISCESv2gas float64 `json:"NEMO3.6-PISCESv2-gas (CNRM)""`
	NEMOPISCESIPSL    float64 `json:"NEMO-PISCES (IPSL)"`
	NEMOPlankTOM5     float64 `json:"NEMO-PlankTOM5"`
	NorESMOC          float64 `json:"NorESM-OC"`
	Landschutzer      float64 `json:"Landschützer"`
	Rodenbeck         float64 `json:"Rödenbeck"`
}

func getOceanData() {
	http.HandleFunc("/ocean/data", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)

		_, err := authorizeUser(r)
		if err != nil {
			log.Println("could not retrieve user from the database, quitting")
			responseHandler(w, StatusBadRequest)
			return
		}

		data, err := ioutil.ReadFile("data/json_data/ocean_sink.json")
		if err != nil {
			log.Println(err)
			responseHandler(w, StatusInternalServerError)
			return
		}

		var x map[string]OceanDataPrelim
		err = json.Unmarshal(data, &x)
		if err != nil {
			log.Println(err)
			responseHandler(w, StatusInternalServerError)
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
			y[utils.ItoS(values.Year)] = temp
		}
		MarshalSend(w, y)
	})
}

type CarbonDataPrelim struct {
	Year                   int `json:"Year"`
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
		checkGet(w, r)
		checkOrigin(w, r)

		_, err := authorizeUser(r)
		if err != nil {
			log.Println("could not retrieve user from the database, quitting")
			responseHandler(w, StatusBadRequest)
			return
		}

		data, err := ioutil.ReadFile("data/json_data/global_carbon_budget.json")
		if err != nil {
			log.Println(err)
			responseHandler(w, StatusInternalServerError)
			return
		}

		var x map[string]CarbonDataPrelim
		err = json.Unmarshal(data, &x)
		if err != nil {
			log.Println(err)
			responseHandler(w, StatusInternalServerError)
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
			y[utils.ItoS(values.Year)] = temp
		}
		MarshalSend(w, y)
	})
}

var NazcaURL = "https://nazcaapiprod.howoco.com/handlers/countrystakeholders.ashx?countryid="

type NazcaResponse struct {
	EntityID       string `json:"entityID"`
	EntityName     string `json:"entityName"`
	CountryName    string `json:"countryName"`
	EntityTypeName string `json:"entityTypeName"`
	Actions        []struct {
		ActionType  string `json:"actionType"`
		Name        string `json:"name"`
		Description string `json:"description"`
		Adhoc       string `json:"adhoc"`
	}
}

func queryNazca() {
	http.HandleFunc("/nazca/data", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)

		_, err := authorizeUser(r)
		if err != nil {
			log.Println("could not retrieve user from the database, quitting")
			responseHandler(w, StatusBadRequest)
			return
		}

		for i := 173; i < 174; i++ {
			apiUrl := "https://nazcaapiprod.howoco.com/handlers/countrystakeholders.ashx?countryid=" + utils.ItoS(i)
			data, err := Get(apiUrl)
			if err != nil {
				log.Println("country: ", i, "not queryable", err)
				time.Sleep(1 * time.Second)
				continue
			}
			var x []NazcaResponse
			err = json.Unmarshal(data, &x)
			if err != nil {
				log.Println("could not unmarshal data, quitting", err)
				responseHandler(w, StatusInternalServerError)
				return
			}
			time.Sleep(1 * time.Second)
			MarshalSend(w, x)
		}
	})
}

func queryNazcaCountry() {
	http.HandleFunc("/nazcacountry/data", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)

		_, err := authorizeUser(r)
		if err != nil {
			log.Println("could not retrieve user from the database, quitting")
			responseHandler(w, StatusBadRequest)
			return
		}

		countryMap := make(map[int]string)
		for i := 1; i < 181; i++ {
			apiUrl := "https://nazcaapiprod.howoco.com/handlers/countrystakeholders.ashx?countryid=" + utils.ItoS(i)
			data, err := Get(apiUrl)
			if err != nil {
				log.Println("country: ", i, "not queryable", err)
				time.Sleep(1 * time.Second)
				continue
			}
			var x []NazcaResponse
			err = json.Unmarshal(data, &x)
			if err != nil {
				log.Println("could not unmarshal data, quitting", err)
				responseHandler(w, StatusInternalServerError)
				return
			}
			if len(x) != 0 {
				countryMap[i] = x[0].CountryName
			}
			time.Sleep(1 * time.Second)
		}

	})
}

type CountryIdResponse struct {
	CountryIds map[int]string
}

func getCountryId() {
	http.HandleFunc("/countries/id", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)

		_, err := authorizeUser(r)
		if err != nil {
			log.Println("could not retrieve user from the database, quitting")
			responseHandler(w, StatusBadRequest)
			return
		}

		countryIds := database.InitUSStates()
		var x CountryIdResponse
		x.CountryIds = countryIds
		MarshalSend(w, x)
	})
}
