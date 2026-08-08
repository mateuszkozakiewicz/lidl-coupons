package lidl

import "errors"

type Promotion struct {
	ID                    string   `json:"id"`
	Image                 Image    `json:"image"`
	Offer                 string   `json:"offer"`
	Title                 string   `json:"title"`
	Description           string   `json:"description"`
	Channel               string   `json:"channel"`
	Channels              []string `json:"channels"`
	IsActivated           bool     `json:"isActivated"`
	IsProcessing          bool     `json:"isProcessing"`
	IsHappyHour           bool     `json:"isHappyHour"`
	IsTest                bool     `json:"isTest"`
	Special               Special  `json:"special,omitempty"`
	StartValidityDate     string   `json:"startValidityDate"`
	EndValidityDate       string   `json:"endValidityDate"`
	DiscountScope         int      `json:"discountScope"`
	NavigationURL         string   `json:"navigationUrl"`
	RedemptionChannels    []string `json:"redemptionChannels"`
	VisualizationChannels []string `json:"visualizationChannels"`
}

type Image struct {
	URL     string `json:"url"`
	AltText string `json:"altText"`
}

type Special struct {
	Tag       string `json:"tag"`
	Color     string `json:"color"`
	FontColor string `json:"fontColor"`
}

var ErrPromotionExpired = errors.New("promotion is expired")
var ErrPromotionNotYetValid = errors.New("promotion is not yet valid")
