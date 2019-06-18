package main

import (
  "io/ioutil"
  "os"
  "encoding/csv"
  "fmt"
  "log"
  "strings"
)

func main() {
  fmt.Println("Hello World!")

  file, err := os.Open("old/data/country_definitions.csv")
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

  fmt.Print(records[16][2])
}
