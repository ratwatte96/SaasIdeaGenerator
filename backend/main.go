package main

import (
	"log"
	"net/http"

	"saasideagenerator/backend/internal/config"
	"saasideagenerator/backend/internal/db"
	"saasideagenerator/backend/internal/server"
)

func main() {
	cfg := config.Load()
	database, err := db.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("db connection failed: %v", err)
	}
	defer database.Close()

	srv := server.New(database)
	addr := ":" + cfg.APIPort
	log.Printf("api listening on %s", addr)
	log.Fatal(http.ListenAndServe(addr, srv.Routes()))
}
