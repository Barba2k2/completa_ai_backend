package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"github.com/completaai/backend/internal/handlers"
	"github.com/completaai/backend/internal/middleware"
)

// New wires up the HTTP routes and returns the top level router.
func New(authHandler *handlers.AuthHandler, collectionHandler *handlers.CollectionHandler, shareHandler *handlers.ShareHandler, statsHandler *handlers.StatsHandler, auth *middleware.AuthMiddleware) http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	r.Route("/auth", func(r chi.Router) {
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)
		r.With(auth.Authenticate).Get("/me", authHandler.Me)
	})

	r.Route("/collection", func(r chi.Router) {
		r.Use(auth.Authenticate)
		r.Get("/", collectionHandler.Get)
		r.Put("/", collectionHandler.Update)
		r.Post("/sync", collectionHandler.Sync)
	})

	r.Route("/share", func(r chi.Router) {
		r.Use(auth.Authenticate)
		r.Get("/missing", shareHandler.GetMissing)
		r.Get("/duplicates", shareHandler.GetDuplicates)
		r.Get("/both", shareHandler.GetBoth)
	})

	r.Route("/stats", func(r chi.Router) {
		r.Use(auth.Authenticate)
		r.Get("/sections", statsHandler.GetSections)
		r.Get("/", statsHandler.GetStats)
	})

	return r
}
