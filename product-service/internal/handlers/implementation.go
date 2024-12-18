package handlers

import (
	"context"
	"net"
	"net/http"
	"os"
	"product-service/internal/service"
	"time"
)

type Implementation struct {
	httpServer *http.Server
}

func NewServer(svc *service.Service) *Implementation {
	return &Implementation{
		httpServer: &http.Server{
			ReadTimeout: 3 * time.Second,
			Handler:     newRouter(svc),
		},
	}
}

func (s *Implementation) Run(_ context.Context) error {
	l, err := net.Listen("tcp", os.Getenv("HTTP_PORT"))
	if err != nil {
		return err
	}

	if err = s.httpServer.Serve(l); err != nil {
		return err
	}

	return nil
}

func (s *Implementation) Stop() error {
	err := s.httpServer.Shutdown(context.Background())
	if err != nil {
		return err
	}

	return nil
}
