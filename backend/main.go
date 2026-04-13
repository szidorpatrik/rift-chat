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
	apiConfig := &api.ApiConfig{}

	if err := apiConfig.LoadEnvs(); err != nil {
		log.Fatalf("Failed to load environment variables: %v", err)
	}

	if err := apiConfig.RunMigrations(); err != nil {
        log.Fatalf("Migration failed: %v", err)
    }

	db, err := sql.Open("pgx", apiConfig.DbUrl)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return
	} else {
		log.Println("Successfully connected to the database")
	}
	defer db.Close()

	dbQueries := database.New(db)
	apiConfig.DbQueries = dbQueries

	mux := http.NewServeMux()

	// Healthcheck endpoint
	mux.HandleFunc("/healthz", apiConfig.HealthzHandler)

	// Admin endpoints

	// Auth endpoints

	// User endpoints

	httpServer := &http.Server{
		Handler: mux,
		Addr: apiConfig.FullAddress,
	}
	log.Printf("rift-chat api listening on :%v\n", httpServer.Addr)
	log.Fatal(httpServer.ListenAndServe())
}


