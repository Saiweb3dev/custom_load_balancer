package main

import (
	"log"
	"path/filepath"
	"os"
	"simple_load_balancer/config"
	"simple_load_balancer/internal/server"
)

func main() {
	cwd, err := os.Getwd()
	if err != nil {
			log.Fatalf("Failed to get current working directory: %v", err)
	}
	log.Printf("Current working directory: %s", cwd)

	// Get the absolute path to the config file
	configPath, err := filepath.Abs("config/config.json")
	if err != nil {
			log.Fatalf("Failed to get absolute path to config file: %v", err)
	}
	log.Printf("Config file path: %s", configPath)

	cfg, err := config.Load(configPath)
	if err != nil {
			log.Fatalf("Failed to load configuration: %v", err)
	}

	s := server.New(cfg)
	if err := s.Start(); err != nil {
			log.Fatalf("Server failed to start: %v", err)
	}
}