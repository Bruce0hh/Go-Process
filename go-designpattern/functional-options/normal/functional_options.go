package main

import (
	"crypto/tls"
	"fmt"
	"time"
)

// 陈皓《Go 编程模式》中的 Functional Options例子

type Server struct {
	Addr     string
	Port     int
	MaxConns int
	Protocol string
	TLS      *tls.Config
	Timeout  time.Duration
}

type Option func(*Server)

func Protocol(p string) Option {
	return func(s *Server) {
		s.Protocol = p
	}
}

func Timeout(timeout time.Duration) Option {
	return func(s *Server) {
		s.Timeout = timeout
	}
}

func MaxConns(maxconns int) Option {
	return func(s *Server) {
		s.MaxConns = maxconns
	}
}

func TLS(tls *tls.Config) Option {
	return func(s *Server) {
		s.TLS = tls
	}
}

func NewServer(addr string, port int, options ...Option) (*Server, error) {
	srv := Server{
		Addr:     addr,
		Port:     port,
		MaxConns: 1000,
		Protocol: "tcp",
		TLS:      nil,
		Timeout:  30 * time.Second,
	}

	for _, option := range options {
		option(&srv)
	}

	//...
	return &srv, nil
}

func main() {
	s1, _ := NewServer("localhost", 1024)
	s2, _ := NewServer("localhost", 2048, Protocol("udp"))
	s3, _ := NewServer("localhost", 8989, Timeout(300*time.Second), MaxConns(1000))

	fmt.Printf("s1: %+v\n", s1)
	fmt.Printf("s2: %+v\n", s2)
	fmt.Printf("s3: %+v\n", s3)
}
