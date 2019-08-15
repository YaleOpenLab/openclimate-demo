package oracle

import (
	// "log"
	"github.com/YaleOpenLab/openclimate/ipfs"
	"github.com/YaleOpenLab/openclimate/blockchain"
	"github.com/pkg/errors"
	// "reflect"
)

// Functions clean the data and return it in the correct format.
// To verify, oracle will check if the methodology used is valid and
// if the values make sense.

func VerifyAtmosCO2(data map[string][]float64) (interface{}, error) {
	var empty interface{}
	return empty, nil
}

func VerifyGlobalTemp() {

}

func VerifyEmissions(data interface{}) (ipfs.Emissions, error) {
	var verifiedData ipfs.Emissions
	return verifiedData, nil
}

func VerifyMitigation(data interface{}) (ipfs.Mitigation, error) {
	var verifiedData ipfs.Mitigation
	return verifiedData, nil
}

func VerifyAdaptation(data interface{}) (ipfs.Adaptation, error) {
	var verifiedData ipfs.Adaptation
	return verifiedData, nil
}

// Calls the relevant verify helper-function to process the data,
// then commits the verified data to IPFS and then returns the hash.
func Verify(reportType string, entityType string, entityID int, data map[string][]float64) error {

	var ipfsHash string
	var err error

	// Invoking a goroutine for report() so that the verification
	// function can run concurrently and doesn't hold up the server.
	go func() {

		var verifiedData interface{}
		
		switch reportType {
		case "AtmosCO2":
			verifiedData, err = VerifyAtmosCO2(data)
		case "Emissions":
			verifiedData, err = VerifyEmissions(data)
		case "Mitigation":
			verifiedData, err = VerifyMitigation(data)
		case "Adaptation":
			verifiedData, err = VerifyAdaptation(data)
		}
		if err != nil {
			err = errors.Wrap(err, "oracle.Verify() failed")
		}

		// Committing to IPFS may not be necessary. We can commit this data
		// directly on to the blockchain if it is small enough. However, once
		// companies start to report a lot of data relating to their assets,
		// IPFS is needed to minimize the amount of blockchain storage required.

		ipfsHash, err = ipfs.IpfsCommitData(verifiedData)
		if err != nil {
			err = errors.Wrap(err, "oracle.Verify() failed")
		}
	}()
	if err != nil {
		return errors.Wrap(err, "oracle.Verify() failed")
	}

	err = blockchain.CommitToChain(ipfsHash)
	if err != nil {
		return errors.Wrap(err, "oracle.Verify() failed")
	}

	return nil
}
