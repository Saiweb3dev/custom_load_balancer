package main

import (
	"encoding/json"
	"fmt"
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
		Method: "GET",
		URL:    "http://localhost:8080/users",
	})
	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics
	for res := range attacker.Attack(targeter, rate, duration, "Big Bang!") {
		metrics.Add(res)
		if res.Code == http.StatusOK {
			var user User
			json.Unmarshal(res.Body, &user)
			fmt.Printf("Retrieved user: %+v\n", user)
		}
	}
	metrics.Close()

	fmt.Printf("99th percentile: %s\n", metrics.Latencies.P99)
	fmt.Printf("Requests: %d\n", metrics.Requests)
	fmt.Printf("Success rate: %.2f%%\n", metrics.Success*100)
}