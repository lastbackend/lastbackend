package client

import (
	"fmt"
	"net"
)

type Server struct {
	port        int
	listener    net.Listener
	connections map[net.Conn]bool
}

func NewTCPServer(port int) *Server {
	var server = Server{
		port:        port,
		connections: make(map[net.Conn]bool),
	}
	return &server
}

func (s *Server) Start() (err error) {
	s.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	return err
}

func (s *Server) Accept(cb func(net.Conn)) {
	for {
		connection, err := s.listener.Accept()
		if err != nil {
			fmt.Printf("Accept failed, %v\n", err)
			cb(nil)
		} else {
			s.connections[connection] = true
			cb(connection)
		}
	}
}

func (s *Server) Send(b []byte) {
	for connection := range s.connections {
		connection.Write(b)
	}
}

func (s *Server) Close() error {
	for connection := range s.connections {
		connection.Close()
	}

	return s.listener.Close()
}
