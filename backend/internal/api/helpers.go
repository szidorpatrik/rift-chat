package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"unicode"

	"github.com/google/uuid"
)

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	res := RespondError{
		Error:      msg,
		StatusCode: code,
	}
	payload, err := json.Marshal(res)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(payload)
}

func RespondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "something went wrong")
		return
	}
	w.Write(payloadBytes)
}

func RespondNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

func (cfg *ApiConfig) LogError(r *http.Request, err error, event EventType, extraArgs ...any) {
	args := append([]any{
		"event", event.String(),
		"ip", GetClientIP(r),
		"error", err.Error(),
	}, extraArgs...)
	cfg.Logger.Error(err.Error(), args...)
}

func (cfg *ApiConfig) LogInfo(r *http.Request, msg string, userId uuid.UUID, event EventType, extraArgs ...any) {
	args := append([]any{
		"event", event.String(),
		"user_id", userId.String(),
		"ip", GetClientIP(r),
	}, extraArgs...)
	cfg.Logger.Info(msg, args...)
}

func (cfg *ApiConfig) LogWarn(r *http.Request, msg string, userId uuid.UUID, event EventType, extraArgs ...any) {
	args := append([]any{
		"event", event.String(),
		"user_id", userId.String(),
		"ip", GetClientIP(r),
	}, extraArgs...)
	cfg.Logger.Warn(msg, args...)
}

func GetClientIP(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return strings.TrimSpace(xri)
	}

	return r.RemoteAddr
}

func sanitizeInput(input, replaceChar string) string {
	// Normalize weird spaces to regular space first
	input = strings.Map(func(r rune) rune {
		// Remove invisible formatting chars
		if unicode.Is(unicode.Cf, r) {
			return ' ' // replace with space
		}
		// Normalize all space separators to plain space
		if unicode.Is(unicode.Zs, r) {
			return ' '
		}
		return r
	}, input)

	input = strings.TrimSpace(input)
	input = whitespace.ReplaceAllString(input, replaceChar)

	return input
}
