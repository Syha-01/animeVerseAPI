package main

import (
	"flag"
	"log/slog"
	"os"
)

const version = "1.0.0"

type configuration struct {
	port int
	env  string
}

type application struct {
	config configuration
	logger *slog.Logger
}

func main() {
	cfg := loadConfig()
	logger := setupLogger()
	app := &application{
		config: cfg,
		logger: logger,
	}

	err := app.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func loadConfig() configuration {
	var cfg configuration
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment(development|staging|production)")
	flag.Parse()
	return cfg
}

func setupLogger() *slog.Logger {
	var logger *slog.Logger
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	return logger
}
