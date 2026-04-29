package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/szidorpatrik/rift-chat/backend/rift-chat-api/internal/auth"
	"github.com/szidorpatrik/rift-chat/backend/rift-chat-api/internal/database"
)

func (cfg *ApiConfig) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	req, code, err := decodeAndValidate(r, validateCreateUserRequest)
	if err != nil {
		msg := err.Error()
		if code != http.StatusBadRequest {
			msg = "registration validation failed"
		}
		cfg.handleAuthError(w, r, code, msg, err, EventRegisterFailed, false, true)
		return
	}

	if code, err := cfg.checkUserExists(req, r.Context()); err != nil {
		if code == http.StatusConflict {
			cfg.handleAuthError(w, r, code, "username or email already exists", err, EventRegisterFailed, false, true, "username", req.Username, "email", req.Email)
			return
		} else {
			cfg.handleAuthError(w, r, http.StatusInternalServerError, "database error checking user existence", err, EventDbError, true, false, "username", req.Username, "email", req.Email)
			return
		}
	}

	userId, err := uuid.NewV7()
	if err != nil {
		cfg.handleAuthError(w, r, http.StatusInternalServerError, "failed to generate UUID", err, EventRegisterFailed, true, false, "username", req.Username, "email", req.Email)
		return
	}

	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		cfg.handleAuthError(w, r, http.StatusInternalServerError, "failed to hash password", err, EventRegisterFailed, true, false, "username", req.Username, "email", req.Email)
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
		cfg.handleAuthError(w, r, http.StatusInternalServerError, "failed to create user in database", err, EventDbError, true, false, "username", req.Username, "email", req.Email)
		return
	}

	cfg.LogInfo(r, "user registered successfully", createdUser.ID, EventRegister)

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
	req, code, err := decodeAndValidate(r, cfg.validateLoginRequest)
	if err != nil {
		msg := "invalid request body"
		if code != http.StatusBadRequest {
			msg = "login validation failed"
		}
		cfg.handleAuthError(w, r, code, msg, err, EventLoginFailed, false, true)
		return
	}

	// Check if user exists
	user, err := cfg.DbQueries.GetUserByUsername(r.Context(), req.Username)
	if err == sql.ErrNoRows {
		// Against timing attacks
		auth.HashPassword("just a silly time waster")
		cfg.handleAuthError(w, r, http.StatusUnauthorized, "user not found", fmt.Errorf("incorrect username or password"), EventLoginFailed, false, true, "username", req.Username)
		return
	} else if err != nil {
		cfg.handleAuthError(w, r, http.StatusInternalServerError, "database error fetching user during login", err, EventDbError, true, false, "username", req.Username)
		return
	}

	// Check password
	if matching, err := auth.CheckPasswordHash(req.Password, user.PasswordHash); err != nil {
		cfg.handleAuthError(w, r, http.StatusUnauthorized, "password hash comparison error", err, EventLoginFailed, true, false, "user_id", user.ID)
		return
	} else if !matching {
		cfg.handleAuthError(w, r, http.StatusUnauthorized, "incorrect password", fmt.Errorf("incorrect username or password"), EventLoginFailed, false, true, "user_id", user.ID)
		return
	}

	refreshToken, hashedRefreshToken := auth.MakeRefreshToken()
	params := database.CreateRefreshTokenParams{
		UserID: user.ID,
		Token:  hashedRefreshToken,
	}
	if _, err := cfg.DbQueries.CreateRefreshToken(r.Context(), params); err != nil {
		cfg.handleAuthError(w, r, http.StatusInternalServerError, "failed to save refresh token", err, EventDbError, true, false, "user_id", user.ID)
		return
	}

	accessToken, err := auth.MakeJWT(user.ID, cfg.JwtSecret, accessTokenExpiry)
	if err != nil {
		cfg.handleAuthError(w, r, http.StatusInternalServerError, "failed to sign access token", err, EventAccessTokenFailedSign, true, false, "user_id", user.ID)
		return
	}

	cfg.setAuthCookie(w, cookieAccessToken, accessToken, accessTokenExpiry)
	cfg.setAuthCookie(w, cookieRefreshToken, refreshToken, refreshTokenExpiry)

	cfg.LogInfo(r, "user logged in", user.ID, EventLogin)

	res := LoginResponse{
		ID:          user.ID.String(),
		Username:    user.Username,
		Email:       user.Email,
		AvatarUrl:   user.AvatarUrl.String,
		Status:      user.Status.String,
		Description: user.Description.String,
	}
	RespondWithJSON(w, http.StatusOK, res)
}

func (cfg *ApiConfig) RefreshTokenHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(cookieRefreshToken)
	if err != nil {
		cfg.handleAuthError(w, r, http.StatusUnauthorized, "refresh token cookie is missing", fmt.Errorf("refresh token cookie is missing"), EventRefreshTokenMissing, false, true)
		return
	}

	rawRefreshToken := cookie.Value
	hashedToken := auth.HashRefreshToken([]byte(rawRefreshToken))
	dbToken, err := cfg.DbQueries.GetRefreshTokenByToken(r.Context(), hashedToken)
	if err != nil {
		if err == sql.ErrNoRows {
			cfg.LogWarn(r, "refresh token not found in database", uuid.Nil, EventRefreshTokenNotFound)
			RespondWithError(w, http.StatusUnauthorized, "invalid refresh token")
			return
		}
		cfg.handleAuthError(w, r, http.StatusUnauthorized, "invalid refresh token", err, EventRefreshTokenInvalid, false, true)
		return
	}

	if code, err := cfg.validateRefreshToken(r.Context(), dbToken.UserID, rawRefreshToken); err != nil {
		cfg.handleAuthError(w, r, code, "refresh token validation failed", err, EventRefreshTokenInvalid, false, true, "user_id", dbToken.UserID)
		return
	}

	newToken, err := auth.MakeJWT(dbToken.UserID, cfg.JwtSecret, accessTokenExpiry)
	if err != nil {
		cfg.handleAuthError(w, r, http.StatusInternalServerError, "failed to sign new access token", err, EventAccessTokenFailedSign, true, false, "user_id", dbToken.UserID)
		return
	}

	newRefreshToken, newHashedRefreshToken := auth.MakeRefreshToken()
	createParams := database.CreateRefreshTokenParams{
		UserID: dbToken.UserID,
		Token:  newHashedRefreshToken,
	}

	// Start a transaction
	tx, err := cfg.DB.BeginTx(r.Context(), nil)
	if err != nil {
		cfg.handleAuthError(w, r, http.StatusInternalServerError, "failed to begin database transaction", err, EventDbError, true, false, "user_id", dbToken.UserID)
		return
	}

	// If a return happens early, it rolls back
	defer tx.Rollback()

	// Bind the sqlc queries to this specific transaction
	qtx := cfg.DbQueries.WithTx(tx)

	if _, err := qtx.CreateRefreshToken(r.Context(), createParams); err != nil {
		cfg.handleAuthError(w, r, http.StatusInternalServerError, "failed to create new refresh token during rotation", err, EventDbError, true, false, "user_id", dbToken.UserID)
		return
	}

	if err := qtx.RevokeRefreshToken(r.Context(), hashedToken); err != nil {
		cfg.handleAuthError(w, r, http.StatusInternalServerError, "failed to revoke old refresh token during rotation", err, EventDbError, true, false, "user_id", dbToken.UserID)
		return
	}

	// If both queries succeed, commit to the database
	if err := tx.Commit(); err != nil {
		cfg.handleAuthError(w, r, http.StatusInternalServerError, "failed to commit token rotation transaction", err, EventDbError, true, false, "user_id", dbToken.UserID)
		return
	}

	cfg.setAuthCookie(w, cookieAccessToken, newToken, accessTokenExpiry)
	cfg.setAuthCookie(w, cookieRefreshToken, newRefreshToken, refreshTokenExpiry)

	cfg.LogInfo(r, "access token refreshed", dbToken.UserID, EventAccessTokenRefreshed)
	RespondNoContent(w)
}

func (cfg *ApiConfig) RevokeTokenHandler(w http.ResponseWriter, r *http.Request) {
	refreshTokenCookie, err := r.Cookie(cookieRefreshToken)
	if err == nil {
		rawRefreshToken := refreshTokenCookie.Value
		hashedToken := auth.HashRefreshToken([]byte(rawRefreshToken))

		if err := cfg.DbQueries.RevokeRefreshToken(r.Context(), hashedToken); err != nil {
			cfg.LogError(r, err, EventDbError)
		}
	}

	// Clear cookies
	cfg.clearAuthCookie(w, cookieAccessToken)
	cfg.clearAuthCookie(w, cookieRefreshToken)

	cfg.LogInfo(r, "refresh token revoked", uuid.Nil, EventRefreshTokenRevoked)

	RespondNoContent(w)
}

func (cfg *ApiConfig) SessionCheckHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := r.Context().Value(userContextKey).(uuid.UUID)
	if !ok {
		cfg.handleAuthError(w, r, http.StatusUnauthorized, "invalid session", fmt.Errorf("invalid session"), EventSessionCheckFailed, false, true)
		return
	}

	user, err := cfg.DbQueries.GetUserById(r.Context(), userID)
	if err != nil {
		if err == sql.ErrNoRows {
			cfg.LogWarn(r, "session check failed: user not found", userID, EventSessionCheckFailed)
			RespondWithError(w, http.StatusUnauthorized, "invalid session")
			return
		}
		cfg.handleAuthError(w, r, http.StatusInternalServerError, "database error during session check", err, EventDbError, true, false, "user_id", userID)
		return
	}

	cfg.LogInfo(r, "session check successful", user.ID, EventSessionCheckSuccess)
	res := LoginResponse{
		ID:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
		AvatarUrl:   user.AvatarUrl.String,
		Status:      user.Status.String,
		Description: user.Description.String,
	}
	RespondWithJSON(w, http.StatusOK, res)
}
