package api

import (
	"encoding/json"
	"net/http"
)

func RespondWithError(w http.ResponseWriter, code int, msg string) {
	res := RespondError{
		Error: msg,
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
