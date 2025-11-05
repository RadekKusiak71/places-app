package server

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type APIServer struct {
	Port   string
	DB     *sql.DB
	Router *chi.Mux
}

func NewAPIServer(port string, db *sql.DB) *APIServer {
	return &APIServer{
		Port:   port,
		DB:     db,
		Router: chi.NewRouter(),
	}
}

func (s *APIServer) Start() error {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", s.Port),
		Handler:      s.Router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Printf("Starting HTTP server on port %s\n", s.Port)
	return server.ListenAndServe()
}
