package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/revenantio/revenant-backend/internal/config"
	"github.com/revenantio/revenant-backend/internal/database"
	"github.com/revenantio/revenant-backend/internal/logger"
	"github.com/revenantio/revenant-backend/internal/server"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize logger
	log := logger.NewLogger(cfg.Environment)
	defer log.Sync()

	// Connect to database
	db, err := database.NewConnection(cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database", map[string]interface{}{
			"error": err.Error(),
		})
	}
	defer db.Close()

	// Run migrations
	if err := database.RunMigrations(db, cfg.Database); err != nil {
		log.Fatal("Failed to run migrations", map[string]interface{}{
			"error": err.Error(),
		})
	}

	// Initialize server
	srv := server.New(cfg, log, db)

	// Start server in a goroutine
	go func() {
		addr := fmt.Sprintf(":%d", cfg.Server.Port)
		log.Info("Starting server", map[string]interface{}{
			"address": addr,
		})

		if err := srv.Run(addr); err != nil && err != http.ErrServerClosed {
			log.Fatal("Server error", map[string]interface{}{
				"error": err.Error(),
			})
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down server", nil)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server shutdown error", map[string]interface{}{
			"error": err.Error(),
		})
	}

	log.Info("Server stopped", nil)
}
