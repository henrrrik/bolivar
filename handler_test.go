package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testServer(t *testing.T) *Server {
	t.Helper()
	db := testDB(t)
	return &Server{
		db:          db,
		mapboxToken: "test-token",
		webhookUser: "testuser",
		webhookPass: "testpass",
		mapTemplate: nil, // not tested here
	}
}

func TestWebhookAuth(t *testing.T) {
	srv := testServer(t)

	// No auth
	req := httptest.NewRequest("POST", "/webhook", strings.NewReader("message=hello"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	w := httptest.NewRecorder()
	srv.handleWebhook(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("no auth: got %d, want 401", w.Code)
	}

	// Wrong auth
	req = httptest.NewRequest("POST", "/webhook", strings.NewReader("message=hello"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth("wrong", "creds")
	w = httptest.NewRecorder()
	srv.handleWebhook(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("wrong auth: got %d, want 401", w.Code)
	}

	// Correct auth, no Garmin URL — should return 200
	req = httptest.NewRequest("POST", "/webhook", strings.NewReader("message=just+a+text"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth("testuser", "testpass")
	w = httptest.NewRecorder()
	srv.handleWebhook(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("valid auth no garmin: got %d, want 200", w.Code)
	}
}

func TestMessagesEndpoint(t *testing.T) {
	srv := testServer(t)

	req := httptest.NewRequest("GET", "/api/messages", nil)
	w := httptest.NewRecorder()
	srv.handleMessages(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("got %d, want 200", w.Code)
	}

	var msgs []Message
	if err := json.NewDecoder(w.Body).Decode(&msgs); err != nil {
		t.Fatal(err)
	}
	if len(msgs) != 0 {
		t.Errorf("expected empty array, got %d messages", len(msgs))
	}
}

func TestHealthEndpoint(t *testing.T) {
	srv := testServer(t)

	req := httptest.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	srv.handleHealth(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("got %d, want 200", w.Code)
	}

	var resp map[string]string
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}
	if resp["status"] != "ok" {
		t.Errorf("status = %q, want ok", resp["status"])
	}
}
