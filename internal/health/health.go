package health

import (
	"net"
	"time"
	"simple_load_balancer/internal/registry"
)

type HealthChecker struct {
	registry *registry.Registry
}

func New(registry *registry.Registry) *HealthChecker {
	return &HealthChecker{
		registry: registry,
	}
}

func (h *HealthChecker) Start() {
	go h.checkLoop()
}

func (h *HealthChecker) checkLoop() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		h.checkBackends()
	}
}

func (h *HealthChecker) checkBackends() {
	for _, backend := range h.registry.GetAll() {
		go h.checkBackend(backend)
	}
}

func (h *HealthChecker) checkBackend(backend registry.Backend) {
	conn, err := net.DialTimeout("tcp", backend.Address, 5*time.Second)
	if err != nil {
		// Backend is unhealthy, remove it from registry
		h.registry.Remove(backend.Address)
		return
	}
	conn.Close()
}