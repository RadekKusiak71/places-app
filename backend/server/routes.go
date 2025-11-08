package server

import (
	"github.com/RadekKusiak71/places-app/internal/middlewares"
	"time"

	"github.com/RadekKusiak71/places-app/internal/handlers"
	"github.com/RadekKusiak71/places-app/internal/services"
	"github.com/RadekKusiak71/places-app/internal/stores"
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

	v1Router := chi.NewRouter()

	// Stores
	userStore := stores.NewUserStore(s.DB)
	rtStore := stores.NewRefreshTokenStore(s.DB)
	placesStore := stores.NewPlacesStore(s.DB)

	// Service
	authService := services.NewAuthService(userStore, rtStore)
	placesService := services.NewPlacesService(placesStore)

	// Handlers
	authHandler := handlers.NewAuthHandler(authService)
	placesHandler := handlers.NewPlacesHandler(placesService)

	v1Router.Route("/auth", func(r chi.Router) {
		r.Post("/register", utils.MakeHandlerFunc(authHandler.Register))
		r.Post("/token", utils.MakeHandlerFunc(authHandler.ObtainTokensPair))
		r.Post("/token/refresh", utils.MakeHandlerFunc(authHandler.RefreshTokensPair))
	})

	v1Router.Route("/places", func(r chi.Router) {
		r.Get("/", utils.MakeHandlerFunc(middlewares.AuthMiddleware(placesHandler.ListPlaces, userStore)))
		r.Post("/", utils.MakeHandlerFunc(middlewares.AuthMiddleware(placesHandler.CreatePlace, userStore)))
		r.Get("/{placeID}", utils.MakeHandlerFunc(middlewares.AuthMiddleware(placesHandler.RetrievePlace, userStore)))
		r.Delete("/{placeID}", utils.MakeHandlerFunc(middlewares.AuthMiddleware(placesHandler.DeletePlace, userStore)))

		r.Route("/{placeID}/photos", func(r chi.Router) {
			r.Get("/", nil)
			r.Post("/", nil)
			r.Delete("/{photoID}", nil)
		})
	})

	v1Router.Route("/users", func(r chi.Router) {
		r.Get("/me", nil)
	})

	s.Router.Mount("/api/v1", v1Router)
}
