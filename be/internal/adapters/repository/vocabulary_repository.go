package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"gopa/internal/core/domain"
)

type VocabularyRepository struct{ db *sqlx.DB }

func NewVocabularyRepository(db *sqlx.DB) *VocabularyRepository { return &VocabularyRepository{db: db} }
func (r *VocabularyRepository) List(ctx context.Context, userID uuid.UUID, limit int) ([]domain.Vocabulary, error) {
	const query = `SELECT id, user_id, language, word, reading, meaning, example_sentence, current_box, next_review_at, created_at, updated_at FROM vocabularies WHERE user_id = $1 ORDER BY next_review_at, created_at DESC LIMIT $2`
	vocabularies := make([]domain.Vocabulary, 0)
	if err := r.db.SelectContext(ctx, &vocabularies, query, userID, limit); err != nil {
		return nil, fmt.Errorf("list vocabularies: %w", err)
	}
	return vocabularies, nil
}
func (r *VocabularyRepository) ListDue(ctx context.Context, userID uuid.UUID, limit int, now time.Time) ([]domain.Vocabulary, error) {
	const query = `SELECT id, user_id, language, word, reading, meaning, example_sentence, current_box, next_review_at, created_at, updated_at FROM vocabularies WHERE user_id = $1 AND next_review_at <= $2 ORDER BY next_review_at LIMIT $3`
	vocabularies := make([]domain.Vocabulary, 0)
	if err := r.db.SelectContext(ctx, &vocabularies, query, userID, now, limit); err != nil {
		return nil, fmt.Errorf("list due vocabularies: %w", err)
	}
	return vocabularies, nil
}
func (r *VocabularyRepository) Get(ctx context.Context, userID, id uuid.UUID) (domain.Vocabulary, error) {
	const query = `SELECT id, user_id, language, word, reading, meaning, example_sentence, current_box, next_review_at, created_at, updated_at FROM vocabularies WHERE id = $1 AND user_id = $2`
	var vocabulary domain.Vocabulary
	if err := r.db.GetContext(ctx, &vocabulary, query, id, userID); err != nil {
		return domain.Vocabulary{}, mapNotFound(err, "get vocabulary")
	}
	return vocabulary, nil
}
func (r *VocabularyRepository) Create(ctx context.Context, vocabulary domain.Vocabulary) (domain.Vocabulary, error) {
	const query = `INSERT INTO vocabularies (id, user_id, language, word, reading, meaning, example_sentence, current_box, next_review_at, created_at, updated_at) VALUES (:id, :user_id, :language, :word, :reading, :meaning, :example_sentence, :current_box, :next_review_at, :created_at, :updated_at)`
	if _, err := r.db.NamedExecContext(ctx, query, vocabulary); err != nil {
		return domain.Vocabulary{}, fmt.Errorf("create vocabulary: %w", err)
	}
	return vocabulary, nil
}
func (r *VocabularyRepository) Update(ctx context.Context, vocabulary domain.Vocabulary) (domain.Vocabulary, error) {
	const query = `UPDATE vocabularies SET language = :language, word = :word, reading = :reading, meaning = :meaning, example_sentence = :example_sentence, updated_at = :updated_at WHERE id = :id AND user_id = :user_id`
	result, err := r.db.NamedExecContext(ctx, query, vocabulary)
	if err != nil {
		return domain.Vocabulary{}, fmt.Errorf("update vocabulary: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return domain.Vocabulary{}, fmt.Errorf("count updated vocabulary: %w", err)
	}
	if rows == 0 {
		return domain.Vocabulary{}, domain.ErrNotFound
	}
	return vocabulary, nil
}
func (r *VocabularyRepository) Delete(ctx context.Context, userID, id uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM vocabularies WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("delete vocabulary: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("count deleted vocabulary: %w", err)
	}
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}
func (r *VocabularyRepository) RecordReview(ctx context.Context, userID, vocabularyID uuid.UUID, quality int16, reviewedAt time.Time, schedule domain.ReviewSchedule) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin vocabulary review: %w", err)
	}
	defer tx.Rollback()
	var exists bool
	if err := tx.GetContext(ctx, &exists, `SELECT EXISTS(SELECT 1 FROM vocabularies WHERE id = $1 AND user_id = $2)`, vocabularyID, userID); err != nil {
		return fmt.Errorf("verify vocabulary ownership: %w", err)
	}
	if !exists {
		return domain.ErrNotFound
	}
	eventID := uuid.New()
	reviewID := uuid.New()
	now := reviewedAt.UTC()
	if _, err := tx.ExecContext(ctx, `INSERT INTO vocabulary_reviews (id, event_id, vocabulary_id, user_id, quality, reviewed_at, created_at) VALUES ($1, $2, $3, $4, $5, $6, $6)`, reviewID, eventID, vocabularyID, userID, quality, now); err != nil {
		return fmt.Errorf("insert vocabulary review: %w", err)
	}
	if _, err := tx.ExecContext(ctx, `UPDATE vocabularies SET current_box = $1, next_review_at = $2, updated_at = $3 WHERE id = $4 AND user_id = $5`, schedule.Box, schedule.NextReviewAt.UTC(), now, vocabularyID, userID); err != nil {
		return fmt.Errorf("schedule vocabulary review: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit vocabulary review: %w", err)
	}
	return nil
}

var _ = sql.ErrNoRows
