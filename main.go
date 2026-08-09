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
	log.Info().Msg("starting lidl-coupons")

	cfg := config.Load()
	zerolog.SetGlobalLevel(cfg.LogLevel)

	lidlClient := lidl.New(cfg.Lidl)
	var promotions []lidl.Promotion
	for i := 0; i < cfg.Retries; i++ {
		var err error
		promotions, err = lidlClient.ActivateAllPromotions()
		if err != nil {
			log.Warn().Err(err).Msgf("failed to activate promotions, retrying (%d/%d)", i+1, cfg.Retries)
		} else {
			break
		}
		if i == cfg.Retries-1 {
			log.Fatal().Err(err).Msg("failed to activate promotions after maximum retries")
			return
		}
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
