package api

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/szidorpatrik/rift-chat/backend/rift-chat-api/sql/schema"
	//"github.com/szidorpatrik/rift-chat/backend/rift-chat-api/internal/models"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func (cfg *ApiConfig) LoadEnvs() error {
	platform := os.Getenv("BACKEND_PLATFORM")

	switch platform {
	case "", "dev", "development":
		envFile := ".env-dev"
		fmt.Printf("=== RUNNING IN DEVELOPMENT MODE, LOADING '%s' ===\n", envFile)
		godotenv.Overload(envFile)
	case "prod", "production":
		fmt.Printf("=== RUNNING IN PRODUCTION MODE ===\n")
	default:
		return fmt.Errorf("=== UNKNOWN BACKEND_PLATFORM '%s': Check './backend/.env' if running locally or './.env' if running via docker compose ===\n", platform)
	}

	requiredEnvs := []string{
		"BACKEND_HTTP_PORT", "BACKEND_QUIC_PORT","BACKEND_LISTEN_ADDR", 
		"BACKEND_JWT_SECRET", "BACKEND_VALKEY_URL", "BACKEND_DATABASE_URL",
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
	valkeyUrl := os.Getenv("BACKEND_VALKEY_URL")

	httpPortInt, _ := strconv.ParseInt(httpPort, 10, 64)
	quicPortInt, _ := strconv.ParseInt(quicPort, 10, 64)

	cfg.HttpPort = int32(httpPortInt)
	cfg.QuicPort = int32(quicPortInt)
	cfg.Platform = platform
	cfg.DbUrl = dbUrl
	cfg.FullAddress = fmt.Sprintf("%v:%v", listenAddr, httpPort)
	cfg.ValkeyUrl = valkeyUrl
	cfg.JwtSecret = jwtSecret

	return nil
}

func (cfg *ApiConfig) RunMigrations() error {
	fmt.Printf("=== RUNNING DATABASE MIGRATIONS ===\n")

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
