package server

import (
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5"
	"go.mongodb.org/mongo-driver/mongo"
		
	"simple_load_balancer/config"
	"simple_load_balancer/internal/registry"
	"simple_load_balancer/internal/balancer"
	"simple_load_balancer/internal/pool"
	"simple_load_balancer/internal/health"
	"simple_load_balancer/internal/listener"
	"simple_load_balancer/internal/database"
	controller "simple_load_balancer/internal/controller"
)

// Server represents the main load balancer server structure
type Server struct {
	config   *config.Config
	registry *registry.Registry
	balancer *balancer.Balancer
	pool     *pool.Pool
	health   *health.HealthChecker
	listener *listener.Listener
	router   *chi.Mux
	db       *mongo.Database
}

// New creates and initializes a new Server instance
func New(cfg *config.Config) *Server {
	db, err := database.ConnectMongoDB(cfg.MongoURI, cfg.MongoDB)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	reg := registry.New(cfg.RegistryFile)
	bal := balancer.New(reg)
	listenerConfig := listener.Config{
		Address:     cfg.ListenAddr,
		TLSCertFile: cfg.TLSCertFile,
		TLSKeyFile:  cfg.TLSKeyFile,
		IdleTimeout: time.Duration(cfg.PoolIdleTimeout),
	}
	lis, err := listener.New(listenerConfig)
	if err != nil {
		log.Fatalf("Failed to create listener: %v", err)
	}
	poolConfig := pool.PoolConfig{
		MaxConns:        cfg.PoolMaxConns,
		IdleTimeout:     time.Duration(cfg.PoolIdleTimeout),
		MaxLifetime:     time.Duration(cfg.PoolMaxLifetime),
		CleanupInterval: time.Duration(cfg.PoolCleanupInterval),
	}
	s := &Server{
		router:   chi.NewRouter(),
		db:       db,
		config:   cfg,
		registry: reg,
		balancer: bal,
		pool:     pool.New(poolConfig),
		health: health.New(
			reg,
			time.Duration(cfg.HealthCheckInterval),
			time.Duration(cfg.HealthCheckTimeout),
			cfg.HealthCheckEndpoint,
		),
		listener: lis,
	}
	lis.SetHandler(s.handleConnection)
	s.setupRoutes()
	return s
}

// Start initializes the server components and begins the main server loop
func (s *Server) Start() error {
	log.Println("Starting load balancer...")

	// Register backend servers
	s.registerBackends()

	// Start the health checker
	s.health.Start()

	// Start periodic logging of server loads
	go s.logServerLoads()

	// Start the listener
	return s.listener.Start()
}

// handleConnection processes a single client connection
func (s *Server) handleConnection(clientConn net.Conn) {
	defer clientConn.Close()

	// Create a new http.Server with our router
	httpServer := &http.Server{
		Handler: s.router,
	}

	// Serve the connection
	httpServer.Serve(&singleConnListener{clientConn})
}

// singleConnListener is a helper type that wraps a single connection
type singleConnListener struct {
	conn net.Conn
}

func (l *singleConnListener) Accept() (net.Conn, error) {
	return l.conn, nil
}

func (l *singleConnListener) Close() error {
	return nil
}

func (l *singleConnListener) Addr() net.Addr {
	return l.conn.LocalAddr()
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

func (s *Server) setupRoutes() {
	userController := controller.NewUserController(s.db)

	s.router.Post("/users", userController.AddUser)
	s.router.Get("/users/last", userController.GetLastUser)

	// Add a catch-all route to forward requests to backend servers
	s.router.HandleFunc("/*", s.forwardToBackend)
}

func (s *Server) forwardToBackend(w http.ResponseWriter, r *http.Request) {
	backend := s.balancer.NextBackend()
	if backend == nil {
		http.Error(w, "No available backend servers", http.StatusServiceUnavailable)
		return
	}

	// Create a reverse proxy
	backendURL, err := url.Parse("http://" + backend.Address)
	if err != nil {
		http.Error(w, "Error parsing backend URL", http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(backendURL)
	proxy.ServeHTTP(w, r)
}

func (s *Server) registerBackends() {
	for _, backendAddr := range s.config.BackendServers {
		s.registry.Add(registry.Backend{Address: backendAddr})
	}
}