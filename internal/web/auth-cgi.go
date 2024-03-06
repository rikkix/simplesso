package web

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rikkix/simplesso/utils/url"
)

// handleAuthCGIAuth handles the auth request.
// If the user is authorized, it will return 200 OK.
// If the user is not authorized, it will return 401 Unauthorized.
func (w *Web) handleAuthCGIAuth(c *fiber.Ctx) error {
	// Get the service for the host
	service := w.config.FindService(c.Hostname())
	// If the service is not found, return 400 Bad Request
	if service == nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// CHeck whether the path is bypassed
	if service.IsBypass(c.Get("X-Forwarded-Uri")) {
		return c.SendStatus(fiber.StatusOK)
	}

	// Get the session and check if the user is authorized and allowed
	ssn := w.session(c)
	if ssn.Authorized && service.IsUserAllowed(ssn.Sub) {
		return c.SendStatus(fiber.StatusOK)
	}

	// Get and check the token from the header
	token := c.Get("X-Authorization")
	if token != "" {
		if service.IsTokenAllowed(strings.TrimSpace(token)) {
			return c.SendStatus(fiber.StatusOK)
		}
	}

	return c.SendStatus(fiber.StatusUnauthorized)
}	

// handleAuthCGILogin handles the login request.
// It will redirect to the SSO page for further actions.
func (w *Web) handleAuthCGILogin(c *fiber.Ctx) error {
	forwardedURI := c.Get("X-Forwarded-Uri")
	rurl := "https://" + c.Hostname() + forwardedURI

	ssourl := "https://" + w.config.Server.SsoHost + "/"

	ssourl = url.AddQueries(ssourl, map[string]string {
		"redirect": rurl,
	})

	return c.Redirect(ssourl)
}

// handleAuthCGICallbackPost handles the callback request from SSO site.
func (w *Web) handleAuthCGICallbackPost(c *fiber.Ctx) error {
	token := c.FormValue("token")
	redirect := c.Query("redirect")
	_,exp,err := w.ssnParser.CGIAuther().ValidateToken(token, c.Hostname())
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	c.Cookie(&fiber.Cookie{
		Name: "auth_token",
		Value: token,
		Expires: exp,
		SameSite: "Lax",
		HTTPOnly: true,
	})

	if !url.SameHost(url.ExtractHost(redirect), c.Hostname()) {
		return c.Redirect("/", 303)
	}

	return c.Redirect(redirect, 303)
}