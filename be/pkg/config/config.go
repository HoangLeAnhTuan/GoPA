package config

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"gopa/pkg/constants"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv          string
	HTTPAddr        string
	PostgresDSN     string
	RedisAddr       string
	RedisPassword   string
	RabbitMQURL     string
	JWTIssuer       string
	JWTAccessSecret string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	LogLevel        string
	WebOrigin       string
	ShutdownTimeout time.Duration
	BcryptCost      int
}

func Load() (Config, error) {
	_ = godotenv.Load()

	accessTTL, err := requiredDuration("ACCESS_TOKEN_TTL")
	if err != nil {
		return Config{}, err
	}
	refreshTTL, err := requiredDuration("REFRESH_TOKEN_TTL")
	if err != nil {
		return Config{}, err
	}
	shutdownTimeout, err := requiredDuration("SHUTDOWN_TIMEOUT")
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		AppEnv:          strings.ToLower(required("APP_ENV")),
		HTTPAddr:        required("HTTP_ADDR"),
		PostgresDSN:     required("POSTGRES_DSN"),
		RedisAddr:       required("REDIS_ADDR"),
		RedisPassword:   os.Getenv("REDIS_PASSWORD"),
		RabbitMQURL:     required("RABBITMQ_URL"),
		JWTIssuer:       required("JWT_ISSUER"),
		JWTAccessSecret: required("JWT_ACCESS_SECRET"),
		AccessTokenTTL:  accessTTL,
		RefreshTokenTTL: refreshTTL,
		LogLevel:        strings.ToLower(required("LOG_LEVEL")),
		WebOrigin:       required("WEB_ORIGIN"),
		ShutdownTimeout: shutdownTimeout,
		BcryptCost:      optionalInt("BCRYPT_COST", constants.DefaultBcryptCost),
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (c Config) Validate() error {
	if c.AppEnv != constants.Development && c.AppEnv != constants.Test && c.AppEnv != constants.Production {
		return fmt.Errorf("APP_ENV must be development, test, or production")
	}
	if c.AccessTokenTTL <= 0 || c.RefreshTokenTTL <= 0 || c.ShutdownTimeout <= 0 {
		return fmt.Errorf("configured durations must be positive")
	}
	if c.BcryptCost < 10 || c.BcryptCost > 14 {
		return fmt.Errorf("BCRYPT_COST must be between 10 and 14")
	}
	if _, err := url.ParseRequestURI(c.PostgresDSN); err != nil {
		return fmt.Errorf("POSTGRES_DSN must be a valid URL")
	}
	if _, err := url.ParseRequestURI(c.RabbitMQURL); err != nil {
		return fmt.Errorf("RABBITMQ_URL must be a valid URL")
	}
	if c.AppEnv == constants.Production {
		if c.JWTAccessSecret == constants.ExampleJWTSecret || len(c.JWTAccessSecret) < 32 {
			return fmt.Errorf("JWT_ACCESS_SECRET must be a non-example value with at least 32 characters in production")
		}
		if c.WebOrigin == "*" {
			return fmt.Errorf("WEB_ORIGIN cannot be a wildcard in production")
		}
	}
	return nil
}

func required(name string) string {
	return strings.TrimSpace(os.Getenv(name))
}

func requiredDuration(name string) (time.Duration, error) {
	value := required(name)
	if value == "" {
		return 0, fmt.Errorf("%s is required", name)
	}
	duration, err := time.ParseDuration(value)
	if err != nil || duration <= 0 {
		return 0, fmt.Errorf("%s must be a positive Go duration", name)
	}
	return duration, nil
}

func optionalInt(name string, fallback int) int {
	value := required(name)
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return parsed
}
