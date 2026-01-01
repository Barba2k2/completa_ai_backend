package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"github.com/completaai/backend/internal/handlers"
	"github.com/completaai/backend/internal/middleware"
)

// New wires up the HTTP routes and returns the top level router.
func New(authHandler *handlers.AuthHandler, collectionHandler *handlers.CollectionHandler, shareHandler *handlers.ShareHandler, statsHandler *handlers.StatsHandler, profileHandler *handlers.ProfileHandler, auth *middleware.AuthMiddleware) http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
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
		r.With(auth.Authenticate).Delete("/account", authHandler.DeleteAccount)
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

	r.Route("/profile", func(r chi.Router) {
		r.Use(auth.Authenticate)
		r.Get("/", profileHandler.GetProfile)
		r.Put("/", profileHandler.UpdateProfile)
		r.Patch("/", profileHandler.UpdateProfile)
	})

	return r
}
