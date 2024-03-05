package web

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rikkix/simplesso/utils/url"
)

func (w *Web) handleAuthCGIAuth(c *fiber.Ctx) error {
	ssn := w.session(c)
	if ssn.Authorized {
		return c.SendStatus(fiber.StatusOK)
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
	})

	if !url.SameHost(url.ExtractHost(redirect), c.Hostname()) {
		return c.Redirect("/", 303)
	}

	return c.Redirect(redirect, 303)
}