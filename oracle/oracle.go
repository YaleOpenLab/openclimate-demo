package oracle

import (
	// "log"
	"github.com/YaleOpenLab/openclimate/blockchain"
	"github.com/YaleOpenLab/openclimate/ipfs"
	"github.com/pkg/errors"
	// "reflect"
)

// Functions clean the data and return it in the correct format.
// To verify, oracle will check if the methodology used is valid and
// if the values make sense.

func VerifyAtmosCO2(data interface{}) (interface{}, error) {
	dataConverted := data.([]GlobalCO2)

	dataSlice := make([]interface{}, len(dataConverted))
	for _, d := range dataConverted {
		dataSlice = append(dataSlice, d.Cycle)
	}

	dVal := dataValue(dataSlice)
	return dVal, nil
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

type BlockChainDataStruct struct {
	EntityType string
	EntityID int
	ReportType string
	Data interface{}
	IpfsHash string
}

// Calls the relevant verify helper-function to process the data,
// then commits the verified data to IPFS and then returns the hash.
func VerifyAndCommit(reportType string, entityType string, entityID int, data interface{}) error {

	var verifiedData interface{}
	var err error

	switch reportType {
	case "AtmosCO2":
		verifiedData, err = VerifyAtmosCO2(data)
	case "Emissions":
		verifiedData, err = VerifyEmissions(data)
	case "Mitigation":
		verifiedData, err = VerifyMitigation(data)
	case "Adaptation":
		verifiedData, err = VerifyAdaptation(data)
	default:
		return errors.New("Verification of this report type is not supported.")
	}

	// Committing to IPFS may not be necessary. We can commit this data
	// directly on to the blockchain if it is small enough. However, once
	// companies start to report a lot of data relating to their assets,
	// IPFS is needed to minimize the amount of blockchain storage required.

	ipfsHash, err := ipfs.IpfsCommitData(verifiedData)
	if err != nil {
		return errors.Wrap(err, "oracle.VerifyAndCommit() failed")
	}

	var bcds BlockChainDataStruct
	bcds.EntityType = entityType
	bcds.EntityID = entityID
	bcds.ReportType = reportType
	bcds.Data = verifiedData
	bcds.IpfsHash = ipfsHash

	err = blockchain.CommitToChain(bcds)
	if err != nil {
		return errors.Wrap(err, "oracle.VerifyAndCommit() failed")
	}

	return nil
}
