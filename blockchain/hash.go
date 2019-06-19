package blockchain

import (
	"log"
)

func CommitToChain(chain string, args ...string) (string, error) {
	log.Println("committing args to chain")

	var inputString string

	for _, strings := range args {
		inputString += strings
	}

	log.Println(inputString)

	if chain == "ethereum" {
		log.Println("built on eth")
	}
	return inputString, nil
}
