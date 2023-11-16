package main

import (
	"log/slog"
	"os"

	"github.com/alexapps/url-shortener/internal/config"
	"github.com/alexapps/url-shortener/internal/lib/sl"
	"github.com/alexapps/url-shortener/internal/storage/sqlite"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// read config
	cfg := config.MustLoad()

	// init logger: slog
	log := setupLogger(cfg.Env)

	log.Info("starting url shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// storage: sqllite
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	// stub
	_ = storage

	// TODO: router: chi, "chi render"

	// TODO: run server
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
