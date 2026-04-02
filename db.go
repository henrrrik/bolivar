package main

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

type Message struct {
	ID        int64     `json:"id"`
	GarminID  string    `json:"garmin_id"`
	Sender    string    `json:"sender"`
	Text      string    `json:"message"`
	Lat       float64   `json:"lat"`
	Lon       float64   `json:"lon"`
	CreatedAt time.Time `json:"created_at"`
}

func OpenDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS messages (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			garmin_id  TEXT,
			sender     TEXT,
			message    TEXT,
			lat        REAL NOT NULL,
			lon        REAL NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_garmin_id ON messages(garmin_id);
	`); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func InsertMessage(db *sql.DB, m Message) error {
	_, err := db.Exec(
		`INSERT OR IGNORE INTO messages (garmin_id, sender, message, lat, lon, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
		m.GarminID, m.Sender, m.Text, m.Lat, m.Lon, m.CreatedAt,
	)
	return err
}

func ListMessages(db *sql.DB) ([]Message, error) {
	rows, err := db.Query(`SELECT id, garmin_id, sender, message, lat, lon, created_at FROM messages ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []Message
	for rows.Next() {
		var m Message
		if err := rows.Scan(&m.ID, &m.GarminID, &m.Sender, &m.Text, &m.Lat, &m.Lon, &m.CreatedAt); err != nil {
			return nil, err
		}
		msgs = append(msgs, m)
	}
	return msgs, rows.Err()
}
