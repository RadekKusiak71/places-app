package server

import (
	"time"

	"github.com/RadekKusiak71/places-app/internal/auth"
	"github.com/RadekKusiak71/places-app/internal/users"
	"github.com/RadekKusiak71/places-app/internal/utils"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (s *APIServer) SetupRouter() {
	s.Router.Use(middleware.RequestID)
	s.Router.Use(middleware.RealIP)
	s.Router.Use(middleware.Logger)
	s.Router.Use(middleware.Recoverer)
	s.Router.Use(middleware.Timeout(60 * time.Second))

	userStore := users.NewStore(s.DB)
	authStore := auth.NewStore(s.DB)

	authService := auth.NewService(userStore, authStore)
	authHandler := auth.NewAuthHandler(authService)

	v1Router := chi.NewRouter()
	v1Router.Route("/auth", func(r chi.Router) {
		r.Post("/register", utils.MakeHandlerFunc(authHandler.Register))
		r.Post("/token", utils.MakeHandlerFunc(authHandler.ObtainJWT))
	})

	s.Router.Mount("/api/v1", v1Router)
}
