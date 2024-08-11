package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/turbekoff/todo/internal/config"
	"github.com/turbekoff/todo/internal/repository/mongo"
	"golang.org/x/exp/slog"
)

func setupLogger(mode config.Mode) *slog.Logger {
	var log *slog.Logger
	mode.RunAt(config.M_LOC, func() {
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	})
	mode.RunAt(config.M_DEV, func() {
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	})
	mode.RunAt(config.M_PROD, func() {
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	})
	mode.RunAt(config.M_NULL, func() {
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	})

	return log.With(slog.String("mode", mode.String()))
}

func Error(err error) slog.Attr {
	return slog.Attr{
		Key:   "error",
		Value: slog.StringValue(err.Error()),
	}
}

func main() {
	cfg, err := config.Load()
	log := setupLogger(cfg.DebugMode)

	if err != nil {
		log.Error("failed to load config", Error(err))
		return
	}

	client, err := mongo.NewConnection(&cfg.Mongo)
	if err != nil {
		log.Error("failed to connect database", Error(err))
		return
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGTSTP)
	<-quit

	ctx, close := context.WithTimeout(context.Background(), 10*time.Second)
	defer close()

	if err := client.Disconnect(ctx); err != nil {
		log.Error("failed to close database connection", Error(err))
	}
}
