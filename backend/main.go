package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/szidorpatrik/rift-chat/backend/rift-chat-api/internal/api"
	"github.com/szidorpatrik/rift-chat/backend/rift-chat-api/internal/database"
)

const apiRoot = "/api"
const apiVersion = "v1"

func main() {
	cfg := &api.ApiConfig{}

	if err := cfg.LoadEnvs(); err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}

	if err := cfg.RunMigrations(); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	db, err := sql.Open("pgx", cfg.DbUrl)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return
	} else {
		log.Println("Successfully connected to the database")
	}
	defer db.Close()

	dbQueries := database.New(db)
	cfg.DbQueries = dbQueries

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

func newEndpoint(httpMethod, path string) string {
	clearedPath := strings.TrimPrefix(path, "/")
	return fmt.Sprintf("%s %s/%s/%s", httpMethod, apiRoot, apiVersion, clearedPath)
}
