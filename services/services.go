package services

import (
	"errors"
)

type Server struct {
	channel chan bool
}
type Arg int
type Result int

type Args struct {
	IPAddress, PortNumber string
}

type Done bool
type Terminated bool

func NewServer() *Server {
	channel := make(chan bool)
	return &Server{channel: channel}
}

func (s *Server) Fibonacci(index Arg, ret *Result) error {
	*ret = s.calculateFib(index)
	if *ret == -1 {
		return errors.New("computation interrupted")
	}
	return nil
}

func (s *Server) Pow(n Arg, ret *Result) error {
	*ret = s.calculatePow(n)
	return nil
}

func (s *Server) StopComputation(done Done, ret *Terminated) error {
	if done == true {
		select {
		case s.channel <- true:
			*ret = true
		default:
			*ret = false
		}
	}
	return nil
}

func (s *Server) calculateFib(n Arg) Result {
	select {
	case <-s.channel:
		return Result(-1)
	default:
		if n <= 1 {
			return Result(n)
		}
		return s.calculateFib(n-1) + s.calculateFib(n-2)
	}
}

func (s *Server) calculatePow(n Arg) Result {
	select {
	case <-s.channel:
		return Result(-1)
	default:
		return Result(n * n)
	}
}
