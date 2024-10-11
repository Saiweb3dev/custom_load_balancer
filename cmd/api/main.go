package main

import (
	"log"
	"simple_load_balancer/config"
	"simple_load_balancer/internal/server"
)

func main() {
	cfg, err := config.Load("config.json")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	s := server.New(cfg)
	if err := s.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}