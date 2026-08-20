package repository

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"gopa/internal/core/domain"
	"gopa/internal/core/ports"
	"gopa/pkg/database"
)

const vocabColumns = `id, user_id, language, word, reading, meaning, example_sentence, example_translation, tags, difficulty_level, audio_url, image_url, ease_factor, interval_days, repetition_count, current_box, next_review_at, mastery_score, created_at, updated_at`

type vocabularyRow struct {
	ID                 uuid.UUID                 `db:"id"`
	UserID             uuid.UUID                 `db:"user_id"`
	Language           domain.VocabularyLanguage `db:"language"`
	Word               string                    `db:"word"`
	Reading            *string                   `db:"reading"`
	Meaning            string                    `db:"meaning"`
	ExampleSentence    *string                   `db:"example_sentence"`
	ExampleTranslation *string                   `db:"example_translation"`
	Tags               pq.StringArray            `db:"tags"`
	DifficultyLevel    string                    `db:"difficulty_level"`
	AudioURL           *string                   `db:"audio_url"`
	ImageURL           *string                   `db:"image_url"`
	EaseFactor         float64                   `db:"ease_factor"`
	IntervalDays       int                       `db:"interval_days"`
	RepetitionCount    int                       `db:"repetition_count"`
	CurrentBox         int16                     `db:"current_box"`
	NextReviewAt       time.Time                 `db:"next_review_at"`
	MasteryScore       int16                     `db:"mastery_score"`
	CreatedAt          time.Time                 `db:"created_at"`
	UpdatedAt          time.Time                 `db:"updated_at"`
}

func (r vocabularyRow) domain() domain.Vocabulary {
	tags := make([]string, len(r.Tags))
	copy(tags, r.Tags)
	return domain.Vocabulary{
		ID: r.ID, UserID: r.UserID, Language: r.Language, Word: r.Word, Reading: r.Reading,
		Meaning: r.Meaning, ExampleSentence: r.ExampleSentence, ExampleTranslation: r.ExampleTranslation,
		Tags: tags, DifficultyLevel: r.DifficultyLevel, AudioURL: r.AudioURL, ImageURL: r.ImageURL,
		EaseFactor: r.EaseFactor, IntervalDays: r.IntervalDays, RepetitionCount: r.RepetitionCount,
		CurrentBox: r.CurrentBox, NextReviewAt: r.NextReviewAt, MasteryScore: r.MasteryScore,
		CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt,
	}
}

type VocabularyRepository struct{ db *sqlx.DB }

func NewVocabularyRepository(db *sqlx.DB) *VocabularyRepository { return &VocabularyRepository{db: db} }

func (r *VocabularyRepository) List(ctx context.Context, userID uuid.UUID, limit int) ([]domain.Vocabulary, error) {
	rows := make([]vocabularyRow, 0)
	if err := r.db.SelectContext(ctx, &rows, `SELECT `+vocabColumns+` FROM vocabularies WHERE user_id = $1 ORDER BY next_review_at, created_at DESC LIMIT $2`, userID, limit); err != nil {
		return nil, fmt.Errorf("list vocabularies: %w", err)
	}
	return vocabRowsToDomain(rows), nil
}

func (r *VocabularyRepository) ListDue(ctx context.Context, userID uuid.UUID, limit int, now time.Time) ([]domain.Vocabulary, error) {
	rows := make([]vocabularyRow, 0)
	if err := r.db.SelectContext(ctx, &rows, `SELECT `+vocabColumns+` FROM vocabularies WHERE user_id = $1 AND next_review_at <= $2 ORDER BY next_review_at LIMIT $3`, userID, now.UTC(), limit); err != nil {
		return nil, fmt.Errorf("list due vocabularies: %w", err)
	}
	return vocabRowsToDomain(rows), nil
}

func (r *VocabularyRepository) Get(ctx context.Context, userID, id uuid.UUID) (domain.Vocabulary, error) {
	var row vocabularyRow
	if err := r.db.GetContext(ctx, &row, `SELECT `+vocabColumns+` FROM vocabularies WHERE id = $1 AND user_id = $2`, id, userID); err != nil {
		return domain.Vocabulary{}, mapNotFound(err, "get vocabulary")
	}
	return row.domain(), nil
}

func (r *VocabularyRepository) Create(ctx context.Context, vocabulary domain.Vocabulary) (domain.Vocabulary, error) {
	tags := pq.StringArray(vocabulary.Tags)
	if _, err := r.db.NamedExecContext(ctx,
		`INSERT INTO vocabularies (`+vocabColumns+`) VALUES (:id, :user_id, :language, :word, :reading, :meaning, :example_sentence, :example_translation, :tags, :difficulty_level, :audio_url, :image_url, :ease_factor, :interval_days, :repetition_count, :current_box, :next_review_at, :mastery_score, :created_at, :updated_at)`,
		vocabToParams(vocabulary, tags)); err != nil {
		return domain.Vocabulary{}, fmt.Errorf("create vocabulary: %w", err)
	}
	return vocabulary, nil
}

// BulkCreate inserts multiple vocabulary items atomically in a single transaction.
func (r *VocabularyRepository) BulkCreate(ctx context.Context, userID uuid.UUID, items []domain.Vocabulary) ([]domain.Vocabulary, error) {
	if len(items) == 0 {
		return nil, nil
	}
	return items, database.WithTx(ctx, r.db, func(tx *sqlx.Tx) error {
		const query = `INSERT INTO vocabularies (` + vocabColumns + `) VALUES (:id, :user_id, :language, :word, :reading, :meaning, :example_sentence, :example_translation, :tags, :difficulty_level, :audio_url, :image_url, :ease_factor, :interval_days, :repetition_count, :current_box, :next_review_at, :mastery_score, :created_at, :updated_at) ON CONFLICT (id) DO NOTHING`
		for i := range items {
			items[i].UserID = userID
			tags := pq.StringArray(items[i].Tags)
			if _, err := tx.NamedExecContext(ctx, query, vocabToParams(items[i], tags)); err != nil {
				return fmt.Errorf("bulk create vocabulary[%d]: %w", i, err)
			}
		}
		return nil
	})
}

func (r *VocabularyRepository) Update(ctx context.Context, vocabulary domain.Vocabulary) (domain.Vocabulary, error) {
	tags := pq.StringArray(vocabulary.Tags)
	result, err := r.db.NamedExecContext(ctx,
		`UPDATE vocabularies SET language = :language, word = :word, reading = :reading, meaning = :meaning, example_sentence = :example_sentence, example_translation = :example_translation, tags = :tags, difficulty_level = :difficulty_level, audio_url = :audio_url, image_url = :image_url, updated_at = :updated_at WHERE id = :id AND user_id = :user_id`,
		vocabToParams(vocabulary, tags))
	if err != nil {
		return domain.Vocabulary{}, fmt.Errorf("update vocabulary: %w", err)
	}
	if err := requireAffected(result, "update vocabulary"); err != nil {
		return domain.Vocabulary{}, err
	}
	return vocabulary, nil
}

func (r *VocabularyRepository) Delete(ctx context.Context, userID, id uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM vocabularies WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("delete vocabulary: %w", err)
	}
	return requireAffected(result, "delete vocabulary")
}

// RecordReview atomically inserts the review event and updates the vocabulary SM-2 state.
func (r *VocabularyRepository) RecordReview(ctx context.Context, review domain.VocabularyReview, vocab domain.Vocabulary) error {
	return database.WithTx(ctx, r.db, func(tx *sqlx.Tx) error {
		// Insert the review record
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO vocabulary_reviews (id, event_id, vocabulary_id, user_id, quality, study_mode, response_time_ms, was_correct, reviewed_at, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $9)`,
			review.ID, review.EventID, review.VocabularyID, review.UserID, review.Quality,
			review.StudyMode, review.ResponseTimeMs, review.WasCorrect, review.ReviewedAt.UTC()); err != nil {
			return fmt.Errorf("insert vocabulary review: %w", err)
		}
		// Update SM-2 state on the vocabulary
		if _, err := tx.ExecContext(ctx,
			`UPDATE vocabularies SET current_box = $1, ease_factor = $2, interval_days = $3, repetition_count = $4, mastery_score = $5, next_review_at = $6, updated_at = $7 WHERE id = $8 AND user_id = $9`,
			vocab.CurrentBox, vocab.EaseFactor, vocab.IntervalDays, vocab.RepetitionCount,
			vocab.MasteryScore, vocab.NextReviewAt.UTC(), review.ReviewedAt.UTC(),
			vocab.ID, vocab.UserID); err != nil {
			return fmt.Errorf("update vocabulary sm2 state: %w", err)
		}
		// Emit outbox event for async processing
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO outbox_events (id, event_type, routing_key, payload, occurred_at, created_at) VALUES ($1, 'learning.review.completed', 'learning.review.completed', $2, $3, $3)`,
			review.EventID, fmt.Sprintf(`{"vocabulary_id":"%s","user_id":"%s","quality":%d}`, vocab.ID, vocab.UserID, review.Quality),
			review.ReviewedAt.UTC()); err != nil {
			return fmt.Errorf("insert review outbox event: %w", err)
		}
		return nil
	})
}

const learningSessionColumns = `id, user_id, language, session_type, items_reviewed, items_correct, duration_seconds, started_at, ended_at`

func (r *VocabularyRepository) CreateLearningSession(ctx context.Context, session domain.LearningSession) (domain.LearningSession, error) {
	if _, err := r.db.NamedExecContext(ctx,
		`INSERT INTO learning_sessions (`+learningSessionColumns+`) VALUES (:id, :user_id, :language, :session_type, :items_reviewed, :items_correct, :duration_seconds, :started_at, :ended_at)`,
		session); err != nil {
		return domain.LearningSession{}, fmt.Errorf("create learning session: %w", err)
	}
	return session, nil
}

func (r *VocabularyRepository) UpdateLearningSession(ctx context.Context, session domain.LearningSession) (domain.LearningSession, error) {
	result, err := r.db.NamedExecContext(ctx,
		`UPDATE learning_sessions SET items_reviewed = :items_reviewed, items_correct = :items_correct, duration_seconds = :duration_seconds, ended_at = :ended_at WHERE id = :id AND user_id = :user_id`,
		session)
	if err != nil {
		return domain.LearningSession{}, fmt.Errorf("update learning session: %w", err)
	}
	if err := requireAffected(result, "update learning session"); err != nil {
		return domain.LearningSession{}, err
	}
	return session, nil
}

// GetStats returns aggregate statistics for a user's vocabulary learning progress.
func (r *VocabularyRepository) GetStats(ctx context.Context, userID uuid.UUID) (domain.VocabularyStats, error) {
	var totals struct {
		Total    int `db:"total"`
		Mastered int `db:"mastered"`
		Due      int `db:"due"`
	}
	if err := r.db.GetContext(ctx, &totals,
		`SELECT COUNT(*) AS total, SUM(CASE WHEN mastery_score >= 80 THEN 1 ELSE 0 END) AS mastered, SUM(CASE WHEN next_review_at <= now() THEN 1 ELSE 0 END) AS due FROM vocabularies WHERE user_id = $1`,
		userID); err != nil {
		return domain.VocabularyStats{}, fmt.Errorf("get vocabulary totals: %w", err)
	}

	boxRows := make([]struct {
		Box   int16 `db:"current_box"`
		Count int   `db:"cnt"`
	}, 0)
	if err := r.db.SelectContext(ctx, &boxRows,
		`SELECT current_box, COUNT(*) AS cnt FROM vocabularies WHERE user_id = $1 GROUP BY current_box`, userID); err != nil {
		return domain.VocabularyStats{}, fmt.Errorf("get box distribution: %w", err)
	}
	boxDist := make(map[int16]int, 5)
	for _, br := range boxRows {
		boxDist[br.Box] = br.Count
	}

	// Last 30-day heatmap using review dates
	heatRows := make([]struct {
		Day   string `db:"day"`
		Count int    `db:"cnt"`
	}, 0)
	if err := r.db.SelectContext(ctx, &heatRows,
		`SELECT to_char(reviewed_at, 'YYYY-MM-DD') AS day, COUNT(*) AS cnt FROM vocabulary_reviews WHERE user_id = $1 AND reviewed_at >= now() - INTERVAL '30 days' GROUP BY day ORDER BY day`, userID); err != nil {
		return domain.VocabularyStats{}, fmt.Errorf("get review heatmap: %w", err)
	}
	heatmap := make(map[string]int, len(heatRows))
	for _, hr := range heatRows {
		heatmap[hr.Day] = hr.Count
	}

	return domain.VocabularyStats{
		TotalWords:      totals.Total,
		MasteredWords:   totals.Mastered,
		DueCount:        totals.Due,
		BoxDistribution: boxDist,
		ReviewHeatmap:   heatmap,
	}, nil
}

func vocabRowsToDomain(rows []vocabularyRow) []domain.Vocabulary {
	result := make([]domain.Vocabulary, 0, len(rows))
	for _, r := range rows {
		result = append(result, r.domain())
	}
	return result
}

func vocabToParams(v domain.Vocabulary, tags pq.StringArray) map[string]any {
	return map[string]any{
		"id": v.ID, "user_id": v.UserID, "language": v.Language, "word": v.Word,
		"reading": v.Reading, "meaning": v.Meaning, "example_sentence": v.ExampleSentence,
		"example_translation": v.ExampleTranslation, "tags": tags,
		"difficulty_level": v.DifficultyLevel, "audio_url": v.AudioURL, "image_url": v.ImageURL,
		"ease_factor": v.EaseFactor, "interval_days": v.IntervalDays, "repetition_count": v.RepetitionCount,
		"current_box": v.CurrentBox, "next_review_at": v.NextReviewAt.UTC(),
		"mastery_score": v.MasteryScore, "created_at": v.CreatedAt.UTC(), "updated_at": v.UpdatedAt.UTC(),
	}
}

// Compile-time interface implementation check
var _ ports.VocabularyRepository = (*VocabularyRepository)(nil)

// silence unused import warning in older file
var _ = strings.TrimSpace
