package web

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/rikkix/simplesso/utils/url"
)

func (w *Web) handleAuthCGIAuth(c *fiber.Ctx) error {
	service := w.config.FindService(c.Hostname())
	if service == nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if service.IsBypass(c.Get("X-Forwarded-Uri")) {
		return c.SendStatus(fiber.StatusOK)
	}

	ssn := w.session(c)
	if ssn.Authorized && service.IsUserAllowed(ssn.Sub) {
		return c.SendStatus(fiber.StatusOK)
	}

	token := c.Get("X-Authorization")
	if token != "" {
		if service.IsTokenAllowed(strings.TrimSpace(token)) {
			return c.SendStatus(fiber.StatusOK)
		}
	}
	return c.SendStatus(fiber.StatusUnauthorized)
}	

func (w *Web) handleAuthCGILogin(c *fiber.Ctx) error {
	forwardedURI := c.Get("X-Forwarded-Uri")
	rurl := "https://" + c.Hostname() + forwardedURI

	ssourl := "https://" + w.config.Server.SsoHost + "/"

	ssourl = url.AddQueries(ssourl, map[string]string {
		"redirect": rurl,
	})

	return c.Redirect(ssourl)
}

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