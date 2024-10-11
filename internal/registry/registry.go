package registry

import (
	"sync"
)

// Backend represents a server that can handle requests
type Backend struct {
	Address string
	// Add more fields as needed (e.g., weight, capacity)
}

// Registry manages a list of backend servers
type Registry struct {
	backends []Backend
	mu       sync.RWMutex
}

// New creates and initializes a new Registry
func New() *Registry {
	return &Registry{
		backends: make([]Backend, 0),
	}
}

// Add appends a new backend to the registry
func (r *Registry) Add(backend Backend) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.backends = append(r.backends, backend)
}

// Remove deletes a backend from the registry based on its address
func (r *Registry) Remove(address string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, b := range r.backends {
		if b.Address == address {
			r.backends = append(r.backends[:i], r.backends[i+1:]...)
			break
		}
	}
}

// GetAll returns a copy of all backends in the registry
func (r *Registry) GetAll() []Backend {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]Backend{}, r.backends...)
}