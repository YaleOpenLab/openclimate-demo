package oracle

import (
	"github.com/YaleOpenLab/openclimate/blockchain"
	"github.com/YaleOpenLab/openclimate/ipfs"
	"github.com/pkg/errors"
)

// Struct defining the data scheme of all reported climate action data
// when it is stored on Ethereum. Contains all metadata necessary to
// locate the correct data and statistic and ensure that it is attached
// to the correct climate actor.
type BlockChainDataStruct struct {
	EntityType string
	EntityID int
	ReportType string
	DataVal float64
	IpfsHash string
}


// Reads data in the form of an array of GlobalCO2 structs (a struct
// that, along with the atmospheric CO2 data measurements themselves,
// also holds metadeta) and computes the "true value".
func VerifyAtmosCO2(data []GlobalCO2) ([]GlobalCO2, float64, error) {

	dataSlice := make([]interface{}, len(data))
	for _, d := range data {
		dataSlice = append(dataSlice, d.Cycle)
	}

	dVal := dataValue(dataSlice)
	return data, dVal, nil
}


// *** TODO: find global temperature data ***
// Reads data in the form of an array of GlobalCO2 structs (a struct
// that, along with the atmospheric CO2 data measurements themselves,
// also holds metadeta) and computes the "true value".
func VerifyGlobalTemp(data []GlobalTemp) ([]GlobalTemp, float64, error) {
	var temp float64
	return data, temp, nil
}


// VerifyAndCommit receives data and depending on what kind of data it is,
// sends it to a helper function to verify the data and compute the "true value".
// Next, it commits the verified data itself to IPFS, receives the IPFS hash,
// then commits the hash and the computed statistic to Ethereum.
func VerifyAndCommit(reportType string, entityType string, entityID int, data interface{}) error {

	var verifiedData interface{}
	var dataVal float64
	var err error

	switch reportType {

	case "Atmospheric CO2":
		verifiedData, dataVal, err = VerifyAtmosCO2(data.([]GlobalCO2))

	case "Global Temperature":
		verifiedData, dataVal, err = VerifyGlobalTemp(data.([]GlobalTemp))

	default:
		return errors.New("Verification of this report type is not supported.")
	}

	// Committing to IPFS may not be necessary. We can commit this data
	// directly on to the blockchain if it is small enough. However, once
	// companies start to report a lot of data relating to their assets,
	// IPFS is needed to minimize blockchain storage overhead required.
	// Here, we commit to IPFS and store the hash on the blockchain to
	// demonstrate the concept.

	ipfsHash, err := ipfs.IpfsCommitData(verifiedData)
	if err != nil {
		return errors.Wrap(err, "oracle.VerifyAndCommit() failed")
	}

	// Marshal the data into the BlockChainDataStruct to ensure we have all
	// the metadata we need to locate the correct IPFS hash and statistic
	// (for display on the front-end).

	var bcds BlockChainDataStruct
	bcds.EntityType = entityType
	bcds.EntityID = entityID
	bcds.ReportType = reportType
	bcds.DataVal = dataVal
	bcds.IpfsHash = ipfsHash

	err = blockchain.CommitToChain(bcds)
	if err != nil {
		return errors.Wrap(err, "oracle.VerifyAndCommit() failed")
	}

	return nil
}
