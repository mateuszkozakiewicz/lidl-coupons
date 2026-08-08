package config

import "os"

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
}

func Load() Config {
	return Config{
		Lidl: &LidlConfig{
			ApiURL:   envOrDefault("LIDL_API_URL", "https://www.lidl.pl/prm/api/v1/PL/"),
			LoginURL: envOrDefault("LIDL_LOGIN_URL", "https://www.lidl.pl/mla/"),
			Login:    os.Getenv("LIDL_LOGIN"),
			Password: os.Getenv("LIDL_PASSWORD"),
			Token:    os.Getenv("LIDL_TOKEN"),
		},
		Notification: &NotificationConfig{
			Username:   os.Getenv("DISCORD_USERNAME"),
			WebhookURL: os.Getenv("DISCORD_WEBHOOK_URL"),
		},
	}
}

func envOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
