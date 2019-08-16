package ipfs

import (
	"encoding/json"
	"github.com/Varunram/essentials/ipfs"
	"log"
)

// Commits data to IPFS and returns IPFS hash
func IpfsCommitData(data interface{}) (string, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return ipfs.IpfsAddBytes(dataBytes)
}
