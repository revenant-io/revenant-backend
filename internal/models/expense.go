package models

import (
	"time"

	"github.com/google/uuid"
)

type Expense struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Currency    string    `json:"currency"`
	Category    string    `json:"category"`
	Date        time.Time `json:"date"`
	SplitType   string    `json:"split_type"`
	CreatedBy   uuid.UUID `json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type ExpenseParticipant struct {
	ID         uuid.UUID `json:"id"`
	ExpenseID  uuid.UUID `json:"expense_id"`
	UserID     uuid.UUID `json:"user_id"`
	Username   string    `json:"username"`
	SplitValue float64   `json:"split_value"`
	CreatedAt  time.Time `json:"created_at"`
}

type ExpenseWithParticipants struct {
	Expense
	Participants []ExpenseParticipant `json:"participants"`
}

type ParticipantInput struct {
	Username   string  `json:"username" validate:"required"`
	SplitValue float64 `json:"split_value"`
}

type CreateExpenseRequest struct {
	Title        string             `json:"title" validate:"required"`
	Description  string             `json:"description"`
	Amount       float64            `json:"amount" validate:"required,gt=0"`
	Currency     string             `json:"currency" validate:"required,len=3"`
	Category     string             `json:"category"`
	Date         string             `json:"date" validate:"required"`
	SplitType    string             `json:"split_type"`
	Participants []ParticipantInput `json:"participants"`
}

type AddParticipantRequest struct {
	Username   string  `json:"username" validate:"required"`
	SplitValue float64 `json:"split_value"`
}

type UpdateExpenseRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	Currency    string  `json:"currency"`
	Category    string  `json:"category"`
	Date        string  `json:"date"`
	SplitType   string  `json:"split_type"`
}
