package exporter

import (
	"fmt"
	"time"

	"github.com/aleksandr-konovalov/whatsapp2telegram/pkg/config"
	"github.com/aleksandr-konovalov/whatsapp2telegram/pkg/telegram"
	"github.com/aleksandr-konovalov/whatsapp2telegram/pkg/whatsapp"
)

// Exporter handles the export process from WhatsApp to Telegram
type Exporter struct {
	config   *config.Config
	whatsapp *whatsapp.Client
	telegram *telegram.Client
}

// NewExporter creates a new exporter instance
func NewExporter(cfg *config.Config) (*Exporter, error) {
	// Create WhatsApp client
	waClient := whatsapp.NewClient(cfg.WhatsApp.SessionFile)

	// Create Telegram client
	tgClient, err := telegram.NewClient(cfg.Telegram.BotToken, cfg.Telegram.ChatID)
	if err != nil {
		return nil, fmt.Errorf("error creating Telegram client: %w", err)
	}

	return &Exporter{
		config:   cfg,
		whatsapp: waClient,
		telegram: tgClient,
	}, nil
}

// Connect connects to WhatsApp
func (e *Exporter) Connect() error {
	return e.whatsapp.Connect()
}

// Login logs in to WhatsApp
func (e *Exporter) Login() error {
	return e.whatsapp.Login()
}

// Disconnect disconnects from WhatsApp
func (e *Exporter) Disconnect() error {
	return e.whatsapp.Disconnect()
}

// ValidateTelegram validates the Telegram connection
func (e *Exporter) ValidateTelegram() error {
	return e.telegram.Validate()
}

// ListChats lists all WhatsApp chats
func (e *Exporter) ListChats() ([]string, error) {
	return e.whatsapp.GetChats()
}

// ExportChat exports a WhatsApp chat to Telegram
func (e *Exporter) ExportChat(chatID string) error {
	fmt.Println("Starting export process...")

	// Parse date range from config
	from, to, err := e.config.Export.ParseDateRange()
	if err != nil {
		return fmt.Errorf("error parsing date range: %w", err)
	}

	// Get messages from WhatsApp
	fmt.Printf("Fetching messages from WhatsApp (from %s to %s)...\n", 
		from.Format("2006-01-02"), 
		to.Format("2006-01-02"))
	
	messages, err := e.whatsapp.GetMessages(chatID, from, to)
	if err != nil {
		return fmt.Errorf("error getting messages: %w", err)
	}

	// Export messages to Telegram
	fmt.Printf("Exporting %d messages to Telegram...\n", len(messages))
	err = e.telegram.ExportMessages(messages, e.config.Export.IncludeMedia)
	if err != nil {
		return fmt.Errorf("error exporting messages: %w", err)
	}

	fmt.Println("Export completed successfully!")
	return nil
}

// ExportAllChats exports all WhatsApp chats to Telegram
func (e *Exporter) ExportAllChats() error {
	// Get all chats
	chats, err := e.ListChats()
	if err != nil {
		return fmt.Errorf("error listing chats: %w", err)
	}

	// Export each chat
	for _, chatID := range chats {
		fmt.Printf("Exporting chat: %s\n", chatID)
		if err := e.ExportChat(chatID); err != nil {
			fmt.Printf("Warning: Error exporting chat %s: %v\n", chatID, err)
		}
		
		// Add a delay between chats to avoid hitting rate limits
		time.Sleep(1 * time.Second)
	}

	return nil
}
