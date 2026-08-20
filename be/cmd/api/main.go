package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gopa/internal/adapters/broker"
	"gopa/internal/adapters/cache"
	httpadapter "gopa/internal/adapters/handler/http"
	"gopa/internal/adapters/repository"
	"gopa/internal/services"
	"gopa/pkg/config"
	"gopa/pkg/database"
	"gopa/pkg/logger"
	"gopa/pkg/utils"

	"github.com/jmoiron/sqlx"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

func main() {
	if err := run(); err != nil {
		slog.Error("api stopped", "error", err)
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
	defer func(db *sqlx.DB) {
		err := db.Close()
		if err != nil {
			slog.Error("failed to close database connection", "error", err)
		}
	}(db)
	redisClient, err := cache.Open(startupCtx, cfg.RedisAddr, cfg.RedisPassword)
	if err != nil {
		return err
	}
	defer func(redisClient *redis.Client) {
		err := redisClient.Close()
		if err != nil {
			slog.Error("failed to close redis connection", "error", err)
		}
	}(redisClient)
	brokerConnection, err := broker.Open(cfg.RabbitMQURL)
	if err != nil {
		return err
	}
	defer func(brokerConnection *amqp091.Connection) {
		err := brokerConnection.Close()
		if err != nil {
			slog.Error("failed to close rabbitmq connection", "error", err)
		}
	}(brokerConnection)

	readiness := services.NewReadinessService(db, redisClient, brokerConnection)
	userRepository := repository.NewUserRepository(db)
	taskRepository := repository.NewTaskRepository(db)
	vocabularyRepository := repository.NewVocabularyRepository(db)
	financeRepository := repository.NewFinanceRepository(db)
	budgetRepository := repository.NewBudgetRepository(db)
	savingsGoalRepository := repository.NewSavingsGoalRepository(db)
	journalRepository := repository.NewJournalRepository(db)
	pomodoroHistoryRepository := repository.NewPomodoroHistoryRepository(db)
	refreshSessions := cache.NewRefreshSessionStore(redisClient)
	rateLimiter := cache.NewRateLimiter(redisClient)
	pomodoroStore := cache.NewPomodoroStore(redisClient)
	tokens := utils.NewTokenManager(cfg.JWTIssuer, cfg.JWTAccessSecret, cfg.AccessTokenTTL)
	authService := services.NewAuthService(userRepository, refreshSessions, tokens, cfg.BcryptCost, cfg.RefreshTokenTTL)
	authHandler := httpadapter.NewAuthHandler(authService, cfg.RefreshTokenTTL, cfg.AppEnv == "production")
	tasksHandler := httpadapter.NewTasksHandler(services.NewTaskService(taskRepository))
	linguisticsHandler := httpadapter.NewLinguisticsHandler(services.NewVocabularyService(vocabularyRepository))
	financeHandler := httpadapter.NewFinanceHandler(services.NewFinanceService(financeRepository, budgetRepository, savingsGoalRepository))
	journalHandler := httpadapter.NewJournalHandler(services.NewJournalService(journalRepository))
	pomodoroHandler := httpadapter.NewPomodoroHandler(services.NewPomodoroService(pomodoroStore, pomodoroHistoryRepository))

	router := httpadapter.NewRouter(
		log,
		cfg.WebOrigin,
		httpadapter.NewHealthHandler(readiness),
		authHandler,
		tasksHandler,
		linguisticsHandler,
		financeHandler,
		journalHandler,
		pomodoroHandler,
		tokens,
		rateLimiter,
	)
	server := &http.Server{Addr: cfg.HTTPAddr, Handler: router, ReadHeaderTimeout: 5 * time.Second}

	shutdownSignal, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	serverErrors := make(chan error, 1)
	go func() { serverErrors <- server.ListenAndServe() }()
	log.Info("api started", "address", cfg.HTTPAddr)

	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	case <-shutdownSignal.Done():
		log.Info("api shutdown started")
		shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
		defer cancelShutdown()
		if err := server.Shutdown(shutdownCtx); err != nil {
			return err
		}
		log.Info("api shutdown complete")
		return nil
	}
}
