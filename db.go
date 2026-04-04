package main

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
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

func OpenDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	if _, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS messages (
			id         BIGSERIAL PRIMARY KEY,
			garmin_id  TEXT,
			sender     TEXT,
			message    TEXT,
			url        TEXT,
			lat        DOUBLE PRECISION,
			lon        DOUBLE PRECISION,
			created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
		)
	`); err != nil {
		db.Close()
		return nil, err
	}

	if _, err := db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_garmin_id ON messages(garmin_id)
	`); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func InsertMessage(db *sql.DB, m Message) error {
	_, err := db.Exec(
		`INSERT INTO messages (garmin_id, sender, message, url, lat, lon, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7)
		 ON CONFLICT (garmin_id) DO NOTHING`,
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
