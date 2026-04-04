package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	mapboxToken := requireEnv("MAPBOX_TOKEN")
	webhookUser := requireEnv("WEBHOOK_USER")
	webhookPass := requireEnv("WEBHOOK_PASS")

	databaseURL := requireEnv("DATABASE_URL")
	port := envOrDefault("PORT", "8080")
	listenAddr := ":" + port

	db, err := OpenDB(databaseURL)
	if err != nil {
		slog.Error("failed to open database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	telegram := NewTelegramNotifier(
		os.Getenv("TELEGRAM_BOT_TOKEN"),
		os.Getenv("TELEGRAM_CHAT_ID"),
	)
	if telegram != nil {
		slog.Info("Telegram notifications enabled")
	}

	srv := NewServer(db, mapboxToken, webhookUser, webhookPass, telegram)

	httpSrv := &http.Server{
		Addr:              listenAddr,
		Handler:           srv.Routes(),
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
	}

	slog.Info("starting server", "addr", listenAddr)
	if err := httpSrv.ListenAndServe(); err != nil {
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
