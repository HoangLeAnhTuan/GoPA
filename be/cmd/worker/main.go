package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gopa/internal/adapters/broker"
	"gopa/internal/adapters/cache"
	"gopa/pkg/config"
	"gopa/pkg/database"
	"gopa/pkg/logger"
)

func main() {
	if err := run(); err != nil {
		slog.Error("worker stopped", "error", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}
	log := logger.New(cfg.AppEnv, cfg.LogLevel)
	startupCtx, cancelStartup := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelStartup()

	db, err := database.Open(startupCtx, cfg.PostgresDSN)
	if err != nil {
		return err
	}
	defer db.Close()
	redisClient, err := cache.Open(startupCtx, cfg.RedisAddr, cfg.RedisPassword)
	if err != nil {
		return err
	}
	defer redisClient.Close()
	brokerConnection, err := broker.Open(cfg.RabbitMQURL)
	if err != nil {
		return err
	}
	defer brokerConnection.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	log.Info("worker started", "status", "idle", "detail", "event consumers are added with their first vertical slice")
	<-ctx.Done()
	log.Info("worker shutdown started")
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancelShutdown()
	select {
	case <-shutdownCtx.Done():
		return shutdownCtx.Err()
	default:
		log.Info("worker shutdown complete")
		return nil
	}
}
