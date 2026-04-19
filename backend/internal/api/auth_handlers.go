package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/szidorpatrik/rift-chat/backend/rift-chat-api/internal/auth"
	"github.com/szidorpatrik/rift-chat/backend/rift-chat-api/internal/database"
)

const accessTokenExpiry = time.Hour

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

func (cfg *ApiConfig) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate user input
	if code, err := cfg.validateLoginRequest(req); err != nil {
		RespondWithError(w, code, err.Error())
		return
	}

	// Check if user exists
	user, err := cfg.DbQueries.GetUserByUsername(r.Context(), req.Username)
	if err == sql.ErrNoRows {
		// Agianst timing attacks
		auth.HashPassword("just a silly time waster")
		RespondWithError(w, http.StatusUnauthorized, "incorrect username or password")
		return
	} else if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	// Check password
	if matching, err := auth.CheckPasswordHash(req.Password, user.PasswordHash); err != nil {
		RespondWithError(w, http.StatusUnauthorized, "something went wrong")
		return
	} else if !matching {
		RespondWithError(w, http.StatusUnauthorized, "incorrect username or password")
		return
	}

	// TODO: Refresh token rotation
	refreshToken, hashedRefreshToken := auth.MakeRefreshToken()
	params := database.CreateRefreshTokenParams{
		UserID: user.ID,
		Token:  hashedRefreshToken,
	}
	if _, err := cfg.DbQueries.CreateRefreshToken(r.Context(), params); err != nil {
		RespondWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.JwtSecret, accessTokenExpiry)
	if err != nil {
		// Failed to create access token
		RespondWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	res := LoginResponse{
		ID:           user.ID.String(),
		Username:     user.Username,
		Email:        user.Email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	RespondWithJSON(w, http.StatusOK, res)
}
