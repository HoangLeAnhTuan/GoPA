package services

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

type ReadinessService struct {
	db     *sqlx.DB
	redis  *redis.Client
	broker *amqp091.Connection
}

func NewReadinessService(db *sqlx.DB, redisClient *redis.Client, brokerConnection *amqp091.Connection) *ReadinessService {
	return &ReadinessService{db: db, redis: redisClient, broker: brokerConnection}
}

func (s *ReadinessService) Check(ctx context.Context) error {
	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("postgres unavailable: %w", err)
	}
	if err := s.redis.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("redis unavailable: %w", err)
	}
	if s.broker.IsClosed() {
		return fmt.Errorf("rabbitmq connection is closed")
	}
	return nil
}
