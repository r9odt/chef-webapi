package httpserver

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/JIexa24/chef-webapi/httpserver/handler"
	"github.com/JIexa24/chef-webapi/logging"
)

// HTTPServer describes the server instance.
type HTTPServer struct {
	address string
	port    string
	logger  logging.Logger
	server  *http.Server
}

// New return new http server object.
func New(address, port string, logger logging.Logger) *HTTPServer {
	return &HTTPServer{
		address: address,
		port:    port,
		logger:  logger,
		server: &http.Server{
			Handler: handler.NewHandler(),
		},
	}
}

// Listen starts serve port.
func (s *HTTPServer) Listen() {
	addr := s.address + ":" + s.port
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		s.logger.Fatal(err)
	}

	s.logger.Infof("Start listening by address: [%s]", addr)

	go func() {
		err := s.server.Serve(listener)
		if err != nil && err != http.ErrServerClosed {
			s.logger.Fatalf("Serve: %v", err)
		}
	}()
}

// Stop server.
func (s *HTTPServer) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Errorf("Can't stop server correctly: %v", err)
	}
}
