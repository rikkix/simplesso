package web

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/json-iterator/go"
	"github.com/rikkix/simplesso/utils/url"
)

func (w *Web) handleOauthGitHub(c *fiber.Ctx) error {
	redirect := c.Query("redirect", "/")
	code := c.Query("code")
	ghurl := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s", 
		w.config.GitHub.ClientID, w.config.GitHub.ClientSecret, code)

	req, err := http.NewRequest(http.MethodPost, ghurl, nil)
	if err != nil {
		w.logger.Error("Error creating request: %s", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	req.Header.Add("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		w.logger.Error("Error making request: %s", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	defer resp.Body.Close()

	var ghtoken struct {
		AccessToken string `json:"access_token"`
	}

	err = jsoniter.NewDecoder(resp.Body).Decode(&ghtoken)
	if err != nil {
		w.logger.Error("Error decoding response: %s", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	req, err = http.NewRequest(http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		w.logger.Error("Error creating request: %s", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	req.Header.Add("Authorization", "Bearer "+ghtoken.AccessToken)
	req.Header.Add("Accept", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		w.logger.Error("Error making request: %s", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	defer resp.Body.Close()

	var ghuser struct {
		Login string `json:"login"`
	}

	err = jsoniter.NewDecoder(resp.Body).Decode(&ghuser)
	if err != nil {
		w.logger.Error("Error decoding response: %s", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	user := w.config.FindUserByGitHub(strings.ToLower(ghuser.Login))
	if user == nil {
		return c.Redirect(url.AddQueries("/login", map[string]string{
			"msg": "User not found.",
		}))
	}

	exp := time.Now().Add(24 * time.Hour)

	token, err := w.ssnParser.SSOAuther().GenerateToken(user.Name, w.config.Server.SsoHost, exp)
	if err != nil {
		w.logger.Error("Error generating token: %s", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.Cookie(
		&fiber.Cookie{
			Name:     "auth_token",
			Value:    token,
			Expires:  exp,
			SameSite: "Lax",
			HTTPOnly: true,
		})

	return c.Redirect(url.AddQueries("/", map[string]string{
		"redirect": redirect,
	}))
}