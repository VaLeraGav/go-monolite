package server

import (
	"errors"
	"go-monolite/internal/config"
	"go-monolite/internal/store"
	"go-monolite/pkg/logger"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

type Server struct {
	config *config.Config
	router *chi.Mux
	store  *store.Store
}

func NewServer(config *config.Config, store *store.Store) *Server {
	s := &Server{
		config: config,
		store:  store,
		router: chi.NewRouter(),
	}
	return s
}

func (s *Server) StartServer() *http.Server {
	logger.Info("starting server", "address", s.config.HTTPServer.Address)

	s.ConfigureRouting()

	server := &http.Server{
		Addr:         s.config.HTTPServer.Address,
		Handler:      s.router,
		ReadTimeout:  time.Duration(s.config.HTTPServer.Timeout) * time.Second, // Таймаут на чтение
		WriteTimeout: time.Duration(s.config.HTTPServer.Timeout) * time.Second, // Таймаут на запись
		IdleTimeout:  120 * time.Second,                                        // Таймаут для неактивных соединений
	}

	go func() {
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal(err, "HTTP server error")
			return
		}

		logger.Info("stopped serving new connections")
	}()

	return server
}
