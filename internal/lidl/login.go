package lidl

import (
	"errors"
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

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
	})
	if err != nil {
		log.Error().Msg(err.Error())
		return errors.New("could not launch browser")
	}

	ctx, err := browser.NewContext(playwright.BrowserNewContextOptions{
		IgnoreHttpsErrors: playwright.Bool(true),
	})
	if err != nil {
		log.Error().Msg(err.Error())
		return errors.New("could not create context")
	}

	page, err := ctx.NewPage()
	if err != nil {
		log.Error().Msg(err.Error())
		return errors.New("could not create page")
	}

	if _, err = page.Goto(c.cfg.LoginURL); err != nil {
		log.Error().Msg(err.Error())
		return errors.New("could not goto login URL")
	}
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
	log.Info().Msg("email submitted")

	if err = page.GetByTestId(passwordInputTestId).Fill(c.cfg.Password); err != nil {
		log.Error().Msg(err.Error())
		return errors.New("could not fill password")
	}
	log.Info().Msg("password input filled")

	if err = page.GetByTestId(passwordNextButtonTestId).Click(); err != nil {
		log.Error().Msg(err.Error())
		return errors.New("could not press button:next with password")
	}
	log.Info().Msg("password submitted")

	// after login there user will be redirected several times and land back on the original login URL
	// we need to wait for the final URL to grab cookies
	if err := page.WaitForURL(c.cfg.LoginURL + "**"); err != nil {
		log.Error().Msg(err.Error())
		return errors.New("after login url did not match expected pattern: " + err.Error())
	}

	cookies, err := page.Context().Cookies()
	if err != nil {
		return errors.New("could not get cookies")
	}
	for _, cookie := range cookies {
		if cookie.Name == "authToken" {
			c.cfg.Token = cookie.Value
			log.Info().Msg("login successful")
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
