package outbox

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

type OutboxRecord struct {
	ID          uuid.UUID       `db:"id"`
	EventType   string          `db:"event_type"`
	RoutingKey  string          `db:"routing_key"`
	Payload     json.RawMessage `db:"payload"`
	OccurredAt  time.Time       `db:"occurred_at"`
	PublishedAt *time.Time      `db:"published_at"`
	CreatedAt   time.Time       `db:"created_at"`
}

type OutboxRelay struct {
	db           *sqlx.DB
	conn         *amqp091.Connection
	exchange     string
	pollInterval time.Duration
	batchSize    int
	logger       *slog.Logger
	ch           *amqp091.Channel
	mu           sync.Mutex
}

func NewOutboxRelay(
	db *sqlx.DB,
	conn *amqp091.Connection,
	logger *slog.Logger,
) *OutboxRelay {
	return &OutboxRelay{
		db:           db,
		conn:         conn,
		exchange:     constants.ExchangeEvents,
		pollInterval: 500 * time.Millisecond,
		batchSize:    50,
		logger:       logger,
	}
}

func (r *OutboxRelay) SetPollInterval(interval time.Duration) {
	r.pollInterval = interval
}

func (r *OutboxRelay) SetBatchSize(size int) {
	if size > 0 {
		r.batchSize = size
	}
}

func (r *OutboxRelay) getChannel() (*amqp091.Channel, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if r.ch != nil && !r.ch.IsClosed() {
		return r.ch, nil
	}

	if r.conn == nil || r.conn.IsClosed() {
		return nil, fmt.Errorf("rabbitmq connection is closed or nil")
	}

	ch, err := r.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("open rabbitmq channel: %w", err)
	}

	if err := ch.ExchangeDeclare(
		r.exchange,
		"topic",
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,
	); err != nil {
		_ = ch.Close()
		return nil, fmt.Errorf("declare exchange %s: %w", r.exchange, err)
	}

	r.ch = ch
	return r.ch, nil
}

func (r *OutboxRelay) ProcessBatch(ctx context.Context) (int, error) {
	if r.db == nil {
		return 0, fmt.Errorf("database connection is nil")
	}

	var records []OutboxRecord
	query := `
		SELECT id, event_type, routing_key, payload, occurred_at, published_at, created_at
		FROM outbox_events
		WHERE published_at IS NULL
		ORDER BY created_at ASC
		LIMIT $1
	`
	if err := r.db.SelectContext(ctx, &records, query, r.batchSize); err != nil {
		return 0, fmt.Errorf("query pending outbox records: %w", err)
	}

	if len(records) == 0 {
		return 0, nil
	}

	ch, err := r.getChannel()
	if err != nil {
		return 0, fmt.Errorf("get rabbitmq channel: %w", err)
	}

	publishedIDs := make([]uuid.UUID, 0, len(records))
	now := time.Now().UTC()

	for _, record := range records {
		routingKey := record.RoutingKey
		if routingKey == "" {
			routingKey = record.EventType
		}

		msg := amqp091.Publishing{
			DeliveryMode: amqp091.Persistent,
			ContentType:  "application/json",
			MessageId:    record.ID.String(),
			Type:         record.EventType,
			Timestamp:    record.OccurredAt.UTC(),
			Body:         record.Payload,
		}

		if err := ch.PublishWithContext(ctx, r.exchange, routingKey, false, false, msg); err != nil {
			if r.logger != nil {
				r.logger.Error("failed to publish outbox event", "id", record.ID, "type", record.EventType, "error", err)
			}
			break
		}

		publishedIDs = append(publishedIDs, record.ID)
	}

	if len(publishedIDs) > 0 {
		updateQuery, args, inErr := sqlx.In(
			`UPDATE outbox_events SET published_at = ? WHERE id IN (?) AND published_at IS NULL`,
			now,
			publishedIDs,
		)
		if inErr != nil {
			return len(publishedIDs), fmt.Errorf("build update query: %w", inErr)
		}
		updateQuery = r.db.Rebind(updateQuery)
		if _, execErr := r.db.ExecContext(ctx, updateQuery, args...); execErr != nil {
			return len(publishedIDs), fmt.Errorf("mark events published: %w", execErr)
		}
		if r.logger != nil {
			r.logger.Debug("published outbox batch", "count", len(publishedIDs))
		}
	}

	return len(publishedIDs), nil
}

func (r *OutboxRelay) Start(ctx context.Context) {
	if r.logger != nil {
		r.logger.Info("outbox relay daemon started", "exchange", r.exchange, "poll_interval", r.pollInterval)
	}

	ticker := time.NewTicker(r.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			r.mu.Lock()
			if r.ch != nil && !r.ch.IsClosed() {
				_ = r.ch.Close()
			}
			r.mu.Unlock()
			if r.logger != nil {
				r.logger.Info("outbox relay daemon stopped")
			}
			return

		case <-ticker.C:
			for {
				count, err := r.ProcessBatch(ctx)
				if err != nil {
					if r.logger != nil {
						r.logger.Warn("outbox relay process error", "error", err)
					}
					break
				}
				if count < r.batchSize || ctx.Err() != nil {
					break
				}
			}
		}
	}
}
