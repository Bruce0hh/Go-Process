package main

import (
	"fmt"
	"time"
)

//Uber推荐的Functional Options模式

type Server struct {
	Addr     string
	Port     int
	MaxConns int
	Timeout  time.Duration
}

type Option interface {
	setOption(*Server)
}

type maxConnsOption int

func (m maxConnsOption) setOption(s *Server) {
	s.MaxConns = int(m)
}

func WithMaxConns(m int) Option {
	return maxConnsOption(m)
}

type timeout time.Duration

func (t timeout) setOption(s *Server) {
	s.Timeout = time.Duration(t)
}

func WithTimeout(t time.Duration) Option {
	return timeout(t)
}

func NewServer(addr string, port int, opts ...Option) (*Server, error) {
	s := &Server{Addr: addr, Port: port}

	for _, o := range opts {
		o.setOption(s)
	}
	return s, nil
}

func main() {

	s1, _ := NewServer("localhost", 1024)
	s2, _ := NewServer("localhost", 2048, WithMaxConns(2000))
	s3, _ := NewServer("localhost", 8989, WithTimeout(30*time.Second))

	fmt.Printf("s1: %+v\n", s1)
	fmt.Printf("s2: %+v\n", s2)
	fmt.Printf("s3: %+v\n", s3)
}
