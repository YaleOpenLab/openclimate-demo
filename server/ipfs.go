package server

import (
	"log"
	"net/http"
	"strconv"

	// "github.com/Varunram/essentials/ipfs"
	eipfs "github.com/Varunram/essentials/ipfs"
	erpc "github.com/Varunram/essentials/rpc"
	"github.com/YaleOpenLab/openclimate/ipfs"
)

func setupIpfsHandlers() {
	RetrieveFromIpfs()
	RetrieveAllFromIpfs()
	getIpfsHash()
}

/*
	Request & retrieve data for a specific actor that has been committed to IPFS.
	HTTP request to our API will provide the actor type and actor id; then
	RetrieveFromIpfs() will look at the smart contract to retrieve all IPFS hashes
	related to that type/ID pair. It then queries IPFS for that data using the hashes
	and then makes that data available on the openclimate API.

	URL parameters:
	- "report_type": the type of climate action data that was reported. can be either
		emissions, mitigation, or adaptation.
	- "actor_type": either city, country, region, state, company, etc.
	- "actor_id": the ID assigned to the actor in the database.
*/
func RetrieveFromIpfs() {
	http.HandleFunc("/ipfs/retrieve", func(w http.ResponseWriter, r *http.Request) {

		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		reportType := r.URL.Query()["report_type"][0]
		actorType := r.URL.Query()["actor_type"][0]
		actorID, err := strconv.Atoi(r.URL.Query()["actor_id"][0])
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		/*
			blockchain.GetFromIpfs() is not a real function yet. The function
			will receive the actor type and the actor id, then search our smart
			contract for all the IPFS hashes that are associated with that actor
			type and actor id. The function will then retrieve the corresponding
			data from IPFS using those hash content addresses and give it to us here.

			For more information, see blockchain/retrieve.go
		*/
		data, err := ipfs.GetFromIpfs(reportType, actorType, actorID)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, data)
	})
}

/*
	Retrieve all data related to a given actor that has been committed to IPFS.
	This includes both emissions, mitigation, and adaptation data.
*/
func RetrieveAllFromIpfs() {
	http.HandleFunc("/ipfs/request", func(w http.ResponseWriter, r *http.Request) {

		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		// reportType := r.URL.Query()["report_type"][0]
		actorType := r.URL.Query()["actor_type"][0]
		actorID, err := strconv.Atoi(r.URL.Query()["actor_id"][0])
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		/*
			blockchain.GetFromIpfs() is not a real function yet. The function
			will receive the actor type and the actor id, then search our smart
			contract for all the IPFS hashes that are associated with that actor
			type and actor id. The function will then retrieve the corresponding
			data from IPFS using those hash content addresses and give it to us here.

			For more information, see blockchain/retrieve.go
		*/
		data, err := ipfs.GetAllFromIpfs(actorType, actorID)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, data)
	})
}

// getIpfsHash gets the ipfs hash of the passed string
func getIpfsHash() {
	http.HandleFunc("/ipfs/hash", func(w http.ResponseWriter, r *http.Request) {
		err := erpc.CheckGet(w, r)
		if err != nil {
			return
		}

		if !checkReqdParams(w, r, "string") {
			return
		}

		hashString := r.URL.Query()["string"][0]
		hash, err := eipfs.IpfsAddString(hashString)
		if err != nil {
			log.Println("did not add string to ipfs", err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		hashCheck, err := eipfs.IpfsGetString(hash)
		if err != nil || hashCheck != hashString {
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
			return
		}

		erpc.MarshalSend(w, hash)
	})
}
