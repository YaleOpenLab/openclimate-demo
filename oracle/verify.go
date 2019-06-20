package main

import (
	"fmt"
)

// Function type that returns a function to test whether data meets user-specified criteria
type UserDefinedIntegrity func(int) bool

// Examples of user-defined integrity statements
func valid_year(year int) bool {
	if year > 2019 {
		return false
	} else {
		return true
	}
}

// Given a set of numbers, score the quality of the data
func data_score(dat []int, params ...UserDefinedIntegrity) int {
	score := 0
	for _, param := range params {
		for _, num := range dat {
			if param(num) {
				score += 1
			} else {
				score -= 1
			}
		}
	}
	return score
}

// Given a set of data, determine the most accurate data value
func data_value(dat []int) int {
	return 0
}

func main() {
	// Test cases
	dat := []int{2011, 2010, 2019, 2020, 2040, 78}
	fmt.Println(data_score(dat, valid_year))
	fmt.Println(data_value(dat))
	fmt.Println(valid_year(2017))
	fmt.Println(valid_year(2020))
}
