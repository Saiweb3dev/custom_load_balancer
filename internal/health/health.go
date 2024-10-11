package health

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
	"simple_load_balancer/internal/registry"
)

// HealthChecker periodically checks the health of backend servers
type HealthChecker struct {
	registry       *registry.Registry
	checkInterval  time.Duration
	timeout        time.Duration
	healthEndpoint string
}

// HealthCheckResult represents the result of a health check
type HealthCheckResult struct {
	Healthy bool
	Latency time.Duration
	Error   error
}

// New creates and initializes a new HealthChecker
func New(registry *registry.Registry, checkInterval, timeout time.Duration, healthEndpoint string) *HealthChecker {
	return &HealthChecker{
		registry:       registry,
		checkInterval:  checkInterval,
		timeout:        timeout,
		healthEndpoint: healthEndpoint,
	}
}

// Start begins the health checking loop in a separate goroutine
func (h *HealthChecker) Start() {
	go h.checkLoop()
}

// checkLoop runs the health checks at regular intervals
func (h *HealthChecker) checkLoop() {
	ticker := time.NewTicker(h.checkInterval)
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

// checkBackend performs a comprehensive health check on a single backend
func (h *HealthChecker) checkBackend(backend registry.Backend) {
	result := h.performHealthCheck(backend)

	if !result.Healthy {
		log.Printf("Backend %s is unhealthy: %v", backend.Address, result.Error)
		h.registry.Remove(backend.Address)
	} else {
		log.Printf("Backend %s is healthy (latency: %v)", backend.Address, result.Latency)
		// Optionally, update the backend's status in the registry
		// h.registry.UpdateStatus(backend.Address, result.Latency)
	}
}

// performHealthCheck conducts a series of health checks on a backend
func (h *HealthChecker) performHealthCheck(backend registry.Backend) HealthCheckResult {
	start := time.Now()

	// 1. TCP Connection Check
	conn, err := net.DialTimeout("tcp", backend.Address, h.timeout)
	if err != nil {
		return HealthCheckResult{Healthy: false, Error: fmt.Errorf("TCP connection failed: %v", err)}
	}
	conn.Close()

	// 2. HTTP Health Endpoint Check
	url := fmt.Sprintf("http://%s%s", backend.Address, h.healthEndpoint)
	ctx, cancel := context.WithTimeout(context.Background(), h.timeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return HealthCheckResult{Healthy: false, Error: fmt.Errorf("failed to create HTTP request: %v", err)}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return HealthCheckResult{Healthy: false, Error: fmt.Errorf("HTTP health check failed: %v", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return HealthCheckResult{Healthy: false, Error: fmt.Errorf("HTTP health check returned non-200 status: %d", resp.StatusCode)}
	}

	// 3. Additional checks can be added here (e.g., checking response body, verifying SSL certificates)

	latency := time.Since(start)
	return HealthCheckResult{Healthy: true, Latency: latency}
}