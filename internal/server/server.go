package server

import (
	"context"
	"fmt"
	"net"
)

type Server struct {
	addr    string
	dir     string
	network string
}

func New(addr, dir, network string) *Server {
	return &Server{
		addr:    addr,
		dir:     dir,
		network: network,
	}
}

func (s *Server) Start(ctx context.Context) error {
	l, err := net.Listen(s.network, s.addr)
	if err != nil {
		return fmt.Errorf("bind to port %s: %w", s.addr, err)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			return fmt.Errorf("accept connection: %w", err)
		}

		go handle(conn, s.dir)
	}
}
