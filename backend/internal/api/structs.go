package api

import (
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
