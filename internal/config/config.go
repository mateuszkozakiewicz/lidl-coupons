package config

import (
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
)

type LidlConfig struct {
	LoginURL    string
	ApiURL      string
	StoragePath string
	Timeout     float64
	Login       string
	Password    string
	Token       string
}

type NotificationConfig struct {
	Username   string
	WebhookURL string
}

type Config struct {
	Lidl         *LidlConfig
	Notification *NotificationConfig
	LogLevel     zerolog.Level
	Retries      int
}

func Load() Config {
	return Config{
		Lidl: &LidlConfig{
			ApiURL:      envOrDefault("API_URL", "https://www.lidl.pl/prm/api/v1/PL/"),
			LoginURL:    envOrDefault("LOGIN_URL", "https://www.lidl.pl/mla/"),
			StoragePath: envOrDefault("STORAGE_PATH", "./playwright-data"),
			Timeout:     getDurationMillis(envOrDefault("TIMEOUT", "5s")),
			Login:       os.Getenv("LOGIN"),
			Password:    os.Getenv("PASSWORD"),
			Token:       os.Getenv("TOKEN"),
		},
		Notification: &NotificationConfig{
			Username:   envOrDefault("USERNAME", "Lidl Bot"),
			WebhookURL: os.Getenv("WEBHOOK_URL"),
		},
		LogLevel: getLogLevel(envOrDefault("LOG_LEVEL", "warn")),
		Retries:  getInt(envOrDefault("RETRIES", "10")),
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

func getDurationMillis(value string) float64 {
	duration, err := time.ParseDuration(value)
	if err != nil {
		return 30000
	}
	return float64(duration / time.Millisecond)
}

func getInt(value string) int {
	i, err := strconv.Atoi(value)
	if err != nil {
		return 10
	}
	return i
}
