package main

import (
	"database/sql"
	"testing"
	"time"
)

func testDB(t *testing.T) *sql.DB {
	t.Helper()
	db, err := OpenDB(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func pf(v float64) *float64 { return &v }

func TestInsertAndList(t *testing.T) {
	db := testDB(t)

	msg := Message{
		GarminID:  "test123",
		Sender:    "Henrik",
		Text:      "Hello from the trail!",
		Lat:       pf(59.387261),
		Lon:       pf(18.054143),
		CreatedAt: time.Now().UTC(),
	}

	if err := InsertMessage(db, msg); err != nil {
		t.Fatal(err)
	}

	msgs, err := ListMessages(db)
	if err != nil {
		t.Fatal(err)
	}
	if len(msgs) != 1 {
		t.Fatalf("got %d messages, want 1", len(msgs))
	}
	if msgs[0].GarminID != "test123" {
		t.Errorf("garmin_id = %q, want test123", msgs[0].GarminID)
	}
	if !msgs[0].HasLocation() || *msgs[0].Lat != 59.387261 {
		t.Errorf("lat = %v, want 59.387261", msgs[0].Lat)
	}
}

func TestInsertWithoutLocation(t *testing.T) {
	db := testDB(t)

	msg := Message{
		GarminID:  "nocoords",
		Sender:    "Henrik",
		Text:      "No location available",
		CreatedAt: time.Now().UTC(),
	}

	if err := InsertMessage(db, msg); err != nil {
		t.Fatal(err)
	}

	msgs, err := ListMessages(db)
	if err != nil {
		t.Fatal(err)
	}
	if len(msgs) != 1 {
		t.Fatalf("got %d messages, want 1", len(msgs))
	}
	if msgs[0].HasLocation() {
		t.Error("expected no location")
	}
}

func TestIdempotentInsert(t *testing.T) {
	db := testDB(t)

	msg := Message{
		GarminID:  "dup456",
		Sender:    "Henrik",
		Text:      "First message",
		Lat:       pf(59.0),
		Lon:       pf(18.0),
		CreatedAt: time.Now().UTC(),
	}

	if err := InsertMessage(db, msg); err != nil {
		t.Fatal(err)
	}

	// Insert same garmin_id again — should be silently ignored
	msg.Text = "Duplicate message"
	if err := InsertMessage(db, msg); err != nil {
		t.Fatal(err)
	}

	msgs, err := ListMessages(db)
	if err != nil {
		t.Fatal(err)
	}
	if len(msgs) != 1 {
		t.Fatalf("got %d messages, want 1 (duplicate should be ignored)", len(msgs))
	}
	if msgs[0].Text != "First message" {
		t.Errorf("text = %q, want 'First message' (original should be kept)", msgs[0].Text)
	}
}

func TestListEmpty(t *testing.T) {
	db := testDB(t)

	msgs, err := ListMessages(db)
	if err != nil {
		t.Fatal(err)
	}
	if msgs != nil {
		t.Errorf("expected nil for empty table, got %v", msgs)
	}
}
