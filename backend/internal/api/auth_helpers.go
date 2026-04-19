package api

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/szidorpatrik/rift-chat/backend/rift-chat-api/internal/auth"
	"github.com/szidorpatrik/rift-chat/backend/rift-chat-api/internal/database"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

const minEntropyBits = 60

func (cfg *ApiConfig) validateLoginRequest(req LoginRequest) (int, error) {
	if req.Username == "" {
		return http.StatusBadRequest, fmt.Errorf("username is required")
	}
	if req.Password == "" {
		return http.StatusBadRequest, fmt.Errorf("password is required")
	}
	if !usernameRegex.MatchString(req.Username) {
		return http.StatusBadRequest, fmt.Errorf("username must be 3-30 characters and contain only letters, numbers, and underscores")
	}
	return http.StatusOK, nil
}

func (cfg *ApiConfig) checkUserExists(req CreateUserRequest, ctx context.Context) (int, error) {
	params := database.GetUserByUsernameOrEmailParams{
		Username: req.Username,
		Email:    req.Email,
	}
	existingUser, err := cfg.DbQueries.GetUserByUsernameOrEmail(ctx, params)
	if err == nil {
		if existingUser.Username == req.Username {
			return http.StatusConflict, fmt.Errorf("username already exists")
		}
		if existingUser.Email == req.Email {
			return http.StatusConflict, fmt.Errorf("email already exists")
		}
	} else if err != sql.ErrNoRows {
		return http.StatusInternalServerError, fmt.Errorf("something went wrong")
	}
	return http.StatusOK, nil
}

func validateCreateUserRequest(req CreateUserRequest) (int, error) {
	if req.Username == "" || req.Password == "" || req.Email == "" {
		return http.StatusBadRequest, fmt.Errorf("username, email, and password are required")
	}
	if !usernameRegex.MatchString(req.Username) {
		return http.StatusBadRequest, fmt.Errorf("username must be 3-30 characters and contain only letters, numbers, and underscores")
	}
	if !emailRegex.MatchString(req.Email) {
		return http.StatusBadRequest, fmt.Errorf("invalid email address format")
	}
	if err := passwordvalidator.Validate(req.Password, minEntropyBits); err != nil {
		return http.StatusBadRequest, fmt.Errorf("password does not meet complexity requirements")
	}
	return http.StatusOK, nil
}

func (cfg *ApiConfig) validateRefreshToken(ctx context.Context, userId uuid.UUID, rawRefreshToken string) (database.RefreshToken, int, error) {
	if rawRefreshToken == "" {
		return database.RefreshToken{}, http.StatusUnauthorized, fmt.Errorf("invalid refresh token")
	}

	hashedToken := auth.HashRefreshToken([]byte(rawRefreshToken))
	token, err := cfg.DbQueries.GetRefreshTokenByToken(ctx, hashedToken)
	if err != nil {
		return database.RefreshToken{}, http.StatusInternalServerError, fmt.Errorf("something went wrong")
	}

    if !checkRefreshTokenValid(token) {
        return database.RefreshToken{}, http.StatusUnauthorized, fmt.Errorf("refresh token expired or revoked")
    }
    return token, http.StatusOK, nil
}

func checkRefreshTokenValid(token database.RefreshToken) bool {
	return !token.RevokedAt.Valid && !token.ExpiresAt.Before(time.Now())
}
