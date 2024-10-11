package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
	"github.com/tsenart/vegeta/v12/lib"
)

type User struct {
	Name        string  `json:"name"`
	Salary      float64 `json:"salary"`
	Department  string  `json:"department"`
	Country     string  `json:"country"`
	Description string  `json:"description"`
}

func main() {
	rate := vegeta.Rate{Freq: 10000, Per: time.Minute}
	duration := 5 * time.Minute
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "POST",
		URL:    "http://localhost:8080/users",
	})
	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Big Bang!") {
		metrics.Add(res)
	}
	metrics.Close()

	fmt.Printf("99th percentile: %s\n", metrics.Latencies.P99)
	fmt.Printf("Requests: %d\n", metrics.Requests)
	fmt.Printf("Success rate: %.2f%%\n", metrics.Success*100)
}

func generateRandomUser() User {
	departments := []string{"IT", "HR", "Finance", "Marketing", "Sales"}
	countries := []string{"USA", "UK", "Canada", "Australia", "Germany"}

	return User{
		Name:        fmt.Sprintf("User%d", rand.Intn(1000)),
		Salary:      float64(30000 + rand.Intn(70000)),
		Department:  departments[rand.Intn(len(departments))],
		Country:     countries[rand.Intn(len(countries))],
		Description: fmt.Sprintf("Description for User%d", rand.Intn(1000)),
	}
}

func sendRequest(user User) {
	jsonUser, err := json.Marshal(user)
	if err != nil {
		log.Printf("Error marshaling user: %v", err)
		return
	}

	resp, err := http.Post("http://localhost:8080/users", "application/json", bytes.NewBuffer(jsonUser))
	if err != nil {
		log.Printf("Error sending request: %v", err)
		return
	}
	defer resp.Body.Close()

	log.Printf("Sent user: %s, Response status: %s", user.Name, resp.Status)
}