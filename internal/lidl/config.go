package lidl

import (
	"github.com/mateuszkozakiewicz/lidl-coupons/internal/config"
)

type Client struct {
	cfg *config.LidlConfig
}

func New(cfg *config.LidlConfig) *Client {
	return &Client{
		cfg: cfg,
	}
}
