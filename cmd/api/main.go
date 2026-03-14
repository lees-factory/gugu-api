package main

import (
	"context"
	"log"
	"net/http"

	httpapi "github.com/ljj/gugu-api/internal/core/api"
	"github.com/ljj/gugu-api/internal/storage/dbcore"
	"github.com/ljj/gugu-api/internal/support/config"
)

func main() {
	cfg := config.Load()

	if cfg.DatabaseURL != "" {
		db, err := dbcore.OpenSupabase(context.Background(), cfg.DatabaseURL)
		if err != nil {
			log.Fatalf("connect database: %v", err)
		}
		defer db.Close()
	}

	server, err := httpapi.NewServer(cfg)
	if err != nil {
		log.Fatalf("build server: %v", err)
	}

	log.Printf("gugu api listening on %s", cfg.HTTPAddress)
	if err := http.ListenAndServe(cfg.HTTPAddress, server.Handler()); err != nil {
		log.Fatalf("serve http: %v", err)
	}
}
