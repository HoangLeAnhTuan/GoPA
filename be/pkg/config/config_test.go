package config

import (
	"strings"
	"testing"
	"time"

	"gopa/pkg/constants"
)

func TestConfigValidate_RejectsExampleSecretInProduction(t *testing.T) {
	cfg := validConfig()
	cfg.AppEnv = constants.Production
	cfg.JWTAccessSecret = constants.ExampleJWTSecret
	cfg.PostgresDSN = "postgres://gopa:gopa@db.example/gopa?sslmode=require"

	if err := cfg.Validate(); err == nil || !strings.Contains(err.Error(), "JWT_ACCESS_SECRET") {
		t.Fatalf("expected validation to reject the example production secret, got %v", err)
	}
}

func TestConfigValidate_RejectsMissingRequiredValues(t *testing.T) {
	tests := []struct {
		name      string
		fieldName string
		mutate    func(*Config)
	}{
		{name: "http address", fieldName: "HTTP_ADDR", mutate: func(cfg *Config) { cfg.HTTPAddr = "" }},
		{name: "redis address", fieldName: "REDIS_ADDR", mutate: func(cfg *Config) { cfg.RedisAddr = "" }},
		{name: "JWT issuer", fieldName: "JWT_ISSUER", mutate: func(cfg *Config) { cfg.JWTIssuer = "" }},
		{name: "JWT secret", fieldName: "JWT_ACCESS_SECRET", mutate: func(cfg *Config) { cfg.JWTAccessSecret = "" }},
		{name: "web origin", fieldName: "WEB_ORIGIN", mutate: func(cfg *Config) { cfg.WebOrigin = "" }},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			cfg := validConfig()
			test.mutate(&cfg)
			err := cfg.Validate()
			if err == nil || !strings.Contains(err.Error(), test.fieldName) {
				t.Fatalf("expected %s validation error, got %v", test.fieldName, err)
			}
		})
	}
}

func TestConfigValidate_RejectsUnsafeProductionConnections(t *testing.T) {
	cfg := validConfig()
	cfg.AppEnv = constants.Production
	cfg.PostgresDSN = "postgres://gopa:gopa@db.example/gopa?sslmode=disable"

	if err := cfg.Validate(); err == nil || !strings.Contains(err.Error(), "TLS") {
		t.Fatalf("expected production database TLS validation error, got %v", err)
	}
}

func TestConfigValidate_RejectsOriginWithPath(t *testing.T) {
	cfg := validConfig()
	cfg.WebOrigin = "https://gopa.example/app"

	if err := cfg.Validate(); err == nil || !strings.Contains(err.Error(), "WEB_ORIGIN") {
		t.Fatalf("expected web origin validation error, got %v", err)
	}
}

func validConfig() Config {
	return Config{
		AppEnv:          constants.Production,
		HTTPAddr:        ":8081",
		PostgresDSN:     "postgres://gopa:gopa@db.example/gopa?sslmode=require",
		RedisAddr:       "localhost:6379",
		RabbitMQURL:     "amqp://gopa:gopa@localhost:5672/",
		JWTIssuer:       "gopa",
		JWTAccessSecret: "a-production-secret-with-at-least-32-characters",
		AccessTokenTTL:  time.Minute,
		RefreshTokenTTL: time.Hour,
		LogLevel:        "info",
		WebOrigin:       "https://gopa.example",
		ShutdownTimeout: time.Second,
		BcryptCost:      constants.DefaultBcryptCost,
	}
}
