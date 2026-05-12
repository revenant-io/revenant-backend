package server

import (
	"context"
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/revenantio/revenant-backend/internal/config"
	"github.com/revenantio/revenant-backend/internal/logger"
	"github.com/revenantio/revenant-backend/internal/server/handlers"
	"github.com/revenantio/revenant-backend/internal/server/middleware"
)

type Server struct {
	*http.Server
	router *gin.Engine
}

func New(cfg *config.Config, log *logger.Logger, db *sql.DB) *Server {
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Apply global middleware
	router.Use(middleware.LoggingMiddleware(log))
	router.Use(gin.Recovery())

	// Initialize handlers
	registerRoutes(router, log, db, cfg)

	httpServer := &http.Server{
		Handler: router,
	}

	return &Server{
		Server: httpServer,
		router: router,
	}
}

func registerRoutes(router *gin.Engine, log *logger.Logger, db *sql.DB, cfg *config.Config) {
	// Health check
	router.GET("/health", handlers.HealthCheck())

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		// Auth routes
		auth := v1.Group("/auth")
		{
			auth.POST("/register", handlers.Register(db, log))
			auth.POST("/login", handlers.Login(db, log, cfg))
		}

		// Protected routes (require auth)
		protected := v1.Group("")
		protected.Use(middleware.AuthMiddleware(cfg))
		{
			protected.GET("/users/search", handlers.SearchUser(db, log))
			protected.GET("/users/:id", handlers.GetUser(db, log))

			expenses := protected.Group("/expenses")
			{
				expenses.POST("", handlers.CreateExpense(db, log))
				expenses.GET("", handlers.ListExpenses(db, log))
				expenses.GET("/:id", handlers.GetExpense(db, log))
				expenses.PUT("/:id", handlers.UpdateExpense(db, log))
				expenses.DELETE("/:id", handlers.DeleteExpense(db, log))
				expenses.POST("/:id/participants", handlers.AddParticipant(db, log))
			}
		}
	}
}

func (s *Server) Run(addr string) error {
	s.Server.Addr = addr
	return s.Server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}
