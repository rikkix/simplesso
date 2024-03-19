package main

import (
	"flag"
	"fmt"
	"os"

	flog "github.com/gofiber/fiber/v2/log"
	"github.com/rikkix/simplesso/internal/config"
	"github.com/rikkix/simplesso/internal/web"
)

func main() {
	// Parse the command line flags.
	configPath := flag.String("config", "config.toml", "Path to the configuration file")
	templatePath := flag.String("templates", "templates", "Path to the templates directory")
	logLevel := flag.String("loglevel", "info", "Log level (trace, debug, info, warn, error, fatal, panic)")
	reloadTemplates := flag.Bool("reload", false, "Reload templates on change")
	listenAddr := flag.String("listen", "", "Address to listen on (overrides configuration if set)")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [options]\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	log := flog.DefaultLogger()

	var ll flog.Level
	switch *logLevel {
	case "trace":
		ll = flog.LevelTrace
	case "debug":
		ll = flog.LevelDebug
	case "info":
		ll = flog.LevelInfo
	case "warn":
		ll = flog.LevelWarn
	case "error":
		ll = flog.LevelError
	case "fatal":
		ll = flog.LevelFatal
	case "panic":
		ll = flog.LevelPanic
	default:
		log.Fatalf("Invalid log level: %s", *logLevel)
	}
	log.SetLevel(ll)

	var err error

	// Load the configuration.
	config, err := config.FromFile(*configPath)
	if err != nil {
		log.Fatalf("Error loading configuration: %s", err)
	}

	if *listenAddr != "" {
		config.Server.ListenAddress = *listenAddr
	}

	// Create a new web server.
	web := web.New(config, *templatePath , *reloadTemplates ,log, nil)

	// Start the web server.
	web.Start()
}