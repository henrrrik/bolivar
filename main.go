package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
)

func main() {
	mapboxToken := requireEnv("MAPBOX_TOKEN")
	webhookUser := requireEnv("WEBHOOK_USER")
	webhookPass := requireEnv("WEBHOOK_PASS")

	dbPath := envOrDefault("DB_PATH", "bolivar.db")
	listenAddr := envOrDefault("LISTEN_ADDR", ":8080")

	db, err := OpenDB(dbPath)
	if err != nil {
		slog.Error("failed to open database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	srv := NewServer(db, mapboxToken, webhookUser, webhookPass)

	slog.Info("starting server", "addr", listenAddr)
	if err := http.ListenAndServe(listenAddr, srv.Routes()); err != nil {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}
}

func requireEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		fmt.Fprintf(os.Stderr, "required environment variable %s is not set\n", key)
		os.Exit(1)
	}
	return v
}

func envOrDefault(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
