package api

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/szidorpatrik/rift-chat/backend/rift-chat-api/internal/database"
	"github.com/szidorpatrik/rift-chat/backend/rift-chat-api/internal/logger"
	"github.com/szidorpatrik/rift-chat/backend/rift-chat-api/sql/schema"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func (cfg *ApiConfig) RegisterRoutes(mux *http.ServeMux) {
	cfg.registerSystemRoutes(mux)
	cfg.registerAuthRoutes(mux)
	cfg.registerServerRoutes(mux)
	cfg.registerChannelRoutes(mux)
	cfg.registerWebtransportRoutes(mux)
}

func (cfg *ApiConfig) LoadEnvs() error {
	platform := os.Getenv("BACKEND_PLATFORM")

	switch platform {
	case "", "dev", "development":
		envFile := ".env-dev"
		cfg.CertPath = "./certs/dev"
		cfg.KeyPath = "./certs/dev"
		log.Default().Printf("=== RUNNING IN DEVELOPMENT MODE, LOADING '%s' ===\n", envFile)
		godotenv.Overload(envFile)
	case "prod", "production":
		cfg.CertPath = "./certs"
		cfg.KeyPath = "./certs"
		log.Default().Println("=== RUNNING IN PRODUCTION MODE ===")
	default:
		return fmt.Errorf("=== UNKNOWN BACKEND_PLATFORM '%s': Check './backend/.env' if running locally or './.env' if running via docker compose ===\n", platform)
	}

	requiredEnvs := []string{
		"BACKEND_HTTP_PORT", "BACKEND_QUIC_PORT", "BACKEND_LISTEN_ADDR",
		"BACKEND_JWT_SECRET", "BACKEND_DATABASE_URL", "BACKEND_VALKEY_HOST",
		"BACKEND_VALKEY_PORT", "BACKEND_CERT_NAME", "BACKEND_KEY_NAME",
	}

	for _, envKey := range requiredEnvs {
		if os.Getenv(envKey) == "" {
			return fmt.Errorf("missing required environment variable: %s", envKey)
		}
	}

	httpPort := os.Getenv("BACKEND_HTTP_PORT")
	if _, err := strconv.ParseInt(httpPort, 10, 64); err != nil {
		return fmt.Errorf("invalid value for BACKEND_HTTP_PORT: %v", err)
	}

	quicPort := os.Getenv("BACKEND_QUIC_PORT")
	if _, err := strconv.ParseInt(quicPort, 10, 64); err != nil {
		return fmt.Errorf("invalid value for BACKEND_QUIC_PORT: %v", err)
	}

	listenAddr := os.Getenv("BACKEND_LISTEN_ADDR")
	jwtSecret := os.Getenv("BACKEND_JWT_SECRET")
	dbUrl := os.Getenv("BACKEND_DATABASE_URL")
	httpPortInt, _ := strconv.ParseInt(httpPort, 10, 64)
	quicPortInt, _ := strconv.ParseInt(quicPort, 10, 64)

	cfg.HttpPort = int32(httpPortInt)
	cfg.QuicPort = int32(quicPortInt)
	cfg.Platform = platform
	cfg.DbUrl = dbUrl
	cfg.FullAddress = fmt.Sprintf("%v:%v", listenAddr, httpPort)
	cfg.FullAddressQuic = fmt.Sprintf("%v:%v", listenAddr, quicPort)
	cfg.JwtSecret = jwtSecret
	cfg.CertPath = fmt.Sprintf("%s/%s", cfg.CertPath, os.Getenv("BACKEND_CERT_NAME"))
	cfg.KeyPath = fmt.Sprintf("%s/%s", cfg.KeyPath, os.Getenv("BACKEND_KEY_NAME"))
	return nil
}

func (cfg *ApiConfig) RunMigrations() error {
	log.Default().Println("=== RUNNING DATABASE MIGRATIONS ===")
	db, err := sql.Open("pgx", cfg.DbUrl)
	if err != nil {
		return fmt.Errorf("failed to open database for migrations: %w", err)
	}
	defer db.Close()

	goose.SetBaseFS(schema.MigrationsFS)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to set goose dialect: %w", err)
	}

	if err := goose.Up(db, "."); err != nil {
		return fmt.Errorf("failed to run goose up: %w", err)
	}
	return nil
}

func (cfg *ApiConfig) InitializeDatabase() error {
	db, err := sql.Open("pgx", cfg.DbUrl)
	if err != nil {
		return fmt.Errorf("Failed to connect to database: %v", err)
	}
	cfg.DB = db
	cfg.DbQueries = database.New(db)
	return nil
}

func (cfg *ApiConfig) InitializeLogger() {
	logBatcher := logger.NewLogBatcher(cfg.DbQueries, logBatcherSize, logBatcherFlushTimer)
	go logBatcher.StartWorker(context.Background())
	cfg.LogBatcher = logBatcher

	dbWriter := logger.NewBatchWriter(cfg.LogBatcher)

	// Send logs to the terminal and the DB
	multiWriter := io.MultiWriter(os.Stdout, dbWriter)

	cfg.Logger = slog.New(slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
		AddSource: false,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Force UTC timestamp
			if a.Key == slog.TimeKey && len(groups) == 0 {
				if t, ok := a.Value.Any().(time.Time); ok {
					return slog.Time(slog.TimeKey, t.UTC())
				}
			}
			// Strip port from IP addresses
			if a.Key == "ip" && a.Value.Kind() == slog.KindString {
				if host, _, err := net.SplitHostPort(a.Value.String()); err == nil {
					return slog.String("ip", host)
				}
			}
			return a
		},
	}))
}

func (cfg *ApiConfig) InitializeValkey() error {
	host := os.Getenv("BACKEND_VALKEY_HOST")
	port := os.Getenv("BACKEND_VALKEY_PORT")

	addr := fmt.Sprintf("%s:%s", host, port)

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	// Test the connection
	if err := client.Ping(context.Background()).Err(); err != nil {
		return fmt.Errorf("failed to connect to valkey: %w", err)
	}
	log.Default().Printf("connected to valkey: %v", client.NodeAddress())

	cfg.Valkey = client
	return nil
}

// Returns a string in the format "METHOD /api/VERSION/path"
func newEndpoint(httpMethod, path string) string {
	clearedPath := strings.TrimPrefix(path, "/")
	return fmt.Sprintf("%s %s/v%d/%s", httpMethod, apiRoot, apiVersion, clearedPath)
}
