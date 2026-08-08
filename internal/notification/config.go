package notification

import (
	"github.com/mateuszkozakiewicz/lidl-coupons/internal/config"
)

type Notifier struct {
	cfg *config.NotificationConfig
}

func New(cfg *config.NotificationConfig) *Notifier {
	return &Notifier{
		cfg: cfg,
	}
}
