package services

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/revenantio/revenant-backend/internal/models"
)

func getParticipantsForExpense(ctx context.Context, db *sql.DB, expenseID uuid.UUID) ([]models.ExpenseParticipant, error) {
	query := `
		SELECT ep.id, ep.expense_id, ep.user_id, u.username, ep.split_value, ep.created_at
		FROM expense_participants ep
		JOIN users u ON u.id = ep.user_id
		WHERE ep.expense_id = $1
	`

	rows, err := db.QueryContext(ctx, query, expenseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []models.ExpenseParticipant
	for rows.Next() {
		var p models.ExpenseParticipant
		if err := rows.Scan(
			&p.ID,
			&p.ExpenseID,
			&p.UserID,
			&p.Username,
			&p.SplitValue,
			&p.CreatedAt,
		); err != nil {
			return nil, err
		}
		participants = append(participants, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	if participants == nil {
		participants = []models.ExpenseParticipant{}
	}
	return participants, nil
}

func scanExpense(row *sql.Row) (*models.Expense, error) {
	e := &models.Expense{}
	err := row.Scan(
		&e.ID,
		&e.Title,
		&e.Description,
		&e.Amount,
		&e.Currency,
		&e.Category,
		&e.Date,
		&e.SplitType,
		&e.CreatedBy,
		&e.CreatedAt,
		&e.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func CreateExpense(ctx context.Context, db *sql.DB, req *models.CreateExpenseRequest, createdBy uuid.UUID) (*models.ExpenseWithParticipants, error) {
	if len(req.Participants) > 0 && req.SplitType == "" {
		return nil, errors.New("split_type is required when participants are present")
	}

	if req.SplitType == "percentage" {
		var total float64
		for _, p := range req.Participants {
			total += p.SplitValue
		}
		if total > 100.0001 {
			return nil, errors.New("percentage splits must not exceed 100%")
		}
	}

	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, errors.New("invalid date format, expected YYYY-MM-DD")
	}

	splitType := req.SplitType
	if splitType == "" {
		splitType = "personal"
	}

	expense := &models.Expense{
		ID:          uuid.New(),
		Title:       req.Title,
		Description: req.Description,
		Amount:      req.Amount,
		Currency:    req.Currency,
		Category:    req.Category,
		Date:        date,
		SplitType:   splitType,
		CreatedBy:   createdBy,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	insertExpenseQuery := `
		INSERT INTO expenses (id, title, description, amount, currency, category, date, split_type, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err = tx.ExecContext(ctx, insertExpenseQuery,
		expense.ID,
		expense.Title,
		expense.Description,
		expense.Amount,
		expense.Currency,
		expense.Category,
		expense.Date,
		expense.SplitType,
		expense.CreatedBy,
		expense.CreatedAt,
		expense.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	var participants []models.ExpenseParticipant
	for _, pi := range req.Participants {
		user, err := GetUserByUsername(ctx, db, pi.Username)
		if err != nil {
			return nil, err
		}
		if user == nil {
			return nil, errors.New("participant user not found: " + pi.Username)
		}

		p := models.ExpenseParticipant{
			ID:         uuid.New(),
			ExpenseID:  expense.ID,
			UserID:     user.ID,
			Username:   user.Username,
			SplitValue: pi.SplitValue,
			CreatedAt:  time.Now(),
		}

		insertParticipantQuery := `
			INSERT INTO expense_participants (id, expense_id, user_id, split_value, created_at)
			VALUES ($1, $2, $3, $4, $5)
		`
		_, err = tx.ExecContext(ctx, insertParticipantQuery,
			p.ID,
			p.ExpenseID,
			p.UserID,
			p.SplitValue,
			p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		participants = append(participants, p)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	if participants == nil {
		participants = []models.ExpenseParticipant{}
	}

	return &models.ExpenseWithParticipants{
		Expense:      *expense,
		Participants: participants,
	}, nil
}

func ListExpenses(ctx context.Context, db *sql.DB, userID uuid.UUID, expenseType string) ([]models.ExpenseWithParticipants, error) {
	var query string
	switch expenseType {
	case "personal":
		query = `
			SELECT id, title, description, amount, currency, category, date, split_type, created_by, created_at, updated_at
			FROM expenses
			WHERE created_by = $1 AND split_type = 'personal'
			ORDER BY date DESC, created_at DESC
		`
	case "shared":
		query = `
			SELECT DISTINCT e.id, e.title, e.description, e.amount, e.currency, e.category, e.date, e.split_type, e.created_by, e.created_at, e.updated_at
			FROM expenses e
			LEFT JOIN expense_participants ep ON ep.expense_id = e.id
			WHERE (e.created_by = $1 OR ep.user_id = $1) AND e.split_type != 'personal'
			ORDER BY e.date DESC, e.created_at DESC
		`
	default:
		query = `
			SELECT DISTINCT e.id, e.title, e.description, e.amount, e.currency, e.category, e.date, e.split_type, e.created_by, e.created_at, e.updated_at
			FROM expenses e
			LEFT JOIN expense_participants ep ON ep.expense_id = e.id
			WHERE e.created_by = $1 OR ep.user_id = $1
			ORDER BY e.date DESC, e.created_at DESC
		`
	}

	rows, err := db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []models.ExpenseWithParticipants
	for rows.Next() {
		var e models.Expense
		if err := rows.Scan(
			&e.ID,
			&e.Title,
			&e.Description,
			&e.Amount,
			&e.Currency,
			&e.Category,
			&e.Date,
			&e.SplitType,
			&e.CreatedBy,
			&e.CreatedAt,
			&e.UpdatedAt,
		); err != nil {
			return nil, err
		}
		results = append(results, models.ExpenseWithParticipants{Expense: e})
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i := range results {
		participants, err := getParticipantsForExpense(ctx, db, results[i].ID)
		if err != nil {
			return nil, err
		}
		results[i].Participants = participants
	}

	if results == nil {
		results = []models.ExpenseWithParticipants{}
	}
	return results, nil
}

func GetExpenseByID(ctx context.Context, db *sql.DB, id uuid.UUID, userID uuid.UUID) (*models.ExpenseWithParticipants, error) {
	query := `
		SELECT id, title, description, amount, currency, category, date, split_type, created_by, created_at, updated_at
		FROM expenses
		WHERE id = $1
	`

	expense, err := scanExpense(db.QueryRowContext(ctx, query, id))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Check if user is creator or participant
	if expense.CreatedBy != userID {
		var participantID uuid.UUID
		checkQuery := `SELECT user_id FROM expense_participants WHERE expense_id = $1 AND user_id = $2`
		err := db.QueryRowContext(ctx, checkQuery, id, userID).Scan(&participantID)
		if err == sql.ErrNoRows {
			return nil, errors.New("forbidden")
		}
		if err != nil {
			return nil, err
		}
	}

	participants, err := getParticipantsForExpense(ctx, db, expense.ID)
	if err != nil {
		return nil, err
	}

	return &models.ExpenseWithParticipants{
		Expense:      *expense,
		Participants: participants,
	}, nil
}

func UpdateExpense(ctx context.Context, db *sql.DB, id uuid.UUID, req *models.UpdateExpenseRequest, userID uuid.UUID) (*models.Expense, error) {
	// Fetch existing expense to verify ownership
	fetchQuery := `
		SELECT id, title, description, amount, currency, category, date, split_type, created_by, created_at, updated_at
		FROM expenses
		WHERE id = $1
	`
	existing, err := scanExpense(db.QueryRowContext(ctx, fetchQuery, id))
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if existing.CreatedBy != userID {
		return nil, errors.New("forbidden")
	}

	// Apply updates, keeping existing values for empty fields
	title := existing.Title
	if req.Title != "" {
		title = req.Title
	}

	description := existing.Description
	if req.Description != "" {
		description = req.Description
	}

	amount := existing.Amount
	if req.Amount > 0 {
		amount = req.Amount
	}

	currency := existing.Currency
	if req.Currency != "" {
		currency = req.Currency
	}

	category := existing.Category
	if req.Category != "" {
		category = req.Category
	}

	date := existing.Date
	if req.Date != "" {
		parsed, err := time.Parse("2006-01-02", req.Date)
		if err != nil {
			return nil, errors.New("invalid date format, expected YYYY-MM-DD")
		}
		date = parsed
	}

	splitType := existing.SplitType
	if req.SplitType != "" {
		splitType = req.SplitType
	}

	updatedAt := time.Now()

	updateQuery := `
		UPDATE expenses
		SET title = $1, description = $2, amount = $3, currency = $4, category = $5, date = $6, split_type = $7, updated_at = $8
		WHERE id = $9
	`
	_, err = db.ExecContext(ctx, updateQuery,
		title,
		description,
		amount,
		currency,
		category,
		date,
		splitType,
		updatedAt,
		id,
	)
	if err != nil {
		return nil, err
	}

	existing.Title = title
	existing.Description = description
	existing.Amount = amount
	existing.Currency = currency
	existing.Category = category
	existing.Date = date
	existing.SplitType = splitType
	existing.UpdatedAt = updatedAt

	return existing, nil
}

func DeleteExpense(ctx context.Context, db *sql.DB, id uuid.UUID, userID uuid.UUID) error {
	fetchQuery := `SELECT created_by FROM expenses WHERE id = $1`
	var createdBy uuid.UUID
	err := db.QueryRowContext(ctx, fetchQuery, id).Scan(&createdBy)
	if err == sql.ErrNoRows {
		return nil
	}
	if err != nil {
		return err
	}

	if createdBy != userID {
		return errors.New("forbidden")
	}

	_, err = db.ExecContext(ctx, `DELETE FROM expenses WHERE id = $1`, id)
	return err
}

func AddParticipant(ctx context.Context, db *sql.DB, expenseID uuid.UUID, req *models.AddParticipantRequest, requestingUserID uuid.UUID) (*models.ExpenseParticipant, error) {
	// Verify the expense exists and the requesting user is the creator
	fetchQuery := `SELECT created_by FROM expenses WHERE id = $1`
	var createdBy uuid.UUID
	err := db.QueryRowContext(ctx, fetchQuery, expenseID).Scan(&createdBy)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if createdBy != requestingUserID {
		return nil, errors.New("forbidden")
	}

	// Look up participant by username
	user, err := GetUserByUsername(ctx, db, req.Username)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("participant user not found: " + req.Username)
	}

	p := models.ExpenseParticipant{
		ID:         uuid.New(),
		ExpenseID:  expenseID,
		UserID:     user.ID,
		Username:   user.Username,
		SplitValue: req.SplitValue,
		CreatedAt:  time.Now(),
	}

	insertQuery := `
		INSERT INTO expense_participants (id, expense_id, user_id, split_value, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err = db.ExecContext(ctx, insertQuery,
		p.ID,
		p.ExpenseID,
		p.UserID,
		p.SplitValue,
		p.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	return &p, nil
}
