package lidl

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/mxschmitt/playwright-go"
	"github.com/rs/zerolog/log"
)

func (c *Client) login() error {
	var emailInputTestId = "input-email"
	var emailNextButtonTestId = "login-or-register-submit-button"
	var passwordInputTestId = "login-input-password"
	var passwordNextButtonTestId = "button-primary"

	if checkTokenValid(c.cfg.Token) {
		return nil
	}

	if c.cfg.Login == "" || c.cfg.Password == "" {
		return errors.New("login or password not set")
	}

	pw, err := playwright.Run()
	if err != nil {
		log.Error().Msg(err.Error())
		return errors.New("could not start Playwright")
	}
	defer pw.Stop()

	ctx, err := pw.Chromium.LaunchPersistentContext(c.cfg.StoragePath, playwright.BrowserTypeLaunchPersistentContextOptions{
		Headless: playwright.Bool(true),
		Channel:  playwright.String("chromium"),
	})
	if err != nil {
		log.Error().Msg(err.Error())
		return errors.New("could not launch browser")
	}
	defer ctx.Close()
	ctx.SetDefaultTimeout(c.cfg.Timeout)
	ctx.SetDefaultNavigationTimeout(c.cfg.Timeout)

	page, err := ctx.NewPage()
	if err != nil {
		log.Error().Msg(err.Error())
		return errors.New("could not create page")
	}

	if _, err = page.Goto(c.cfg.LoginURL); err != nil {
		log.Error().Msg(err.Error())
		return errors.New("could not goto login URL")
	}

	if strings.Contains(page.URL(), "accounts.lidl.com") {
		log.Debug().Msg("redirected to accounts.lidl.com, login required")
		err = page.GetByTestId(emailInputTestId).Fill(c.cfg.Login)
		if err != nil {
			log.Error().Msg(err.Error())
			return errors.New("could not fill email")
		}
		log.Debug().Msg("email input filled")

		if err = page.GetByTestId(emailNextButtonTestId).Click(); err != nil {
			log.Error().Msg(err.Error())
			return errors.New("could not press button:next with email")
		}
		log.Debug().Msg("email submitted")

		if err = page.GetByTestId(passwordInputTestId).Fill(c.cfg.Password); err != nil {
			log.Error().Msg(err.Error())
			return errors.New("could not fill password")
		}
		log.Debug().Msg("password input filled")

		if err = page.GetByTestId(passwordNextButtonTestId).Click(); err != nil {
			log.Error().Msg(err.Error())
			return errors.New("could not press button:next with password")
		}
		log.Debug().Msg("password submitted")

		// after login there user will be redirected several times and land back on the original login URL
		// we need to wait for the final URL to grab cookies
		if err := page.WaitForURL(c.cfg.LoginURL + "**"); err != nil {
			log.Error().Msg(err.Error())
			return errors.New("after login url did not match expected pattern: " + err.Error())
		}
		log.Info().Msg("login successful")
	} else {
		log.Debug().Msg("reusing previous session, no login required")
	}

	cookies, err := page.Context().Cookies()
	if err != nil {
		return errors.New("could not get cookies")
	}
	for _, cookie := range cookies {
		if cookie.Name == "authToken" {
			c.cfg.Token = cookie.Value
			log.Info().Msg("token retrieved from cookies")
			return nil
		}
	}
	return errors.New("authToken not found in cookies")
}

func checkTokenValid(token string) bool {
	claims := jwt.MapClaims{}
	_, _, err := jwt.NewParser().ParseUnverified(token, claims)
	if err != nil {
		return false
	}
	exp, err := claims.GetExpirationTime()
	if err != nil {
		return false
	}
	if time.Now().Unix() > int64(exp.Unix()) {
		return false
	}
	log.Debug().Msg("existing token is valid")
	return true
}
