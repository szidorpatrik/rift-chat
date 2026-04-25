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

// Blocks any request that does not have a valid user ID in its context.
func (cfg *ApiConfig) RequireAuth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        userId, ok := r.Context().Value(userContextKey).(uuid.UUID)
        
        // If the ID is missing or is the Nil UUID, kicks them out
        if !ok || userId == uuid.Nil {
            RespondWithError(w, http.StatusUnauthorized, "unauthorized")
            return
        }

        next.ServeHTTP(w, r)
    })
}

// Attempts to extract the user ID from a JWT cookie
// and attach it to the request context. It does NOT block unauthenticated requests.
func (cfg *ApiConfig) IdentityMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        cookie, err := r.Cookie("access_token")
        if err == nil {
            // If the cookie exists, try to validate it
            if userId, err := auth.ValidateJWT(cookie.Value, cfg.JwtSecret); err == nil {
                // Attach to context and replace the request pointer locally
                ctx := context.WithValue(r.Context(), userContextKey, userId)
                r = r.WithContext(ctx)
            }
        }
        next.ServeHTTP(w, r)
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

		cfg.Logger.Info("http request",
			"event", EventHTTPRequest.String(),
			"ip", r.RemoteAddr,
			"user_id", userId,
			"method", r.Method,
			"path", r.URL.Path,
			"status", recorder.status,
		)
	})
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}
