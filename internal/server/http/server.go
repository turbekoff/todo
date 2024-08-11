package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/turbekoff/todo/internal/config"
)

type Server struct {
	http *http.Server
}

func New(handler http.Handler, config *config.HTTPConfig) *Server {
	return &Server{
		http: &http.Server{
			Handler:        handler,
			Addr:           fmt.Sprintf("%s:%d", config.Host, config.Port),
			ReadTimeout:    config.ReadTimeout,
			WriteTimeout:   config.WriteTimeout,
			MaxHeaderBytes: config.MaxHeaderBytes,
		},
	}
}

func (s *Server) Run() error {
	return s.http.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.http.Shutdown(ctx)
}
