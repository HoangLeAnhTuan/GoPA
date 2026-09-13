package repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopa/internal/core/domain"
)

func TestJournalRepo_FullTextSearch_GINIndexAndRank(t *testing.T) {
	db, userID := setupTestDB(t)
	repo := NewJournalRepository(db)
	ctx := context.Background()

	now := time.Now().UTC()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	// Entry 1: "Golang" in TITLE (Weighted 'A')
	journal1 := domain.Journal{
		ID:            uuid.New(),
		UserID:        userID,
		Title:         "Mastering Golang Concurrency Patterns",
		Content:       "Deep dive into goroutines, channels, and worker pools for high-throughput microservices.",
		Tags:          []string{"golang", "concurrency"},
		PublishedDate: today,
		CreatedAt:     now,
		UpdatedAt:     now,
	}

	// Entry 2: "Golang" only in BODY (Weighted 'B')
	journal2 := domain.Journal{
		ID:            uuid.New(),
		UserID:        userID,
		Title:         "Weekly Work Reflection and Tasks",
		Content:       "Today I spent three hours refactoring legacy code into idiomatic Golang.",
		Tags:          []string{"work", "weekly"},
		PublishedDate: today.AddDate(0, 0, -1),
		CreatedAt:     now.Add(-24 * time.Hour),
		UpdatedAt:     now.Add(-24 * time.Hour),
	}

	// Entry 3: Completely unrelated content
	journal3 := domain.Journal{
		ID:            uuid.New(),
		UserID:        userID,
		Title:         "Japanese Grammar N2 Preparation",
		Content:       "Studying advanced kanji readings and nuanced sentence particles.",
		Tags:          []string{"japanese", "jlpt"},
		PublishedDate: today.AddDate(0, 0, -2),
		CreatedAt:     now.Add(-48 * time.Hour),
		UpdatedAt:     now.Add(-48 * time.Hour),
	}

	for _, j := range []domain.Journal{journal1, journal2, journal3} {
		_, err := repo.Create(ctx, j)
		require.NoError(t, err, "insert journal")
	}

	t.Cleanup(func() {
		cleanupCtx := context.Background()
		_, _ = db.ExecContext(cleanupCtx, `DELETE FROM journals WHERE user_id = $1`, userID)
	})

	// 1. Search for "Golang"
	results, err := repo.Search(ctx, userID, "Golang")
	require.NoError(t, err)
	require.Len(t, results, 2, "Expected exactly 2 journals containing 'Golang'")

	// Because Title is weighted 'A' and Content is weighted 'B' in journals_tsvector_trigger,
	// journal1 (title match) must be ranked higher than journal2 (body match only).
	assert.Equal(t, journal1.ID, results[0].ID, "Title match should have higher ts_rank than body match")
	assert.Equal(t, journal2.ID, results[1].ID)

	// 2. Search for "Japanese"
	jpResults, err := repo.Search(ctx, userID, "Japanese")
	require.NoError(t, err)
	require.Len(t, jpResults, 1)
	assert.Equal(t, journal3.ID, jpResults[0].ID)

	// 3. Search for non-existent term
	emptyResults, err := repo.Search(ctx, userID, "NonExistentKeywordXYZ")
	require.NoError(t, err)
	assert.Empty(t, emptyResults)
}
