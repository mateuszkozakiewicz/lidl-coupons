package notification

import (
	"strings"

	"github.com/gtuk/discordwebhook"
	"github.com/mateuszkozakiewicz/lidl-coupons/internal/lidl"
)

func (n *Notifier) send(message discordwebhook.Message) error {
	if n.cfg.WebhookURL == "" {
		return ErrWebhookURLNotSet
	}
	err := discordwebhook.SendMessage(n.cfg.WebhookURL, message)
	if err != nil {
		return err
	}
	return nil
}

func (n *Notifier) NotifyDiscord(promotions []lidl.Promotion) error {
	message := discordwebhook.Message{
		Username: &n.cfg.Username,
		Content:  formatMessage(promotions),
	}
	return n.send(message)
}

func formatMessage(promotions []lidl.Promotion) *string {
	var message string
	for _, p := range promotions {
		details := strings.ReplaceAll(p.Offer, "*", "")
		details = strings.ReplaceAll(details, "\n", " ")
		details = strings.ReplaceAll(details, "_", " ")
		if len(p.Title) > 3 {
			title := strings.ReplaceAll(p.Title, "*", "")
			title = strings.ReplaceAll(title, "\n", " ")
			title = strings.ReplaceAll(title, "_", " ")
			details += " _(" + title + ")_"
		}
		message += "**" + p.Description + "**" + ": " + details + "\n"
	}
	return &message
}
