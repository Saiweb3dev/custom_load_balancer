package listener

import (
	"net"
	"log"
)

// Listener handles incoming network connections
type Listener struct {
	address string
}

// New creates and initializes a new Listener
func New(address string) *Listener {
	return &Listener{
		address: address,
	}
}

// Start begins listening for incoming connections and handles them
func (l *Listener) Start() error {
	listener, err := net.Listen("tcp", l.address)
	if err != nil {
		return err
	}
	defer listener.Close()

	log.Printf("Listening on %s", l.address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v", err)
			continue
		}
		go l.handleConnection(conn)
	}
}

// handleConnection processes a single client connection
func (l *Listener) handleConnection(conn net.Conn) {
	// Here you would implement the logic to forward the connection
	// to a backend server using the balancer, pool, etc.
	defer conn.Close()
	log.Printf("Received connection from %s", conn.RemoteAddr())
}