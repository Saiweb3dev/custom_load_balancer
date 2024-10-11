package registry

import (
	"sync"
)

type Backend struct {
	Address string
	// Add more fields as needed (e.g., weight, capacity)
}

type Registry struct {
	backends []Backend
	mu       sync.RWMutex
}

func New() *Registry {
	return &Registry{
		backends: make([]Backend, 0),
	}
}

func (r *Registry) Add(backend Backend) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.backends = append(r.backends, backend)
}

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

func (r *Registry) GetAll() []Backend {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return append([]Backend{}, r.backends...)
}