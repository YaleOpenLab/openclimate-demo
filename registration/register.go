package main

import (
<<<<<<< HEAD
	"../database"
	"log"
	"strings"
=======
  "fmt"
  "log"
  "strings"
  "bufio"
  "os"
  queue "github.com/jupp0r/go-priority-queue"
>>>>>>> Created repl for finding user name
)

// Compute Levenshtein Edit Distance (LED) between w1 and w2
func LED(w1 string, w2 string) int {
<<<<<<< HEAD
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
=======
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
      if w1[i - 1] == w2[j - 1] {
        dif = 0
      } else {
        dif = 1
      }
      if dist[i - 1][j - 1] + dif < newval {
        newval = dist[i - 1][j - 1] + dif
      }
      if dist[i - 1][j] + 1 < newval {
        newval = dist[i - 1][j] + 1
      }
      if dist[i][j - 1] + 1 < newval {
        newval = dist[i][j - 1] + 1
      }
      dist[i][j] = newval
    }
  }

  return dist[l1 - 1][l2 - 1]
}

// Get the official user name corresponding to input "name"
func fetchUser(name string) (*queue.PriorityQueue, error) {
  //users, err := database.RetrieveAllUsers()
  users := [...]string{"New York Greater City Area", "New York, Co.", "New York City", "New York State"}
  pq := queue.New()
  for _, user := range users {
    if strings.Contains(user, name) {
      pq.Insert(user, float64(LED(name, user)))
      continue
    }
    dist := LED(name, user)
    if dist < 10 {
      pq.Insert(user, float64(dist))
    }
  }
  return &pq, nil
}

func main() {
  // Tests
  reader := bufio.NewReader(os.Stdin)
  fmt.Println("What account would you like to claim?")

  name, _ := reader.ReadString('\n')
  name = strings.Replace(name, "\n", "", -1)
  pq, err := fetchUser(name)
  if err != nil {
    log.Fatal(err)
  }
  for i := 0; i < 5; i++ {
    val, err := pq.Pop()
    if err != nil {
      fmt.Println("It looks like we have no record of this account...")
      fmt.Println("Let's register it now!")
      break
    }
    fmt.Println("Did you mean", val.(string), "?")
    reply, _ := reader.ReadString('\n')
    reply = strings.Replace(reply, "\n", "", -1)
    if reply == "y" {
      fmt.Println("Great! Let's continue.")
      break
    }
  }
>>>>>>> Created repl for finding user name
}
