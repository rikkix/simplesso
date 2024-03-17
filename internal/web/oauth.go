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

const (
	// GitHubOAuthURL is the URL for GitHub OAuth.
	GitHubOAuthURL = "https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s"
	// GitHubUserURL is the URL for GitHub user info.
	GitHubUserURL = "https://api.github.com/user"
)

// handleOauthGitHub handles the GitHub OAuth callback.
// GitHub OAuth ref: https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/authorizing-oauth-apps
func (w *Web) handleOauthGitHub(c *fiber.Ctx) error {
	// Get the redirect URL and the code from the query
	redirect := c.Query("redirect", "/")
	code := c.Query("code")

	// Generate the URL for requesting the access token
	ghurl := fmt.Sprintf(GitHubOAuthURL, 
		w.config.GitHub.ClientID, w.config.GitHub.ClientSecret, code)

	// Request the access token
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

	// Decode the access token
	var ghtoken struct {
		AccessToken string `json:"access_token"`
	}
	err = jsoniter.NewDecoder(resp.Body).Decode(&ghtoken)
	if err != nil {
		w.logger.Error("Error decoding response: %s", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Get the current user info
	req, err = http.NewRequest(http.MethodGet, GitHubUserURL, nil)
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

	// Decode the username of the current user
	var ghuser struct {
		Login string `json:"login"`
	}
	err = jsoniter.NewDecoder(resp.Body).Decode(&ghuser)
	if err != nil {
		w.logger.Error("Error decoding response: %s", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Find the user by the GitHub username
	user := w.config.FindUserByGitHub(strings.ToLower(ghuser.Login))
	if user == nil {
		return c.Redirect(url.AddQueries("/login", map[string]string{
			"msg": "User not found.",
		}))
	}

	// Generate a new JWT token
	exp := time.Now().Add(24 * time.Hour)
	token, err := w.ssnParser.SSOAuther().GenerateToken(user.Name, w.config.Server.SsoHost, exp)
	if err != nil {
		w.logger.Error("Error generating token: %s", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Set the token in the cookie
	c.Cookie(
		&fiber.Cookie{
			Name:     "auth_token",
			Value:    token,
			Expires:  exp,
			SameSite: "Lax",
			HTTPOnly: true,
			Secure: true,
		})

	return c.Redirect(url.AddQueries("/", map[string]string{
		"redirect": redirect,
	}))
}