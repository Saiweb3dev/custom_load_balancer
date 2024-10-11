package balancer

import (
	"sync/atomic"
	"simple_load_balancer/internal/registry"
)

type Balancer struct {
	registry *registry.Registry
	current  uint64
}

func New(registry *registry.Registry) *Balancer {
	return &Balancer{
		registry: registry,
		current:  0,
	}
}

func (b *Balancer) NextBackend() *registry.Backend {
	backends := b.registry.GetAll()
	if len(backends) == 0 {
		return nil
	}
	
	next := atomic.AddUint64(&b.current, 1) % uint64(len(backends))
	return &backends[next]
}