package telegram

import (
	"fmt"
	"time"

	"github.com/aleksandr-konovalov/whatsapp2telegram/pkg/whatsapp"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Client represents a Telegram client
type Client struct {
	bot    *tgbotapi.BotAPI
	chatID string
}

// NewClient creates a new Telegram client
func NewClient(token, chatID string) (*Client, error) {
	// Create new bot instance
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("error creating Telegram bot: %w", err)
	}

	// Set debug mode to false in production
	bot.Debug = false

	return &Client{
		bot:    bot,
		chatID: chatID,
	}, nil
}

// SendMessage sends a text message to Telegram
func (c *Client) SendMessage(text string) error {
	msg := tgbotapi.NewMessage(parseChatID(c.chatID), text)
	_, err := c.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("error sending message: %w", err)
	}
	return nil
}

// SendMedia sends media to Telegram
func (c *Client) SendMedia(mediaType string, data []byte, caption string) error {
	chatID := parseChatID(c.chatID)

	fileBytes := tgbotapi.FileBytes{
		Name:  fmt.Sprintf("media_%d", time.Now().Unix()),
		Bytes: data,
	}

	var msg tgbotapi.Chattable
	switch mediaType {
	case "image":
		photoMsg := tgbotapi.NewPhoto(chatID, fileBytes)
		photoMsg.Caption = caption
		msg = photoMsg
	case "video":
		videoMsg := tgbotapi.NewVideo(chatID, fileBytes)
		videoMsg.Caption = caption
		msg = videoMsg
	case "audio":
		audioMsg := tgbotapi.NewAudio(chatID, fileBytes)
		audioMsg.Caption = caption
		msg = audioMsg
	case "document":
		docMsg := tgbotapi.NewDocument(chatID, fileBytes)
		docMsg.Caption = caption
		msg = docMsg
	default:
		return fmt.Errorf("unsupported media type: %s", mediaType)
	}

	_, err := c.bot.Send(msg)
	if err != nil {
		return fmt.Errorf("error sending media: %w", err)
	}
	return nil
}

// ExportMessages exports WhatsApp messages to Telegram
func (c *Client) ExportMessages(messages []whatsapp.Message, includeMedia bool) error {
	// Sort messages by timestamp
	// In a real application, you would implement a proper sort here
	
	// Group messages by day for better organization
	msgHeader := fmt.Sprintf("📱 WhatsApp Export (%d messages)", len(messages))
	if err := c.SendMessage(msgHeader); err != nil {
		return err
	}

	// Process each message
	for _, msg := range messages {
		formattedMsg := formatMessage(msg)
		
		// Send the text message
		if err := c.SendMessage(formattedMsg); err != nil {
			return err
		}
		
		// If media is included and this message has media, send it
		if includeMedia && msg.MediaData != nil && len(msg.MediaData) > 0 {
			if err := c.SendMedia(msg.MediaType, msg.MediaData, ""); err != nil {
				fmt.Printf("Warning: Failed to send media for message %s: %v\n", msg.ID, err)
			}
		}
		
		// Throttle to avoid hitting Telegram rate limits
		time.Sleep(100 * time.Millisecond)
	}
	
	return nil
}

// Validate checks if the Telegram bot token and chat ID are valid
func (c *Client) Validate() error {
	// Check if we can get bot information
	_, err := c.bot.GetMe()
	if err != nil {
		return fmt.Errorf("invalid bot token: %w", err)
	}

	// Try to send a test message to verify the chat ID
	testMsg := tgbotapi.NewMessage(parseChatID(c.chatID), "🔄 Testing connection...")
	_, err = c.bot.Send(testMsg)
	if err != nil {
		return fmt.Errorf("invalid chat ID: %w", err)
	}

	return nil
}

// formatMessage formats a WhatsApp message for Telegram display
func formatMessage(msg whatsapp.Message) string {
	timestamp := msg.Timestamp.Format("2006-01-02 15:04:05")
	
	// Format the message
	formatted := fmt.Sprintf("From: %s (%s)\nTime: %s\n\n%s", 
		msg.FromName,
		msg.FromPhone,
		timestamp,
		msg.Text)
	
	// Add media info if present
	if msg.MediaType != "" {
		formatted += fmt.Sprintf("\n\n[Contains %s]", msg.MediaType)
	}
	
	return formatted
}

// parseChatID parses a string chat ID to int64
func parseChatID(chatID string) int64 {
	var id int64
	fmt.Sscanf(chatID, "%d", &id)
	return id
}
