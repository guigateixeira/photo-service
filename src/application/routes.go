package application

import (
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"photo-service/src/handler"
)

func loadRoutes(photoHandler *handler.PhotoHandler) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	// Starndard health check endpoint
	v1Router.Get("/health", handler.HandlerReadiness)

	v1Router.Route("/photos", func(router chi.Router) {
		loadPhotoRoutes(router, photoHandler)
	})

	router.Mount("/v1", v1Router)

	return router
}

func loadPhotoRoutes(router chi.Router, photoHandler *handler.PhotoHandler) {
	router.Post("/upload", photoHandler.CreatePhoto)
}
