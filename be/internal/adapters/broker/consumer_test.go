package broker

import (
	"context"
	"encoding/json"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopa/internal/constants"
)

func TestParseDelivery_EnvelopeFormat(t *testing.T) {
	eventID := uuid.New()
	userID := uuid.New()
	now := time.Now().UTC()

	envelope := EventEnvelope{
		EventID:       eventID,
		EventType:     "learning.review.completed",
		SchemaVersion: 1,
		OccurredAt:    now,
		UserID:        userID,
		Payload:       json.RawMessage(`{"quality":5}`),
	}

	body, err := json.Marshal(envelope)
	require.NoError(t, err)

	delivery := amqp091.Delivery{
		Body: body,
	}

	parsed, err := ParseDelivery(delivery)
	require.NoError(t, err)
	assert.Equal(t, eventID, parsed.EventID)
	assert.Equal(t, "learning.review.completed", parsed.EventType)
	assert.Equal(t, userID, parsed.UserID)
}

func TestParseDelivery_RawPayloadFormat(t *testing.T) {
	eventID := uuid.New()
	userID := uuid.New()
	now := time.Now().UTC()

	rawBody := []byte(`{"user_id":"` + userID.String() + `","amount":50000}`)

	delivery := amqp091.Delivery{
		MessageId:  eventID.String(),
		Type:       "finance.budget.alert",
		Timestamp:  now,
		RoutingKey: "finance.budget.alert",
		Body:       rawBody,
	}

	parsed, err := ParseDelivery(delivery)
	require.NoError(t, err)
	assert.Equal(t, eventID, parsed.EventID)
	assert.Equal(t, "finance.budget.alert", parsed.EventType)
	assert.Equal(t, userID, parsed.UserID)
	assert.Equal(t, rawBody, []byte(parsed.Payload))
}

func TestIdempotencyStore_NilDB(t *testing.T) {
	store := NewIdempotencyStore(nil)

	acquired, err := store.TryAcquire(context.Background(), uuid.New(), "test.event")
	assert.NoError(t, err)
	assert.True(t, acquired)

	err = store.Release(context.Background(), uuid.New())
	assert.NoError(t, err)
}

func TestConsumer_RegistrationAndConfiguration(t *testing.T) {
	consumer := NewConsumer(nil, nil, slog.Default())

	assert.Equal(t, 10, consumer.prefetch)

	consumer.SetPrefetch(20)
	assert.Equal(t, 20, consumer.prefetch)

	called := false
	consumer.Register(constants.QueueSRS, func(ctx context.Context, delivery amqp091.Delivery, event EventEnvelope) error {
		called = true
		return nil
	})

	assert.NotNil(t, consumer.handlers[constants.QueueSRS])
	err := consumer.handlers[constants.QueueSRS](context.Background(), amqp091.Delivery{}, EventEnvelope{})
	assert.NoError(t, err)
	assert.True(t, called)
}

func TestSetupTopology_ClosedConnection(t *testing.T) {
	err := SetupTopology(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "rabbitmq connection is closed or nil")
}
