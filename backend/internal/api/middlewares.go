package api

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/szidorpatrik/rift-chat/backend/rift-chat-api/internal/auth"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

const userContextKey = "userID"
var publicEndpoints = []string{
	"/favicon.ico",
	"/healthz",
}

// AuthMiddleware checks for a valid JWT access token in the Authorization header and extracts the user ID
func (cfg *ApiConfig) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			RespondWithError(w, http.StatusForbidden, err.Error())
			return
		}
		userId, err := auth.ValidateJWT(token, cfg.JwtSecret)
		if err != nil {
			RespondWithError(w, http.StatusForbidden, err.Error())
			return
		}
		ctx := context.WithValue(r.Context(), userContextKey, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// LoggingMiddleware logs incoming requests with IP, user ID (if authenticated), method, path, and response status
func (cfg *ApiConfig) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, path := range publicEndpoints {
			if strings.HasSuffix(r.URL.Path, path) {
				next.ServeHTTP(w, r)
				return
			}
		}

		recorder := &statusRecorder{
			ResponseWriter: w,
			status:         http.StatusOK,
		}
		next.ServeHTTP(recorder, r)

		userId, ok := r.Context().Value(userContextKey).(uuid.UUID)
		if !ok {
			userId = uuid.Nil
		}

		cfg.Logger.Info("request", "ip", r.RemoteAddr, "user", userId, "method", r.Method, "path", r.URL.Path, "status", recorder.status)
	})
}
