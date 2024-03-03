package web

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/template/html/v2"
	"github.com/rikkix/simplesso/internal/config"
	"github.com/rikkix/simplesso/internal/web/loginreq"
	"github.com/rikkix/simplesso/internal/web/session"
	"github.com/rikkix/simplesso/internal/web/tg"
)


type Web struct {
	// Config is the configuration for the web server.
	config *config.Config
	// Logger is the logger for the web server.
	logger log.Logger
	// Server is the web server.
	server *fiber.App
	
	// ssnParser is the session parser for the web server.
	ssnParser *session.SessionParser
	// loginreqdb is the login request database for the web server.
	loginreqdb *loginreq.MemDB
	// tgbot is the telegram bot for the web server.
	tgbot *tg.TG

	// route_registered is a flag to indicate if the routes have been registered.
	routeRegistered bool
}

// New creates a new Web instance.
func New(c *config.Config, l log.Logger, a *fiber.App) *Web {
	l.Info("Creating new web server...")
	if a == nil {
		l.Info("Empty fiber app, creating new one...")
		a = fiber.New(fiber.Config{
			Views: html.New(c.Server.WebPath + "layouts", ".html"),
		})
	}
	lrq := loginreq.NewMemDB(30, 5 * time.Minute)
	tbot, err := tg.New(c.Server.TelegramToken, lrq, l)
	if err != nil {
		l.Fatal("Error creating telegram bot: %s", err)
	}
	return &Web{
		config: c,
		logger: l,
		server: a,
		ssnParser: session.NewSessionParser(&c.Server),
		loginreqdb: lrq,
		tgbot: tbot,
		routeRegistered: false,
	}
}

// RegisterRoutes registers the routes for the web server.
func (w *Web) RegisterRoutes() {
	if w.routeRegistered {
		return
	}
	w.logger.Info("Registering routes...")

	// Register the routes here.
	w.server.Get("/", w.handleIndex)
	w.server.Get("/login", w.handleLogin)
	w.server.Post("/login", w.handleLoginPost)
	w.server.Post("/verify", w.handleVerifyPost)
	w.server.Get("/logout", w.handleLogout)
}

// Start starts the web server.
func (w *Web) Start() {
	w.logger.Warn("Starting telegram bot...")
	go w.tgbot.StartPolling()

	w.logger.Warn("Starting web server...")
	err := w.server.Listen(w.config.Server.ListenAddress)
	if err != nil {
		w.logger.Fatal("Error starting web server: %s", err)
	}
}

func (w *Web) session(c *fiber.Ctx) *session.Session {
	return w.ssnParser.Parse(c)
}