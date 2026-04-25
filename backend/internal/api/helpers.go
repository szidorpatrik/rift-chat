package api

import (
	"encoding/json"
	"net/http"

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
	w.WriteHeader(code)
	w.Header().Add("Content-Type", "application/json")
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
		"ip", r.RemoteAddr,
		"error", err.Error(),
	}, extraArgs...)
	cfg.Logger.Error(err.Error(), args...)
}

func (cfg *ApiConfig) LogInfo(r *http.Request, msg string, userId uuid.UUID, event EventType, extraArgs ...any) {
	args := append([]any{
		"event", event.String(),
		"user_id", userId.String(),
		"ip", r.RemoteAddr,
	}, extraArgs...)
	cfg.Logger.Info(msg, args...)
}

func (cfg *ApiConfig) LogWarn(r *http.Request, msg string, userId uuid.UUID, event EventType, extraArgs ...any) {
	args := append([]any{
		"event", event.String(),
		"user_id", userId.String(),
		"ip", r.RemoteAddr,
	}, extraArgs...)
	cfg.Logger.Warn(msg, args...)
}
