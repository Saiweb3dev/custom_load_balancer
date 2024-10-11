package pool

import (
	"net"
	"sync"
	"time"
	"errors"
)

// ConnectionWrapper wraps a net.Conn with creation time information
type ConnectionWrapper struct {
	conn       net.Conn
	createdAt  time.Time
	lastUsedAt time.Time
}

// Pool manages a pool of reusable network connections
type Pool struct {
	connections     map[string][]*ConnectionWrapper
	mu              sync.Mutex
	maxConns        int
	idleTimeout     time.Duration
	maxLifetime     time.Duration
	cleanupInterval time.Duration
}

// PoolConfig holds configuration for the connection pool
type PoolConfig struct {
	MaxConns        int
	IdleTimeout     time.Duration
	MaxLifetime     time.Duration
	CleanupInterval time.Duration
}

// New creates and initializes a new connection Pool
func New(config PoolConfig) *Pool {
	p := &Pool{
		connections:     make(map[string][]*ConnectionWrapper),
		maxConns:        config.MaxConns,
		idleTimeout:     config.IdleTimeout,
		maxLifetime:     config.MaxLifetime,
		cleanupInterval: config.CleanupInterval,
	}
	go p.periodicCleanup()
	return p
}

// Get retrieves a connection from the pool or creates a new one if none are available
func (p *Pool) Get(address string) (net.Conn, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	now := time.Now()

	if conns, ok := p.connections[address]; ok && len(conns) > 0 {
		for i, wrapper := range conns {
			if p.isConnectionValid(wrapper, now) {
				// Remove the connection from the slice
				p.connections[address] = append(conns[:i], conns[i+1:]...)
				wrapper.lastUsedAt = now
				return wrapper.conn, nil
			}
		}
	}

	// If no valid connection is available, create a new one
	if len(p.connections[address]) >= p.maxConns {
		return nil, ErrPoolExhausted
	}

	conn, err := net.DialTimeout("tcp", address, 10*time.Second)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// Put adds a connection back to the pool for reuse
func (p *Pool) Put(address string, conn net.Conn) {
	p.mu.Lock()
	defer p.mu.Unlock()

	now := time.Now()
	wrapper := &ConnectionWrapper{
		conn:       conn,
		createdAt:  now,
		lastUsedAt: now,
	}

	if len(p.connections[address]) < p.maxConns {
		p.connections[address] = append(p.connections[address], wrapper)
	} else {
		conn.Close()
	}
}

// isConnectionValid checks if a connection is still valid based on idle timeout and max lifetime
func (p *Pool) isConnectionValid(wrapper *ConnectionWrapper, now time.Time) bool {
	if now.Sub(wrapper.lastUsedAt) > p.idleTimeout {
		wrapper.conn.Close()
		return false
	}
	if now.Sub(wrapper.createdAt) > p.maxLifetime {
		wrapper.conn.Close()
		return false
	}
	return true
}

// periodicCleanup removes expired connections from the pool
func (p *Pool) periodicCleanup() {
	ticker := time.NewTicker(p.cleanupInterval)
	defer ticker.Stop()

	for range ticker.C {
		p.mu.Lock()
		now := time.Now()
		for address, conns := range p.connections {
			validConns := make([]*ConnectionWrapper, 0, len(conns))
			for _, wrapper := range conns {
				if p.isConnectionValid(wrapper, now) {
					validConns = append(validConns, wrapper)
				}
			}
			p.connections[address] = validConns
		}
		p.mu.Unlock()
	}
}

// Close closes all connections in the pool
func (p *Pool) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()

	for _, conns := range p.connections {
		for _, wrapper := range conns {
			wrapper.conn.Close()
		}
	}
	p.connections = make(map[string][]*ConnectionWrapper)
}

// ErrPoolExhausted is returned when the pool has reached its maximum number of connections
var ErrPoolExhausted = errors.New("connection pool exhausted")