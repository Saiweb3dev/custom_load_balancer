package listener

import (
	"net"
	"log"
)

type Listener struct {
	address string
}

func New(address string) *Listener {
	return &Listener{
		address: address,
	}
}

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

func (l *Listener) handleConnection(conn net.Conn) {
	// Here you would implement the logic to forward the connection
	// to a backend server using the balancer, pool, etc.
	defer conn.Close()
	log.Printf("Received connection from %s", conn.RemoteAddr())
}