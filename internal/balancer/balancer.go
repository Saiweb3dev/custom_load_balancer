package balancer

import (
	"sync"
	"time"
	"simple_load_balancer/internal/registry"
)

// Balancer implements a load-aware balancing algorithm
type Balancer struct {
	registry *registry.Registry
	mu       sync.RWMutex
	serverLoads map[string]float64
	lastUpdate time.Time
}

// New creates and initializes a new Balancer
func New(registry *registry.Registry) *Balancer {
	b := &Balancer{
		registry: registry,
		serverLoads: make(map[string]float64),
		lastUpdate: time.Now(),
	}
	go b.periodicLoadUpdate()
	return b
}

// NextBackend selects the next backend server based on current load
func (b *Balancer) NextBackend() *registry.Backend {
	b.mu.RLock()
	defer b.mu.RUnlock()

	backends := b.registry.GetAll()
	if len(backends) == 0 {
		return nil
	}

	var leastLoadedBackend *registry.Backend
	minLoad := float64(101) // Initialize with a value higher than possible load percentage

	for _, backend := range backends {
		load, exists := b.serverLoads[backend.Address]
		if !exists {
			// If we don't have load info, assume 50% as a neutral value
			load = 50
		}

		if load < minLoad {
			minLoad = load
			leastLoadedBackend = &backend
		}
	}

	return leastLoadedBackend
}

// UpdateServerLoad updates the load information for a specific server
func (b *Balancer) UpdateServerLoad(serverAddress string, load float64) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.serverLoads[serverAddress] = load
	b.lastUpdate = time.Now()
}

// periodicLoadUpdate simulates periodic load updates from servers
// In a real-world scenario, this would be replaced by actual server metrics
func (b *Balancer) periodicLoadUpdate() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		b.mu.Lock()
		for _, backend := range b.registry.GetAll() {
			// Simulate load update. In reality, you'd get this data from the server
			b.serverLoads[backend.Address] = float64(50 + (time.Now().UnixNano() % 50))
		}
		b.lastUpdate = time.Now()
		b.mu.Unlock()
	}
}

// GetServerLoads returns the current load information for all servers
func (b *Balancer) GetServerLoads() map[string]float64 {
	b.mu.RLock()
	defer b.mu.RUnlock()

	loads := make(map[string]float64)
	for k, v := range b.serverLoads {
		loads[k] = v
	}
	return loads
}