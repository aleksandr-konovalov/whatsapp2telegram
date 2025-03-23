package main

import (
	"fmt"
	"os"

	"github.com/aleksandr-konovalov/whatsapp2telegram/pkg/config"
	"github.com/aleksandr-konovalov/whatsapp2telegram/pkg/exporter"
	"github.com/spf13/cobra"
)

var (
	configFile string
	chatName   string
)

func main() {
	// Define root command
	rootCmd := &cobra.Command{
		Use:   "whatsapp2telegram",
		Short: "Export WhatsApp chats to Telegram",
		Long:  `A command-line tool to export WhatsApp chats, including messages and media, to Telegram channels or groups.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Ensure config file exists
			return config.EnsureConfigFile(configFile)
		},
	}

	// Add global flags
	rootCmd.PersistentFlags().StringVar(&configFile, "config", "config", "Config file name without extension")

	// Login command
	loginCmd := &cobra.Command{
		Use:   "login",
		Short: "Login to WhatsApp",
		Long:  `Login to WhatsApp by scanning the QR code with your phone.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return loginToWhatsApp()
		},
	}

	// Export command
	exportCmd := &cobra.Command{
		Use:   "export",
		Short: "Export WhatsApp chat to Telegram",
		Long:  `Export a specific WhatsApp chat to Telegram.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return exportChat()
		},
	}
	exportCmd.Flags().StringVarP(&chatName, "chat", "c", "", "Chat name to export (required)")
	exportCmd.MarkFlagRequired("chat")

	// Add commands to root
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(exportCmd)

	// Execute
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// loginToWhatsApp handles logging in to WhatsApp
func loginToWhatsApp() error {
	// Load config
	cfg, err := config.Load(configFile)
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	// Create exporter
	exp, err := exporter.NewExporter(cfg)
	if err != nil {
		return fmt.Errorf("error creating exporter: %w", err)
	}

	// Connect to WhatsApp
	if err := exp.Connect(); err != nil {
		return fmt.Errorf("error connecting to WhatsApp: %w", err)
	}

	// Login to WhatsApp
	fmt.Println("Logging in to WhatsApp...")
	if err := exp.Login(); err != nil {
		return fmt.Errorf("error logging in to WhatsApp: %w", err)
	}

	fmt.Println("Login successful!")
	return nil
}

// exportChat handles exporting a WhatsApp chat to Telegram
func exportChat() error {
	// Load config
	cfg, err := config.Load(configFile)
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	// Validate config
	if cfg.Telegram.BotToken == "" {
		return fmt.Errorf("Telegram bot token not configured. Please set it in the config file")
	}
	if cfg.Telegram.ChatID == "" {
		return fmt.Errorf("Telegram chat ID not configured. Please set it in the config file")
	}

	// Create exporter
	exp, err := exporter.NewExporter(cfg)
	if err != nil {
		return fmt.Errorf("error creating exporter: %w", err)
	}

	// Connect to WhatsApp
	if err := exp.Connect(); err != nil {
		return fmt.Errorf("error connecting to WhatsApp: %w", err)
	}
	defer exp.Disconnect()

	// Validate Telegram connection
	if err := exp.ValidateTelegram(); err != nil {
		return fmt.Errorf("error validating Telegram connection: %w", err)
	}

	// Export chat
	fmt.Printf("Exporting chat '%s' to Telegram...\n", chatName)
	if err := exp.ExportChat(chatName); err != nil {
		return fmt.Errorf("error exporting chat: %w", err)
	}

	return nil
} 