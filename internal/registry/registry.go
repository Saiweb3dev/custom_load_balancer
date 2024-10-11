package registry

import (
	"encoding/json"
	"io/ioutil"
	"os"
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
	filePath string
}

// New creates and initializes a new Registry
func New(filePath string) *Registry {
	r := &Registry{
		backends: make([]Backend, 0),
		filePath: filePath,
	}
	r.load() // Load existing data from file
	return r
}

// Add appends a new backend to the registry
func (r *Registry) Add(backend Backend) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.backends = append(r.backends, backend)
	r.save() // Save changes to file
}

// Remove deletes a backend from the registry based on its address
func (r *Registry) Remove(address string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for i, b := range r.backends {
		if b.Address == address {
			r.backends = append(r.backends[:i], r.backends[i+1:]...)
			r.save() // Save changes to file
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

// save writes the current state of the registry to a file
func (r *Registry) save() error {
	data, err := json.Marshal(r.backends)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(r.filePath, data, 0644)
}

// load reads the registry state from a file
func (r *Registry) load() error {
	data, err := ioutil.ReadFile(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // It's okay if the file doesn't exist yet
		}
		return err
	}
	return json.Unmarshal(data, &r.backends)
}