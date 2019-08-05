package server

import (
	// "github.com/pkg/errors"
	"log"
	// "math/big"
	"net/http"
	"strconv"

	// ipfs "github.com/Varunram/essentials/ipfs"
	erpc "github.com/Varunram/essentials/rpc"
	"github.com/YaleOpenLab/openclimate/database"
)


// URL parameter: candidate ID
func VerifyUser() {
	http.HandleFunc("/user/admin/verify", func(w http.ResponseWriter, r *http.Request) {

		var candidate database.User

		_, err := CheckPostAdmin(w, r) // Check if the person is an admin/has authority to verify users
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		if r.URL.Query()["candidate"] == nil {
			log.Println("Candidate for verification not specified")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		id, err := strconv.Atoi(r.URL.Query()["candidate"][0])
		if err != nil {
			log.Println("Failed to typecast candidate id from string to int")
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		candidate, err = database.RetrieveUser(id)
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
