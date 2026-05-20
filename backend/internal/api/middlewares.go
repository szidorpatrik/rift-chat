package api

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/szidorpatrik/rift-chat/backend/rift-chat-api/internal/auth"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

var publicEndpoints = []string{
	"/favicon.ico",
	"/healthz",
	// "/wt",
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

func (cfg *ApiConfig) ServerMembershipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Safe type assertion to prevent panics
		userId, ok := r.Context().Value(userContextKey).(uuid.UUID)
		if !ok || userId == uuid.Nil {
			RespondWithError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		// 2. Get server_id from path parameters
		serverIdStr := r.PathValue(serverContextKey)
		if serverIdStr == "" {
			cfg.LogWarn(r, fmt.Sprintf("missing %s path parameter", serverContextKey), userId, EventMissingParameter)
			RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("%s path parameter is required", serverContextKey))
			return
		}

		serverId, err := uuid.Parse(serverIdStr)
		if err != nil {
			cfg.LogWarn(r, fmt.Sprintf("invalid %s format", serverContextKey), userId, EventInvalidParameter)
			RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("invalid %s format", serverContextKey))
			return
		}

		// 3. Verify membership
		servers, err := cfg.DbQueries.GetServersByUserId(r.Context(), userId)
		if err != nil {
			cfg.LogError(r, err, EventDbError, serverContextKey, serverId)
			RespondWithError(w, http.StatusInternalServerError, "failed to verify server membership")
			return
		}

		isMember := false
		for _, server := range servers {
			if server.ID == serverId {
				isMember = true
				break
			}
		}

		if !isMember {
			cfg.LogWarn(r, "user is not a member of the server", userId, EventServerNotMember, serverContextKey, serverId)
			RespondWithError(w, http.StatusForbidden, "you are not a member of this server")
			return
		}

		// 4. Attach server ID to context and proceed
		ctx := context.WithValue(r.Context(), serverContextKey, serverId)
		r = r.WithContext(ctx)
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
			"event", EventHTTPRequest,
			"ip", GetClientIP(r),
			"user_id", userId,
			"method", r.Method,
			"path", r.URL.Path,
			"status", recorder.status,
		)
	})
}

func (r *statusRecorder) Unwrap() http.ResponseWriter {
	return r.ResponseWriter
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}
