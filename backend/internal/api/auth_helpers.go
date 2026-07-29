package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/szidorpatrik/rift-chat/backend/rift-chat-api/internal/auth"
	"github.com/szidorpatrik/rift-chat/backend/rift-chat-api/internal/database"
	passwordvalidator "github.com/wagslane/go-password-validator"
)

const minEntropyBits = 60

func (cfg *ApiConfig) validateLoginRequest(req LoginRequest) (status_code int, err error) {
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

func (cfg *ApiConfig) checkUserExists(req CreateUserRequest, ctx context.Context) (status_code int, err error) {
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

func validateCreateUserRequest(req CreateUserRequest) (status_code int, err error) {
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

func (cfg *ApiConfig) validateRefreshToken(ctx context.Context, userId uuid.UUID, rawRefreshToken string) (status_code int, err error) {
	if rawRefreshToken == "" {
		return http.StatusUnauthorized, fmt.Errorf("invalid refresh token")
	}

	hashedToken := auth.HashRefreshToken([]byte(rawRefreshToken))
	token, err := cfg.DbQueries.GetRefreshTokenByToken(ctx, hashedToken)
	if err != nil {
		if err == sql.ErrNoRows {
			return http.StatusUnauthorized, fmt.Errorf("invalid refresh token")
		}
		return http.StatusInternalServerError, fmt.Errorf("something went wrong")
	}

	if token.UserID != userId {
		return http.StatusUnauthorized, fmt.Errorf("invalid refresh token")
	}

	if !checkRefreshTokenValid(token) {
		return http.StatusUnauthorized, fmt.Errorf("refresh token expired or revoked")
	}

	return http.StatusOK, nil
}

func checkRefreshTokenValid(token database.RefreshToken) bool {
	return !token.RevokedAt.Valid && !token.ExpiresAt.Before(time.Now())
}

func buildRefreshTokenParams(userID uuid.UUID, token string) database.CreateRefreshTokenParams {
	return database.CreateRefreshTokenParams{
		UserID: userID,
		Token:  token,
	}
}

// Creates a cookie for JWTs and refresh tokens
func (cfg *ApiConfig) setAuthCookie(w http.ResponseWriter, name, value string, duration time.Duration) {
	isSecure := cfg.Platform == "production" || cfg.Platform == "prod"

	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		MaxAge:   int(duration.Seconds()),
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: http.SameSiteStrictMode,
	})
}

func (cfg *ApiConfig) setAuthCookiePair(w http.ResponseWriter, accessToken, refreshToken string) {
	cfg.setAuthCookie(w, cookieAccessToken, accessToken, accessTokenExpiry)
	cfg.setAuthCookie(w, cookieRefreshToken, refreshToken, refreshTokenExpiry)
}

// Overwrites an existing cookie to instantly expire it
func (cfg *ApiConfig) clearAuthCookie(w http.ResponseWriter, name string) {
	isSecure := cfg.Platform == "production" || cfg.Platform == "prod"

	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: http.SameSiteStrictMode,
	})
}

func (cfg *ApiConfig) clearAuthCookies(w http.ResponseWriter) {
	cfg.clearAuthCookie(w, cookieAccessToken)
	cfg.clearAuthCookie(w, cookieRefreshToken)
}

func (cfg *ApiConfig) handleAuthError(w http.ResponseWriter, r *http.Request, code int, msg string, err error, event EventType, isError bool, shouldRespondErrMsg bool, extraArgs ...any) {
	args := append([]any{
		"event", event.String(),
		"ip", GetClientIP(r),
		"error", err.Error(),
	}, extraArgs...)

	if isError {
		cfg.Logger.Error(msg, args...)
	} else {
		cfg.Logger.Warn(msg, args...)
	}

	if EventType(code) == EventDbError || !shouldRespondErrMsg {
		msg = "something went wrong"
	}

	RespondWithError(w, code, msg)
}

// Reads the JSON body into a generic struct T, then runs the provided validation function.
func decodeAndValidate[T any](r *http.Request, validator func(T) (int, error)) (T, int, error) {
	var req T

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return req, http.StatusBadRequest, fmt.Errorf("malformed json: %w", err)
	}

	if code, err := validator(req); err != nil {
		return req, code, err
	}

	return req, http.StatusOK, nil
}
