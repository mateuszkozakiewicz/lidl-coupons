package config

import (
	"os"

	"github.com/rs/zerolog"
)

type LidlConfig struct {
	LoginURL string
	ApiURL   string
	Login    string
	Password string
	Token    string
}

type NotificationConfig struct {
	Username   string
	WebhookURL string
}

type Config struct {
	Lidl         *LidlConfig
	Notification *NotificationConfig
	LogLevel     zerolog.Level
}

func Load() Config {
	return Config{
		Lidl: &LidlConfig{
			ApiURL:   envOrDefault("API_URL", "https://www.lidl.pl/prm/api/v1/PL/"),
			LoginURL: envOrDefault("LOGIN_URL", "https://www.lidl.pl/mla/"),
			Login:    os.Getenv("LOGIN"),
			Password: os.Getenv("PASSWORD"),
			Token:    os.Getenv("TOKEN"),
		},
		Notification: &NotificationConfig{
			Username:   envOrDefault("USERNAME", "Lidl Bot"),
			WebhookURL: os.Getenv("WEBHOOK_URL"),
		},
		LogLevel: getLogLevel(envOrDefault("LOG_LEVEL", "warn")),
	}
}

func envOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getLogLevel(level string) zerolog.Level {
	parsed, err := zerolog.ParseLevel(level)
	if err != nil {
		return zerolog.InfoLevel
	}
	return parsed
}
