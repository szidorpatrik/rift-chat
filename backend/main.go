package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/szidorpatrik/rift-chat/backend/rift-chat-api/internal/api"
)

const apiRoot = "/api"
const apiVersion = "v1"

func main() {
	cfg := &api.ApiConfig{}

	if err := cfg.LoadEnvs(); err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}

	if err := cfg.InitializeDatabase(); err != nil {
		log.Fatalf("Database initialization failed:\n\t %v", err)
	} else {
		defer cfg.DB.Close()
	}

	if err := cfg.RunMigrations(); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	cfg.InitializeLogger()

	mux := http.NewServeMux()

	// SYSTEM ENDPOINTS
	mux.HandleFunc("GET /healthz", cfg.HealthzHandler)

	// PUBLIC API ENDPOINTS
	mux.HandleFunc(newEndpoint("POST", "/users"), cfg.CreateUserHandler)
	mux.HandleFunc(newEndpoint("POST", "/login"), cfg.LoginHandler)
	mux.HandleFunc(newEndpoint("POST", "/refresh"), cfg.RefreshTokenHandler)
	mux.HandleFunc(newEndpoint("POST", "/revoke"), cfg.RevokeTokenHandler)

	// PROTECTED API ENDPOINTS

	// SERVER SETUP
	loggedRouter := cfg.LoggingMiddleware(mux)
	httpServer := &http.Server{
		Handler: loggedRouter,
		Addr:    cfg.FullAddress,
	}

	// Start the server
	log.Printf("rift-chat api listening on: %v\n", httpServer.Addr)
	log.Fatal(httpServer.ListenAndServe())
}

// Returns a string in the format "METHOD /api/VERSION/path"
func newEndpoint(httpMethod, path string) string {
	clearedPath := strings.TrimPrefix(path, "/")
	return fmt.Sprintf("%s %s/%s/%s", httpMethod, apiRoot, apiVersion, clearedPath)
}
