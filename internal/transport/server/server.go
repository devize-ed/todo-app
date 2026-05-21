package http

import (
	"context"
	"fmt"
	"net/http"

	"github.com/devize-ed/todo-app/internal/config"
	"github.com/devize-ed/todo-app/internal/logger"
	"github.com/devize-ed/todo-app/internal/transport/handlers"
)

type Server struct {
	*http.Server
	cfg *config.Config
}

func NewServer(cfg *config.Config, h *handlers.Handler) *Server {
	return &Server{
		Server: &http.Server{
			Addr:    cfg.Server.Addr,
			Handler: h.NewRouter(),
		},
		cfg: cfg,
	}
}

func (s *Server) Start(ctx context.Context) error {

	errChan := make(chan error)
	go func() {
		if err := s.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
		close(errChan)
	}()
	select {
	case err := <-errChan:
		if err != nil {
			return fmt.Errorf("failed to listen: %w", err)
		}
	case <-ctx.Done():
		return s.stop()
	}
	return nil
}

func (s *Server) stop() error {
	logger.Warn("shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), s.cfg.Server.Timeout)
	defer cancel()

	if err := s.Server.Shutdown(shutdownCtx); err != nil {

		s.Close()
		return fmt.Errorf("shutdown server error: %w", err)
	}

	logger.Info("server shutted down")
	return nil
}
