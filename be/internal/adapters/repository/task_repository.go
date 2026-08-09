package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"gopa/internal/core/domain"
)

type TaskRepository struct{ db *sqlx.DB }

func NewTaskRepository(db *sqlx.DB) *TaskRepository { return &TaskRepository{db: db} }

func (r *TaskRepository) List(ctx context.Context, userID uuid.UUID) ([]domain.Task, error) {
	const query = `SELECT id, user_id, title, description, status, priority, category, due_date, position, created_at, updated_at FROM tasks WHERE user_id = $1 ORDER BY status, position, created_at DESC`
	tasks := make([]domain.Task, 0)
	if err := r.db.SelectContext(ctx, &tasks, query, userID); err != nil {
		return nil, fmt.Errorf("list tasks: %w", err)
	}
	return tasks, nil
}
func (r *TaskRepository) Get(ctx context.Context, userID, id uuid.UUID) (domain.Task, error) {
	const query = `SELECT id, user_id, title, description, status, priority, category, due_date, position, created_at, updated_at FROM tasks WHERE id = $1 AND user_id = $2`
	var task domain.Task
	if err := r.db.GetContext(ctx, &task, query, id, userID); err != nil {
		return domain.Task{}, mapNotFound(err, "get task")
	}
	return task, nil
}
func (r *TaskRepository) Create(ctx context.Context, task domain.Task) (domain.Task, error) {
	const query = `INSERT INTO tasks (id, user_id, title, description, status, priority, category, due_date, position, created_at, updated_at) VALUES (:id, :user_id, :title, :description, :status, :priority, :category, :due_date, :position, :created_at, :updated_at)`
	if _, err := r.db.NamedExecContext(ctx, query, task); err != nil {
		return domain.Task{}, fmt.Errorf("create task: %w", err)
	}
	return task, nil
}
func (r *TaskRepository) Update(ctx context.Context, task domain.Task) (domain.Task, error) {
	const query = `UPDATE tasks SET title = :title, description = :description, status = :status, priority = :priority, category = :category, due_date = :due_date, position = :position, updated_at = :updated_at WHERE id = :id AND user_id = :user_id`
	result, err := r.db.NamedExecContext(ctx, query, task)
	if err != nil {
		return domain.Task{}, fmt.Errorf("update task: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return domain.Task{}, fmt.Errorf("count updated task: %w", err)
	}
	if rows == 0 {
		return domain.Task{}, domain.ErrNotFound
	}
	return task, nil
}
func (r *TaskRepository) Delete(ctx context.Context, userID, id uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM tasks WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("delete task: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("count deleted task: %w", err)
	}
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}
func (r *TaskRepository) NextPosition(ctx context.Context, userID uuid.UUID, status domain.TaskStatus) (int64, error) {
	var next int64
	if err := r.db.GetContext(ctx, &next, `SELECT COALESCE(MAX(position), -1) + 1 FROM tasks WHERE user_id = $1 AND status = $2`, userID, status); err != nil {
		return 0, fmt.Errorf("get next task position: %w", err)
	}
	return next, nil
}
func mapNotFound(err error, operation string) error {
	if errors.Is(err, sql.ErrNoRows) {
		return domain.ErrNotFound
	}
	return fmt.Errorf("%s: %w", operation, err)
}
