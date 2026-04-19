package logger

import (
	"encoding/json"
	"time"

	"github.com/szidorpatrik/rift-chat/backend/rift-chat-api/internal/database"
)

// LogEntry holds the raw JSON bytes and timestamp for database insertion
type LogEntry struct {
	CreatedAt time.Time
	Data      json.RawMessage
}

// LogBatcher manages the async queue and database flushing
type LogBatcher struct {
	logChannel chan string
	db         *database.Queries
	batchSize  int
	timeout    time.Duration
}
