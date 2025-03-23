package config

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

// Config represents the application configuration
type Config struct {
	WhatsApp WhatsAppConfig `mapstructure:"whatsapp"`
	Telegram TelegramConfig `mapstructure:"telegram"`
	Export   ExportConfig   `mapstructure:"export"`
}

// WhatsAppConfig holds WhatsApp-specific configuration
type WhatsAppConfig struct {
	SessionFile string `mapstructure:"session_file"`
}

// TelegramConfig holds Telegram-specific configuration
type TelegramConfig struct {
	BotToken string `mapstructure:"bot_token"`
	ChatID   string `mapstructure:"chat_id"`
}

// ExportConfig holds export-specific configuration
type ExportConfig struct {
	IncludeMedia bool   `mapstructure:"include_media"`
	DateFrom     string `mapstructure:"date_from"`
	DateTo       string `mapstructure:"date_to"`
}

// ParseDateRange parses the configured date range
func (e *ExportConfig) ParseDateRange() (from, to time.Time, err error) {
	if e.DateFrom != "" {
		from, err = time.Parse("2006-01-02", e.DateFrom)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid date_from format: %w", err)
		}
	}

	if e.DateTo != "" {
		to, err = time.Parse("2006-01-02", e.DateTo)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid date_to format: %w", err)
		}
	} else {
		to = time.Now()
	}

	return from, to, nil
}

// Load loads configuration from file and environment variables
func Load(configFile string) (*Config, error) {
	viper.SetConfigName(configFile)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	// Set default values
	viper.SetDefault("whatsapp.session_file", "whatsapp_session.gob")
	viper.SetDefault("export.include_media", true)

	// Load from environment variables (optional)
	viper.AutomaticEnv()

	// Read config
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found, create a default one
			defaultCfg := &Config{
				WhatsApp: WhatsAppConfig{
					SessionFile: "whatsapp_session.gob",
				},
				Telegram: TelegramConfig{},
				Export: ExportConfig{
					IncludeMedia: true,
				},
			}
			
			return defaultCfg, nil
		}
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config: %w", err)
	}

	return &cfg, nil
}

// Save saves the configuration to file
func (c *Config) Save(configFile string) error {
	viper.SetConfigName(configFile)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	for k, v := range map[string]interface{}{
		"whatsapp.session_file": c.WhatsApp.SessionFile,
		"telegram.bot_token":    c.Telegram.BotToken,
		"telegram.chat_id":      c.Telegram.ChatID,
		"export.include_media":  c.Export.IncludeMedia,
		"export.date_from":      c.Export.DateFrom,
		"export.date_to":        c.Export.DateTo,
	} {
		viper.Set(k, v)
	}

	return viper.WriteConfig()
}

// EnsureConfigFile makes sure that the config file exists
func EnsureConfigFile(configFile string) error {
	if _, err := os.Stat(configFile + ".yaml"); os.IsNotExist(err) {
		// Create default config file
		cfg := &Config{
			WhatsApp: WhatsAppConfig{
				SessionFile: "whatsapp_session.gob",
			},
			Telegram: TelegramConfig{},
			Export: ExportConfig{
				IncludeMedia: true,
			},
		}
		
		return cfg.Save(configFile)
	}
	return nil
}
