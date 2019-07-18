package oracle

import (
	"encoding/json"
	"github.com/Varunram/essentials/ipfs"
	"log"
)

func IpfsCommitData(data interface{}) (string, error) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return ipfs.IpfsAddBytes(dataBytes)
}
