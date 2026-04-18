package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/google/uuid"
	"github.com/szidorpatrik/rift-chat/backend/rift-chat-api/internal/auth"
	"github.com/szidorpatrik/rift-chat/backend/rift-chat-api/internal/database"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

const minEntropyBits = 60

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,30}$`)
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func (cfg *ApiConfig) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if code, err := validateCreateUserRequest(req); err != nil {
		RespondWithError(w, code, err.Error())
		return
	}

	if code, err := cfg.checkUserExists(req, r.Context()); err != nil {
		RespondWithError(w, code, err.Error())
		return
	}

	userId, err := uuid.NewV7()
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	params := database.CreateUserParams{
		ID:           userId,
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passwordHash,
	}

	createdUser, err := cfg.DbQueries.CreateUser(r.Context(), params)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	res := CreateUserResponse{
		ID:        createdUser.ID.String(),
		Username:  createdUser.Username,
		Email:     createdUser.Email,
		CreatedAt: createdUser.CreatedAt.Time.UTC(),
		UpdatedAt: createdUser.UpdatedAt.Time.UTC(),
	}
	RespondWithJSON(w, http.StatusCreated, res)
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

func (cfg *ApiConfig) checkUserExists(req CreateUserRequest, ctx context.Context) (int, error) {
	// Check if Username exists
	if _, err := cfg.DbQueries.GetUserByUsername(ctx, req.Username); err == nil {
		return http.StatusConflict, fmt.Errorf("username already exists")
	} else if err != sql.ErrNoRows {
		// Something went wrong with the DB.
		return http.StatusInternalServerError, fmt.Errorf("something went wrong")
	}

	// Check if Email exists
	if _, err := cfg.DbQueries.GetUserByEmail(ctx, req.Email); err == nil {
		return http.StatusConflict, fmt.Errorf("email already exists")
	} else if err != sql.ErrNoRows {
		// Something went wrong with the DB.
		return http.StatusInternalServerError, fmt.Errorf("something went wrong")
	}
	return http.StatusOK, nil
}
