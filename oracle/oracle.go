package oracle

import (
	"log"
)

// Helper functions check if the methodology used is valid. Functions should
// also clean the data and return it in the correct format, as determined by
// the structs defined in datastructs.go

func VerifyEmissions(data map[string]string) (Emissions, error) {
	var verifiedData Emissions
	return verifiedData, nil
}

func VerifyPledge(data map[string]string) (Pledges, error) {
	var verifiedData Pledges
	return verifiedData, nil
}

func VerifyMitigation(data map[string]string) (Mitigation, error) {
	var verifiedData Mitigation
	return verifiedData, nil
}

func VerifyAdaptation(data map[string]string) (Adaptation, error) {

	var verifiedData Adaptation
	return verifiedData, nil
}

// Calls the relevant verify helper-function to process the data,
// then commits the data to IPFS and returns the hash
func Verify(data map[string]string, reportType string) (string, error) {

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
