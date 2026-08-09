package repository

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"gopa/internal/core/domain"
)

type PomodoroHistoryRepository struct{ db *sqlx.DB }

func NewPomodoroHistoryRepository(db *sqlx.DB) *PomodoroHistoryRepository {
	return &PomodoroHistoryRepository{db: db}
}
func (r *PomodoroHistoryRepository) Create(ctx context.Context, history domain.PomodoroHistory) error {
	_, err := r.db.NamedExecContext(ctx, `INSERT INTO pomodoro_history (id, event_id, user_id, task_id, started_at, ended_at, duration_seconds, created_at) VALUES (:id, :event_id, :user_id, :task_id, :started_at, :ended_at, :duration_seconds, :created_at) ON CONFLICT (event_id) DO NOTHING`, history)
	if err != nil {
		return fmt.Errorf("create pomodoro history: %w", err)
	}
	return nil
}
