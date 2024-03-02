package main

import (
	flog "github.com/gofiber/fiber/v2/log"
	"github.com/rikkix/simplesso/internal/config"
	"github.com/rikkix/simplesso/internal/web"

	"github.com/qinains/fastergoding"
)

const (
	DefaultConfigPath = "config.toml"
)

func main() {
	fastergoding.Run("./cmd/simplesso") // hot reload

	log := flog.DefaultLogger()
	log.SetLevel(flog.LevelDebug)

	var err error

	// Load the configuration.
	config, err := config.FromFile(DefaultConfigPath)
	if err != nil {
		log.Fatalf("Error loading configuration: %s", err)
	}

	// Create a new web server.
	web := web.New(config, log, nil)

	// Register the routes.
	web.RegisterRoutes()

	// Start the web server.
	web.Start()
}