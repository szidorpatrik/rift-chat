package api

import (
	"context"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/szidorpatrik/rift-chat/backend/rift-chat-api/internal/auth"
	"github.com/szidorpatrik/rift-chat/backend/rift-chat-api/internal/database"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.status = code
	rec.ResponseWriter.WriteHeader(code)
}

type contextKey struct{}

var userContextKey = contextKey{}

var publicEndpoints = []string{
	"/favicon.ico",
	"/healthz",
}

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
		currentUser, err := cfg.DbQueries.GetUserById(r.Context(), userId)
		if err != nil {
			RespondWithError(w, http.StatusUnauthorized, "user not found")
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, currentUser)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

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

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}

		user, ok := r.Context().Value(userContextKey).(database.User)
		if !ok {
			user = database.User{Username: "Guest"}
		}

		log.Printf("IP: %s | User: %s | Method: %s | Path: %s | Status: %d",
			ip, user.Username, r.Method, r.URL.Path, recorder.status)
	})
}
