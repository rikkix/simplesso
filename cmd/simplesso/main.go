package main

import (
	"fmt"
	"os"

	flog "github.com/gofiber/fiber/v2/log"
	"github.com/rikkix/simplesso/internal/config"
	"github.com/rikkix/simplesso/internal/web"
)

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Usage: simplesso <config file>")
		os.Exit(1)
	}
	configPath := args[1]

	log := flog.DefaultLogger()
	log.SetLevel(flog.LevelWarn)

	var err error

	// Load the configuration.
	config, err := config.FromFile(configPath)
	if err != nil {
		log.Fatalf("Error loading configuration: %s", err)
	}

	// Create a new web server.
	web := web.New(config, log, nil)

	// Start the web server.
	web.Start()
}