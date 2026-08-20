package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"gopa/internal/core/domain"
	"gopa/internal/core/ports"
	"gopa/pkg/database"
)

type PomodoroHistoryRepository struct{ db *sqlx.DB }

func NewPomodoroHistoryRepository(db *sqlx.DB) *PomodoroHistoryRepository {
	return &PomodoroHistoryRepository{db: db}
}

func (r *PomodoroHistoryRepository) Create(ctx context.Context, history domain.PomodoroHistory) error {
	return database.WithTx(ctx, r.db, func(tx *sqlx.Tx) error {
		_, err := tx.NamedExecContext(ctx, `INSERT INTO pomodoro_history (id, event_id, user_id, task_id, started_at, ended_at, duration_seconds, created_at) VALUES (:id, :event_id, :user_id, :task_id, :started_at, :ended_at, :duration_seconds, :created_at) ON CONFLICT (event_id) DO NOTHING`, history)
		if err != nil {
			return fmt.Errorf("create pomodoro history: %w", err)
		}
		var taskIDStr string
		if history.TaskID != nil {
			taskIDStr = fmt.Sprintf(`"%s"`, history.TaskID.String())
		} else {
			taskIDStr = "null"
		}
		payload := fmt.Sprintf(`{"history_id":"%s","event_id":"%s","user_id":"%s","task_id":%s,"duration_seconds":%d}`, history.ID, history.EventID, history.UserID, taskIDStr, history.DurationSeconds)
		if _, err := tx.ExecContext(ctx,
			`INSERT INTO outbox_events (id, event_type, routing_key, payload, occurred_at, created_at) VALUES ($1, 'pomodoro.finished', 'pomodoro.finished', $2, $3, $3)`,
			history.EventID, payload, history.EndedAt.UTC()); err != nil {
			return fmt.Errorf("insert pomodoro outbox event: %w", err)
		}
		return nil
	})
}

// List returns the most recent pomodoro history entries for a user, ordered by
// started_at descending. The cursor parameter is an ISO-8601 timestamp; only
// entries started before that time are returned (exclusive), enabling
// keyset pagination.
func (r *PomodoroHistoryRepository) List(ctx context.Context, userID uuid.UUID, limit int, cursor string) ([]domain.PomodoroHistory, error) {
	rows := make([]domain.PomodoroHistory, 0)
	var err error
	if cursor == "" {
		err = r.db.SelectContext(ctx, &rows,
			`SELECT id, event_id, user_id, task_id, started_at, ended_at, duration_seconds, created_at FROM pomodoro_history WHERE user_id = $1 ORDER BY started_at DESC LIMIT $2`,
			userID, limit)
	} else {
		err = r.db.SelectContext(ctx, &rows,
			`SELECT id, event_id, user_id, task_id, started_at, ended_at, duration_seconds, created_at FROM pomodoro_history WHERE user_id = $1 AND started_at < $2::timestamptz ORDER BY started_at DESC LIMIT $3`,
			userID, cursor, limit)
	}
	if err != nil {
		return nil, fmt.Errorf("list pomodoro history: %w", err)
	}
	return rows, nil
}

var _ ports.PomodoroHistoryRepository = (*PomodoroHistoryRepository)(nil)
