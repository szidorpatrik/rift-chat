-- name: BulkInsertJSONLogs :exec
INSERT INTO api_logs (created_at, log_data)
SELECT 
    (COALESCE(parsed.log_json->>'time', NOW()::text))::timestamptz, 
    parsed.log_json
FROM (
    SELECT CAST(unnest(sqlc.arg(log_records)::text[]) AS jsonb) AS log_json
) parsed;
