package main

import (
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/szidorpatrik/rift-chat/backend/rift-chat-api/internal/api"
)

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
	cfg.RegisterRoutes(mux)

	// SERVER SETUP
	loggedRouter := cfg.LoggingMiddleware(mux)
	identifiedRouter := cfg.IdentityMiddleware(loggedRouter)

	httpServer := &http.Server{
		Handler: identifiedRouter,
		Addr:    cfg.FullAddress,
	}

	// Start the server
	log.Printf("rift-chat api listening on: %v\n", httpServer.Addr)
	log.Fatal(httpServer.ListenAndServe())
}
