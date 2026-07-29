package auth

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestGetAPIKey(t *testing.T) {
	headers := http.Header{}
	headers.Set("Authorization", "ApiKey secret-123")

	got, err := GetAPIKey(headers)
	if err != nil {
		t.Fatalf("GetAPIKey() error = %v", err)
	}
	if got != "secret-123" {
		t.Fatalf("GetAPIKey() = %q, want %q", got, "secret-123")
	}

	headers.Set("Authorization", "Bearer secret-123")
	if _, err := GetAPIKey(headers); err == nil {
		t.Fatal("GetAPIKey() expected error for malformed header")
	}
}

func TestHashRefreshToken(t *testing.T) {
	hash1 := HashRefreshToken([]byte("token"))
	hash2 := HashRefreshToken([]byte("token"))
	if hash1 != hash2 {
		t.Fatalf("HashRefreshToken() not stable: %q != %q", hash1, hash2)
	}
	if len(hash1) != 64 {
		t.Fatalf("HashRefreshToken() length = %d, want 64", len(hash1))
	}
}

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"

	token, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT() error = %v", err)
	}

	got, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("ValidateJWT() error = %v", err)
	}
	if got != userID {
		t.Fatalf("ValidateJWT() = %s, want %s", got, userID)
	}
}

func TestValidateJWTRejectsBadSecret(t *testing.T) {
	token, err := MakeJWT(uuid.New(), "secret-a", time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT() error = %v", err)
	}

	_, err = ValidateJWT(token, "secret-b")
	if err == nil {
		t.Fatal("ValidateJWT() expected error with wrong secret")
	}
	if !strings.Contains(err.Error(), "signature") && !strings.Contains(err.Error(), "token") {
		t.Fatalf("ValidateJWT() error = %v, want signature/token error", err)
	}
}
