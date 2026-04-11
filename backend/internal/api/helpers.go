package api

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func (cfg *ApiConfig) LoadEnvs() error {
	// TODO: Remove in production
	godotenv.Load()
	
	requiredEnvs := []string{
		"HTTP_PORT", "QUIC_PORT", "PLATFORM",
		"LISTEN_ADDR", "JWT_SECRET", "VALKEY_URL",
		"DB_URL",
	}
	
	for _, envKey := range requiredEnvs {
		if os.Getenv(envKey) == "" {
			return fmt.Errorf("missing required environment variable: %s", envKey)
		}
	}
	
	httpPort := os.Getenv("HTTP_PORT")
	if _, err := strconv.ParseInt(httpPort, 10, 64); err != nil {
		return fmt.Errorf("invalid value for HTTP_PORT: %v", err)
	}
	
	quicPort := os.Getenv("QUIC_PORT")
	if _, err := strconv.ParseInt(quicPort, 10, 64); err != nil {
		return fmt.Errorf("invalid value for QUIC_PORT: %v", err)
	}
	
	platform := os.Getenv("PLATFORM")
	listenAddr := os.Getenv("LISTEN_ADDR")
	jwtSecret := os.Getenv("JWT_SECRET")
	dbUrl := os.Getenv("DB_URL")
	valkeyUrl := os.Getenv("VALKEY_URL")

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
