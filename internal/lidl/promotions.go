package lidl

import (
	"github.com/rs/zerolog/log"
)

func (c *Client) ActivateAllPromotions() (promotions []Promotion, err error) {
	if err := c.login(); err != nil {
		return nil, err
	}
	promotions, err = c.getPromotions()
	if err != nil {
		return nil, err
	}
	activatedCount := 0
	for _, promotion := range promotions {
		promotionDetails := promotion.Description + " - " + promotion.Offer
		if len(promotion.Title) > 3 {
			promotionDetails += " (" + promotion.Title + ")"
		}
		if !promotion.IsActivated {
			if err := c.activatePromotion(promotion); err != nil {
				log.Warn().Err(err).Msgf("failed to activate promotion id:%s details:%s", promotion.ID, promotionDetails)
			} else {
				activatedCount++
			}
		}
	}
	log.Info().Msgf("promotions total: %d, activated: %d", len(promotions), activatedCount)
	return promotions, nil
}
