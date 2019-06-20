package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// Function type that returns a function to test whether data meets user-specified criteria
type UserDefinedIntegrity func(int) bool

// Variance allowed for data
var sigma float64 = 1.5

// Examples of user-defined integrity statements
func valid_record_year(year int) bool {
	// For a record to be valid, it must be measured before the current time
	dt := time.Now()
	current_year, _ := strconv.Atoi(strings.Split(dt.Format("01-02-2006"), "-")[2])

	// Check that the record year is smaller than the current year
	if year > current_year {
		return false
	} else {
		return true
	}
}

// Given a set of numbers, score the quality of the data from 0 to 5
func data_score(dat []int, params ...UserDefinedIntegrity) float64 {
	avg := average(dat)
	stdev := stdev(dat)
	n := float64(len(dat))
	p := float64(len(params))
	score := 0.0
	for _, num := range dat {
		for _, param := range params {
			if param(num) {
				score += (1.0 / (n * p))
			}
		}

		if math.Abs(float64(num)-avg) < sigma*stdev {
			score += (1.0 / n)
		}
	}

	return score * 5.0 / (p + 1.0)
}

// Given a set of data, determine the most accurate data value
func data_value(dat []int) int {
	// Take average of values within 1 standard deviation of the mean
	avg := average(dat)
	stdev := stdev(dat)
	val := 0.0
	ctr := 0.0
	for _, num := range dat {
		if math.Abs(float64(num)-avg) < sigma*stdev {
			val += float64(num)
			ctr += 1.0
		}
	}
	val = val / ctr

	return int(val)
}

// Given a set of data, determine the precision of the data
func stdev(dat []int) float64 {
	avg := average(dat)
	res := 0.0
	for _, num := range dat {
		res += math.Pow(avg-float64(num), 2)
	}
	res = res / float64(len(dat))
	res = math.Sqrt(res)

	return res
}

// Computes the average of a set of data
func average(dat []int) float64 {
	avg := 0.0
	for _, num := range dat {
		avg += float64(num)
	}
	avg = avg / float64(len(dat))

	return avg
}

func main() {
	// Example use case: dat holds measurements of temperature taken from IoT devices
	// Devices are in the same general region and measurements are taken within the same hour
	dat := []int{0, 70, 71, 72, 71, 75, 89, 71}
	fmt.Println(data_value(dat)) // 74, which is approximately the average temperature from these
	fmt.Println(data_score(dat)) // 4.4 out of 5, which makes sense because only 2 values are "off"
}
