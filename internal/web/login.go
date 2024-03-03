package web

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/rikkix/simplesso/internal/web/loginreq"
	"github.com/rikkix/simplesso/utils/url"
)

func (w *Web) handleIndex(c *fiber.Ctx) error {
	ssn := w.session(c)
	redirect := c.Query("redirect", "/")

	if !ssn.Authorized {
		args := map[string]string{
			"msg": "Please login to continue.",
			"redirect": redirect,
		}

		return c.Redirect(
			url.AddQueries("/login", args),
		)
	}

	return c.SendString("Welcome, "+ssn.Sub+"!")
}

func (w *Web) handleLogin(c *fiber.Ctx) error {
	msg := c.Query("msg")
	redirect := c.Query("redirect")

	return c.Render("login", fiber.Map{
		"msg": msg,
		"redirect": redirect,
	})
}

func (w *Web) handleLoginPost(c *fiber.Ctx) error {
	redirect := c.Query("redirect")
	username := c.FormValue("username")
	username = strings.ToLower(strings.TrimSpace(username))
	user := w.config.FindUser(username)

	reqid := ""
	if user != nil {
		reqid = w.loginreqdb.NewReq(utils.CopyString(username), 20)
		go w.tgbot.SendConfirmaion(user.TelegramId, reqid)
	} else {
		reqid = loginreq.NewReqID()
	}

	return c.Render("login_verify", fiber.Map{
		"redirect": redirect,
		"reqid": reqid,
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
			"msg": "Error verifying code. Please try again.",
		}))
	}
	exp := time.Now().Add(time.Second * time.Duration(loginreq.Dur))
	fmt.Println("loginreq.Username: ", loginreq.Username)
	token, err := w.ssnParser.SSOAuther().GenerateToken(loginreq.Username, w.config.Server.SsoHost, exp)
	if err != nil {
		w.logger.Error("Error generating token: %s", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	c.Cookie(
		&fiber.Cookie{
			Name: "auth_token",
			Value: token,
			Expires: exp,
			SameSite: "Lax",
			HTTPOnly: true,
			Secure: w.config.Server.CookieSecure,
		})

	return c.Redirect(redirect)
}
