package main

import (
	"errors"

	"github.com/mateuszkozakiewicz/lidl-coupons/internal/config"
	"github.com/mateuszkozakiewicz/lidl-coupons/internal/lidl"
	"github.com/mateuszkozakiewicz/lidl-coupons/internal/notification"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	cfg := config.Load()
	zerolog.SetGlobalLevel(cfg.LogLevel)

	log.Info().Msg("starting lidl-coupons")

	lidlClient := lidl.New(cfg.Lidl)
	promotions, err := lidlClient.ActivateAllPromotions()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to activate promotions")
	}

	if len(promotions) > 0 {
		notifier := notification.New(cfg.Notification)
		if err := notifier.NotifyDiscord(promotions); err != nil {
			if errors.Is(err, notification.ErrWebhookURLNotSet) {
				log.Warn().Msg("webhook URL not set, skipping notification")
			} else {
				log.Error().Err(err).Msg("failed to send notification")
			}
		}
	} else {
		log.Info().Msg("no coupons activated")
	}
}
