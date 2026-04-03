package main

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"html/template"
	"log/slog"
	"net/http"
	"time"
)

//go:embed static/index.html
var indexHTML string

type Server struct {
	db           *sql.DB
	mapboxToken  string
	webhookUser  string
	webhookPass  string
	mapTemplate  *template.Template
	telegram     *TelegramNotifier
}

func NewServer(db *sql.DB, mapboxToken, webhookUser, webhookPass string, telegram *TelegramNotifier) *Server {
	tmpl := template.Must(template.New("index").Parse(indexHTML))
	return &Server{
		db:          db,
		mapboxToken: mapboxToken,
		webhookUser: webhookUser,
		webhookPass: webhookPass,
		mapTemplate: tmpl,
		telegram:    telegram,
	}
}

func (s *Server) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /webhook", s.handleWebhook)
	mux.HandleFunc("GET /api/messages", s.handleMessages)
	mux.HandleFunc("GET /health", s.handleHealth)
	mux.HandleFunc("GET /", s.handleMap)
	return mux
}

func (s *Server) handleWebhook(w http.ResponseWriter, r *http.Request) {
	user, pass, ok := r.BasicAuth()
	if !ok || user != s.webhookUser || pass != s.webhookPass {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err := r.ParseForm(); err != nil {
		slog.Error("failed to parse form", "error", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	msgText := r.FormValue("message")
	smsFrom := r.FormValue("from")

	garminLink, extID, found := ExtractGarminURL(msgText)
	if !found {
		slog.Warn("no Garmin URL found in message", "from", smsFrom, "message", msgText)
		w.WriteHeader(http.StatusOK)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	msg := Message{
		GarminID:  extID,
		Text:      StripGarminURL(msgText),
		URL:       garminLink,
		CreatedAt: time.Now().UTC(),
	}

	result, err := FetchGarminLocation(ctx, garminLink)
	if err != nil {
		slog.Error("failed to fetch Garmin location", "url", garminLink, "error", err)
	} else {
		msg.Sender = result.Sender
		msg.Lat = &result.Lat
		msg.Lon = &result.Lon
	}

	if err := InsertMessage(s.db, msg); err != nil {
		slog.Error("failed to insert message", "error", err)
	}

	if err := s.telegram.Send(ctx, msg); err != nil {
		slog.Error("failed to send Telegram notification", "error", err)
	}

	w.WriteHeader(http.StatusOK)
}

func (s *Server) handleMessages(w http.ResponseWriter, r *http.Request) {
	msgs, err := ListMessages(s.db)
	if err != nil {
		slog.Error("failed to list messages", "error", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}
	if msgs == nil {
		msgs = []Message{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(msgs)
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	err := s.db.Ping()
	status := "ok"
	if err != nil {
		status = "error"
		w.WriteHeader(http.StatusServiceUnavailable)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": status})
}

func (s *Server) handleMap(w http.ResponseWriter, r *http.Request) {
	s.mapTemplate.Execute(w, map[string]string{
		"MapboxToken": s.mapboxToken,
	})
}
