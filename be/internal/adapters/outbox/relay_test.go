package outbox

import (
	"context"
	"encoding/json"
	"log/slog"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopa/internal/constants"
)

func TestOutboxRecord_Parsing(t *testing.T) {
	eventID := uuid.New()
	occurredAt := time.Now().UTC()
	payloadJSON := []byte(`{"task_id":"123","user_id":"456"}`)

	record := OutboxRecord{
		ID:         eventID,
		EventType:  "pomodoro.finished",
		RoutingKey: "pomodoro.finished",
		Payload:    payloadJSON,
		OccurredAt: occurredAt,
		CreatedAt:  occurredAt,
	}

	assert.Equal(t, eventID, record.ID)
	assert.Equal(t, "pomodoro.finished", record.EventType)
	assert.Nil(t, record.PublishedAt)

	var parsed map[string]string
	err := json.Unmarshal(record.Payload, &parsed)
	require.NoError(t, err)
	assert.Equal(t, "123", parsed["task_id"])
}

func TestOutboxRelay_Configuration(t *testing.T) {
	relay := NewOutboxRelay(nil, nil, slog.Default())

	assert.Equal(t, constants.ExchangeEvents, relay.exchange)
	assert.Equal(t, 50, relay.batchSize)
	assert.Equal(t, 500*time.Millisecond, relay.pollInterval)

	relay.SetBatchSize(100)
	assert.Equal(t, 100, relay.batchSize)

	relay.SetPollInterval(2 * time.Second)
	assert.Equal(t, 2*time.Second, relay.pollInterval)
}

func TestOutboxRelay_ChannelClosedHandling(t *testing.T) {
	relay := NewOutboxRelay(nil, nil, slog.Default())

	ch, err := relay.getChannel()
	assert.Error(t, err)
	assert.Nil(t, ch)
	assert.Contains(t, err.Error(), "rabbitmq connection is closed or nil")
}

func TestOutboxRelay_ContextCancellation(t *testing.T) {
	relay := NewOutboxRelay(nil, nil, slog.Default())
	relay.SetPollInterval(10 * time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())

	done := make(chan struct{})
	go func() {
		relay.Start(ctx)
		close(done)
	}()

	// Wait briefly and cancel
	time.Sleep(25 * time.Millisecond)
	cancel()

	select {
	case <-done:
		// Clean exit
	case <-time.After(1 * time.Second):
		t.Fatal("relay did not stop upon context cancellation")
	}
}
