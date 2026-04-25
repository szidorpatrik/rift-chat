package api

import (
	"database/sql"
	"log/slog"
	"time"

	"github.com/szidorpatrik/rift-chat/backend/rift-chat-api/internal/database"
	"github.com/szidorpatrik/rift-chat/backend/rift-chat-api/internal/logger"
)

type ApiConfig struct {
	HttpPort    int32
	QuicPort    int32
	Platform    string
	DbUrl       string
	DbQueries   *database.Queries
	DB          *sql.DB
	FullAddress string
	JwtSecret   string
	Logger      *slog.Logger
	LogBatcher  *logger.LogBatcher
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

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email,omitempty"`
}
