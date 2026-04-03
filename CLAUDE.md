# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

Bolivar receives SMS messages from a Garmin inReach device via 46elks webhooks, parses geolocation from the embedded Garmin map link, stores messages in SQLite, and displays them on a Mapbox Outdoors map.

red/green tdd

## Commands

```bash
go build ./...          # build
go test ./...           # run all tests
go test ./... -v        # verbose test output
go test -run TestName   # run a single test
go vet ./...            # static analysis
```

Run the server (requires env vars):
```bash
export MAPBOX_TOKEN=... WEBHOOK_USER=... WEBHOOK_PASS=...
go run .
```

## Architecture

Single `main` package, flat layout. No frameworks — stdlib `net/http` with Go 1.22+ routing.

- `main.go` — config from env vars, route wiring, server start
- `handler.go` — HTTP handlers: `POST /webhook` (Basic Auth), `GET /api/messages`, `GET /health`, `GET /`
- `garmin.go` — fetches Garmin inReach page HTML, extracts lat/lon with regex
- `db.go` — SQLite via `modernc.org/sqlite` (pure Go), schema migration, insert (idempotent on `garmin_id`), list. Lat/lon are nullable (`*float64`) for messages without GPS coordinates
- `telegram.go` — optional Telegram notifications (text message + location pin). Nil-safe — disabled when env vars are unset
- `static/index.html` — Mapbox GL JS map page, embedded via `go:embed`, templated with `html/template` to inject Mapbox token

## Data Flow

1. 46elks POSTs form-encoded SMS to `/webhook` (authenticated via Basic Auth)
2. Handler extracts Garmin URL from message text
3. Garmin parser fetches the page and extracts lat/lon via regex
4. Message stored in SQLite (duplicate `garmin_id` silently ignored)
5. Map page fetches `/api/messages` JSON and renders markers

## Configuration

All via environment variables. Required: `MAPBOX_TOKEN`, `WEBHOOK_USER`, `WEBHOOK_PASS`. Optional: `TELEGRAM_BOT_TOKEN`, `TELEGRAM_CHAT_ID`, `DB_PATH` (default: `bolivar.db`), `LISTEN_ADDR` (default: `:8080`).

## Workflow

- Create a PR for all changes — do not push directly to master.
- CI runs tests and deploys to Fly.io on merge to master.
