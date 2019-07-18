package oracle

import (
	"log"
	"encoding/json"
	"github.com/Varunram/essentials/ipfs"
)


func IpfsCommitData(data interface{}) (string, error) {

	dataBytes, err := json.Marshal(data)
	hash, err := ipfs.IpfsAddBytes(dataBytes)
	if err != nil {
		log.Println(err)
	}
	return hash, err
}

