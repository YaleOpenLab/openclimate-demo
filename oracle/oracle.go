package oracle

import (
	"log"
	// "reflect"
)

// Functions clean the data and return it in the correct format, as determined by
// the structs defined in datastructs.go
// To verify, oracle will check if the methodology used is valid and if the
// values make sense.

func VerifyEmissions(data interface{}) (Emissions, error) {
	var verifiedData Emissions
	return verifiedData, nil
}

// func VerifyPledge(data interface{}) (Pledges, error) {
// 	var actorPledges Pledges

// 	log.Println(reflect.TypeOf(data))

// 	// actorPledges.UserID = data["UserID"]
// 	// actorPledges.EntityType = data["EntityType"]

// 	return actorPledges, nil
// }

func VerifyMitigation(data interface{}) (Mitigation, error) {
	var verifiedData Mitigation
	return verifiedData, nil
}

func VerifyAdaptation(data interface{}) (Adaptation, error) {
	var verifiedData Adaptation
	return verifiedData, nil
}

// Calls the relevant verify helper-function to process the data,
// then commits the data to IPFS and returns the hash
func Verify(data interface{}, reportType string) (string, error) {
	var ipfsHash string
	var err error

	var verifiedData interface{}
	switch reportType {
	case "Earth":
		verifiedData, err = VerifyEarth(data)
	case "Emissions":
		verifiedData, err = VerifyEmissions(data)
	// case "Pledges":
	// 	log.Println("hit")
	// 	verifiedData, err = VerifyPledge(data)
	case "Mitigation":
		verifiedData, err = VerifyMitigation(data)
	case "Adaptation":
		verifiedData, err = VerifyAdaptation(data)
	}

	if err != nil {
		log.Println("failed to verify data")
		return ipfsHash, err
	}

	ipfsHash, err = IpfsCommitData(verifiedData)
	if err != nil {
		log.Println("Failed to commit data to IPFS")
		return ipfsHash, err
	}

	return ipfsHash, err
}
