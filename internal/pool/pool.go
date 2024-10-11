package pool

import (
	"net"
	"sync"
)

type Pool struct {
	connections map[string][]net.Conn
	mu          sync.Mutex
}

func New() *Pool {
	return &Pool{
		connections: make(map[string][]net.Conn),
	}
}

func (p *Pool) Get(address string) net.Conn {
	p.mu.Lock()
	defer p.mu.Unlock()

	if conns, ok := p.connections[address]; ok && len(conns) > 0 {
		conn := conns[len(conns)-1]
		p.connections[address] = conns[:len(conns)-1]
		return conn
	}
	
	// If no connection available, create a new one
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return nil
	}
	return conn
}

func (p *Pool) Put(address string, conn net.Conn) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.connections[address] = append(p.connections[address], conn)
}