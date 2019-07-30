package wrappers

import (
	"encoding/json"
	erpc "github.com/Varunram/essentials/rpc"
	"github.com/Varunram/essentials/utils"
	// "io/ioutil"
	"log"
	"net/http"
	"time"
)

/**************************/
/* NAZCA DATA API HANDLER */
/**************************/

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
		_, err := CheckGetAuth(w, r)
		if err != nil {
			return
		}

		for i := 173; i < 174; i++ {
			iString, err := utils.ToString(i)
			if err != nil {
				log.Println(err)
				erpc.ResponseHandler(w, erpc.StatusInternalServerError)
				return
			}
			apiUrl := "https://nazcaapiprod.howoco.com/handlers/countrystakeholders.ashx?countryid=" + iString
			// alt url: https://nazcaapiprod.howoco.com/handlers/countrystakeholders.ashx?countryid=173&entitytypeid=3
			// coutnry code for the US is 173, hence the weird loop here. Once we have the frontend ready for
			// the US, we can add the other countries here
			data, err := erpc.GetRequest(apiUrl)
			if err != nil {
				log.Println("country: ", i, "not queryable", err)
				time.Sleep(1 * time.Second)
				continue
			}
			var x []NazcaResponse
			err = json.Unmarshal(data, &x)
			if err != nil {
				log.Println("could not unmarshal data, quitting", err)
				erpc.ResponseHandler(w, erpc.StatusInternalServerError)
				return
			}
			time.Sleep(1 * time.Second)
			erpc.MarshalSend(w, x)
		}
	})
}

func queryNazcaCountry() {
	http.HandleFunc("/nazcacountry/data", func(w http.ResponseWriter, r *http.Request) {
		_, err := CheckGetAuth(w, r)
		if err != nil {
			return
		}

		countryMap := make(map[int]string)
		for i := 1; i < 181; i++ {
			iString, err := utils.ToString(i)
			if err != nil {
				log.Println(err)
				erpc.ResponseHandler(w, erpc.StatusInternalServerError)
				return
			}
			apiUrl := "https://nazcaapiprod.howoco.com/handlers/countrystakeholders.ashx?countryid=" + iString
			data, err := erpc.GetRequest(apiUrl)
			if err != nil {
				log.Println("country: ", i, "not queryable", err)
				time.Sleep(1 * time.Second)
				continue
			}
			var x []NazcaResponse
			err = json.Unmarshal(data, &x)
			if err != nil {
				log.Println("could not unmarshal data, quitting", err)
				erpc.ResponseHandler(w, erpc.StatusInternalServerError)
				return
			}
			if len(x) != 0 {
				countryMap[i] = x[0].CountryName
			}
			time.Sleep(1 * time.Second)
		}
	})
}
