package server

import (
	"github.com/pkg/errors"
	"log"
	"math/big"
	"net/http"

	// ipfs "github.com/Varunram/essentials/ipfs"
	erpc "github.com/Varunram/essentials/rpc"
	"github.com/YaleOpenLab/openclimate/database"
)

// Calls all database handlers
func setupUser() {
	newUser()
	retrieveUser()
	retrieveAllUsers()
	deleteUser()
	updateUser()

}

/*****************/
/* USER HANDLERS */
/*****************/

// setupPingHandler is a ping route for remote callers to check if the platform is up
func newUser() {
	http.HandleFunc("/user/new", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		if !checkReqdParams(w, r, "username", "pwhash", "email", "entity_type") {
			return
		}

		username := r.URL.Query()["username"][0]
		pwhash := r.URL.Query()["pwhash"][0]
		email := r.URL.Query()["email"][0]
		entityType := r.URL.Query()["entity_type"][0]
		entityName := r.URL.Query()["entity_name"][0]
		entityParent := r.URL.Query()["entity_parent"][0]

		user, err := database.NewUser(username, pwhash, email, entityType, entityName, entityParent)
		if err != nil {
			log.Println("couldn't create new user", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, user)
	})
}

// func newChild() {
// 	http.HandleFunc("/user/add/child", func(w http.ResponseWriter, r *http.Request) {
// 		err := erpc.CheckGet(w, r)
// 		if err != nil {
// 			return
// 		}

// 		if r.URL.Query()["child"] == nil {
// 			log.Println("required param child missing")
// 			erpc.ResponseHandler(w, erpc.StatusBadRequest)
// 			return
// 		}

// 		username := r.URL.Query()["username"][0]
// 		child := r.URL.Query()["child"][0]

// 		user, err := database.RetrieveUserByUsername(username)
// 		if err != nil {
// 			log.Println("failed to retrieve user, quitting")
// 			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
// 		}

// 		err = user.AddChild(child)
// 		if err != nil {
// 			log.Println("failed to add child, quitting")
// 			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
// 		}

// 		erpc.MarshalSend(w, user)
// 	})
// }

func CheckGetAuth(w http.ResponseWriter, r *http.Request) (database.User, error) {
	var user database.User
	err := erpc.CheckGet(w, r)
	if err != nil {
		return user, errors.Wrap(err, "could not checkgetauth")
	}

	if !checkReqdParams(w, r, "username", "pwhash") {
		return user, nil
	}

	username := r.URL.Query()["username"][0]
	pwhash := r.URL.Query()["pwhash"][0]

	user, err = database.ValidateUser(username, pwhash)
	if err != nil {
		log.Println("could not retrieve user from the database, quitting")
		erpc.ResponseHandler(w, erpc.StatusBadRequest)
		return user, errors.New("user not found in database, quitting")
	}
	return user, nil
}

func CheckPostAuth(w http.ResponseWriter, r *http.Request) (database.User, error) {
	var user database.User
	err := erpc.CheckPost(w, r)
	if err != nil {
		return user, errors.Wrap(err, "could not checkpostauth")
	}

	if !checkReqdParams(w, r, "username", "pwhash") {
		return user, nil
	}

	username := r.URL.Query()["username"][0]
	pwhash := r.URL.Query()["pwhash"][0]

	user, err = database.ValidateUser(username, pwhash)
	if err != nil {
		erpc.ResponseHandler(w, erpc.StatusBadRequest)
		return user, errors.New("user not found in database, quitting")
	}

	if user.Verified == false {
		erpc.ResponseHandler(w, erpc.StatusBadRequest)
		return user, errors.New("user is not verified for this entity, quitting")
	}

	return user, nil
}

func CheckPostAdmin(w http.ResponseWriter, r *http.Request) (database.User, error) {

	user, err := CheckPostAuth(w, r)
	if err != nil {
		return user, errors.Wrap(err, "CheckPostAdmin failed while calling CheckPostAuth")
	}

	if !user.Admin {
		return user, errors.New("User is not an admin.")
	}

	return user, nil
}

func retrieveUser() {
	http.HandleFunc("/user/retrieve", func(w http.ResponseWriter, r *http.Request) {
		user, err := CheckGetAuth(w, r)
		if err != nil {
			return
		}

		erpc.MarshalSend(w, user)
	})
}

func retrieveAllUsers() {
	http.HandleFunc("/user/retrieve/all", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		users, err := database.RetrieveAllUsers()
		if err != nil {
			log.Println("could not retrieve user from the database, quittting")
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, users)
	})
}

func deleteUser() {
	http.HandleFunc("/user/delete", func(w http.ResponseWriter, r *http.Request) {
		user, err := CheckGetAuth(w, r)
		if err != nil {
			return
		}

		err = database.DeleteKeyFromBucket(user.Index, database.UserBucket)
		if err != nil {
			log.Println("could not delete user from database, quittting", err)
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

func updateUser() {
	http.HandleFunc("/user/update", func(w http.ResponseWriter, r *http.Request) {
		user, err := CheckGetAuth(w, r)
		if err != nil {
			return
		}

		if r.URL.Query()["email"] != nil {
			user.Email = r.URL.Query()["email"][0]
		} else if r.URL.Query()["newpwhash"] != nil {
			user.Pwhash = r.URL.Query()["newpwhash"][0]
		} else if r.URL.Query()["newusername"] != nil {
			user.Username = r.URL.Query()["newusername"][0]
		} else {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		err = user.Save()
		if err != nil {
			log.Println("error while savingt user to database")
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, user)
	})
}

/*********************/
/* ETHEREUM HANDLERS */
/*********************/

func sendEth() {
	http.HandleFunc("/user/sendeth", func(w http.ResponseWriter, r *http.Request) {
		user, err := CheckGetAuth(w, r)
		if err != nil {
			return
		}

		if !checkReqdParams(w, r, "address", "amount") {
			return
		}

		address := r.URL.Query()["address"][0]
		amountStr := r.URL.Query()["amount"][0] // convert this to bigint

		var amount big.Int
		_, boolErr := amount.SetString(amountStr, 10)
		if !boolErr {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		txhash, err := user.SendEthereumTx(address, amount)
		if err != nil {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		log.Println("user: ", user.Username, "has sent tx with txhash: ", txhash)
		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

// func getAllRequests() {
// 	http.HandleFunc("/requests/all", func(w http.ResponseWriter, r *http.Request) {
// 		err := erpc.CheckGet(w, r)
// 		if err != nil {
// 			return
// 		}

// 		requests, err := database.RetrieveAllRequests()
// 		if err != nil {
// 			log.Println("error while retrieving all requests, quitting")
// 			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
// 			return
// 		}

// 		erpc.MarshalSend(w, requests)
// 	})
// }
