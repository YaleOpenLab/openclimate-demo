package main

import (
  "log"
  "../database"
  "github.com/texttheater/golang-levenshtein/levenshtein"
)

func fetchUser(name string) (string, error) {
  user, err := database.RetrieveAllUsers()
  log.Println(users)
  return name, nil
}

func main() {
  // Tests
  log.Println(DistanceForStrings([]rune("a"), []rune("abb"), DefaultOptions)
  user, err := fetchUser("Target")
  if err != nil {
    log.Fatal(err)
  }
  log.Println(user)
}
