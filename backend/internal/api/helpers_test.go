package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/szidorpatrik/rift-chat/backend/rift-chat-api/internal/database"
)

func TestNewEndpoint(t *testing.T) {
	got := newEndpoint("GET", "/servers/{server_id}/channels")
	want := "GET /api/v1/servers/{server_id}/channels"
	if got != want {
		t.Fatalf("newEndpoint() = %q, want %q", got, want)
	}
}

func TestSanitizeInput(t *testing.T) {
	got := sanitizeInput("\u200b  hello\t\nworld  ", "-")
	want := "hello-world"
	if got != want {
		t.Fatalf("sanitizeInput() = %q, want %q", got, want)
	}
}

func TestGetClientIP(t *testing.T) {
	t.Run("xff wins", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set("X-Forwarded-For", "1.2.3.4, 5.6.7.8")
		if got := GetClientIP(r); got != "1.2.3.4" {
			t.Fatalf("GetClientIP() = %q, want %q", got, "1.2.3.4")
		}
	})

	t.Run("real ip fallback", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.Header.Set("X-Real-IP", "9.9.9.9")
		if got := GetClientIP(r); got != "9.9.9.9" {
			t.Fatalf("GetClientIP() = %q, want %q", got, "9.9.9.9")
		}
	})

	t.Run("remote addr fallback", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.RemoteAddr = "127.0.0.1:12345"
		if got := GetClientIP(r); got != "127.0.0.1:12345" {
			t.Fatalf("GetClientIP() = %q, want %q", got, "127.0.0.1:12345")
		}
	})
}

func TestRespondWithError(t *testing.T) {
	recorder := httptest.NewRecorder()
	RespondWithError(recorder, http.StatusUnauthorized, "nope")

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusUnauthorized)
	}
	if got := recorder.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("content-type = %q, want %q", got, "application/json")
	}

	var res RespondError
	if err := json.Unmarshal(recorder.Body.Bytes(), &res); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if res.Error != "nope" || res.StatusCode != http.StatusUnauthorized {
		t.Fatalf("response = %+v, want error=nope status=%d", res, http.StatusUnauthorized)
	}
}

func TestRespondWithJSON(t *testing.T) {
	recorder := httptest.NewRecorder()
	RespondWithJSON(recorder, http.StatusCreated, map[string]string{"ok": "yes"})

	if recorder.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusCreated)
	}
	if got := recorder.Header().Get("Content-Type"); got != "application/json" {
		t.Fatalf("content-type = %q, want %q", got, "application/json")
	}
	if !bytes.Contains(recorder.Body.Bytes(), []byte(`"ok":"yes"`)) {
		t.Fatalf("body = %q, want JSON payload", recorder.Body.String())
	}
}

func TestDecodeAndValidate(t *testing.T) {
	type payload struct {
		Username string `json:"username"`
	}

	t.Run("valid payload", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"username":"bob"}`))
		decoded, code, err := decodeAndValidate(req, func(in payload) (int, error) {
			if in.Username == "" {
				return http.StatusBadRequest, errors.New("missing username")
			}
			return http.StatusOK, nil
		})

		if err != nil || code != http.StatusOK || decoded.Username != "bob" {
			t.Fatalf("decodeAndValidate() = %+v, %d, %v", decoded, code, err)
		}
	})

	t.Run("malformed json", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString(`{"username":`))
		_, code, err := decodeAndValidate(req, func(in payload) (int, error) {
			return http.StatusOK, nil
		})

		if err == nil || code != http.StatusBadRequest {
			t.Fatalf("decodeAndValidate() = code %d err %v, want bad request error", code, err)
		}
	})
}

func TestValidateCreateUserRequest(t *testing.T) {
	if code, err := validateCreateUserRequest(CreateUserRequest{Username: "bob", Email: "bob@example.com", Password: "xY9!qW7#eR5$uT3&vI1@pL8^nM6?"}); err != nil || code != http.StatusOK {
		t.Fatalf("valid request = code %d err %v", code, err)
	}

	if code, err := validateCreateUserRequest(CreateUserRequest{Username: "", Email: "", Password: ""}); code != http.StatusBadRequest || err == nil {
		t.Fatalf("empty request = code %d err %v, want bad request", code, err)
	}
}

func TestValidateLoginRequest(t *testing.T) {
	cfg := &ApiConfig{}
	if code, err := cfg.validateLoginRequest(LoginRequest{Username: "bob", Password: "secret"}); err != nil || code != http.StatusOK {
		t.Fatalf("valid request = code %d err %v", code, err)
	}

	if code, err := cfg.validateLoginRequest(LoginRequest{Username: "b", Password: "secret"}); code != http.StatusBadRequest || err == nil {
		t.Fatalf("short username = code %d err %v, want bad request", code, err)
	}
}

func TestCheckRefreshTokenValid(t *testing.T) {
	token := database.RefreshToken{
		RevokedAt: sql.NullTime{Valid: false},
		ExpiresAt: time.Now().Add(time.Hour),
	}
	if !checkRefreshTokenValid(token) {
		t.Fatal("expected token to be valid")
	}

	token.RevokedAt = sql.NullTime{Valid: true}
	if checkRefreshTokenValid(token) {
		t.Fatal("expected revoked token to be invalid")
	}
}

func TestServerJoinRequestRoundTrip(t *testing.T) {
	serverID := uuid.New()
	req := ServerJoinRequest{ServerId: serverID}
	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if !bytes.Contains(data, []byte(serverID.String())) {
		t.Fatalf("marshal = %s, want server id", data)
	}
}
