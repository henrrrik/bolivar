# Bolivar

![](https://vignette.wikia.nocookie.net/disney/images/c/c3/Tumblr_n1scdfSDfE1qhcrb0o1_1280.jpg/revision/latest/scale-to-width-down/185?cb=20140302031630)

Bolivar receives Garmin inReach satellite text messages through [46elks](https://46elks.se) and plots them on a nicer map from [Mapbox](https://www.mapbox.com).

<img width="662" height="436" alt="Skärmavbild 2026-04-03 kl  00 31 06" src="https://github.com/user-attachments/assets/13509eb1-9bbd-49a8-bd6c-67d37041e28e" />

## Environment Variables

| Variable | Required | Default | Description |
|---|---|---|---|
| `MAPBOX_TOKEN` | Yes | — | Mapbox GL JS access token |
| `WEBHOOK_USER` | Yes | — | Basic auth username for the 46elks webhook |
| `WEBHOOK_PASS` | Yes | — | Basic auth password for the 46elks webhook |
| `TELEGRAM_BOT_TOKEN` | No | — | Telegram bot API token for notifications |
| `TELEGRAM_CHAT_ID` | No | — | Telegram group chat ID for notifications |
| `DB_PATH` | No | `bolivar.db` | Path to the SQLite database file |
| `LISTEN_ADDR` | No | `:8080` | Address to listen on |
