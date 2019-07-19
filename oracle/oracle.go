package oracle

import (
	"log"
)

// Helper functions check if the methodology used is valid. Functions should
// also clean the data and return it in the correct format, as determined by
// the structs defined in datastructs.go

func VerifyEmissions(data interface{}) (Emissions, error) {
	var verifiedData Emissions
	return verifiedData, nil
}

func VerifyPledge(data interface{}) (Pledges, error) {
	var verifiedData Pledges
	return verifiedData, nil
}

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
	case "Emissions":
		verifiedData, err = VerifyEmissions(data)
	case "Pledges":
		verifiedData, err = VerifyPledge(data)
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

	// COMMIT TO CHAIN

	





	return ipfsHash, err
}
