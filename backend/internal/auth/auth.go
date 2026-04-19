package auth

import (
	"crypto"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func HashPassword(password string) (string, error) {
	return argon2id.CreateHash(password, argon2id.DefaultParams)
}

func CheckPasswordHash(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}

// Extracts Bearer token from the Authorization header in the format "Bearer <key>"
func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("Authorization header missing")
	}

	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", fmt.Errorf("Malformed authorization header")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return "", fmt.Errorf("Bearer token missing")
	}

	return token, nil
}

// Extracts the API key from the Authorization header in the format "ApiKey <key>"
func GetAPIKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("Authorization header missing")
	}
	if !strings.HasPrefix(authHeader, "ApiKey ") {
		return "", fmt.Errorf("Malformed authorization header")
	}
	apiKey := strings.TrimPrefix(authHeader, "ApiKey ")
	if apiKey == "" {
		return "", fmt.Errorf("API key missing")
	}
	return apiKey, nil
}

// Returns a new Refresh token and its hash
func MakeRefreshToken() (token string, hash string) {
	slice := make([]byte, 32)
	rand.Read(slice)
	return hex.EncodeToString(slice), HashRefreshToken(slice)
}

// Hashes the refresh token using SHA256
func HashRefreshToken(token []byte) string {
	hasher := crypto.SHA256.New()
	hasher.Write([]byte(token))
	return hex.EncodeToString(hasher.Sum(nil))
}

// Creates a JWT access token with the user ID as the subject and the specified expiration time
func MakeJWT(userId uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "rift-chat-access",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn).UTC()),
		Subject:   userId.String(),
	})
	return token.SignedString([]byte(tokenSecret))
}

// Validates the JWT and returns the user ID if valid, or an error if invalid
func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.Nil, err
	}

	if !token.Valid {
		return uuid.Nil, fmt.Errorf("invalid token")
	}

	userId, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid token subject")
	}
	return userId, nil
}
