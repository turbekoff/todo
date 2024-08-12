package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/turbekoff/todo/internal/config"
	"github.com/turbekoff/todo/internal/delivery/rest"
	"github.com/turbekoff/todo/internal/repository/mongo"
	server "github.com/turbekoff/todo/internal/server/http"
	"github.com/turbekoff/todo/internal/service"
	"github.com/turbekoff/todo/pkg/hash"
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

	database := client.Database(cfg.Mongo.Database)
	userRepository := mongo.NewUserRepositry(database)
	taskRepository := mongo.NewTaskRepository(database)
	sessionRepository := mongo.NewSessionRepository(database)
	hasher := hash.NewArgon2idHasher(cfg.PasswordPepper)

	userService := service.NewUserService(hasher, userRepository, taskRepository, sessionRepository)
	taskService := service.NewTaskService(userRepository, taskRepository)
	sessionService := service.NewSessionService(hasher, userRepository, sessionRepository, &cfg.JWT)

	router := rest.NewRouter(log, userService, taskService, sessionService)
	server := server.New(router, &cfg.HTTP)

	go func() {
		if err := server.Run(); !errors.Is(err, http.ErrServerClosed) {
			log.Error(fmt.Sprintf("failed to start HTTP server on %s:%d", cfg.HTTP.Host, cfg.HTTP.Port), Error(err))
		}
	}()

	log.Info(fmt.Sprintf("Starting HTTP server on %s:%d", cfg.HTTP.Host, cfg.HTTP.Port))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGTSTP)
	<-quit

	ctx, close := context.WithTimeout(context.Background(), 10*time.Second)
	defer close()

	if err := server.Shutdown(context.Background()); err != nil {
		log.Error("failed to shutdown HTTP server", Error(err))
	}

	if err := client.Disconnect(ctx); err != nil {
		log.Error("failed to close database connection", Error(err))
	}
}
