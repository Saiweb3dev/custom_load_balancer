package server

import (
	"log"
	"net"
	"time"
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
func New(cfg *config.Config) *Server {
	reg := registry.New("./registry_data.json") 
	bal := balancer.New(reg)
	listenerConfig := listener.Config{
		Address:     cfg.ListenAddr,
		TLSCertFile: cfg.TLSCertFile, // Add these to your config
		TLSKeyFile:  cfg.TLSKeyFile,  // Add these to your config
		IdleTimeout: 5 * time.Minute,
	}
	lis, err := listener.New(listenerConfig)
	if err != nil {
		log.Fatalf("Failed to create listener: %v", err)
	}
	poolConfig := pool.PoolConfig{
		MaxConns:        100,  // Adjust as needed
		IdleTimeout:     5 * time.Minute,
		MaxLifetime:     30 * time.Minute,
		CleanupInterval: 1 * time.Minute,
	}
	s := &Server{
		config:   cfg,
		registry: reg,
		balancer: bal,
		pool:     pool.New(poolConfig),
		health: health.New(
			reg,
			10*time.Second,  // Check interval
			5*time.Second,   // Timeout
			"/health",       // Health endpoint
		),
		listener: lis,
	}
	lis.SetHandler(s.handleConnection)
	return s
}

// Start initializes the server components and begins the main server loop
func (s *Server) Start() error {
	log.Println("Starting load balancer...")

	// Start the health checker
	s.health.Start()

	// Start periodic logging of server loads
	go s.logServerLoads()

	// Start the listener
	return s.listener.Start()
}

// handleConnection processes a single client connection
func (s *Server) handleConnection(clientConn net.Conn) {
	// Ensure the client connection is closed when we're done
	defer clientConn.Close()

	// Get the next available backend based on the load-aware algorithm
	backend := s.balancer.NextBackend()
	if backend == nil {
		log.Println("No available backend servers")
		return
	}
	// Get a connection to the selected backend from the connection pool
	backendConn, err := s.pool.Get(backend.Address)
if err != nil {
    if err == pool.ErrPoolExhausted {
        log.Printf("Connection pool exhausted for backend %s", backend.Address)
        // Handle the exhausted pool (e.g., return a 503 Service Unavailable to the client)
    } else {
        log.Printf("Failed to get connection to backend %s: %v", backend.Address, err)
    }
    return
}
defer s.pool.Put(backend.Address, backendConn)

	// Log the connection forwarding
	log.Printf("Forwarding connection from %s to backend %s", clientConn.RemoteAddr(), backend.Address)

	// In a real implementation, you would use io.Copy to forward data between connections
	// go io.Copy(backendConn, clientConn)
	// go io.Copy(clientConn, backendConn)
}

// logServerLoads periodically logs the current load of all servers
func (s *Server) logServerLoads() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		loads := s.balancer.GetServerLoads()
		log.Printf("Current server loads: %v", loads)
	}
}

// UpdateBackendLoad updates the load for a specific backend server
func (s *Server) UpdateBackendLoad(address string, load float64) {
	s.balancer.UpdateServerLoad(address, load)
}