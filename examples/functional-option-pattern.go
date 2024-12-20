package main

import (
	"fmt"
)

// Server represents a server with configurable pattern
type Server struct {
	Host string
	Port int
	TLS  bool
}

// Option type: A function that modifies a Server instance
type Option func(*Server)

// Functional options to configure the server
func WithHost(host string) Option {
	return func(s *Server) {
		s.Host = host
	}
}

func WithPort(port int) Option {
	return func(s *Server) {
		s.Port = port
	}
}

func WithTLS(enabled bool) Option {
	return func(s *Server) {
		s.TLS = enabled
	}
}

// Constructor function that applies options
func NewServer(opts ...Option) *Server {
	// Default configuration
	server := &Server{
		Host: "localhost",
		Port: 80,
		TLS:  false,
	}

	for _, opt := range opts {
		opt(server)
	}

	return server
}

func main() {
	// Use functional options to configure the server
	server := NewServer(WithHost("127.0.0.1"), WithPort(8080), WithTLS(true))
	fmt.Printf("Server config: Host=%s, Port=%d, TLS=%v\n",
		server.Host, server.Port, server.TLS)
}
