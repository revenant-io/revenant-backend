package services

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"

	"github.com/revenantio/revenant-backend/internal/models"
	"github.com/revenantio/revenant-backend/internal/utils/hash"
)

func RegisterUser(ctx context.Context, db *sql.DB, req *models.CreateUserRequest) (*models.User, error) {
	hashedPassword, err := hash.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:        uuid.New(),
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	query := `
		INSERT INTO users (id, email, password, first_name, last_name, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err = db.ExecContext(ctx, query,
		user.ID,
		user.Email,
		user.Password,
		user.FirstName,
		user.LastName,
		user.CreatedAt,
		user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	// Return user without password
	user.Password = ""
	return user, nil
}

func GetUserByID(ctx context.Context, db *sql.DB, id uuid.UUID) (*models.User, error) {
	user := &models.User{}

	query := `
		SELECT id, email, first_name, last_name, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	err := db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

func GetUserByEmail(ctx context.Context, db *sql.DB, email string) (*models.User, error) {
	user := &models.User{}

	query := `
		SELECT id, email, password, first_name, last_name, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	err := db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}
