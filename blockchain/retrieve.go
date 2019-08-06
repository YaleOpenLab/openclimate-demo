package blockchain

/*

	GetFromIpfs() is not a real function yet. The function will receive the
	actor type and the actor id, then search our smart contract for all the
	IPFS hashes that are associated with that actor type and actor id. The
	function will then retrieve the corresponding data from IPFS using those
	hash content addresses and return it to us here.

	arguments:
	- "actorType": the type of the actor (company, city, region, country, etc.)
	- "actorID": the ID assigned to the actor in the database

	return type:
	A map that maps the type of the reported climate data to a struct
	of the data itself. For example, "emissions" would map to the actor's
	actual reported emissions data, formatted as a struct (to be defined).
*/

func GetFromIpfs(actorType string, actorID int) (map[string]interface{}, error) {
	var empty map[string]interface{}
	return empty, nil
}
