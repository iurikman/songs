package rest

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
)

const (
	gracefulShutdownTimeout = 5 * time.Second
	readHeaderTimeout       = 5 * time.Second
	maxHeaderBytes          = 1 << 20
)

type SrvConfig struct {
	BindAddr string
}

type Server struct {
	config SrvConfig
	router *chi.Mux
	server *http.Server
	svc    service
}

func NewServer(cfg SrvConfig, svc service) (*Server, error) {
	router := chi.NewRouter()

	srv := &http.Server{
		Addr:           cfg.BindAddr,
		Handler:        router,
		ReadTimeout:    readHeaderTimeout,
		MaxHeaderBytes: maxHeaderBytes,
	}

	return &Server{
		config: cfg,
		router: router,
		server: srv,
		svc:    svc,
	}, nil
}

func (s *Server) Start(ctx context.Context) error {
	s.configRouter()

	go func() {
		<-ctx.Done()

		ctxWithTimeout, cancel := context.WithTimeout(ctx, gracefulShutdownTimeout)
		defer cancel()

		if err := s.server.Shutdown(ctxWithTimeout); err != nil {
			log.Warnf("failed to shutdown gracefully %s", err)
		}
	}()

	err := s.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("s.server.ListenAndServe() err: %w", err)
	}

	return nil
}

func (s *Server) configRouter() {
	s.router.Route("/api", func(r chi.Router) {
		r.Route("/v1", func(r chi.Router) {
			r.Route("/songs", func(r chi.Router) {
				r.Post("/", s.createSong)
			})
		})
	})
}
