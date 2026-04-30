package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/revenantio/revenant-backend/internal/config"
	"github.com/revenantio/revenant-backend/internal/models"
	"github.com/revenantio/revenant-backend/internal/utils/hash"
	"github.com/revenantio/revenant-backend/internal/utils/jwt"
)

func Login(ctx context.Context, db *sql.DB, req *models.LoginRequest, cfg *config.Config) (string, *models.User, error) {
	user, err := GetUserByEmail(ctx, db, req.Email)
	if err != nil {
		return "", nil, err
	}

	if user == nil {
		return "", nil, errors.New("user not found")
	}

	if !hash.VerifyPassword(user.Password, req.Password) {
		return "", nil, errors.New("invalid password")
	}

	token, err := jwt.GenerateToken(user.ID, cfg.JWT.Secret, cfg.JWT.Expiration)
	if err != nil {
		return "", nil, err
	}

	user.Password = ""
	return token, user, nil
}
