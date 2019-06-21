package main

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

// Variance allowed for data
var sigma float64 = 1.5

// Function type that returns a bool to test whether data meets user-specified criteria
type UserDefinedIntegrity func(interface{}) bool

// Examples of user-defined integrity functions
func valid_record_year(year interface{}) bool {
  switch y := year.(type) {
  case int:
	  // For a record to be valid, it must be measured before the current time
	  dt := time.Now()
	  current_year, _ := strconv.Atoi(strings.Split(dt.Format("01-02-2006"), "-")[2])

	  // Check that the record year is smaller than the current year
	  if y <= current_year {
		  return true
	  } else {
      return false
    }
  default:
    return false
  }
}

func valid_record_temperature(temp interface{}) bool {
  switch t := temp.(type) {
  case float64:
    return t > -50.0 && t < 150.0
  default:
    return false
  }
}

// Convert []interface to []float64
func interface_to_float64(dat []interface{}) []float64 {
  ret := make([]float64, len(dat))
  for i, val := range dat {
    switch v := val.(type) {
    case int:
      ret[i] = float64(v)
    case float64:
      ret[i] = v
    default:
      ret[i] = 0.0
    }
  }

  return ret
}

// Given a set of numbers, score the quality of the data from 0 to 5
func data_score(dat []interface{}, params ...UserDefinedIntegrity) float64 {
  d := interface_to_float64(dat)
  avg := average(d)
  stdev := stdev(d)
	l := float64(len(dat))
	p := float64(len(params))
	score := 0.0
	for _, num := range dat {
	  for _, param := range params {
      if param(num) {
        score += (1.0 / (l * p))
		  }
	  }

	  if math.Abs(num.(float64)-avg) < sigma*stdev {
		  score += (1.0 / l)
    }
	}

	return score * 5.0 / (p + 1.0)
}

// Given a set of data, determine the most accurate data value
func data_value(dat []interface{}) float64 {
  d := interface_to_float64(dat)
	// Take average of values within 1 standard deviation of the mean
	avg := average(d)
	stdev := stdev(d)
	val := 0.0
	ctr := 0.0
	for _, num := range d {
		if math.Abs(num-avg) < sigma*stdev {
			val += num
			ctr += 1.0
		}
	}
	val = val / ctr

	return val
}

// Computes the precision of a set of data
func stdev(dat []float64) float64 {
	avg := average(dat)
	res := 0.0
	for _, num := range dat {
		res += math.Pow(avg-num, 2)
	}
	res = res / float64(len(dat))
	res = math.Sqrt(res)

	return res
}

// Computes the average of a set of data
func average(dat []float64) float64 {
	avg := 0.0
	for _, num := range dat {
		avg += num
	}
	avg = avg / float64(len(dat))

	return avg
}

func main() {
	// Example use case: dat holds measurements of temperature taken from IoT devices
	// Devices are in the same general region and measurements are taken within the same hour
	dat := []interface{}{0.0, 70.1, 70.2, 74.0, 75.1, 73.3, 72.0, 71.6, 75.0, 89.0, 71.1, 200.3}
	fmt.Println(data_value(dat)) // 74, which is approximately the average temperature from these
	fmt.Println(data_score(dat, valid_record_temperature)) // 4.4 out of 5, which makes sense because only 2 values are "off"
}
