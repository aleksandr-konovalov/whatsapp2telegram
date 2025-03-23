# WhatsApp to Telegram Exporter

A Golang tool for exporting data from WhatsApp and importing it to Telegram.

## Features

- Export messages from WhatsApp individual chats or groups
- Export media (images, videos, documents) from WhatsApp
- Import the exported data to Telegram channels or groups
- Configurable export options (date range, media types, etc.)

## Prerequisites

- Go 1.21 or higher
- WhatsApp account with access to the chats you want to export
- Telegram Bot API token

## Installation

```bash
# Clone the repository
git clone https://github.com/aleksandr-konovalov/whatsapp2telegram.git
cd whatsapp2telegram

# Initialize and download dependencies
go mod tidy

# Build the application
go build -o whatsapp2telegram cmd/whatsapp2telegram/main.go
```

If you encounter dependency errors during the build, try initializing the dependencies explicitly:

```bash
# Explicitly download dependencies
go get github.com/Rhymen/go-whatsapp@v0.1.1
go get github.com/go-telegram-bot-api/telegram-bot-api/v5@v5.5.1
go get github.com/spf13/cobra@v1.7.0
go get github.com/spf13/viper@v1.16.0

# Then build
go build -o whatsapp2telegram cmd/whatsapp2telegram/main.go
```

## Configuration

Create a `config.yaml` file in the root directory with the following content:

```yaml
whatsapp:
  session_file: "whatsapp_session.gob"
telegram:
  bot_token: "YOUR_TELEGRAM_BOT_TOKEN"
  chat_id: "YOUR_TELEGRAM_CHAT_ID"
export:
  include_media: true
  date_from: "2023-01-01"
  date_to: "2023-12-31"
```

## Usage

```bash
# Log in to WhatsApp (only needed once)
./whatsapp2telegram login

# Export WhatsApp chat to Telegram
./whatsapp2telegram export --chat "Chat Name"

# Show help
./whatsapp2telegram --help
```

## Troubleshooting

### Missing go.sum Entries

If you see errors like these:

```
missing go.sum entry for module providing package github.com/spf13/viper
missing go.sum entry for module providing package github.com/Rhymen/go-whatsapp
```

These errors occur when Go can't find checksums for dependencies in the go.sum file. To fix this:

```bash
# Option 1: Use go mod tidy to download dependencies and update go.mod/go.sum
go mod tidy

# Option 2: If that doesn't work, try getting each dependency explicitly
go get github.com/Rhymen/go-whatsapp@v0.1.1
go get github.com/go-telegram-bot-api/telegram-bot-api/v5@v5.5.1
go get github.com/spf13/cobra@v1.7.0
go get github.com/spf13/viper@v1.16.0
```

### Checksum Mismatch

If you encounter a checksum mismatch error like this:

```
verifying github.com/spf13/viper@v1.16.0: checksum mismatch
        downloaded: h1:rGGH0XDZhdUOryiDWjmIvUSWpbNqisK8Wk0Vyefw8hc=
        go.sum:     h1:vawHUee0VqjUf7VdnU56IQBpO2qr4qk7s8S3tzJvTW8=
```

Fix it by removing the go.sum file and letting Go regenerate it:

```bash
rm go.sum
go mod download
```

### WhatsApp Connection Issues

If you encounter issues connecting to WhatsApp, make sure your phone and WhatsApp are up to date. The WhatsApp library used has some limitations with accessing historical data. You might need to explore alternative APIs or techniques for specific use cases.

## License

MIT
