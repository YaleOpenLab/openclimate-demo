package server

import (
	"github.com/YaleOpenLab/openclimate/database"
	"log"
	"net/http"
)

func setupDBHandlers() {
	newUser()
	retrieveUser()
	retrieveAllUsers()
	deleteUser()
	updateUser()
}

// setupPingHandler is a ping route for remote callers to check if the platform is up
func newUser() {
	http.HandleFunc("/user/new", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)

		if r.URL.Query()["username"] == nil || r.URL.Query()["pwhash"] == nil || r.URL.Query()["email"] == nil {
			log.Println("required params - username, pwhash, email missing")
			log.Println(r.URL.Query()["username"])
			log.Println(r.URL.Query()["pwhash"])
			log.Println(r.URL.Query()["email"])
			responseHandler(w, StatusBadRequest)
			return
		}

		username := r.URL.Query()["username"][0]
		pwhash := r.URL.Query()["pwhash"][0]
		email := r.URL.Query()["email"][0]

		user, err := database.NewUser(username, pwhash, email)
		if err != nil {
			log.Println("couldn't create new user", err)
			responseHandler(w, StatusInternalServerError)
			return
		}

		MarshalSend(w, user)
	})
}

func authorizeUser(r *http.Request) bool {
	username := r.URL.Query()["username"][0]
	pwhash := r.URL.Query()["pwhash"][0]

	return database.AuthUser(username, pwhash)
}

func retrieveUser() {
	http.HandleFunc("/user/retrieve", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)

		if !authorizeUser(r) {
			log.Println("user has invalid credentials")
			responseHandler(w, StatusUnauthorized)
			return
		}

		username := r.URL.Query()["username"][0]
		pwhash := r.URL.Query()["pwhash"][0]

		user, err := database.RetrieveUser(username, pwhash)
		if err != nil {
			log.Println("could not retrieve user from the database, quittting")
			responseHandler(w, StatusInternalServerError)
			return
		}

		MarshalSend(w, user)
	})
}

func retrieveAllUsers() {
	http.HandleFunc("/user/retrieve/all", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)

		if !authorizeUser(r) {
			log.Println("user has invalid credentials")
			responseHandler(w, StatusUnauthorized)
			return
		}

		users, err := database.RetrieveAllUsers()
		if err != nil {
			log.Println("could not retrieve user from the database, quittting")
			responseHandler(w, StatusInternalServerError)
			return
		}

		MarshalSend(w, users)
	})
}

func deleteUser() {
	http.HandleFunc("/user/delete", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)

		if !authorizeUser(r) {
			log.Println("user has invalid credentials")
			responseHandler(w, StatusUnauthorized)
			return
		}

		username := r.URL.Query()["username"][0]
		pwhash := r.URL.Query()["pwhash"][0]

		err := database.DeleteUser(username, pwhash)
		if err != nil {
			log.Println("could not delete user from database, quittting", err)
			responseHandler(w, StatusBadRequest)
			return
		}

		responseHandler(w, StatusOK)
	})
}

func updateUser() {
	http.HandleFunc("/user/update", func(w http.ResponseWriter, r *http.Request) {
		checkGet(w, r)
		checkOrigin(w, r)

		if !authorizeUser(r) {
			log.Println("user has invalid credentials")
			responseHandler(w, StatusUnauthorized)
			return
		}

		username := r.URL.Query()["username"][0]
		pwhash := r.URL.Query()["pwhash"][0]

		user, err := database.RetrieveUser(username, pwhash)
		if err != nil {
			log.Println("could not retrieve user from the database, quittting")
			responseHandler(w, StatusInternalServerError)
			return
		}

		if r.URL.Query()["email"] != nil {
			user.Email = r.URL.Query()["email"][0]
		} else if r.URL.Query()["newpwhash"] != nil {
			user.Pwhash = r.URL.Query()["newpwhash"][0]
		} else if r.URL.Query()["newusername"] != nil {
			user.Name = r.URL.Query()["newusername"][0]
		} else {
			responseHandler(w, StatusBadRequest)
			return
		}

		err = user.Save()
		if err != nil {
			log.Println("error while savingt user to database")
			responseHandler(w, StatusInternalServerError)
			return
		}

		MarshalSend(w, user)
	})
}
