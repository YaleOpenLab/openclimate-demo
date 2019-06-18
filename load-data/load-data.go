package main

import (
  "io/ioutil"
  "os"
  "encoding/csv"
  "fmt"
  "log"
  "strings"
)

func create_csv_dictionary(path_to_file string) [][]string {
  file, err := os.Open(path_to_file)
  if err != nil {
    log.Fatal(err)
  }
  defer file.Close()

  in, err := ioutil.ReadAll(file)
  r := csv.NewReader(strings.NewReader(string(in)))
  records, err := r.ReadAll()
  if err != nil {
    log.Fatal(err)
  }

  return records
}

func main() {
  // Example call to create_csv_dictionary
  fmt.Println(create_csv_dictionary("data/country_definitions.csv"))
}
