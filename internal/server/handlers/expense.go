package handlers

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/revenantio/revenant-backend/internal/logger"
	"github.com/revenantio/revenant-backend/internal/models"
	"github.com/revenantio/revenant-backend/internal/services"
	"github.com/revenantio/revenant-backend/internal/utils/validator"
)

func CreateExpense(db *sql.DB, log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		createdBy, ok := userID.(uuid.UUID)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		var req models.CreateExpenseRequest
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

		expense, err := services.CreateExpense(c.Request.Context(), db, &req, createdBy)
		if err != nil {
			log.Error("Failed to create expense", map[string]any{
				"user_id": createdBy,
				"error":   err.Error(),
			})
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, expense)
	}
}

func ListExpenses(db *sql.DB, log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		uid, ok := userID.(uuid.UUID)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		expenseType := c.Query("type")

		expenses, err := services.ListExpenses(c.Request.Context(), db, uid, expenseType)
		if err != nil {
			log.Error("Failed to list expenses", map[string]any{
				"user_id": uid,
				"error":   err.Error(),
			})
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to list expenses",
			})
			return
		}

		c.JSON(http.StatusOK, expenses)
	}
}

func GetExpense(db *sql.DB, log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		uid, ok := userID.(uuid.UUID)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		expenseID := c.Param("id")
		id, err := uuid.Parse(expenseID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid expense id",
			})
			return
		}

		expense, err := services.GetExpenseByID(c.Request.Context(), db, id, uid)
		if err != nil {
			if err.Error() == "forbidden" {
				c.JSON(http.StatusForbidden, gin.H{
					"error": "access denied",
				})
				return
			}
			log.Error("Failed to get expense", map[string]any{
				"expense_id": expenseID,
				"error":      err.Error(),
			})
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to get expense",
			})
			return
		}

		if expense == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "expense not found",
			})
			return
		}

		c.JSON(http.StatusOK, expense)
	}
}

func UpdateExpense(db *sql.DB, log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		uid, ok := userID.(uuid.UUID)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		expenseID := c.Param("id")
		id, err := uuid.Parse(expenseID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid expense id",
			})
			return
		}

		var req models.UpdateExpenseRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid request body",
			})
			return
		}

		expense, err := services.UpdateExpense(c.Request.Context(), db, id, &req, uid)
		if err != nil {
			if err.Error() == "forbidden" {
				c.JSON(http.StatusForbidden, gin.H{
					"error": "access denied",
				})
				return
			}
			log.Error("Failed to update expense", map[string]any{
				"expense_id": expenseID,
				"error":      err.Error(),
			})
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to update expense",
			})
			return
		}

		if expense == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "expense not found",
			})
			return
		}

		c.JSON(http.StatusOK, expense)
	}
}

func DeleteExpense(db *sql.DB, log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		uid, ok := userID.(uuid.UUID)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		expenseID := c.Param("id")
		id, err := uuid.Parse(expenseID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid expense id",
			})
			return
		}

		err = services.DeleteExpense(c.Request.Context(), db, id, uid)
		if err != nil {
			if err.Error() == "forbidden" {
				c.JSON(http.StatusForbidden, gin.H{
					"error": "access denied",
				})
				return
			}
			log.Error("Failed to delete expense", map[string]any{
				"expense_id": expenseID,
				"error":      err.Error(),
			})
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "failed to delete expense",
			})
			return
		}

		c.Status(http.StatusNoContent)
	}
}

func AddParticipant(db *sql.DB, log *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, ok := c.Get("user_id")
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		uid, ok := userID.(uuid.UUID)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		expenseID := c.Param("id")
		id, err := uuid.Parse(expenseID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid expense id",
			})
			return
		}

		var req models.AddParticipantRequest
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

		participant, err := services.AddParticipant(c.Request.Context(), db, id, &req, uid)
		if err != nil {
			if err.Error() == "forbidden" {
				c.JSON(http.StatusForbidden, gin.H{
					"error": "access denied",
				})
				return
			}
			log.Error("Failed to add participant", map[string]any{
				"expense_id": expenseID,
				"error":      err.Error(),
			})
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		if participant == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "expense not found",
			})
			return
		}

		c.JSON(http.StatusCreated, participant)
	}
}
