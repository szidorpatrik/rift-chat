package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// MessageMetadata holds the flexible JSONB data
type MessageMetadata struct {
	Reactions   map[string]int `json:"reactions,omitempty"`
	Attachments []string       `json:"attachments,omitempty"`
}

// Value makes MessageMetadata implement the driver.Valuer interface
func (m MessageMetadata) Value() (driver.Value, error) {
	return json.Marshal(m)
}

// Scan makes MessageMetadata implement the sql.Scanner interface
func (m *MessageMetadata) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &m)
}
