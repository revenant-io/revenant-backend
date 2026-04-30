package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/revenantio/revenant-backend/internal/config"
	"github.com/revenantio/revenant-backend/internal/logger"
	"github.com/revenantio/revenant-backend/internal/models"
	"github.com/revenantio/revenant-backend/internal/services"
	"github.com/revenantio/revenant-backend/internal/utils/validator"
)

func Register(db *sql.DB, log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.CreateUserRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body",
			})
			return
		}

		if err := validator.Validate(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		user, err := services.RegisterUser(c.Request.Context(), db, &req)
		if err != nil {
			log.Error("Failed to register user", map[string]interface{}{
				"email": req.Email,
				"error": err.Error(),
			})
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to register user",
			})
			return
		}

		c.JSON(http.StatusCreated, user)
	}
}

func Login(db *sql.DB, log *logger.Logger, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req models.LoginRequest

		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body",
			})
			return
		}

		if err := validator.Validate(req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		token, user, err := services.Login(c.Request.Context(), db, &req, cfg)
		if err != nil {
			log.Warn("Failed login attempt", map[string]interface{}{
				"email": req.Email,
			})
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "invalid credentials",
			})
			return
		}

		c.JSON(http.StatusOK, models.AuthResponse{
			Token: token,
			User:  user,
		})
	}
}
