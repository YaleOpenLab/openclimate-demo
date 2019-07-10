package server

import (
	"github.com/pkg/errors"
	"log"
	"math/big"
	"net/http"

	ipfs "github.com/Varunram/essentials/ipfs"
	erpc "github.com/Varunram/essentials/rpc"
	"github.com/YaleOpenLab/openclimate/database"
)

func setupDBHandlers() {
	newUser()
	retrieveUser()
	retrieveAllUsers()
	deleteUser()
	updateUser()
	getIpfsHash()
	getAllCompanies()
	getCompany()
	getAllRegions()
	getRegion()
	getAllCities()
	getCity()
	newRegion("Connecticut", "USA")
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

		if r.URL.Query()["username"] == nil || r.URL.Query()["pwhash"] == nil || r.URL.Query()["email"] == nil {
			log.Println("required params - username, pwhash, email missing")
			log.Println(r.URL.Query()["username"])
			log.Println(r.URL.Query()["pwhash"])
			log.Println(r.URL.Query()["email"])
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		username := r.URL.Query()["username"][0]
		pwhash := r.URL.Query()["pwhash"][0]
		email := r.URL.Query()["email"][0]

		user, err := database.NewUser(username, pwhash, email)
		if err != nil {
			log.Println("couldn't create new user", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, user)
	})
}

func CheckGetAuth(w http.ResponseWriter, r *http.Request) (database.User, error) {
	var user database.User
	err := erpc.CheckGet(w, r)
	if err != nil {
		return user, errors.Wrap(err, "could not checkgetauth")
	}

	if r.URL.Query()["username"] == nil || r.URL.Query()["pwhash"] == nil {
		log.Println("missing params in call")
		erpc.ResponseHandler(w, erpc.StatusBadRequest)
		return user, errors.New("missing params in call")
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
			user.Name = r.URL.Query()["newusername"][0]
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

/*****************/
/* IPFS HANDLERS */
/*****************/

// getIpfsHash gets the ipfs hash of the passed string
func getIpfsHash() {
	http.HandleFunc("/ipfs/hash", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		if r.URL.Query()["string"] == nil {
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		hashString := r.URL.Query()["string"][0]
		hash, err := ipfs.AddStringToIpfs(hashString)
		if err != nil {
			log.Println("did not add string to ipfs", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		hashCheck, err := ipfs.GetStringFromIpfs(hash)
		if err != nil || hashCheck != hashString {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, hash)
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

		if r.URL.Query()["address"] == nil || r.URL.Query()["amount"] == nil {
			log.Println("address or amount missing, quitting")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
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

		log.Println("user: ", user.Name, "has sent tx with txhash: ", txhash)
		erpc.ResponseHandler(w, erpc.StatusOK)
	})
}

/*******************/
/* REGION HANDLERS */
/*******************/

func getAllRegions() {
	http.HandleFunc("/region/all", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		regions, err := database.RetrieveAllRegions()
		if err != nil {
			log.Println("Error while retrieving all regions, quitting")
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, regions)
	})
}

func getRegion() {
	http.HandleFunc("/region", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		if r.URL.Query()["region_name"] == nil || r.URL.Query()["region_country"] == nil {
			log.Println("Region_name or region_country not passed, quitting")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
		}

		name := r.URL.Query()["region_name"][0]
		country := r.URL.Query()["region_country"][0]
		region, err := database.RetrieveRegionByName(name, country) //************ STOP ***********
		if err != nil {
			log.Println("Error while retrieving all regions, quitting")
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, region)
	})
}

/*****************/
/* CITY HANDLERS */
/*****************/

func getAllCities() {
	http.HandleFunc("/city/all", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		cities, err := database.RetrieveAllCities()
		if err != nil {
			log.Println("Error while retrieving all cities, quitting")
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, cities)
	})
}

func getCity() {
	http.HandleFunc("/city", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		if r.URL.Query()["city_name"] == nil || r.URL.Query()["city_region"] == nil {
			log.Println("City name or city region not passed, quitting")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
		}

		name := r.URL.Query()["city_name"][0]
		region := r.URL.Query()["city_region"][0]
		city, err := database.RetrieveCityByName(name, region) //************ STOP ***********
		if err != nil {
			log.Println("Error while retrieving all cities, quitting")
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, city)
	})
}

/********************/
/* COMPANY HANDLERS */
/********************/

func getAllCompanies() {
	http.HandleFunc("/companies/all", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		companies, err := database.RetrieveAllCompanies()
		if err != nil {
			log.Println("error while retrieving all companies, quitting")
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, companies)
	})
}

func getCompany() {
	http.HandleFunc("/company", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		if r.URL.Query()["company_name"] == nil || r.URL.Query()["company_country"] == nil {
			log.Println("company name or country not passed, quitting")
			erpc.ResponseHandler(w, erpc.StatusBadRequest)
			return
		}

		name := r.URL.Query()["company_name"][0]
		country := r.URL.Query()["company_country"][0]
		company, err := database.RetrieveCompanyByName(name, country)
		if err != nil {
			log.Println("error while retrieving all companies, quitting")
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, company)
	})
}
