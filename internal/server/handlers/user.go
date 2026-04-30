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
			log.Error("Failed to get user", map[string]interface{}{
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
