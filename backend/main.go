package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/szidorpatrik/rift-chat/backend/rift-chat-api/internal/api"
	"github.com/szidorpatrik/rift-chat/backend/rift-chat-api/internal/database"
	_ "github.com/jackc/pgx/v5/stdlib"
)

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

	mux := http.NewServeMux()

	// Healthcheck endpoint
	mux.HandleFunc("/healthz", cfg.HealthzHandler)

	// Admin endpoints

	// Auth endpoints
	usersHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("User created successfully!\n"))
	})
	//mux.Handle("POST /users", cfg.AuthMiddleware(http.HandlerFunc(cfg.CreateUserHandler)))
	mux.Handle("POST /users", usersHandler)
	// User endpoints

	// Middlewares
	loggedRouter := cfg.LoggingMiddleware(mux)
	httpServer := &http.Server{
		Handler: loggedRouter, 
		Addr:    cfg.FullAddress,
	}

	// Start the server
	log.Printf("rift-chat api listening on :%v\n", httpServer.Addr)
	log.Fatal(httpServer.ListenAndServe())
}
