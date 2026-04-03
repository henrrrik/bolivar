package main

import (
	"context"
	"testing"
	"time"
)

func TestNilNotifierSend(t *testing.T) {
	var n *TelegramNotifier
	err := n.Send(context.Background(), Message{
		Sender:    "Test",
		Text:      "Hello",
		Lat:       pf(59.0),
		Lon:       pf(18.0),
		CreatedAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("nil notifier should no-op, got: %v", err)
	}
}

func TestNewTelegramNotifierEmpty(t *testing.T) {
	if n := NewTelegramNotifier("", "123"); n != nil {
		t.Error("expected nil with empty token")
	}
	if n := NewTelegramNotifier("token", ""); n != nil {
		t.Error("expected nil with empty chatID")
	}
	if n := NewTelegramNotifier("token", "123"); n == nil {
		t.Error("expected non-nil with both set")
	}
}
