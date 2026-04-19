-- name: BulkInsertJSONLogs :exec
INSERT INTO api_logs (created_at, log_data)
SELECT 
    unnest(sqlc.arg(times)::timestamptz[]), 
    CAST(unnest(sqlc.arg(log_records)::text[]) AS jsonb);
