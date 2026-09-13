package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/rabbitmq/amqp091-go"
	"gopa/internal/adapters/broker"
	"gopa/internal/core/domain"
	"gopa/internal/core/domain/linguistics"
	"gopa/internal/core/ports"
)

// ─── Pomodoro Finished ────────────────────────────────────────────────────────

// pomodoroFinishedPayload mirrors the JSON written to outbox_events by
// PomodoroHistoryRepository.Create.
type pomodoroFinishedPayload struct {
	HistoryID       string  `json:"history_id"`
	EventID         string  `json:"event_id"`
	UserID          string  `json:"user_id"`
	TaskID          *string `json:"task_id"`
	DurationSeconds int     `json:"duration_seconds"`
}

// HandlePomodoroFinished returns a broker.HandlerFunc that:
//  1. Parses the pomodoro.finished event payload.
//  2. Persists the pomodoro_history record (idempotent via event_id ON CONFLICT).
//  3. Increments completed_pomodoros on the linked task, if one is present.
//
// The db parameter is used only for the task counter increment; history
// persistence delegates through PomodoroHistoryRepository which owns its own
// transaction.
func HandlePomodoroFinished(history ports.PomodoroHistoryRepository, db *sqlx.DB, logger *slog.Logger) broker.HandlerFunc {
	return func(ctx context.Context, delivery amqp091.Delivery, event broker.EventEnvelope) error {
		payloadBytes := event.Payload
		if len(payloadBytes) == 0 {
			payloadBytes = delivery.Body
		}
		var p pomodoroFinishedPayload
		if err := json.Unmarshal(payloadBytes, &p); err != nil {
			return fmt.Errorf("parse pomodoro.finished payload: %w", err)
		}

		histID, err := uuid.Parse(p.HistoryID)
		if err != nil {
			return fmt.Errorf("invalid history_id %q: %w", p.HistoryID, err)
		}
		eventID, err := uuid.Parse(p.EventID)
		if err != nil {
			return fmt.Errorf("invalid event_id %q: %w", p.EventID, err)
		}
		userID, err := uuid.Parse(p.UserID)
		if err != nil {
			return fmt.Errorf("invalid user_id %q: %w", p.UserID, err)
		}

		var taskID *uuid.UUID
		if p.TaskID != nil && *p.TaskID != "" && *p.TaskID != "null" {
			tid, err := uuid.Parse(*p.TaskID)
			if err != nil {
				return fmt.Errorf("invalid task_id %q: %w", *p.TaskID, err)
			}
			taskID = &tid
		}

		now := time.Now().UTC()
		h := domain.PomodoroHistory{
			ID:              histID,
			EventID:         eventID,
			UserID:          userID,
			TaskID:          taskID,
			StartedAt:       event.OccurredAt,
			EndedAt:         now,
			DurationSeconds: p.DurationSeconds,
			CreatedAt:       now,
		}

		if err := history.Create(ctx, h); err != nil {
			return fmt.Errorf("persist pomodoro history: %w", err)
		}

		if taskID != nil {
			if _, err := db.ExecContext(
				ctx,
				`UPDATE tasks SET completed_pomodoros = completed_pomodoros + 1, updated_at = now()
				 WHERE id = $1 AND user_id = $2`,
				taskID, userID,
			); err != nil {
				// Non-fatal: task counter is best-effort; log and continue.
				logger.Warn("failed to increment task completed_pomodoros",
					"task_id", taskID, "user_id", userID, "error", err)
			}
		}

		logger.Info("pomodoro.finished processed",
			"event_id", eventID, "user_id", userID, "task_id", taskID,
			"duration_seconds", p.DurationSeconds)
		return nil
	}
}

// ─── Review Completed ─────────────────────────────────────────────────────────

// reviewCompletedPayload mirrors the outbox payload emitted after a vocabulary
// review is submitted via the SRS HTTP endpoint.
type reviewCompletedPayload struct {
	ReviewID     string `json:"review_id"`
	EventID      string `json:"event_id"`
	UserID       string `json:"user_id"`
	VocabularyID string `json:"vocabulary_id"`
	Quality      int16  `json:"quality"`
	StudyMode    string `json:"study_mode"`
}

// HandleReviewCompleted returns a broker.HandlerFunc that applies an SM-2 SRS
// schedule update for the reviewed vocabulary item.
// This mirrors VocabularyService.Review but is driven by the async event
// rather than a synchronous HTTP request, allowing cross-service decoupling.
func HandleReviewCompleted(vocabularies ports.VocabularyRepository, logger *slog.Logger) broker.HandlerFunc {
	return func(ctx context.Context, delivery amqp091.Delivery, event broker.EventEnvelope) error {
		payloadBytes := event.Payload
		if len(payloadBytes) == 0 {
			payloadBytes = delivery.Body
		}
		var p reviewCompletedPayload
		if err := json.Unmarshal(payloadBytes, &p); err != nil {
			return fmt.Errorf("parse learning.review_completed payload: %w", err)
		}

		userID, err := uuid.Parse(p.UserID)
		if err != nil {
			return fmt.Errorf("invalid user_id %q: %w", p.UserID, err)
		}
		vocabID, err := uuid.Parse(p.VocabularyID)
		if err != nil {
			return fmt.Errorf("invalid vocabulary_id %q: %w", p.VocabularyID, err)
		}

		vocabulary, err := vocabularies.Get(ctx, userID, vocabID)
		if err != nil {
			return fmt.Errorf("get vocabulary %s: %w", vocabID, err)
		}

		quality := p.Quality
		if quality < 0 {
			quality = 0
		} else if quality > 5 {
			quality = 5
		}

		now := time.Now().UTC()
		result := linguistics.CalculateNextReview(linguistics.VocabularyState{
			CurrentBox:      vocabulary.CurrentBox,
			EaseFactor:      vocabulary.EaseFactor,
			IntervalDays:    vocabulary.IntervalDays,
			RepetitionCount: vocabulary.RepetitionCount,
			MasteryScore:    vocabulary.MasteryScore,
		}, quality, now)

		vocabulary.CurrentBox = result.NextBox
		vocabulary.EaseFactor = result.EaseFactor
		vocabulary.IntervalDays = result.IntervalDays
		vocabulary.RepetitionCount = result.RepetitionCount
		vocabulary.MasteryScore = result.MasteryScore
		vocabulary.NextReviewAt = result.NextReviewAt

		reviewID, _ := uuid.Parse(p.ReviewID)
		eventID, _ := uuid.Parse(p.EventID)
		if reviewID == uuid.Nil {
			reviewID = uuid.New()
		}
		if eventID == uuid.Nil {
			eventID = event.EventID
		}

		studyMode := p.StudyMode
		if studyMode == "" {
			studyMode = "FLASHCARD"
		}

		review := domain.VocabularyReview{
			ID:             reviewID,
			EventID:        eventID,
			VocabularyID:   vocabID,
			UserID:         userID,
			Quality:        quality,
			StudyMode:      studyMode,
			ResponseTimeMs: 0,
			WasCorrect:     quality >= 3,
			ReviewedAt:     event.OccurredAt,
			CreatedAt:      now,
		}

		if err := vocabularies.RecordReview(ctx, review, vocabulary); err != nil {
			return fmt.Errorf("record review for vocabulary %s: %w", vocabID, err)
		}

		logger.Info("learning.review_completed processed",
			"event_id", event.EventID, "user_id", userID,
			"vocabulary_id", vocabID, "quality", quality,
			"next_box", result.NextBox, "next_review_at", result.NextReviewAt)
		return nil
	}
}

// ─── Budget Alert ─────────────────────────────────────────────────────────────

// budgetAlertPayload mirrors the finance.budget_alert outbox event emitted
// when a transaction causes a budget's utilization to cross its alert threshold.
type budgetAlertPayload struct {
	BudgetID        string `json:"budget_id"`
	UserID          string `json:"user_id"`
	UtilizationRate string `json:"utilization_rate"` // decimal string e.g. "0.85"
	AlertThreshold  string `json:"alert_threshold"`  // decimal string e.g. "0.80"
}

// HandleBudgetAlert returns a broker.HandlerFunc that:
//  1. Parses the finance.budget_alert payload.
//  2. Retrieves the live BudgetStatus from the repository to confirm the alert.
//  3. Emits a structured warning log with budget and utilization context.
//
// No state mutation is performed; alert state is derived on-demand by
// GetBudgetStatus, which computes utilization from live transaction data.
func HandleBudgetAlert(budgets ports.BudgetRepository, logger *slog.Logger) broker.HandlerFunc {
	return func(ctx context.Context, delivery amqp091.Delivery, event broker.EventEnvelope) error {
		payloadBytes := event.Payload
		if len(payloadBytes) == 0 {
			payloadBytes = delivery.Body
		}
		var p budgetAlertPayload
		if err := json.Unmarshal(payloadBytes, &p); err != nil {
			return fmt.Errorf("parse finance.budget_alert payload: %w", err)
		}

		userID, err := uuid.Parse(p.UserID)
		if err != nil {
			return fmt.Errorf("invalid user_id %q: %w", p.UserID, err)
		}
		budgetID, err := uuid.Parse(p.BudgetID)
		if err != nil {
			return fmt.Errorf("invalid budget_id %q: %w", p.BudgetID, err)
		}

		status, err := budgets.GetBudgetStatus(ctx, userID, budgetID)
		if err != nil {
			return fmt.Errorf("get budget status for %s: %w", budgetID, err)
		}

		logger.Warn("budget alert triggered",
			"event_id", event.EventID,
			"user_id", userID,
			"budget_id", budgetID,
			"budget_name", status.Budget.Name,
			"budget_amount", status.Budget.Amount,
			"spent_amount", status.SpentAmount,
			"utilization_rate", status.UtilizationRate,
			"alert_threshold", status.Budget.AlertThreshold,
			"is_alert_triggered", status.IsAlertTriggered,
		)

		return nil
	}
}
