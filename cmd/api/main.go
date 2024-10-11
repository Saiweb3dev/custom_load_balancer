package main

import (
	"log"
	"simple_load_balancer/internal/server"
	"simple_load_balancer/config"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create and start the load balancer server
	s := server.New(cfg)
	if err := s.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}