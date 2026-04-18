package api

import (
	"time"

	"github.com/szidorpatrik/rift-chat/backend/rift-chat-api/internal/database"
)

type ApiConfig struct {
	HttpPort    int32
	QuicPort    int32
	Platform    string
	DbUrl       string
	DbQueries   *database.Queries
	FullAddress string
	ValkeyUrl   string
	JwtSecret   string
}

type RespondError struct {
	Error      string `json:"error"`
	StatusCode int    `json:"status_code,omitempty"`
}

type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CreateUserResponse struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Bearer token is required in Authorization header
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email,omitempty"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token,omitempty"`
}
