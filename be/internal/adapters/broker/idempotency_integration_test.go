package broker

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopa/pkg/database"
)

func getTestDSN() string {
	if dsn := os.Getenv("POSTGRES_DSN"); dsn != "" {
		return dsn
	}
	return "postgres://gopa:gopa@localhost:5432/gopa?sslmode=disable"
}

func TestIdempotencyStore_PostgreSQLIntegration(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := database.Open(ctx, getTestDSN())
	if err != nil {
		t.Skipf("Skipping postgres idempotency test: %v", err)
		return
	}
	defer db.Close()

	store := NewIdempotencyStore(db)
	eventID := uuid.New()
	eventType := "learning.review.completed"

	t.Cleanup(func() {
		cleanupCtx := context.Background()
		_, _ = db.ExecContext(cleanupCtx, `DELETE FROM processed_events WHERE event_id = $1`, eventID)
	})

	// First attempt: should acquire successfully (true)
	acquired, err := store.TryAcquire(ctx, eventID, eventType)
	require.NoError(t, err)
	assert.True(t, acquired, "First attempt must acquire lock")

	// Duplicate attempt with same eventID: should be rejected (false)
	dupAcquired, err := store.TryAcquire(ctx, eventID, eventType)
	require.NoError(t, err)
	assert.False(t, dupAcquired, "Duplicate attempt must be rejected by idempotency constraint")

	// Release event lock
	err = store.Release(ctx, eventID)
	require.NoError(t, err)

	// After release: should be able to acquire again
	reacquired, err := store.TryAcquire(ctx, eventID, eventType)
	require.NoError(t, err)
	assert.True(t, reacquired, "Should be able to acquire lock after release")
}
