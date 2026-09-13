package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
	"gopa/internal/adapters/broker"
	"gopa/internal/adapters/cache"
	"gopa/internal/adapters/repository"
	"gopa/internal/constants"
	"gopa/internal/services"
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

	// ── Infrastructure ──────────────────────────────────────────────────────
	db, err := database.Open(startupCtx, cfg.PostgresDSN)
	if err != nil {
		return err
	}
	defer func(db *sqlx.DB) {
		if err := db.Close(); err != nil {
			slog.Error("failed to close database", "error", err)
		}
	}(db)

	redisClient, err := cache.Open(startupCtx, cfg.RedisAddr, cfg.RedisPassword)
	if err != nil {
		return err
	}
	defer func(r *redis.Client) {
		if err := r.Close(); err != nil {
			slog.Error("failed to close redis", "error", err)
		}
	}(redisClient)

	brokerConnection, err := broker.Open(cfg.RabbitMQURL)
	if err != nil {
		return err
	}
	defer func(conn *amqp091.Connection) {
		if err := conn.Close(); err != nil {
			slog.Error("failed to close rabbitmq", "error", err)
		}
	}(brokerConnection)

	// ── Repositories ────────────────────────────────────────────────────────
	pomodoroHistoryRepo := repository.NewPomodoroHistoryRepository(db)
	vocabularyRepo := repository.NewVocabularyRepository(db)
	budgetRepo := repository.NewBudgetRepository(db)

	// ── Consumer setup ──────────────────────────────────────────────────────
	consumer := broker.NewConsumer(brokerConnection, db, log)

	consumer.Register(
		constants.QueuePomodoro,
		services.HandlePomodoroFinished(pomodoroHistoryRepo, db, log),
	)
	consumer.Register(
		constants.QueueSRS,
		services.HandleReviewCompleted(vocabularyRepo, log),
	)
	consumer.Register(
		constants.QueueFinance,
		services.HandleBudgetAlert(budgetRepo, log),
	)

	// ── Signal context (blocks until SIGINT / SIGTERM) ────────────────────
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Info("worker starting", "queues", []string{
		constants.QueuePomodoro,
		constants.QueueSRS,
		constants.QueueFinance,
	})

	if err := consumer.Start(ctx); err != nil {
		return err
	}

	log.Info("worker started — consuming events")
	<-ctx.Done()

	log.Info("worker shutdown started — draining in-flight messages")
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancelShutdown()

	done := make(chan struct{})
	go func() {
		consumer.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Info("worker shutdown complete")
	case <-shutdownCtx.Done():
		log.Warn("worker shutdown timed out; some messages may not have been acknowledged")
	}

	return nil
}
