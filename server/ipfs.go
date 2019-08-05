package server

import (

	"log"
	"net/http"

	erpc "github.com/Varunram/essentials/rpc"

)




/* 
	Request & retrieve data for a specific actor that has been committed to IPFS.
	HTTP request to our API will provide the actor type and actor id; then
	RetrieveFromIpfs() will look at the smart contract to retrieve all IPFS hashes 
	related to that type/ID pair. It then queries IPFS for that data using the hashes
	and then makes that data available on the openclimate API.
	
	URL parameters:
	- "actor_type":
	- "actor_id":

*/
func RetrieveFromIpfs() {
	http.HandleFunc("/ipfs/request", func(w http.ResponseWriter, r *http.Request) {

		err := erpc.CheckGet(w, r)
		if err != nil {
			log.Println(err)
			erpc.ResponseHandler(w, erpc.StatusInternalServerError)
		}


	})
}