package server

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	erpc "github.com/Varunram/essentials/rpc"
)

func setupSwytchApis() {
	getAccessToken()
	getRefreshToken()
	getSwytchUser()
	getAssets()
	getEnergy()
	getEnergyAttribution()
}

type getAccessTokenDataHelper struct {
	Access_token  string `json:"access_token"`
	Issued_at     int64  `json:"issued_at"`
	Refresh_token string `json:"refresh_token"`
	Token_type    string `json:"token_type"`
	Expires_in    int64  `json:"expires_in"`
}

type GetAccessTokenData struct {
	Data []getAccessTokenDataHelper `json:"data"`
}

func getAccessToken() {
	http.HandleFunc("/swytch/accessToken", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		if r.URL.Query()["clientId"] == nil || r.URL.Query()["clientSecret"] == nil ||
			r.URL.Query()["username"] == nil || r.URL.Query()["password"] == nil {
			erpc.MarshalSend(w, StatusBadRequest)
			return
		}
		url := "https://platformapi-staging.swytch.io/v1/oauth/token"
		pwd := "password"

		clientId := r.URL.Query()["clientId"][0]
		clientSecret := r.URL.Query()["clientSecret"][0]
		username := r.URL.Query()["username"][0]
		password := r.URL.Query()["password"][0]

		a := `{
			"grant_type":"` + pwd + `",
			"client_id":"` + clientId + `",
			"client_secret":"` + clientSecret + `",
			"username":"` + username + `",
			"password":"` + password + `"
		}`
		log.Println(a)
		reqbody := strings.NewReader(a)
		req, err := http.NewRequest("POST", url, reqbody)
		if err != nil {
			log.Println(err)
			erpc.MarshalSend(w, StatusInternalServerError)
			return
		}

		req.Header.Add("content-type", "application/json")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println(err)
			erpc.MarshalSend(w, StatusInternalServerError)
			return
		}

		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)

		var x GetAccessTokenData
		err = json.Unmarshal(body, &x)
		if err != nil {
			log.Println(err)
			erpc.MarshalSend(w, StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, x)
	})
}

func getRefreshToken() {
	http.HandleFunc("/swytch/refreshToken", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		if r.URL.Query()["clientId"] == nil || r.URL.Query()["clientSecret"] == nil ||
			r.URL.Query()["refreshToken"] == nil {
			erpc.MarshalSend(w, StatusBadRequest)
			return
		}

		url := "https://platformapi-staging.swytch.io/v1/oauth/token"
		pwd := "refresh_token"

		clientId := r.URL.Query()["clientId"][0]
		clientSecret := r.URL.Query()["clientSecret"][0]
		refreshToken := r.URL.Query()["refreshToken"][0]

		a := `
		{
			"grant_type":"` + pwd + `",
			"client_id":"` + clientId + `",
			"client_secret":"` + clientSecret + `",
			"refresh_token": "` + refreshToken + `"
		}`

		reqbody := strings.NewReader(a)
		req, err := http.NewRequest("POST", url, reqbody)
		if err != nil {
			log.Println(err)
			erpc.MarshalSend(w, StatusInternalServerError)
			return
		}

		req.Header.Add("content-type", "application/json")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println(err)
			erpc.MarshalSend(w, StatusInternalServerError)
			return
		}

		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)

		var x GetAccessTokenData
		err = json.Unmarshal(body, &x)
		if err != nil {
			log.Println(err)
			erpc.MarshalSend(w, StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, x)
	})
}

type getSwytchUserStructToken struct {
	Valid_token_balance bool   `json:"valid_token_balance"`
	Checked_on          string `json:"checked_on"`
	Token_hash          string `json:"token_hash"`
}

type getSwytchUserStructHelper struct {
	Id            string                   `json:"id"`
	First_name    string                   `json:"first_name"`
	Last_name     string                   `json:"last_name"`
	Name          string                   `json:"name"`
	Email         string                   `json:"email"`
	Username      string                   `json:"username"`
	Roles         []string                 `json:"roles"`
	Token_staking getSwytchUserStructToken `json:"token_staking"`
	Wallet        string                   `json:"wallet"`
}

type GetSwytchUserStruct struct {
	Data []getSwytchUserStructHelper `json:"data"`
}

func getSwytchUser() {
	http.HandleFunc("/swytch/getuser", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		if r.URL.Query()["authToken"] == nil {
			erpc.MarshalSend(w, StatusBadRequest)
			return
		}

		url := "https://platformapi-staging.swytch.io/v1/auth/user"
		auth_token := r.URL.Query()["authToken"][0]

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Println(err)
			erpc.MarshalSend(w, StatusInternalServerError)
		}

		req.Header.Add("authorization", "Bearer "+auth_token)
		req.Header.Add("cache-control", "no-cache")

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println(err)
			erpc.MarshalSend(w, StatusInternalServerError)
		}

		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println(err)
			erpc.MarshalSend(w, StatusInternalServerError)
		}

		log.Println(res)

		var x GetSwytchUserStruct
		err = json.Unmarshal(body, &x)
		if err != nil {
			log.Println(err)
			erpc.MarshalSend(w, StatusInternalServerError)
		}

		erpc.MarshalSend(w, x)
	})
}

type gA2Meta struct {
	Manufacturer    string  `json:"manufacturer"`
	NameplateRating float64 `json:"nameplateRating"`
	SerialNO        string  `json:"serialNO"`
	ThingName       string  `json:"thingName"`
	ThingArn        string  `json:"thingArn"`
	ThingId         string  `json:"thingId"`
}

type gA2Position struct {
	Id          string    `json:"_id"`
	Coordinates []float64 `json:"coordinates"`
	Type        string    `json:"type"`
}

type gA2 struct {
	Id         string      `json:"_id"`
	UpdatedAt  string      `json:"updatedAt"`
	CreatedAt  string      `json:"createdAt"`
	Position   gA2Position `json:"position"`
	Arn        string      `json:"arn"`
	Asset_id   string      `json:"asset_id"`
	Owner_id   string      `json:"owner_id"`
	Name       string      `json:"name"`
	Type       string      `json:"type"`
	Location   string      `json:"location"`
	Meta       gA2Meta     `json:"meta"`
	Country    string      `json:"country"`
	Status     string      `json:"status"`
	Generating bool        `json:"generating"`
	Node_type  string      `json:"node_type"`
}

type GetAssetStruct struct {
	Data []gA2 `json:"data"`
}

func getAssets() {
	http.HandleFunc("/swytch/getassets", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		if r.URL.Query()["authToken"] == nil || r.URL.Query()["userId"] == nil {
			erpc.MarshalSend(w, StatusBadRequest)
			return
		}

		auth_token := r.URL.Query()["authToken"][0]
		userId := r.URL.Query()["userId"][0]

		url := "https://platformapi-staging.swytch.io/v1/users/" + userId + "/assets"

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Println(err)
			erpc.MarshalSend(w, StatusInternalServerError)
		}

		req.Header.Add("authorization", "Bearer "+auth_token)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println(err)
			erpc.MarshalSend(w, StatusInternalServerError)
		}

		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println(err)
			erpc.MarshalSend(w, StatusInternalServerError)
		}

		log.Println(string(body))
		var x GetAssetStruct
		err = json.Unmarshal(body, &x)
		if err != nil {
			log.Println(err)
			erpc.MarshalSend(w, StatusInternalServerError)
		}

		erpc.MarshalSend(w, x)
	})
}

type gEMetadata struct {
	Status           string  `json:"status"`
	Elevation        string  `json:"elevation"`
	Longitude        string  `json:"longitude"`
	Latitude         string  `json:"latitude"`
	Field8           string  `json:"field8"`
	Field7           string  `json:"field7"`
	Field6           string  `json:"field6"`
	Field5           string  `json:"field5"`
	Field4           float64 `json:"field4"`
	Field3           string  `json:"field3"`
	Field2           string  `json:"field2"`
	Field1           string  `json:"field1"`
	Entry_id         float64 `json:"entry_id"`
	Created_at       string  `json:"created_at"`
	Manufacturer     string  `json:"manufacturer"`
	NameplateRating  float64 `json:"nameplateRating"`
	SerialNO         string  `json:"serialNO"`
	ThingName        string  `json:"thingName"`
	ThingArn         string  `json:"thingArn"`
	ThingId          string  `json:"thingId"`
	Source_timestamp string  `json:"source_timestamp"`
}

type getEnergyHelper struct {
	Id               string     `json:"_id"`
	Asset_id         string     `json:"asset_id"`
	Asset_type       string     `json:"asset_type"`
	Source           string     `json:"source"`
	Value            float64    `json:"value"`
	Unit             string     `json:"unit"`
	Lat              string     `json:"lat"`
	Lng              string     `json:"lng"`
	Energy_timestamp string     `json:"energy_timestamp"`
	Timestamp        string     `json:"timestamp"`
	Metadata         gEMetadata `json:"meta"`
	Hash             string     `json:"hash"`
	Block_id         string     `json:"block_id"`
	Block_hash       string     `json:"block_hash"`
	Block_time       string     `json:"block_time"`
	CreatedAt        string     `json:"createdAt"`
	UpdatedAt        string     `json:"updatedAt"`
}

type GetEnergyStruct struct {
	Data []getEnergyHelper `json:"data"`
}

func getEnergy() {
	http.HandleFunc("/swytch/getenergy", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		if r.URL.Query()["authToken"] == nil || r.URL.Query()["assetId"] == nil {
			erpc.MarshalSend(w, StatusBadRequest)
			return
		}

		auth_token := r.URL.Query()["authToken"][0]
		assetId := r.URL.Query()["assetId"][0]

		url := "https://platformapi-staging.swytch.io/v1/assets/" + assetId + "/energy?limit=100&offset=0"

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Println(err)
			erpc.MarshalSend(w, StatusInternalServerError)
		}

		req.Header.Add("authorization", "Bearer "+auth_token)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println(err)
			erpc.MarshalSend(w, StatusInternalServerError)
		}

		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println(err)
			erpc.MarshalSend(w, StatusInternalServerError)
		}

		var x GetEnergyStruct
		err = json.Unmarshal(body, &x)
		if err != nil {
			log.Println(err)
			erpc.MarshalSend(w, StatusInternalServerError)
		}

		erpc.MarshalSend(w, x)
	})
}

type getEnergyAttributionOrigin struct {
	Id          string    `json:"_id"`
	Coordinates []float64 `json:"coordinates"`
	Type        string    `json:"type"`
}

type getEnergyAttributionI struct {
	NextTokenFactorEffectiveAfter string
	NextTokenFactor               string
	CurrentTokenFactor            float64
}

type getEnergyAttributionInputs struct {
	TokenFactor   getEnergyAttributionI `json:"tokenFactor"`
	CarbonOffsets []string              `json:"carbonOffsets"`
}

type getEnergyAttributionHelper struct {
	Id                     string                     `json:"_id"`
	Asset_id               string                     `json:"asset_id"`
	Attribution_holder     string                     `json:"attribution_holder"`
	Carbon_offset          string                     `json:"carbon_offset"`
	Energy_produced        string                     `json:"energy_produced"`
	Actual_energy_produced string                     `json:"actual_energy_produced"`
	Token_award            string                     `json:"token_award"`
	Version                string                     `json:"version"`
	Origin                 getEnergyAttributionOrigin `json:"origin"`
	Asset_type             string                     `json:"asset_type"`
	Production_period      string                     `json:"production_period"`
	Timestamp              string                     `json:"timestamp"`
	Epoch                  string                     `json:"epoch"`
	Block_hash             string                     `json:"block_hash"`
	Validation_authority   string                     `json:"validation_authority"`
	Signature              string                     `json:"signature"`
	Token_id               float64                    `json:"token_id"`
	CreatedAt              string                     `json:"createdAt"`
	UpdatedAt              string                     `json:"updatedAt"`
	Inputs                 getEnergyAttributionInputs `json:"inputs"`
	Tags                   string                     `json:"tags"`
	Tx_history             string                     `json:"tx_history"`
	Processing_status      string                     `json:"processing_status"`
	Transactions           []string                   `json:"transactions"`
	Redeemable             bool                       `json:"redeemable"`
	Claimed                bool                       `json:"claimed"`
	Confirmed              bool                       `json:"confirmed"`
}

type GetEnergyAttributionData struct {
	Data []getEnergyAttributionHelper `json:"data"`
}

func getEnergyAttribution() {
	http.HandleFunc("/swytch/geteattributes", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		if r.URL.Query()["authToken"] == nil || r.URL.Query()["assetId"] == nil {
			erpc.MarshalSend(w, StatusBadRequest)
			return
		}

		auth_token := r.URL.Query()["authToken"][0]
		assetId := r.URL.Query()["assetId"][0]

		url := "https://platformapi-staging.swytch.io/v1/assets/" + assetId + "/attributions?limit=100&offset=0"

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Println(err)
			erpc.MarshalSend(w, StatusInternalServerError)
		}

		req.Header.Add("authorization", "Bearer "+auth_token)

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Println(err)
			erpc.MarshalSend(w, StatusInternalServerError)
		}

		defer res.Body.Close()
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			log.Println(err)
			erpc.MarshalSend(w, StatusInternalServerError)
		}

		log.Println(string(body))
		var x GetEnergyAttributionData
		err = json.Unmarshal(body, &x)
		if err != nil {
			log.Println(err)
			erpc.MarshalSend(w, StatusInternalServerError)
		}

		erpc.MarshalSend(w, x)
	})
}
