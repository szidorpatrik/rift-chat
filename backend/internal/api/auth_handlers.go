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
		cfg.Logger.Warn("invalid JSON body in registration", "ip", r.RemoteAddr, "error", err.Error())
		RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if code, err := validateCreateUserRequest(req); err != nil {
		cfg.Logger.Warn("registration validation failed", "ip", r.RemoteAddr, "error", err.Error())
		RespondWithError(w, code, err.Error())
		return
	}

	if code, err := cfg.checkUserExists(req, r.Context()); err != nil {
		if code == http.StatusConflict {
			cfg.Logger.Warn("registration conflict", "ip", r.RemoteAddr, "error", err.Error())
		} else {
			cfg.Logger.Error("database error checking user existence", "error", err.Error())
		}
		RespondWithError(w, code, err.Error())
		return
	}

	userId, err := uuid.NewV7()
	if err != nil {
		cfg.Logger.Error("failed to generate UUID", "error", err.Error())
		RespondWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		cfg.Logger.Error("failed to hash password", "error", err.Error())
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
		cfg.Logger.Error("failed to create user in database", "error", err.Error())
		RespondWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	cfg.Logger.Info("user registered successfully",
		"event", "user_registration",
		"user_id", createdUser.ID,
		"ip", r.RemoteAddr,
	)

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
		cfg.Logger.Warn("invalid JSON body in login", "ip", r.RemoteAddr, "error", err.Error())
		RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate user input
	if code, err := cfg.validateLoginRequest(req); err != nil {
		cfg.Logger.Warn("login validation failed", "ip", r.RemoteAddr, "error", err.Error())
		RespondWithError(w, code, err.Error())
		return
	}

	// Check if user exists
	user, err := cfg.DbQueries.GetUserByUsername(r.Context(), req.Username)
	if err == sql.ErrNoRows {
		// Against timing attacks
		auth.HashPassword("just a silly time waster")
		cfg.Logger.Warn("failed login attempt - user not found", "username", req.Username, "ip", r.RemoteAddr)
		RespondWithError(w, http.StatusUnauthorized, "incorrect username or password")
		return
	} else if err != nil {
		cfg.Logger.Error("database error fetching user during login", "username", req.Username, "error", err.Error())
		RespondWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	// Check password
	if matching, err := auth.CheckPasswordHash(req.Password, user.PasswordHash); err != nil {
		cfg.Logger.Error("password hash comparison error", "user_id", user.ID, "error", err.Error())
		RespondWithError(w, http.StatusUnauthorized, "something went wrong")
		return
	} else if !matching {
		cfg.Logger.Warn("failed login attempt - bad password", "username", req.Username, "ip", r.RemoteAddr)
		RespondWithError(w, http.StatusUnauthorized, "incorrect username or password")
		return
	}

	// Refresh token rotation
	refreshToken, hashedRefreshToken := auth.MakeRefreshToken()
	params := database.CreateRefreshTokenParams{
		UserID: user.ID,
		Token:  hashedRefreshToken,
	}
	if _, err := cfg.DbQueries.CreateRefreshToken(r.Context(), params); err != nil {
		cfg.Logger.Error("failed to save refresh token", "user_id", user.ID, "error", err.Error())
		RespondWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.JwtSecret, accessTokenExpiry)
	if err != nil {
		cfg.Logger.Error("failed to sign access token", "user_id", user.ID, "error", err.Error())
		RespondWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	cfg.Logger.Info("user logged in successfully",
		"event", "user_login",
		"ip", r.RemoteAddr,
		"user_id", user.ID,
	)

	res := LoginResponse{
		ID:           user.ID.String(),
		Username:     user.Username,
		Email:        user.Email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	RespondWithJSON(w, http.StatusOK, res)
}
