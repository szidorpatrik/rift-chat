-- +goose Up
CREATE TABLE api_logs (
    id BIGSERIAL PRIMARY KEY,
    created_at TIMESTAMPTZ NOT NULL,
    log_data JSONB NOT NULL
);

-- GIN index allows blazingly fast queries INSIDE the JSON object
CREATE INDEX idx_api_logs_data ON api_logs USING GIN (log_data);

-- +goose Down
DROP TABLE api_logs;

DROP INDEX idx_api_logs_data;
