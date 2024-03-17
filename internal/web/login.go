package web

import (
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
	"github.com/rikkix/simplesso/internal/web/loginreq"
	"github.com/rikkix/simplesso/utils/url"
)

// handleIndex handles the index page.
// If the user is not authenticated, it will redirect to the login page.
// If the user is authenticated, it will redirect to the redirect URL if provided.
func (w *Web) handleIndex(c *fiber.Ctx) error {
	// Get the session and check if the user is authenticated
	ssn := w.session(c)
	authed := ssn.Authorized

	// Get the redirect URL if provided
	redirect := c.Query("redirect", "")

	// If the user is not authenticated, redirect to the login page
	if !authed {
		return c.Redirect(
			url.AddQueries("/login", map[string]string{
				"msg":      "Please login to continue.",
				"redirect": redirect,
			}),
		)
	}

	// Extract the host of the redirect url
	re_host := url.ExtractHost(redirect)
	// Check the corresponding service for the host
	service := w.config.FindService(re_host)

	// If redirect is not empty, validate the redirect URL
	if redirect != "" {
		if re_host == "" && redirect[0] != '/' {
			redirect = "/"
		}
		if re_host != "" && service == nil {
			redirect = "/"
		}
	}

	// If the redirect URL is empty, return a welcome message
	if redirect == "" {
		return c.SendString("Welcome, " + ssn.Sub + "!")
	}

	// If the redirect URL is index page, redirect to the index page
	if redirect[0] == '/' {
		return c.Redirect(redirect)
	}

	// If the user is not allowed, return 403 Forbidden
	if !service.IsUserAllowed(ssn.Sub) {
		return c.SendStatus(fiber.StatusForbidden)
	}

	// Redirect to the service page
	callback := "https://" + service.Host + "/auth-cgi/callback"
	callback = url.AddQueries(callback, map[string]string{
		"redirect": redirect,
	})
	// Generate token for the service
	token, err := w.ssnParser.CGIAuther().GenerateToken(ssn.Sub, re_host, ssn.Exp)
	if err != nil {
		w.logger.Error("Error generating token: %s", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	// Redirect to the service callback page
	return c.Render("callback_form", fiber.Map{
		"callback": callback,
		"token":    token,
	})
}

// handleLogin handles the login page.
func (w *Web) handleLogin(c *fiber.Ctx) error {
	// Get the message and redirect URL if provided
	msg := c.Query("msg")
	redirect := c.Query("redirect")

	return c.Render("login", fiber.Map{
		"msg":              msg,
		"redirect":         redirect,
		"github_client_id": w.config.GitHub.ClientID,
	})
}

// handleLoginPost handles the login form submission.
func (w *Web) handleLoginPost(c *fiber.Ctx) error {
	// Get the redirect URL if provided
	redirect := c.Query("redirect")
	// Get the username from the form
	username := c.FormValue("username")
	username = strings.ToLower(strings.TrimSpace(username))
	// Find the user by the username
	user := w.config.FindUser(username)
	// Get the remember me checkbox value
	remember := c.FormValue("remember") == "on"

	// Set the validity duration of the token
	dur := 6 * time.Hour
	if remember {
		dur = 30 * 24 * time.Hour
	}

	// Generate a new login request ID
	// If the user is found, generate a new login request ID and send the code to the user.
	// Otherwise, generate a dummy login request ID
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

// handleVerify handles the verification form request.
func (w *Web) handleVerifyPost(c *fiber.Ctx) error {
	// Get the redirect URL if provided
	redirect := c.Query("redirect", "/")
	// Get the request ID and code from the form
	reqid := c.FormValue("reqid")
	code := c.FormValue("code")

	// Check if the request ID and code are valid
	succ, loginreq := w.loginreqdb.Finish(reqid, code)

	// If the request ID and code are not valid, redirect to the login page with an error message
	if !succ {
		return c.Redirect(url.AddQueries("/login", map[string]string{
			"msg":      "Error verifying code. Please try again.",
			"redirect": redirect,
		}))
	}

	// Generate a new JWT token
	exp := time.Now().Add(loginreq.Dur)
	token, err := w.ssnParser.SSOAuther().GenerateToken(loginreq.Username, w.config.Server.SsoHost, exp)
	if err != nil {
		w.logger.Error("Error generating token: %s", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Set Cookie
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

// handleLogout handles the logout request.
func (w *Web) handleLogout(c *fiber.Ctx) error {
	c.ClearCookie("auth_token")
	return c.Redirect("/")
}
