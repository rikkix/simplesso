package web

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/rikkix/simplesso/internal/web/loginreq"
	"github.com/rikkix/simplesso/utils/url"
)

func (w *Web) handleIndex(c *fiber.Ctx) error {
	ssn := w.session(c)
	authed := ssn.Authorized
	redirect := c.Query("redirect", "")

	if !authed {
		return c.Redirect(
			url.AddQueries("/login", map[string]string{
				"msg":      "Please login to continue.",
				"redirect": redirect,
			}),
		)
	}

	re_host := url.ExtractHost(redirect)
	service := w.config.FindService(re_host)
	if redirect != "" {
		if re_host == "" && redirect[0] != '/' {
			redirect = "/"
		}
		if re_host != "" && service == nil {
			redirect = "/"
		}
	}

	if redirect == "" {
		return c.SendString("Welcome, " + ssn.Sub + "!")
	}

	if redirect[0] == '/' {
		return c.Redirect(redirect)
	}

	callback := "https://" + service.Host + "/auth-cgi/callback"
	token, err := w.ssnParser.CGIAuther().GenerateToken(re_host, ssn.Sub, ssn.Exp)
	if err != nil {
		w.logger.Error("Error generating token: %s", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	return c.Render("callback_form", fiber.Map{
		"callback": callback,
		"token": token,
	})
	
}

func (w *Web) handleLogin(c *fiber.Ctx) error {
	msg := c.Query("msg")
	redirect := c.Query("redirect")

	return c.Render("login", fiber.Map{
		"msg":      msg,
		"redirect": redirect,
	})
}

func (w *Web) handleLoginPost(c *fiber.Ctx) error {
	redirect := c.Query("redirect")
	username := c.FormValue("username")
	remember := c.FormValue("remember") == "on"
	username = strings.ToLower(strings.TrimSpace(username))
	user := w.config.FindUser(username)

	dur := 6 * time.Hour
	if remember {
		dur = 30 * 24 * time.Hour
	}

	reqid := ""
	if user != nil {
		reqid = w.loginreqdb.NewReq(utils.CopyString(username), dur)
		go w.tgbot.SendConfirmaion(user.TelegramId, reqid)
	} else {
		reqid = loginreq.NewReqID()
	}

	return c.Render("login_verify", fiber.Map{
		"redirect": redirect,
		"reqid":    reqid,
		"username": username,
	})
}

func (w *Web) handleVerifyPost(c *fiber.Ctx) error {
	redirect := c.Query("redirect", "/")
	reqid := c.FormValue("reqid")
	code := c.FormValue("code")
	succ, loginreq := w.loginreqdb.Finish(reqid, code)
	if !succ {
		return c.Redirect(url.AddQueries("/login", map[string]string{
			"msg":      "Error verifying code. Please try again.",
			"redirect": redirect,
		}))
	}
	exp := time.Now().Add(loginreq.Dur)
	token, err := w.ssnParser.SSOAuther().GenerateToken(loginreq.Username, w.config.Server.SsoHost, exp)
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

	re_host := url.ExtractHost(redirect)
	if re_host != "" && w.config.FindService(re_host) == nil {
		redirect = "/"
	}

	return c.Redirect(redirect)
}

func (w *Web) handleLogout(c *fiber.Ctx) error {
	c.ClearCookie("auth_token")
	return c.Redirect("/")
}