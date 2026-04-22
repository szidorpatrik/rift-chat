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
const refreshTokenExpiry = time.Hour * 24 * 30 // 30 days

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,30}$`)
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

func (cfg *ApiConfig) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		cfg.Logger.Warn("invalid JSON body",
			"event", EventRegisterFailed.String(),
			"ip", r.RemoteAddr,
			"error", err.Error(),
		)
		RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if code, err := validateCreateUserRequest(req); err != nil {
		cfg.Logger.Warn("registration validation failed",
			"event", EventRegisterFailed.String(),
			"ip", r.RemoteAddr,
			"error", err.Error(),
		)
		RespondWithError(w, code, err.Error())
		return
	}

	if code, err := cfg.checkUserExists(req, r.Context()); err != nil {
		if code == http.StatusConflict {
			cfg.Logger.Warn("registration conflict",
				"event", EventRegisterFailed.String(),
				"username", req.Username,
				"email", req.Email,
				"ip", r.RemoteAddr,
				"error", err.Error(),
			)
		} else {
			cfg.Logger.Error("database error checking user existence",
				"event", EventDbError.String(),
				"username", req.Username,
				"email", req.Email,
				"ip", r.RemoteAddr,
				"error", err.Error(),
			)
		}
		RespondWithError(w, code, err.Error())
		return
	}

	userId, err := uuid.NewV7()
	if err != nil {
		cfg.Logger.Error("failed to generate UUID",
			"event", EventRegisterFailed.String(),
			"user_id", userId,
			"username", req.Username,
			"email", req.Email,
			"ip", r.RemoteAddr,
			"error", err.Error(),
		)
		RespondWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		cfg.Logger.Error("failed to hash password",
			"event", EventRegisterFailed.String(),
			"user_id", userId,
			"username", req.Username,
			"email", req.Email,
			"ip", r.RemoteAddr,
			"error", err.Error(),
		)
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
		cfg.Logger.Error("failed to create user in database",
			"event", EventDbError.String(),
			"user_id", params.ID,
			"username", params.Username,
			"email", params.Email,
			"ip", r.RemoteAddr,
			"error", err.Error(),
		)
		RespondWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	cfg.Logger.Info("user registered successfully",
		"event", EventRegister.String(),
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
		cfg.Logger.Warn("invalid JSON body in login",
			"event", EventLoginFailed.String(),
			"ip", r.RemoteAddr,
			"error", err.Error(),
		)
		RespondWithError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Validate user input
	if code, err := cfg.validateLoginRequest(req); err != nil {
		cfg.Logger.Warn("login validation failed",
			"event", EventLoginFailed.String(),
			"ip", r.RemoteAddr,
			"error", err.Error(),
		)
		RespondWithError(w, code, err.Error())
		return
	}

	// Check if user exists
	user, err := cfg.DbQueries.GetUserByUsername(r.Context(), req.Username)
	if err == sql.ErrNoRows {
		// Against timing attacks
		auth.HashPassword("just a silly time waster")
		cfg.Logger.Warn("failed login attempt - user not found",
			"event", EventDbError.String(),
			"username", req.Username,
			"ip", r.RemoteAddr,
		)
		RespondWithError(w, http.StatusUnauthorized, "incorrect username or password")
		return
	} else if err != nil {
		cfg.Logger.Error("database error fetching user during login",
			"event", EventDbError.String(),
			"username", req.Username,
			"ip", r.RemoteAddr,
			"error", err.Error(),
		)
		RespondWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	// Check password
	if matching, err := auth.CheckPasswordHash(req.Password, user.PasswordHash); err != nil {
		cfg.Logger.Error("password hash comparison error",
			"event", EventLoginFailed.String(),
			"user_id", user.ID,
			"ip", r.RemoteAddr,
			"error", err.Error(),
		)
		RespondWithError(w, http.StatusUnauthorized, "something went wrong")
		return
	} else if !matching {
		cfg.Logger.Warn("failed login attempt - bad password",
			"event", EventLoginFailed.String(),
			"user_id", user.ID,
			"ip", r.RemoteAddr,
		)
		RespondWithError(w, http.StatusUnauthorized, "incorrect username or password")
		return
	}

	refreshToken, hashedRefreshToken := auth.MakeRefreshToken()
	params := database.CreateRefreshTokenParams{
		UserID: user.ID,
		Token:  hashedRefreshToken,
	}
	if _, err := cfg.DbQueries.CreateRefreshToken(r.Context(), params); err != nil {
		cfg.Logger.Error("failed to save refresh token",
			"event", EventDbError.String(),
			"user_id", user.ID,
			"error", err.Error(),
		)
		RespondWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.JwtSecret, accessTokenExpiry)
	if err != nil {
		cfg.Logger.Error("failed to sign access token",
			"event", EventAccessTokenFailedSign.String(),
			"ip", r.RemoteAddr,
			"user_id", user.ID,
			"error", err.Error(),
		)
		RespondWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	isSecure := cfg.Platform == "production" || cfg.Platform == "prod"
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		Path:     "/",
		MaxAge:   int(accessTokenExpiry.Seconds()),
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: http.SameSiteStrictMode, // Protects against CSRF
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   int(refreshTokenExpiry.Seconds()),
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: http.SameSiteStrictMode,
	})

	cfg.Logger.Info("user logged in",
		"event", EventLogin.String(),
		"ip", r.RemoteAddr,
		"user_id", user.ID,
	)

	res := LoginResponse{
		ID:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
	}
	RespondWithJSON(w, http.StatusOK, res)
}

func (cfg *ApiConfig) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		cfg.Logger.Warn("missing refresh token cookie",
			"event", EventRefreshTokenMissing.String(),
			"ip", r.RemoteAddr,
		)
		RespondWithError(w, http.StatusUnauthorized, "missing refresh token")
		return
	}

	rawRefreshToken := cookie.Value
	hashedToken := auth.HashRefreshToken([]byte(rawRefreshToken))
	dbToken, err := cfg.DbQueries.GetRefreshTokenByToken(r.Context(), hashedToken)
	if err != nil {
		if err == sql.ErrNoRows {
			cfg.Logger.Warn("refresh token not found in database",
				"event", EventRefreshTokenNotFound.String(),
				"ip", r.RemoteAddr,
			)
			RespondWithError(w, http.StatusUnauthorized, "invalid refresh token")
			return
		}

		cfg.Logger.Warn("invalid refresh token",
			"event", EventRefreshTokenInvalid.String(),
			"ip", r.RemoteAddr,
		)
		RespondWithError(w, http.StatusUnauthorized, "invalid refresh token")
		return
	}

	if code, err := cfg.validateRefreshToken(r.Context(), dbToken.UserID, rawRefreshToken); err != nil {
		cfg.Logger.Warn("refresh token validation failed",
			"event", EventRefreshTokenInvalid.String(),
			"ip", r.RemoteAddr,
			"user_id", dbToken.UserID,
			"error", err.Error(),
		)
		RespondWithError(w, code, err.Error())
		return
	}

	newToken, err := auth.MakeJWT(dbToken.UserID, cfg.JwtSecret, accessTokenExpiry)
	if err != nil {
		cfg.Logger.Error("failed to sign new access token",
			"event", EventAccessTokenFailedSign.String(),
			"ip", r.RemoteAddr,
			"user_id", dbToken.UserID,
			"error", err.Error(),
		)
		RespondWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	newRefreshToken, newHashedRefreshToken := auth.MakeRefreshToken()
	createParams := database.CreateRefreshTokenParams{
		UserID: dbToken.UserID,
		Token:  newHashedRefreshToken,
	}
	if _, err := cfg.DbQueries.CreateRefreshToken(r.Context(), createParams); err != nil {
		cfg.Logger.Error("failed to create new refresh token during rotation",
			"event", EventDbError.String(),
			"user_id", dbToken.UserID,
			"error", err.Error(),
		)
		RespondWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	if err := cfg.DbQueries.RevokeRefreshToken(r.Context(), hashedToken); err != nil {
		cfg.Logger.Error("failed to revoke old refresh token during rotation",
			"event", EventDbError.String(),
			"user_id", dbToken.UserID,
			"error", err.Error(),
		)
		RespondWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}

	isSecure := cfg.Platform == "production" || cfg.Platform == "prod"
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    newToken,
		Path:     "/",
		MaxAge:   int(accessTokenExpiry.Seconds()),
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: http.SameSiteStrictMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		Path:     "/",
		MaxAge:   int(refreshTokenExpiry.Seconds()),
		HttpOnly: true,
		Secure:   isSecure,
		SameSite: http.SameSiteStrictMode,
	})

	cfg.Logger.Info("access token refreshed",
		"event", EventAccessTokenRefreshed.String(),
		"ip", r.RemoteAddr,
		"user_id", dbToken.UserID,
	)

	RespondWithJSON(w, http.StatusOK, "")
}

func (cfg *ApiConfig) RevokeTokenHandler(w http.ResponseWriter, r *http.Request) {
	refreshTokenCookie, err := r.Cookie("refresh_token")
	if err == nil {
		rawRefreshToken := refreshTokenCookie.Value
		hashedToken := auth.HashRefreshToken([]byte(rawRefreshToken))

		if err := cfg.DbQueries.RevokeRefreshToken(r.Context(), hashedToken); err != nil {
			cfg.Logger.Error("failed to revoke refresh token",
				"event", EventDbError.String(),
				"ip", r.RemoteAddr,
				"error", err.Error(),
			)
		}
	}

	// Clear access token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	// Clear refresh token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	cfg.Logger.Info("refresh token revoked",
		"event", EventRefreshTokenRevoked.String(),
		"ip", r.RemoteAddr,
	)

	RespondWithJSON(w, http.StatusOK, "")
}
