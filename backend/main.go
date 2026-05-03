package main

import (
	"crypto/tls"
	"log"
	"net/http"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/quic-go/quic-go/http3"
	"github.com/quic-go/webtransport-go"
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

	if err := cfg.InitializeValkey(); err != nil {
		log.Fatalf("Valkey initialization failed:\n\t %v", err)
	}

	mux := http.NewServeMux()
	cfg.RegisterRoutes(mux)

	// SERVER SETUP
	loggedRouter := cfg.LoggingMiddleware(mux)
	identifiedRouter := cfg.IdentityMiddleware(loggedRouter)

	wtServer := &webtransport.Server{
		H3: &http3.Server{
			Addr:            cfg.FullAddressQuic,
			Handler:         identifiedRouter,
			EnableDatagrams: true,
		},
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	cfg.WTServer = wtServer

	go func() {
		certFile := "/root/certs/cert.pem"
		keyFile := "/root/certs/key.pem"
		if cfg.Platform == "" || cfg.Platform == "dev" || cfg.Platform == "development" {
			certFile = "./certs/dev/cert.pem"
			keyFile = "./certs/dev/key.pem"
		}
		
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			log.Fatalf("FATAL: Failed to load WebTransport certificates: %v", err)
		}
		
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			NextProtos:   []string{"h3"},
		}
		
		wtServer.H3.TLSConfig = tlsConfig
		
		log.Printf("rift-chat WebTransport listening on: %v (UDP)\n", cfg.FullAddressQuic)
		log.Fatal(wtServer.ListenAndServe())
	}()

	httpServer := &http.Server{
		Handler: identifiedRouter,
		Addr:    cfg.FullAddress,
	}

	// Start the server
	log.Printf("rift-chat api listening on: %v\n", httpServer.Addr)
	log.Fatal(httpServer.ListenAndServe())
}
