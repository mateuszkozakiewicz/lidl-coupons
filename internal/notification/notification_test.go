package notification

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mateuszkozakiewicz/lidl-coupons/internal/config"
	"github.com/mateuszkozakiewicz/lidl-coupons/internal/lidl"
)

func TestFormatMessages(t *testing.T) {
	t.Run("includes description and offer, strips markdown", func(t *testing.T) {
		promotions := []lidl.Promotion{
			{Description: "Store A", Offer: "10% *off*\nsome_thing"},
		}
		messages := *formatMessages(promotions)
		if len(messages) != 1 {
			t.Fatalf("got %d messages, want 1", len(messages))
		}
		if strings.Contains(messages[0], "*off*") || strings.Contains(messages[0], "_") {
			t.Errorf("expected markdown characters to be stripped, got: %q", messages[0])
		}
		if !strings.Contains(messages[0], "**Store A**") {
			t.Errorf("expected description as bold header, got: %q", messages[0])
		}
	})

	t.Run("omits short titles but includes longer ones", func(t *testing.T) {
		promotions := []lidl.Promotion{
			{Description: "Store A", Offer: "Offer A", Title: "AB"},
			{Description: "Store B", Offer: "Offer B", Title: "Long Title"},
		}
		messages := *formatMessages(promotions)
		if len(messages) != 1 {
			t.Fatalf("got %d messages, want 1", len(messages))
		}
		if strings.Contains(messages[0], "AB") {
			t.Errorf("expected short title to be omitted, got: %q", messages[0])
		}
		if !strings.Contains(messages[0], "Long Title") {
			t.Errorf("expected long title to be included, got: %q", messages[0])
		}
	})

	t.Run("splits into multiple messages when exceeding limit", func(t *testing.T) {
		var promotions []lidl.Promotion
		line := strings.Repeat("x", 200)
		for i := 0; i < 20; i++ {
			promotions = append(promotions, lidl.Promotion{Description: "Store", Offer: line})
		}
		messages := *formatMessages(promotions)
		if len(messages) < 2 {
			t.Fatalf("expected messages to be split, got %d message(s)", len(messages))
		}
		for i, msg := range messages {
			if len(msg) > 1900+300 { // allow for the line that pushed it over before splitting
				t.Errorf("message %d too long: %d chars", i, len(msg))
			}
		}
	})
}

func TestNotifyDiscord_WebhookURLNotSet(t *testing.T) {
	n := New(&config.NotificationConfig{Username: "Bot"})
	err := n.NotifyDiscord([]lidl.Promotion{{Description: "Store", Offer: "Offer"}})
	if err != ErrWebhookURLNotSet {
		t.Errorf("got %v, want %v", err, ErrWebhookURLNotSet)
	}
}

func TestNotifyDiscord_Success(t *testing.T) {
	var receivedBody string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		receivedBody = string(body)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	n := New(&config.NotificationConfig{Username: "Bot", WebhookURL: server.URL})
	err := n.NotifyDiscord([]lidl.Promotion{{Description: "Store", Offer: "10% off"}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(receivedBody, "Store") || !strings.Contains(receivedBody, "10% off") {
		t.Errorf("unexpected webhook payload: %s", receivedBody)
	}
}

func TestNotifyDiscord_ServerError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	n := New(&config.NotificationConfig{Username: "Bot", WebhookURL: server.URL})
	err := n.NotifyDiscord([]lidl.Promotion{{Description: "Store", Offer: "10% off"}})
	if err == nil {
		t.Error("expected error when webhook returns non-2xx status")
	}
}
