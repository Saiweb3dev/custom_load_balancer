package server

import (
	"log"
	"simple_load_balancer/config"
	"simple_load_balancer/internal/registry"
	"simple_load_balancer/internal/balancer"
	"simple_load_balancer/internal/pool"
	"simple_load_balancer/internal/health"
	"simple_load_balancer/internal/listener"
)

type Server struct {
	config   *config.Config
	registry *registry.Registry
	balancer *balancer.Balancer
	pool     *pool.Pool
	health   *health.HealthChecker
	listener *listener.Listener
}

func New(cfg *config.Config) *Server {
	return &Server{
		config:   cfg,
		registry: registry.New(),
		balancer: balancer.New(),
		pool:     pool.New(),
		health:   health.New(),
		listener: listener.New(cfg.ListenAddr),
	}
}

func (s *Server) Start() error {
	log.Println("Starting load balancer...")
	// Initialize components and start the server
	// This is where you'd implement the main server loop
	return nil
}