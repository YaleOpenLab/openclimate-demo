package server

import (
	// "github.com/pkg/errors"
	"log"
	// "math/big"
	"net/http"
	"strconv"
	"io/ioutil"
	"encoding/json"

	// ipfs "github.com/Varunram/essentials/ipfs"
	erpc "github.com/Varunram/essentials/rpc"
	db "github.com/YaleOpenLab/openclimate/database"
)


/*
	Function allows admins of a particular entity to "verify" other users who claim to be
	part of the same entity.

	URL parameters: 
	- "candidate_id": the ID of the user who is being considered for verification
*/
func VerifyUser() {
	http.HandleFunc("/user/admin/verify", func(w http.ResponseWriter, r *http.Request) {

		var candidate db.User

		_, err := CheckPostAdmin(w, r) // Check if the person is an admin/has authority to verify users
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		if r.URL.Query()["candidate_id"] == nil {
			log.Println("Candidate for verification not specified")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(r.URL.Query()["candidate_id"][0])
		if err != nil {
			log.Println("Failed to typecast candidate id from string to int")
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		candidate, err = db.RetrieveUser(id)
		if err != nil {
			log.Println("Candidate could not be found in database")
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		candidate.Verified = true
		candidate.Save()

		erpc.MarshalSend(w, candidate)
	})
}

/*
	URL parameters: N/A
	Response body (all strings):
	- name
	- companyID (not required)
	- location
	- type

*/
func AddAsset() {
	http.HandleFunc("/user/assets/add", func(w http.ResponseWriter, r *http.Request) {

		user, err := CheckPostAdmin(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		bytes, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		var asset map[string]string
		err = json.Unmarshal(bytes, &asset)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		name := asset["name"]

		companyID := user.EntityID

		location := asset["location"]
		assetType := asset["type"]

		new, err := db.NewAsset(name, companyID, location, assetType)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, new)
	})
}


/*
	URL parameters:
	- Asset ID

	Response body data (all strings):
	- companyID (not necessary)
	- name
	- location
	- type
*/
func UpdateAsset() {
	http.HandleFunc("/user/assets/update", func(w http.ResponseWriter, r *http.Request) {

		user, err := CheckPostAdmin(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		assetID, err := strconv.Atoi(r.URL.Query()["asset_id"][0])
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		bytes, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		var asset db.Asset
		err = json.Unmarshal(bytes, &asset)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		asset.Index = assetID
		asset.CompanyID = user.EntityID

		err = db.UpdateAsset(assetID, asset)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		erpc.MarshalSend(w, asset)
	})
}


