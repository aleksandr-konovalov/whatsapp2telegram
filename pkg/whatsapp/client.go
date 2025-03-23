package whatsapp

import (
	"encoding/gob"
	"fmt"
	"os"
	"time"

	"github.com/Rhymen/go-whatsapp"
)

// Message represents a WhatsApp message
type Message struct {
	ID        string
	FromName  string
	FromPhone string
	Text      string
	Timestamp time.Time
	MediaType string
	MediaURL  string
	MediaData []byte
}

// Client represents a WhatsApp client
type Client struct {
	conn       *whatsapp.Conn
	sessionFile string
	loggedIn   bool
}

// NewClient creates a new WhatsApp client
func NewClient(sessionFile string) *Client {
	return &Client{
		sessionFile: sessionFile,
	}
}

// Connect connects to WhatsApp
func (c *Client) Connect() error {
	// Create new WhatsApp connection
	conn, err := whatsapp.NewConn(5 * time.Second)
	if err != nil {
		return fmt.Errorf("error creating connection: %w", err)
	}
	c.conn = conn

	// Set handler for connection events
	c.conn.AddHandler(&waHandler{})

	// Try to restore session
	if err := c.restoreSession(); err != nil {
		fmt.Println("Session restore failed, new login required")
		c.loggedIn = false
	} else {
		c.loggedIn = true
	}

	return nil
}

// Login logs in to WhatsApp using QR code
func (c *Client) Login() error {
	if c.conn == nil {
		return fmt.Errorf("connection not established, call Connect first")
	}

	// Request QR code and show to user
	qr := make(chan string)
	go func() {
		fmt.Println("Please scan the QR code with your WhatsApp app:")
		for code := range qr {
			fmt.Println(code)
		}
	}()

	session, err := c.conn.Login(qr)
	if err != nil {
		return fmt.Errorf("error during login: %w", err)
	}

	// Save session
	if err := c.saveSession(session); err != nil {
		return fmt.Errorf("error saving session: %w", err)
	}

	c.loggedIn = true
	return nil
}

// Disconnect disconnects from WhatsApp
func (c *Client) Disconnect() error {
	if c.conn != nil {
		_, err := c.conn.Disconnect()
		return err
	}
	return nil
}

// GetChats gets a list of all WhatsApp chats
func (c *Client) GetChats() ([]string, error) {
	if !c.loggedIn {
		return nil, fmt.Errorf("not logged in")
	}

	// This is a placeholder since go-whatsapp doesn't have a direct method to list chats
	// You would typically need to implement a handler to collect chats from incoming messages
	return []string{}, fmt.Errorf("chat listing not implemented in the current library")
}

// GetMessages gets messages from a chat
func (c *Client) GetMessages(chatID string, from, to time.Time) ([]Message, error) {
	if !c.loggedIn {
		return nil, fmt.Errorf("not logged in")
	}

	// This is a placeholder since go-whatsapp doesn't have a direct method to fetch historical messages
	// You would typically need to implement a handler to collect messages from the chat history
	return []Message{}, fmt.Errorf("message history retrieval not implemented in the current library")
}

// saveSession saves the WhatsApp session to file
func (c *Client) saveSession(session whatsapp.Session) error {
	file, err := os.Create(c.sessionFile)
	if err != nil {
		return fmt.Errorf("error creating session file: %w", err)
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(session); err != nil {
		return fmt.Errorf("error encoding session: %w", err)
	}

	return nil
}

// restoreSession restores a WhatsApp session from file
func (c *Client) restoreSession() error {
	file, err := os.Open(c.sessionFile)
	if err != nil {
		return fmt.Errorf("error opening session file: %w", err)
	}
	defer file.Close()

	var session whatsapp.Session
	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&session); err != nil {
		return fmt.Errorf("error decoding session: %w", err)
	}

	// Restore session
	session, err = c.conn.RestoreWithSession(session)
	if err != nil {
		return fmt.Errorf("error restoring session: %w", err)
	}

	// Save the new session
	if err := c.saveSession(session); err != nil {
		return fmt.Errorf("error saving restored session: %w", err)
	}

	return nil
}

// waHandler implements the WhatsApp message handler interface
type waHandler struct{}

func (h *waHandler) HandleError(err error) {
	fmt.Fprintf(os.Stderr, "WhatsApp error: %v\n", err)
}

func (h *waHandler) HandleTextMessage(message whatsapp.TextMessage) {
	// Placeholder for handling incoming text messages
}

func (h *waHandler) HandleImageMessage(message whatsapp.ImageMessage) {
	// Placeholder for handling incoming image messages
}

func (h *waHandler) HandleDocumentMessage(message whatsapp.DocumentMessage) {
	// Placeholder for handling incoming document messages
}

func (h *waHandler) HandleVideoMessage(message whatsapp.VideoMessage) {
	// Placeholder for handling incoming video messages
}

func (h *waHandler) HandleAudioMessage(message whatsapp.AudioMessage) {
	// Placeholder for handling incoming audio messages
}
