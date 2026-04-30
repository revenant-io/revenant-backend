package middleware

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/revenantio/revenant-backend/internal/logger"
)

func LoggingMiddleware(log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		duration := time.Since(startTime)
		log.Info("HTTP Request", map[string]interface{}{
			"method":       c.Request.Method,
			"path":         c.Request.URL.Path,
			"status":       c.Writer.Status(),
			"duration_ms":  duration.Milliseconds(),
			"client_ip":    c.ClientIP(),
		})
	}
}
