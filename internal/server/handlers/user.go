package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/revenantio/revenant-backend/internal/logger"
	"github.com/revenantio/revenant-backend/internal/services"
)

func GetUser(db *sql.DB, log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Param("id")

		id, err := uuid.Parse(userID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid user id",
			})
			return
		}

		user, err := services.GetUserByID(c.Request.Context(), db, id)
		if err != nil {
			log.Error("Failed to get user", map[string]any{
				"user_id": userID,
				"error":   err.Error(),
			})
			c.JSON(http.StatusNotFound, gin.H{
				"error": "user not found",
			})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}

func SearchUser(db *sql.DB, log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.Query("username")
		if len(username) < 1 {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "username query parameter is required",
			})
			return
		}

		user, err := services.SearchUserByUsername(c.Request.Context(), db, username)
		if err != nil {
			log.Error("Failed to search user", map[string]any{
				"username": username,
				"error":    err.Error(),
			})
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to search user",
			})
			return
		}

		if user == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "user not found",
			})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
