package config

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
)

func unsetEnv(t *testing.T, key string) {
	t.Helper()
	original, wasSet := os.LookupEnv(key)
	if err := os.Unsetenv(key); err != nil {
		t.Fatalf("failed to unset %s: %v", key, err)
	}
	t.Cleanup(func() {
		if wasSet {
			os.Setenv(key, original)
		}
	})
}

func TestEnvOrDefault(t *testing.T) {
	t.Run("returns default when unset", func(t *testing.T) {
		if got := envOrDefault("LIDL_COUPONS_TEST_UNSET", "default"); got != "default" {
			t.Errorf("got %q, want %q", got, "default")
		}
	})

	t.Run("returns env value when set", func(t *testing.T) {
		t.Setenv("LIDL_COUPONS_TEST_SET", "value")
		if got := envOrDefault("LIDL_COUPONS_TEST_SET", "default"); got != "value" {
			t.Errorf("got %q, want %q", got, "value")
		}
	})

	t.Run("returns empty string when explicitly set empty", func(t *testing.T) {
		t.Setenv("LIDL_COUPONS_TEST_EMPTY", "")
		if got := envOrDefault("LIDL_COUPONS_TEST_EMPTY", "default"); got != "" {
			t.Errorf("got %q, want empty string", got)
		}
	})
}

func TestGetLogLevel(t *testing.T) {
	tests := []struct {
		input string
		want  zerolog.Level
	}{
		{"debug", zerolog.DebugLevel},
		{"warn", zerolog.WarnLevel},
		{"error", zerolog.ErrorLevel},
		{"not-a-level", zerolog.InfoLevel},
	}
	for _, tt := range tests {
		if got := getLogLevel(tt.input); got != tt.want {
			t.Errorf("getLogLevel(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestGetDurationMillis(t *testing.T) {
	tests := []struct {
		input string
		want  float64
	}{
		{"5s", 5000},
		{"1500ms", 1500},
		{"2m", 120000},
		{"not-a-duration", 30000},
		{"", 30000},
	}
	for _, tt := range tests {
		if got := getDurationMillis(tt.input); got != tt.want {
			t.Errorf("getDurationMillis(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestGetInt(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"5", 5},
		{"0", 0},
		{"not-a-number", 10},
		{"", 10},
	}
	for _, tt := range tests {
		if got := getInt(tt.input); got != tt.want {
			t.Errorf("getInt(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestLoad(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		for _, key := range []string{
			"API_URL", "LOGIN_URL", "STORAGE_PATH", "TIMEOUT", "LOGIN", "PASSWORD",
			"TOKEN", "USERNAME", "WEBHOOK_URL", "LOG_LEVEL", "RETRIES",
		} {
			unsetEnv(t, key)
		}

		cfg := Load()

		if cfg.Lidl.ApiURL != "https://www.lidl.pl/prm/api/v1/PL/" {
			t.Errorf("unexpected default ApiURL: %q", cfg.Lidl.ApiURL)
		}
		if cfg.Lidl.LoginURL != "https://www.lidl.pl/mla/" {
			t.Errorf("unexpected default LoginURL: %q", cfg.Lidl.LoginURL)
		}
		if cfg.Lidl.StoragePath != "./playwright-data" {
			t.Errorf("unexpected default StoragePath: %q", cfg.Lidl.StoragePath)
		}
		if cfg.Lidl.Timeout != 5000 {
			t.Errorf("unexpected default Timeout: %v", cfg.Lidl.Timeout)
		}
		if cfg.Notification.Username != "Lidl Bot" {
			t.Errorf("unexpected default Username: %q", cfg.Notification.Username)
		}
		if cfg.LogLevel != zerolog.WarnLevel {
			t.Errorf("unexpected default LogLevel: %v", cfg.LogLevel)
		}
		if cfg.Retries != 10 {
			t.Errorf("unexpected default Retries: %v", cfg.Retries)
		}
	})

	t.Run("overrides from env", func(t *testing.T) {
		t.Setenv("API_URL", "https://example.test/api/")
		t.Setenv("LOGIN", "user@example.test")
		t.Setenv("PASSWORD", "hunter2")
		t.Setenv("TOKEN", "sometoken")
		t.Setenv("WEBHOOK_URL", "https://example.test/webhook")
		t.Setenv("LOG_LEVEL", "debug")
		t.Setenv("RETRIES", "3")

		cfg := Load()

		if cfg.Lidl.ApiURL != "https://example.test/api/" {
			t.Errorf("unexpected ApiURL: %q", cfg.Lidl.ApiURL)
		}
		if cfg.Lidl.Login != "user@example.test" {
			t.Errorf("unexpected Login: %q", cfg.Lidl.Login)
		}
		if cfg.Lidl.Password != "hunter2" {
			t.Errorf("unexpected Password: %q", cfg.Lidl.Password)
		}
		if cfg.Lidl.Token != "sometoken" {
			t.Errorf("unexpected Token: %q", cfg.Lidl.Token)
		}
		if cfg.Notification.WebhookURL != "https://example.test/webhook" {
			t.Errorf("unexpected WebhookURL: %q", cfg.Notification.WebhookURL)
		}
		if cfg.LogLevel != zerolog.DebugLevel {
			t.Errorf("unexpected LogLevel: %v", cfg.LogLevel)
		}
		if cfg.Retries != 3 {
			t.Errorf("unexpected Retries: %v", cfg.Retries)
		}
	})
}
