package main

import (
	"../database"
	"log"
	"strings"
)

// Compute Levenshtein Edit Distance (LED) between w1 and w2
func LED(w1 string, w2 string) int {
	w1 = strings.TrimSpace(w1)
	w2 = strings.TrimSpace(w2)
	if w1 == w2 {
		return 0
	}
	l1 := len(w1) + 1
	l2 := len(w2) + 1
	dist := make([][]int, l1)
	for i := 0; i < l1; i++ {
		dist[i] = make([]int, l2)
	}
	for i := 0; i < l1; i++ {
		dist[i][0] = i
	}
	for i := 0; i < l2; i++ {
		dist[0][i] = i
	}
	for j := 1; j < l2; j++ {
		for i := 1; i < l1; i++ {
			var newval int
			if l1 > l2 {
				newval = l1 + 1
			} else {
				newval = l2 + 1
			}
			var dif int
			if w1[i-1] == w2[j-1] {
				dif = 0
			} else {
				dif = 1
			}
			if dist[i-1][j-1]+dif < newval {
				newval = dist[i-1][j-1] + dif
			}
			if dist[i-1][j]+1 < newval {
				newval = dist[i-1][j] + 1
			}
			if dist[i][j-1]+1 < newval {
				newval = dist[i][j-1] + 1
			}
			dist[i][j] = newval
		}
	}
	log.Println(dist)
	return dist[l1-1][l2-1]
}

// Get the official user name corresponding to input "name"
func fetchUser(name string) (string, error) {
	users, err := database.RetrieveAllUsers()
	log.Println(users)
	return name, nil
}

func main() {
	// Tests
	user, err := fetchUser("Target")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(LED("aba", "bb"))
	log.Println(user)
}
