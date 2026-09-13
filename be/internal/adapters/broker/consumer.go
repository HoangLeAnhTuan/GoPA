package broker

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rabbitmq/amqp091-go"
	"gopa/internal/constants"
)

type EventEnvelope struct {
	EventID       uuid.UUID       `json:"event_id"`
	EventType     string          `json:"event_type"`
	SchemaVersion int             `json:"schema_version"`
	OccurredAt    time.Time       `json:"occurred_at"`
	CorrelationID string          `json:"correlation_id,omitempty"`
	UserID        uuid.UUID       `json:"user_id"`
	Payload       json.RawMessage `json:"payload"`
}

func ParseDelivery(d amqp091.Delivery) (EventEnvelope, error) {
	var env EventEnvelope
	if err := json.Unmarshal(d.Body, &env); err == nil && env.EventID != uuid.Nil {
		if env.EventType == "" {
			env.EventType = d.Type
		}
		if env.OccurredAt.IsZero() {
			env.OccurredAt = d.Timestamp
		}
		if len(env.Payload) == 0 {
			env.Payload = d.Body
		}
		return env, nil
	}

	// Fallback when delivery body contains raw JSON payload
	var rawObj map[string]any
	_ = json.Unmarshal(d.Body, &rawObj)

	eventID, err := uuid.Parse(d.MessageId)
	if err != nil {
		if idStr, ok := rawObj["event_id"].(string); ok {
			eventID, _ = uuid.Parse(idStr)
		}
		if eventID == uuid.Nil {
			if idStr, ok := rawObj["id"].(string); ok {
				eventID, _ = uuid.Parse(idStr)
			}
		}
	}

	var userID uuid.UUID
	if idStr, ok := rawObj["user_id"].(string); ok {
		userID, _ = uuid.Parse(idStr)
	}

	eventType := d.Type
	if eventType == "" {
		if et, ok := rawObj["event_type"].(string); ok {
			eventType = et
		} else {
			eventType = d.RoutingKey
		}
	}

	occurredAt := d.Timestamp
	if occurredAt.IsZero() {
		occurredAt = time.Now().UTC()
	}

	return EventEnvelope{
		EventID:       eventID,
		EventType:     eventType,
		SchemaVersion: 1,
		OccurredAt:    occurredAt,
		UserID:        userID,
		Payload:       d.Body,
	}, nil
}

type IdempotencyStore struct {
	db *sqlx.DB
}

func NewIdempotencyStore(db *sqlx.DB) *IdempotencyStore {
	return &IdempotencyStore{db: db}
}

func (s *IdempotencyStore) TryAcquire(ctx context.Context, eventID uuid.UUID, eventType string) (bool, error) {
	if s.db == nil {
		return true, nil
	}
	res, err := s.db.ExecContext(
		ctx,
		`INSERT INTO processed_events (event_id, event_type, processed_at) VALUES ($1, $2, $3) ON CONFLICT (event_id) DO NOTHING`,
		eventID,
		eventType,
		time.Now().UTC(),
	)
	if err != nil {
		return false, fmt.Errorf("record processed event: %w", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("check rows affected: %w", err)
	}
	return rows > 0, nil
}

func (s *IdempotencyStore) Release(ctx context.Context, eventID uuid.UUID) error {
	if s.db == nil {
		return nil
	}
	_, err := s.db.ExecContext(ctx, `DELETE FROM processed_events WHERE event_id = $1`, eventID)
	return err
}

type HandlerFunc func(ctx context.Context, delivery amqp091.Delivery, event EventEnvelope) error

type QueueBinding struct {
	QueueName  string
	RoutingKey string
}

var DefaultTopology = []QueueBinding{
	{QueueName: constants.QueueSRS, RoutingKey: "learning.#"},
	{QueueName: constants.QueuePomodoro, RoutingKey: "pomodoro.#"},
	{QueueName: constants.QueueFinance, RoutingKey: "finance.#"},
	{QueueName: constants.QueueJournal, RoutingKey: "journal.#"},
}

func SetupTopology(conn *amqp091.Connection) error {
	if conn == nil || conn.IsClosed() {
		return fmt.Errorf("rabbitmq connection is closed or nil")
	}
	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("open topology channel: %w", err)
	}
	defer func() {
		_ = ch.Close()
	}()

	// 1. Declare Topic Exchange: gopa.events
	if err := ch.ExchangeDeclare(
		constants.ExchangeEvents,
		"topic",
		true,  // durable
		false, // auto-delete
		false, // internal
		false, // no-wait
		nil,
	); err != nil {
		return fmt.Errorf("declare exchange %s: %w", constants.ExchangeEvents, err)
	}

	// 2. Declare Dead Letter Exchange: gopa.events.dlx
	if err := ch.ExchangeDeclare(
		constants.ExchangeDLX,
		"topic",
		true,  // durable
		false, // auto-delete
		false, // internal
		false, // no-wait
		nil,
	); err != nil {
		return fmt.Errorf("declare exchange %s: %w", constants.ExchangeDLX, err)
	}

	// 3. Declare DLQ and bind to DLX
	if _, err := ch.QueueDeclare(
		constants.QueueDLQ,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	); err != nil {
		return fmt.Errorf("declare queue %s: %w", constants.QueueDLQ, err)
	}
	if err := ch.QueueBind(
		constants.QueueDLQ,
		"#",
		constants.ExchangeDLX,
		false,
		nil,
	); err != nil {
		return fmt.Errorf("bind dlq: %w", err)
	}

	// 4. Declare Domain Queues with DLX configuration and Bind to gopa.events
	dlxArgs := amqp091.Table{
		"x-dead-letter-exchange": constants.ExchangeDLX,
	}

	for _, b := range DefaultTopology {
		if _, err := ch.QueueDeclare(
			b.QueueName,
			true,  // durable
			false, // auto-delete
			false, // exclusive
			false, // no-wait
			dlxArgs,
		); err != nil {
			return fmt.Errorf("declare queue %s: %w", b.QueueName, err)
		}

		if err := ch.QueueBind(
			b.QueueName,
			b.RoutingKey,
			constants.ExchangeEvents,
			false,
			nil,
		); err != nil {
			return fmt.Errorf("bind queue %s to %s: %w", b.QueueName, b.RoutingKey, err)
		}
	}

	return nil
}

type Consumer struct {
	conn        *amqp091.Connection
	idempotency *IdempotencyStore
	logger      *slog.Logger
	handlers    map[string]HandlerFunc
	prefetch    int
	mu          sync.RWMutex
	wg          sync.WaitGroup
}

func NewConsumer(
	conn *amqp091.Connection,
	db *sqlx.DB,
	logger *slog.Logger,
) *Consumer {
	return &Consumer{
		conn:        conn,
		idempotency: NewIdempotencyStore(db),
		logger:      logger,
		handlers:    make(map[string]HandlerFunc),
		prefetch:    10,
	}
}

func (c *Consumer) SetPrefetch(prefetch int) {
	if prefetch > 0 {
		c.prefetch = prefetch
	}
}

func (c *Consumer) Register(queueName string, handler HandlerFunc) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.handlers[queueName] = handler
}

func (c *Consumer) Start(ctx context.Context) error {
	if err := SetupTopology(c.conn); err != nil {
		return fmt.Errorf("setup topology: %w", err)
	}

	c.mu.RLock()
	handlers := make(map[string]HandlerFunc, len(c.handlers))
	for k, v := range c.handlers {
		handlers[k] = v
	}
	c.mu.RUnlock()

	for queueName, handler := range handlers {
		ch, err := c.conn.Channel()
		if err != nil {
			return fmt.Errorf("open channel for %s: %w", queueName, err)
		}

		if err := ch.Qos(c.prefetch, 0, false); err != nil {
			_ = ch.Close()
			return fmt.Errorf("set qos for %s: %w", queueName, err)
		}

		deliveries, err := ch.ConsumeWithContext(
			ctx,
			queueName,
			"",    // consumer tag
			false, // manual ack
			false, // exclusive
			false, // no-local
			false, // no-wait
			nil,
		)
		if err != nil {
			_ = ch.Close()
			return fmt.Errorf("consume %s: %w", queueName, err)
		}

		c.wg.Add(1)
		go func(q string, h HandlerFunc, channel *amqp091.Channel, msgs <-chan amqp091.Delivery) {
			defer c.wg.Done()
			defer func() {
				_ = channel.Close()
			}()
			c.consumeLoop(ctx, q, h, msgs)
		}(queueName, handler, ch, deliveries)
	}

	return nil
}

func (c *Consumer) consumeLoop(ctx context.Context, queueName string, handler HandlerFunc, deliveries <-chan amqp091.Delivery) {
	if c.logger != nil {
		c.logger.Info("consumer listening on queue", "queue", queueName)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case d, ok := <-deliveries:
			if !ok {
				return
			}
			c.handleDelivery(ctx, queueName, handler, d)
		}
	}
}

func (c *Consumer) handleDelivery(ctx context.Context, queueName string, handler HandlerFunc, d amqp091.Delivery) {
	event, err := ParseDelivery(d)
	if err != nil {
		if c.logger != nil {
			c.logger.Error("malformed event delivery, sending to dlq", "queue", queueName, "error", err)
		}
		_ = d.Nack(false, false)
		return
	}

	if event.EventID != uuid.Nil {
		acquired, acquireErr := c.idempotency.TryAcquire(ctx, event.EventID, event.EventType)
		if acquireErr != nil {
			if c.logger != nil {
				c.logger.Error("idempotency check error", "event_id", event.EventID, "error", acquireErr)
			}
			_ = d.Nack(false, true)
			return
		}
		if !acquired {
			if c.logger != nil {
				c.logger.Debug("event already processed, skipping", "event_id", event.EventID, "type", event.EventType)
			}
			_ = d.Ack(false)
			return
		}
	}

	if err := handler(ctx, d, event); err != nil {
		if c.logger != nil {
			c.logger.Error("event handler error", "queue", queueName, "event_id", event.EventID, "type", event.EventType, "error", err)
		}
		if event.EventID != uuid.Nil {
			_ = c.idempotency.Release(ctx, event.EventID)
		}
		_ = d.Nack(false, false)
		return
	}

	_ = d.Ack(false)
}

func (c *Consumer) Wait() {
	c.wg.Wait()
}
