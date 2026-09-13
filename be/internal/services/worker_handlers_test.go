package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
	"gopa/internal/adapters/broker"
	"gopa/internal/core/domain"
)

// ─── Shared test helpers ───────────────────────────────────────────────────────

var discardLogger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError + 10}))

// buildDelivery creates a minimal amqp091.Delivery for testing.
func buildDelivery(routingKey string) amqp091.Delivery {
	return amqp091.Delivery{
		RoutingKey: routingKey,
		Timestamp:  time.Now().UTC(),
	}
}

// buildEnvelope wraps an arbitrary payload map into a broker.EventEnvelope.
func buildEnvelope(eventType string, payload any) broker.EventEnvelope {
	raw, _ := json.Marshal(payload)
	return broker.EventEnvelope{
		EventID:    uuid.New(),
		EventType:  eventType,
		OccurredAt: time.Now().UTC(),
		Payload:    raw,
	}
}

// ─── HandlePomodoroFinished ───────────────────────────────────────────────────

// pomodoroHistoryRecorder records Create calls and exposes them for assertions.
type pomodoroHistoryRecorder struct {
	calls []domain.PomodoroHistory
}

func (r *pomodoroHistoryRecorder) Create(_ context.Context, h domain.PomodoroHistory) error {
	r.calls = append(r.calls, h)
	return nil
}

func (r *pomodoroHistoryRecorder) List(_ context.Context, _ uuid.UUID, _ int, _ string) ([]domain.PomodoroHistory, error) {
	return nil, nil
}

func TestHandlePomodoroFinished_HappyPath(t *testing.T) {
	histID := uuid.New()
	eventID := uuid.New()
	userID := uuid.New()

	rec := &pomodoroHistoryRecorder{}
	// Nil db: task increment is skipped (no taskID provided)
	handler := HandlePomodoroFinished(rec, nil, discardLogger)

	env := buildEnvelope("pomodoro.finished", map[string]any{
		"history_id":       histID.String(),
		"event_id":         eventID.String(),
		"user_id":          userID.String(),
		"duration_seconds": 1500,
	})

	if err := handler(context.Background(), buildDelivery("pomodoro.finished"), env); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(rec.calls) != 1 {
		t.Fatalf("expected 1 history create call, got %d", len(rec.calls))
	}
	if rec.calls[0].EventID != eventID {
		t.Errorf("expected event_id %s, got %s", eventID, rec.calls[0].EventID)
	}
	if rec.calls[0].DurationSeconds != 1500 {
		t.Errorf("expected duration 1500, got %d", rec.calls[0].DurationSeconds)
	}
}

func TestHandlePomodoroFinished_MalformedPayload(t *testing.T) {
	rec := &pomodoroHistoryRecorder{}
	handler := HandlePomodoroFinished(rec, nil, discardLogger)

	env := broker.EventEnvelope{
		EventID:   uuid.New(),
		EventType: "pomodoro.finished",
		Payload:   []byte(`not json`),
	}

	if err := handler(context.Background(), buildDelivery("pomodoro.finished"), env); err == nil {
		t.Fatal("expected an error for malformed payload, got nil")
	}
}

func TestHandlePomodoroFinished_InvalidHistoryID(t *testing.T) {
	rec := &pomodoroHistoryRecorder{}
	handler := HandlePomodoroFinished(rec, nil, discardLogger)

	env := buildEnvelope("pomodoro.finished", map[string]any{
		"history_id":       "not-a-uuid",
		"event_id":         uuid.New().String(),
		"user_id":          uuid.New().String(),
		"duration_seconds": 900,
	})

	if err := handler(context.Background(), buildDelivery("pomodoro.finished"), env); err == nil {
		t.Fatal("expected an error for invalid history_id, got nil")
	}
}

// ─── HandleReviewCompleted ────────────────────────────────────────────────────

// workerVocabRepoStub is a minimal stub for the VocabularyRepository port used
// in worker handler tests. Named distinctly to avoid collision with stubs in
// finance_service_test.go which share the same package.
type workerVocabRepoStub struct {
	vocab         domain.Vocabulary
	reviewCalled  bool
	recordedVocab domain.Vocabulary
	err           error
}

func (s *workerVocabRepoStub) Get(_ context.Context, _, _ uuid.UUID) (domain.Vocabulary, error) {
	return s.vocab, s.err
}

func (s *workerVocabRepoStub) RecordReview(_ context.Context, _ domain.VocabularyReview, updated domain.Vocabulary) error {
	s.reviewCalled = true
	s.recordedVocab = updated
	return s.err
}

// Unimplemented stubs to satisfy the interface.
func (s *workerVocabRepoStub) List(_ context.Context, _ uuid.UUID, _ int) ([]domain.Vocabulary, error) {
	return nil, nil
}
func (s *workerVocabRepoStub) ListDue(_ context.Context, _ uuid.UUID, _ int, _ time.Time) ([]domain.Vocabulary, error) {
	return nil, nil
}
func (s *workerVocabRepoStub) Create(_ context.Context, v domain.Vocabulary) (domain.Vocabulary, error) {
	return v, nil
}
func (s *workerVocabRepoStub) BulkCreate(_ context.Context, _ uuid.UUID, vs []domain.Vocabulary) ([]domain.Vocabulary, error) {
	return vs, nil
}
func (s *workerVocabRepoStub) Update(_ context.Context, v domain.Vocabulary) (domain.Vocabulary, error) {
	return v, nil
}
func (s *workerVocabRepoStub) Delete(_ context.Context, _, _ uuid.UUID) error { return nil }
func (s *workerVocabRepoStub) CreateLearningSession(_ context.Context, ls domain.LearningSession) (domain.LearningSession, error) {
	return ls, nil
}
func (s *workerVocabRepoStub) UpdateLearningSession(_ context.Context, ls domain.LearningSession) (domain.LearningSession, error) {
	return ls, nil
}
func (s *workerVocabRepoStub) GetStats(_ context.Context, _ uuid.UUID) (domain.VocabularyStats, error) {
	return domain.VocabularyStats{}, nil
}

func TestHandleReviewCompleted_HappyPath(t *testing.T) {
	userID := uuid.New()
	vocabID := uuid.New()

	stub := &workerVocabRepoStub{
		vocab: domain.Vocabulary{
			ID:              vocabID,
			UserID:          userID,
			Language:        domain.VocabularyLanguageJapanese,
			Word:            "猫",
			Meaning:         "cat",
			EaseFactor:      2.5,
			IntervalDays:    1,
			RepetitionCount: 1,
			CurrentBox:      2,
			MasteryScore:    20,
			NextReviewAt:    time.Now().UTC(),
			CreatedAt:       time.Now().UTC(),
			UpdatedAt:       time.Now().UTC(),
		},
	}

	handler := HandleReviewCompleted(stub, discardLogger)

	env := buildEnvelope("learning.review_completed", map[string]any{
		"review_id":     uuid.New().String(),
		"event_id":      uuid.New().String(),
		"user_id":       userID.String(),
		"vocabulary_id": vocabID.String(),
		"quality":       4,
		"study_mode":    "FLASHCARD",
	})

	if err := handler(context.Background(), buildDelivery("learning.review_completed"), env); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !stub.reviewCalled {
		t.Fatal("expected RecordReview to be called")
	}
	// Quality ≥ 3 → box should advance
	if stub.recordedVocab.CurrentBox <= 2 {
		t.Errorf("expected box to advance beyond 2, got %d", stub.recordedVocab.CurrentBox)
	}
}

func TestHandleReviewCompleted_LowQualityResetsBox(t *testing.T) {
	userID := uuid.New()
	vocabID := uuid.New()

	stub := &workerVocabRepoStub{
		vocab: domain.Vocabulary{
			ID: vocabID, UserID: userID,
			Language: domain.VocabularyLanguageEnglish, Word: "hello", Meaning: "greeting",
			EaseFactor: 2.5, IntervalDays: 6, RepetitionCount: 2,
			CurrentBox: 3, MasteryScore: 40,
			NextReviewAt: time.Now().UTC(), CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(),
		},
	}

	handler := HandleReviewCompleted(stub, discardLogger)

	env := buildEnvelope("learning.review_completed", map[string]any{
		"review_id":     uuid.New().String(),
		"event_id":      uuid.New().String(),
		"user_id":       userID.String(),
		"vocabulary_id": vocabID.String(),
		"quality":       1, // below threshold → box reset to 1
		"study_mode":    "FLASHCARD",
	})

	if err := handler(context.Background(), buildDelivery("learning.review_completed"), env); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stub.recordedVocab.CurrentBox != 1 {
		t.Errorf("expected box reset to 1, got %d", stub.recordedVocab.CurrentBox)
	}
}

func TestHandleReviewCompleted_RepositoryError(t *testing.T) {
	userID := uuid.New()
	vocabID := uuid.New()

	stub := &workerVocabRepoStub{err: fmt.Errorf("db connection lost")}
	handler := HandleReviewCompleted(stub, discardLogger)

	env := buildEnvelope("learning.review_completed", map[string]any{
		"user_id":       userID.String(),
		"vocabulary_id": vocabID.String(),
		"quality":       4,
	})

	if err := handler(context.Background(), buildDelivery("learning.review_completed"), env); err == nil {
		t.Fatal("expected error from repository failure, got nil")
	}
}

// ─── HandleBudgetAlert ────────────────────────────────────────────────────────

// workerBudgetRepoStub is a minimal stub for the BudgetRepository port used in
// worker handler tests. Named distinctly to avoid collision with the budgetRepoStub
// declared in finance_service_test.go (same package).
type workerBudgetRepoStub struct {
	status domain.BudgetStatus
	err    error
}

func (s *workerBudgetRepoStub) GetBudgetStatus(_ context.Context, _, _ uuid.UUID) (domain.BudgetStatus, error) {
	return s.status, s.err
}
func (s *workerBudgetRepoStub) ListBudgets(_ context.Context, _ uuid.UUID, _ bool) ([]domain.Budget, error) {
	return nil, nil
}
func (s *workerBudgetRepoStub) GetBudget(_ context.Context, _, _ uuid.UUID) (domain.Budget, error) {
	return domain.Budget{}, nil
}
func (s *workerBudgetRepoStub) CreateBudget(_ context.Context, b domain.Budget) (domain.Budget, error) {
	return b, nil
}
func (s *workerBudgetRepoStub) UpdateBudget(_ context.Context, b domain.Budget) (domain.Budget, error) {
	return b, nil
}
func (s *workerBudgetRepoStub) DeleteBudget(_ context.Context, _, _ uuid.UUID) error { return nil }
func (s *workerBudgetRepoStub) ListActiveBudgetsForCategory(_ context.Context, _, _ uuid.UUID) ([]domain.Budget, error) {
	return nil, nil
}

func TestHandleBudgetAlert_HappyPath(t *testing.T) {
	userID := uuid.New()
	budgetID := uuid.New()

	stub := &workerBudgetRepoStub{
		status: domain.BudgetStatus{
			Budget:           domain.Budget{ID: budgetID, UserID: userID, Name: "Monthly Groceries"},
			IsAlertTriggered: true,
		},
	}

	handler := HandleBudgetAlert(stub, discardLogger)

	env := buildEnvelope("finance.budget_alert", map[string]any{
		"budget_id":        budgetID.String(),
		"user_id":          userID.String(),
		"utilization_rate": "0.87",
		"alert_threshold":  "0.80",
	})

	if err := handler(context.Background(), buildDelivery("finance.budget_alert"), env); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestHandleBudgetAlert_RepositoryError(t *testing.T) {
	stub := &workerBudgetRepoStub{err: fmt.Errorf("timeout")}
	handler := HandleBudgetAlert(stub, discardLogger)

	env := buildEnvelope("finance.budget_alert", map[string]any{
		"budget_id": uuid.New().String(),
		"user_id":   uuid.New().String(),
	})

	if err := handler(context.Background(), buildDelivery("finance.budget_alert"), env); err == nil {
		t.Fatal("expected error from repository failure, got nil")
	}
}

func TestHandleBudgetAlert_InvalidBudgetID(t *testing.T) {
	stub := &workerBudgetRepoStub{}
	handler := HandleBudgetAlert(stub, discardLogger)

	env := buildEnvelope("finance.budget_alert", map[string]any{
		"budget_id": "not-a-uuid",
		"user_id":   uuid.New().String(),
	})

	if err := handler(context.Background(), buildDelivery("finance.budget_alert"), env); err == nil {
		t.Fatal("expected error for invalid budget_id, got nil")
	}
}
