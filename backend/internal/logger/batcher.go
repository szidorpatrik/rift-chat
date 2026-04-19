package logger

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/szidorpatrik/rift-chat/backend/rift-chat-api/internal/database"
)

func NewLogBatcher(db *database.Queries, batchSize int, timeout time.Duration) *LogBatcher {
	return &LogBatcher{
		logChannel: make(chan string, batchSize*2),
		db:         db,
		batchSize:  batchSize,
		timeout:    timeout,
	}
}

// StartWorker should be run as a goroutine when the app starts
func (b *LogBatcher) StartWorker(ctx context.Context) {
	var batch []string
	ticker := time.NewTicker(b.timeout)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			if len(batch) > 0 {
				b.flushToDatabase(context.Background(), batch)
			}
			return
		case entry := <-b.logChannel:
			batch = append(batch, entry)
			if len(batch) >= b.batchSize {
				b.flushToDatabase(ctx, batch)
				batch = nil
				ticker.Reset(b.timeout)
			}
		case <-ticker.C:
			if len(batch) > 0 {
				b.flushToDatabase(ctx, batch)
				batch = nil
			}
		}
	}
}

func (b *LogBatcher) flushToDatabase(ctx context.Context, batch []string) {
	err := b.db.BulkInsertJSONLogs(ctx, batch)
	if err != nil {
		log.Println("ERROR: Failed to bulk insert JSON logs:", err)
	}
}

type BatchWriter struct {
	batcher *LogBatcher
}

func NewBatchWriter(batcher *LogBatcher) *BatchWriter {
	return &BatchWriter{
		batcher: batcher,
	}
}

func (bw *BatchWriter) Write(p []byte) (n int, err error) {
	// Casting []byte to string creates a safe copy
	select {
	case bw.batcher.logChannel <- string(p):
		// Success!
	default:
		os.Stdout.Write([]byte("WARNING: DB Log Queue Full! Dropping log to stdout: "))
		os.Stdout.Write(p)
	}

	return len(p), nil
}
