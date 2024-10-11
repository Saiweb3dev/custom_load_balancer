package balancer

import (
	"sync/atomic"
	"simple_load_balancer/internal/registry"
)

// Balancer implements a round-robin load balancing algorithm
type Balancer struct {
	registry *registry.Registry
	current  uint64
}

// New creates and initializes a new Balancer
func New(registry *registry.Registry) *Balancer {
	return &Balancer{
		registry: registry,
		current:  0,
	}
}

// NextBackend selects the next backend server in a round-robin fashion
func (b *Balancer) NextBackend() *registry.Backend {
	backends := b.registry.GetAll()
	if len(backends) == 0 {
		return nil
	}
	
	next := atomic.AddUint64(&b.current, 1) % uint64(len(backends))
	return &backends[next]
}