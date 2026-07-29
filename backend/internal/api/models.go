package api

import (
	"database/sql"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/szidorpatrik/rift-chat/backend/rift-chat-api/internal/database"
	"github.com/szidorpatrik/rift-chat/backend/rift-chat-api/internal/logger"
)

type ApiConfig struct {
	HttpPort        int32
	Platform        string
	DbUrl           string
	DbQueries       *database.Queries
	DB              *sql.DB
	FullAddress     string
	JwtSecret       string
	Logger          *slog.Logger
	LogBatcher      *logger.LogBatcher
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
	ID          string `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Status      string `json:"status"`
	AvatarUrl   string `json:"avatar_url,omitempty"`
	Description string `json:"description,omitempty"`
}

type ServerRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	IsPublic    bool   `json:"is_public,omitempty"`
}

type ServerResponse struct {
	ID            string    `json:"id,omitempty"`
	OwnerId       string    `json:"owner_id,omitempty"`
	Name          string    `json:"name"`
	IconUrl       string    `json:"icon_url,omitempty"`
	BackgroundUrl string    `json:"background_url,omitempty"`
	IsPublic      bool      `json:"is_public,omitempty"`
	Description   string    `json:"description,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
}

type ServerJoinRequest struct {
	ServerId uuid.UUID `json:"server_id"`
}

type ChannelCreateRequest struct {
	Name     string `json:"name"`
	Type     int32  `json:"type"`
	Position int32  `json:"position"`
}

type ChannelResponse struct {
	ID       uuid.UUID `json:"id"`
	ServerID uuid.UUID `json:"server_id"`
	Name     string    `json:"name"`
	Type     int32     `json:"type"`
	Position int32     `json:"position"`
}

type MembersResponse struct {
	ID          uuid.UUID `json:"id"`
	Username    string    `json:"username"`
	Status      string    `json:"status"`
	AvatarUrl   string    `json:"avatar_url,omitempty"`
	Description string    `json:"description,omitempty"`
}
