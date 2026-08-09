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

	readiness := services.NewReadinessService(db, redisClient, brokerConnection)
	userRepository := repository.NewUserRepository(db)
	taskRepository := repository.NewTaskRepository(db)
	vocabularyRepository := repository.NewVocabularyRepository(db)
	vehicleRepository := repository.NewVehicleRepository(db)
	networkRepository := repository.NewNetworkRepository(db)
	financeRepository := repository.NewFinanceRepository(db)
	journalRepository := repository.NewJournalRepository(db)
	pomodoroHistoryRepository := repository.NewPomodoroHistoryRepository(db)
	refreshSessions := cache.NewRefreshSessionStore(redisClient)
	rateLimiter := cache.NewRateLimiter(redisClient)
	pomodoroStore := cache.NewPomodoroStore(redisClient)
	tokens := utils.NewTokenManager(cfg.JWTIssuer, cfg.JWTAccessSecret, cfg.AccessTokenTTL)
	authService := services.NewAuthService(userRepository, refreshSessions, tokens, cfg.BcryptCost, cfg.RefreshTokenTTL)
	authHandler := httpadapter.NewAuthHandler(authService, cfg.RefreshTokenTTL, cfg.AppEnv == "production")
	moduleHandler := httpadapter.NewModuleHandler(services.NewTaskService(taskRepository), services.NewVocabularyService(vocabularyRepository), services.NewAssetService(vehicleRepository, networkRepository), services.NewFinanceService(financeRepository), services.NewJournalService(journalRepository), services.NewPomodoroService(pomodoroStore, pomodoroHistoryRepository))
	router := httpadapter.NewRouter(log, cfg.WebOrigin, httpadapter.NewHealthHandler(readiness), authHandler, moduleHandler, tokens, rateLimiter)
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
