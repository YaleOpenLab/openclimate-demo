package main

import (
	"bytes"
	"encoding/csv"
	"io/ioutil"
	"log"
)

func createCsvDictionary(filePath string) ([][]string, error) {
	var records [][]string
	byteData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return records, err
	}

	csvReader := csv.NewReader(bytes.NewReader(byteData))
	records, err = csvReader.ReadAll()
	if err != nil {
		return records, err
	}

	return records, nil
}

func main() {
	// Example call to createCsvDictionary
	records, err := createCsvDictionary("data/country_definitions.csv")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(records)
}
