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
	messages := formatMessages(promotions)
	for _, msg := range *messages {
		message := discordwebhook.Message{
			Username: &n.cfg.Username,
			Content:  &msg,
		}
		if err := n.send(message); err != nil {
			return err
		}
	}
	return nil
}

func formatMessages(promotions []lidl.Promotion) *[]string {
	var messages []string
	msg := ""
	for _, p := range promotions {
		details := sanitizeString(p.Offer)
		description := sanitizeString(p.Description)
		if len(p.Title) > 3 {
			title := sanitizeString(p.Title)
			details += " _(" + title + ")_"
		}
		line := "**" + description + "**" + ": " + details + "\n"
		if len(msg)+len(line) > 1900 {
			messages = append(messages, msg)
			msg = line
		} else {
			msg += line
		}
	}
	messages = append(messages, msg)
	return &messages
}

func sanitizeString(input string) string {
	sanitized := strings.ReplaceAll(input, "*", "")
	sanitized = strings.ReplaceAll(sanitized, "\n", " ")
	sanitized = strings.ReplaceAll(sanitized, "_", " ")
	return sanitized
}
