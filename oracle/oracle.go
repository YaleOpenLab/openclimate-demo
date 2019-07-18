package oracle

import (
	"log"

)

func VerifyEmissions(data map[string]string) (map[string]string, error) {
	return data, nil
}

func VerifyPledge(data map[string]string) (map[string]string, error) {
	return data, nil
}

func VerifyMitigation(data map[string]string) (map[string]string, error) {
	return data, nil
}

func VerifyAdaptation(data map[string]string) (map[string]string, error) {
	return data, nil
}

func Verify(data map[string]string, reportType string) (string, error) {

	var ipfsHash string
	var err error

	verifiedData := make(map[string]string)

	switch reportType {
		case "Emissions":
			verifiedData, err = VerifyEmissions(data)
			if err != nil {
				log.Println("failed verify emissions")
				return "", err
			}
		case "Pledges":	
			verifiedData, err = VerifyPledge(data)
			if err != nil {
				log.Println("failed verify pledges")
				return "", err
			}
		case "Mitigation":
			verifiedData, err = VerifyMitigation(data)
			if err != nil {
				log.Println("failed verify mitigation")
				return "", err				
			}
		case "Adaptation":
			verifiedData, err = VerifyAdaptation(data)
			if err != nil {
				log.Println("failed verify adaptation")
				return "", err
			}
		}	

	ipfsHash, err = IpfsCommitData(verifiedData)
	if err != nil {
		log.Println("Failed to commit data to IPFS")
		return "", err
	}

	return ipfsHash, err
}