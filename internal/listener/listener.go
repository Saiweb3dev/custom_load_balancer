package listener

import (
	"crypto/tls"
	"log"
	"net"
	"time"
)

// Listener handles incoming network connections
type Listener struct {
	address     string
	tlsConfig   *tls.Config
	handler     func(net.Conn)
	idleTimeout time.Duration
}

// Config holds the configuration for the Listener
type Config struct {
	Address     string
	TLSCertFile string
	TLSKeyFile  string
	IdleTimeout time.Duration
}

// New creates and initializes a new Listener
func New(cfg Config) (*Listener, error) {
	l := &Listener{
			address:     cfg.Address,
			idleTimeout: cfg.IdleTimeout,
	}

	if cfg.TLSCertFile != "" && cfg.TLSKeyFile != "" {
			cert, err := tls.LoadX509KeyPair(cfg.TLSCertFile, cfg.TLSKeyFile)
			if err != nil {
					return nil, err
			}
			l.tlsConfig = &tls.Config{
					Certificates: []tls.Certificate{cert},
			}
	}

	return l, nil
}

// SetHandler sets the connection handler function
func (l *Listener) SetHandler(handler func(net.Conn)) {
	l.handler = handler
}

// Start begins listening for incoming connections and handles them
func (l *Listener) Start() error {
	var listener net.Listener
	var err error

	if l.tlsConfig != nil {
			listener, err = tls.Listen("tcp", l.address, l.tlsConfig)
	} else {
			listener, err = net.Listen("tcp", l.address)
	}

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
	defer conn.Close()

	// Set idle timeout
	if l.idleTimeout > 0 {
		err := conn.SetDeadline(time.Now().Add(l.idleTimeout))
		if err != nil {
			log.Printf("Error setting connection deadline: %v", err)
			return
		}
	}

	// Handle TLS handshake if TLS is enabled
	if l.tlsConfig != nil {
		tlsConn, ok := conn.(*tls.Conn)
		if !ok {
			log.Printf("Error: expected TLS connection")
			return
		}
		err := tlsConn.Handshake()
		if err != nil {
			log.Printf("TLS handshake error: %v", err)
			return
		}
	}

	// Call the user-defined handler
	if l.handler != nil {
		l.handler(conn)
	} else {
		log.Printf("Warning: no handler set for connection")
	}
}