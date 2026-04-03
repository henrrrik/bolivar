package main

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

type Message struct {
	ID        int64    `json:"id"`
	GarminID  string   `json:"garmin_id"`
	Sender    string   `json:"sender"`
	Text      string   `json:"message"`
	URL       string   `json:"url"`
	Lat       *float64 `json:"lat"`
	Lon       *float64 `json:"lon"`
	CreatedAt time.Time `json:"created_at"`
}

func (m Message) HasLocation() bool {
	return m.Lat != nil && m.Lon != nil
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
			url        TEXT,
			lat        REAL,
			lon        REAL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_garmin_id ON messages(garmin_id);
	`); err != nil {
		db.Close()
		return nil, err
	}

	// Migrate existing databases that lack the url column.
	db.Exec(`ALTER TABLE messages ADD COLUMN url TEXT`)

	return db, nil
}

func InsertMessage(db *sql.DB, m Message) error {
	_, err := db.Exec(
		`INSERT OR IGNORE INTO messages (garmin_id, sender, message, url, lat, lon, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		m.GarminID, m.Sender, m.Text, m.URL, m.Lat, m.Lon, m.CreatedAt,
	)
	return err
}

func ListMessages(db *sql.DB) ([]Message, error) {
	rows, err := db.Query(`SELECT id, garmin_id, sender, message, COALESCE(url, ''), lat, lon, created_at FROM messages ORDER BY created_at DESC LIMIT 30`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []Message
	for rows.Next() {
		var m Message
		var lat, lon sql.NullFloat64
		if err := rows.Scan(&m.ID, &m.GarminID, &m.Sender, &m.Text, &m.URL, &lat, &lon, &m.CreatedAt); err != nil {
			return nil, err
		}
		if lat.Valid {
			m.Lat = &lat.Float64
		}
		if lon.Valid {
			m.Lon = &lon.Float64
		}
		msgs = append(msgs, m)
	}
	return msgs, rows.Err()
}
