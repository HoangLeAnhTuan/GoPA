package config

import (
	"testing"

	"gopa/pkg/constants"
)

func TestConfigValidate_RejectsExampleSecretInProduction(t *testing.T) {
	cfg := Config{
		AppEnv:          constants.Production,
		PostgresDSN:     "postgres://gopa:gopa@localhost:5432/gopa?sslmode=disable",
		RedisAddr:       "localhost:6379",
		RabbitMQURL:     "amqp://gopa:gopa@localhost:5672/",
		JWTIssuer:       "gopa",
		JWTAccessSecret: constants.ExampleJWTSecret,
		AccessTokenTTL:  1,
		RefreshTokenTTL: 1,
		LogLevel:        "info",
		WebOrigin:       "https://gopa.example",
		ShutdownTimeout: 1,
		BcryptCost:      constants.DefaultBcryptCost,
	}

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation to reject the example production secret")
	}
}
