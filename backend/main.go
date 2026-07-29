package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/szidorpatrik/rift-chat/backend/rift-chat-api/internal/api"
)

func main() {
	handleHealthcheck()

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

func handleHealthcheck() {
	if len(os.Args) <= 1 || os.Args[1] != "healthcheck" {
		return
	}

	port := os.Getenv("BACKEND_HTTP_PORT")
	if port == "" {
		port = "8081"
	}

	resp, err := http.Get(fmt.Sprintf("http://127.0.0.1:%s/api/v1/healthz", port))
	if err != nil || resp.StatusCode != http.StatusOK {
		os.Exit(1)
	}

	os.Exit(0)
}
