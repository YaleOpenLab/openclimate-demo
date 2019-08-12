package server

import (
	// "github.com/pkg/errors"
	"log"
	// "math/big"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	erpc "github.com/Varunram/essentials/rpc"
	db "github.com/YaleOpenLab/openclimate/database"
	"github.com/YaleOpenLab/openclimate/ipfs"
)

func setupManage() {
	VerifyUser()
	AddAsset()
	UpdateAsset()
	AddPledge()
	UpdatePledge()
	CommitPledge()
}

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
		state := asset["state"]
		assetType := asset["type"]

		new, err := db.NewAsset(name, companyID, location, state, assetType)
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

func AddPledge() {
	http.HandleFunc("/user/pledges/add", func(w http.ResponseWriter, r *http.Request) {
		user, err := CheckPostAuth(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		// Read the data in the body of the response into bytes,
		// which will be parsed into a map[string]string.
		bytes, err := ioutil.ReadAll(r.Body)
		defer r.Body.Close()
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		// Unmarshal the json bytes into the map[string]string data type
		// to be parsed for the arguments to create a new Pledge item and
		// inserted into the pledge bucket.
		var pledge map[string]string
		err = json.Unmarshal(bytes, &pledge)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		// Convert string data into the correct data type to be passed
		// to the NewPledge() function, in order to create the new pledge.

		actorID := user.EntityID

		pledgeType := pledge["pledge_type"]

		baseYear, err := strconv.Atoi(pledge["base_year"])
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		targetYear, err := strconv.Atoi(pledge["target_year"])
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		goal, err := strconv.ParseFloat(pledge["goal"], 64)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		regulatory, err := strconv.ParseBool(pledge["regulatory"])
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		// Call NewPledge() with all the arguments, which have been typecasted
		// into the proper types required by the NewPledge function.
		new, err := db.NewPledge(pledgeType, baseYear, targetYear, goal, regulatory, actorID)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		// Get the entity that the user is associated with from the database
		// so that we can update its Pledges field to include the new pledge.
		entity, err := user.RetrieveUserEntity()
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		// Add the ID of the new pledge to the entity's Pledges field.
		err = entity.AddPledges(new.ID)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		// TODO: Convert pledge into smart contract condition
		erpc.MarshalSend(w, new)
	})
}

func UpdatePledge() {
	http.HandleFunc("/user/pledges/update", func(w http.ResponseWriter, r *http.Request) {

		user, err := CheckPostAdmin(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		pledgeID, err := strconv.Atoi(r.URL.Query()["pledge_id"][0])
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

		var pledge db.Pledge
		err = json.Unmarshal(bytes, &pledge)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		pledge.ID = pledgeID
		pledge.ActorID = user.EntityID

		err = db.UpdatePledge(pledgeID, pledge)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		erpc.MarshalSend(w, pledge)
	})
}

func CommitPledge() {
	http.HandleFunc("user/pledges/commit", func(w http.ResponseWriter, r *http.Request) {
		_, err := CheckGetAuth(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		if r.URL.Query()["pledge_ID"] == nil {
			log.Println("pledge ID not passed, quitting")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		pledgeID, err := strconv.Atoi(r.URL.Query()["pledge_ID"][0])
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}

		pledge, err := db.RetrievePledge(pledgeID)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		ipfsHash, err := ipfs.IpfsCommitData(pledge)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, ipfsHash)
	})
}

func UpdateMRV() {
	http.HandleFunc("user/mrv", func(w http.ResponseWriter, r *http.Request) {
		user, err := CheckGetAuth(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		if r.URL.Query()["MRV"] == nil {
			log.Println("Updated MRV not passed, quitting")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		mrv := r.URL.Query()["MRV"][0]

		actor, err := user.RetrieveUserEntity()
		actor.UpdateMRV(mrv)

		erpc.MarshalSend(w, mrv)
	})
}

