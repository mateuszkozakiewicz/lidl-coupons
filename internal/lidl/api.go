package lidl

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

func (c *Client) request(path string, method string) (*http.Response, error) {
	path = strings.TrimPrefix(path, "/")
	client := &http.Client{}
	req, err := http.NewRequest(method, c.cfg.ApiURL+path, nil)
	if err != nil {
		return nil, err
	}
	if c.cfg.Token == "" {
		return nil, errors.New("authToken is not set, log in first")
	}
	req.Header.Set("Cookie", "authToken="+c.cfg.Token)
	log.Debug().Msgf("requesting %s %s", method, req.URL.String())
	return client.Do(req)
}

func (c *Client) getPromotions() ([]Promotion, error) {
	resp, err := c.request("promotionslist", "GET")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("getting unexpected status code: " + resp.Status)
	}
	var response struct {
		Sections []struct {
			Promotions []Promotion `json:"promotions"`
		} `json:"sections"`
	}
	if err = json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, err
	}
	var promotions []Promotion
	for _, section := range response.Sections {
		promotions = append(promotions, section.Promotions...)
	}
	log.Info().Msgf("promotions fetched successfully")
	return promotions, nil
}

func (c *Client) activatePromotion(p Promotion) error {
	d := time.Now().Format("2006-01-02T15:04:05Z")
	if p.EndValidityDate < d {
		return ErrPromotionExpired
	}
	if p.StartValidityDate > d {
		return ErrPromotionNotYetValid
	}
	resp, err := c.request("promotions/"+p.ID+"/activation", "POST")
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusAccepted {
		return errors.New("unexpected status code: " + resp.Status)
	}
	log.Info().Msgf("activated promotion id:%s description:%s", p.ID, p.Description)
	return nil
}
