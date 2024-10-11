package health

import (
	"net"
	"time"
	"simple_load_balancer/internal/registry"
)

// HealthChecker periodically checks the health of backend servers
type HealthChecker struct {
	registry *registry.Registry
}

// New creates and initializes a new HealthChecker
func New(registry *registry.Registry) *HealthChecker {
	return &HealthChecker{
		registry: registry,
	}
}

// Start begins the health checking loop in a separate goroutine
func (h *HealthChecker) Start() {
	go h.checkLoop()
}

// checkLoop runs the health checks at regular intervals
func (h *HealthChecker) checkLoop() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		h.checkBackends()
	}
}

// checkBackends initiates a health check for all registered backends
func (h *HealthChecker) checkBackends() {
	for _, backend := range h.registry.GetAll() {
		go h.checkBackend(backend)
	}
}

// checkBackend performs a health check on a single backend
func (h *HealthChecker) checkBackend(backend registry.Backend) {
	conn, err := net.DialTimeout("tcp", backend.Address, 5*time.Second)
	if err != nil {
		// Backend is unhealthy, remove it from registry
		h.registry.Remove(backend.Address)
		return
	}
	conn.Close()
}