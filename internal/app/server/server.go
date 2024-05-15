package server

import (
	"context"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

func Listen(ctx context.Context, opts Opts) (*Server, error) {
	var server *Server

	listener, err := net.Listen("tcp", opts.Config.Address())
	if err != nil {
		return server, fmt.Errorf("TCP listen: %w", err)
	}

	server = &Server{
		listener: listener,
		config:   opts.Config,
		logger:   opts.Logger,
		service:  opts.Service,
	}

	server.shutdownWg.Add(1)
	go server.acceptConnections(ctx)

	return server, nil
}

// Opts - options to run server.
type Opts struct {
	Config  Config
	Logger  Logger
	Service Service
}

// Sever - tcp server.
type Server struct {
	listener net.Listener
	config   Config
	logger   Logger
	service  Service

	shutdownWg    sync.WaitGroup
	isShutingDown atomic.Bool
}

// Shutdown - shutdown server gracefully.
func (s *Server) Shutdown() {
	const operationName = "server.Shutdown"

	s.isShutingDown.Store(true)
	s.listener.Close()

	done := make(chan struct{})
	go func() {
		s.shutdownWg.Wait()
		close(done)
	}()

	select {
	case <-done:
		s.logger.Debug("shutdown server gracefully", "operationName", operationName)

		return
	case <-time.After(s.config.ShutdownTimeout()):
		s.logger.Debug("shutdown server by timeout", "operationName", operationName)

		return
	}
}

func (s *Server) acceptConnections(_ context.Context) {
	const operationName = "server.acceptConnections"

	defer s.shutdownWg.Done()

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.isShutingDown.Load() {
				s.logger.Debug("server closed", "operationName", operationName)

				return
			}

			s.logger.Error(err.Error(), "operationName", operationName)

			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	const operationName = "server.handleConnection"

	defer conn.Close()

	err := conn.SetReadDeadline(time.Now().Add(s.config.ConnectionTimeout()))
	if err != nil {
		s.logger.Error(err.Error(), "operationName", operationName)

		return
	}

	if s.isShutingDown.Load() {
		s.logger.Error("server closed", "operationName", operationName)

		return
	}

	s.service.HandleMessages(conn.RemoteAddr().String(), conn)
}
