package main

import (
	"log/slog"
	"net/http"
	"os"

	chi "github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/alexapps/url-shortener/internal/config"
	"github.com/alexapps/url-shortener/internal/http-server/handlers/redirect"
	"github.com/alexapps/url-shortener/internal/http-server/handlers/url/save"
	mwLogger "github.com/alexapps/url-shortener/internal/http-server/middleware/logger"
	"github.com/alexapps/url-shortener/internal/lib/logger/handlers/slogpretty"
	"github.com/alexapps/url-shortener/internal/lib/logger/sl"
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

	err = storage.DeleteURL("google")
	if err != nil {
		log.Error("failed to delete url", sl.Err(err))
		os.Exit(1)
	}

	// router: chi, "chi render"
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mwLogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Post("/url", save.New(log, storage))
	router.Get("/{alias}", redirect.New(log, storage))

	log.Info("starting server", slog.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}
	log.Error("server stopped")

	// router.Use(mi)

	// TODO: run server
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}
	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}
