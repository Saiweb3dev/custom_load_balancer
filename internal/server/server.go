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

// Server represents the main load balancer server structure
type Server struct {
	config   *config.Config
	registry *registry.Registry
	balancer *balancer.Balancer
	pool     *pool.Pool
	health   *health.HealthChecker
	listener *listener.Listener
}

// New creates and initializes a new Server instance
// It takes a configuration object and sets up all the necessary components
func New(cfg *config.Config) *Server {
	reg := registry.New()
	return &Server{
		config:   cfg,
		registry: reg,
		balancer: balancer.New(reg),
		pool:     pool.New(),
		health:   health.New(reg),
		listener: listener.New(cfg.ListenAddr),
	}
}

// Start initializes the server components and begins the main server loop
// It returns an error if there's any issue during the startup process
func (s *Server) Start() error {
	log.Println("Starting load balancer...")
	// Initialize components and start the server
	// This is where you'd implement the main server loop
	return nil
}