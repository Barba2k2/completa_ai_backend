package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/completaai/backend/internal/config"
	"github.com/completaai/backend/internal/handlers"
	"github.com/completaai/backend/internal/middleware"
	"github.com/completaai/backend/internal/repository"
	"github.com/completaai/backend/internal/router"
	"github.com/completaai/backend/pkg/database"
)

func main() {
	cfg := config.Load()

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	client, err := database.NewClient(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer client.Close()

	userRepo := repository.NewUserRepository(client.DB, cfg.JWTSecret)
	collectionRepo := repository.NewCollectionRepository(client.DB)

	authHandler := handlers.NewAuthHandler(userRepo)
	collectionHandler := handlers.NewCollectionHandler(collectionRepo)
	shareHandler := handlers.NewShareHandler(collectionRepo)
	statsHandler := handlers.NewStatsHandler(collectionRepo)
	authMiddleware := middleware.NewAuthMiddleware(cfg.JWTSecret)

	r := router.New(authHandler, collectionHandler, shareHandler, statsHandler, authMiddleware)

	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("server listening on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("server stopped cleanly")
}
