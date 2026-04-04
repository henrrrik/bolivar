# Bolivar

![](https://vignette.wikia.nocookie.net/disney/images/c/c3/Tumblr_n1scdfSDfE1qhcrb0o1_1280.jpg/revision/latest/scale-to-width-down/185?cb=20140302031630)

Garmin inReach devices are really nice, but the backend servie is clunky and kind of ugly (at least the version that supports my inReach mini). Bolivar aims to improve this by providing a nicer way for your loved ones to track your journey through the wilderness. 

Bolivar receives Garmin inReach satellite text messages through an [46elks](https://46elks.se) phone number, strips out the Garmin map URLs and parses the lat/lon from them, plots them on a map from [Mapbox](https://www.mapbox.com), and optionally forwards them to a [Telegram](https://telegram.org) group along with a map pin. Just add whoever is interested in your trip to the Telegram group instead of adding them as Garmin contacts.

## How it looks

### Web view
<img width="662" height="436" alt="Mapbox Outdoor Map" src="https://github.com/user-attachments/assets/13509eb1-9bbd-49a8-bd6c-67d37041e28e" />

### Telegram
<img width="331" alt="Telegram Group" src="https://github.com/user-attachments/assets/f40b8f80-d09c-4d3f-aa40-08f04a79dcb9" />

## Requirements

- A Garmin inReach with a subscription
- An 46elks account with a phone number
- A Mapbox API key
- Somewhere to host this (I use [Runway](https://www.runway.horse))
- (Optionally) A Telegram bot and a group


## Environment Variables

| Variable | Required | Default | Description |
|---|---|---|---|
| `MAPBOX_TOKEN` | Yes | — | Mapbox GL JS access token |
| `WEBHOOK_USER` | Yes | — | Basic auth username for the 46elks webhook |
| `WEBHOOK_PASS` | Yes | — | Basic auth password for the 46elks webhook |
| `TELEGRAM_BOT_TOKEN` | No | — | Telegram bot API token for notifications |
| `TELEGRAM_CHAT_ID` | No | — | Telegram group chat ID for notifications |
| `DATABASE_URL` | Yes | — | PostgreSQL connection string |
| `LISTEN_ADDR` | No | `:8080` | Address to listen on |
