package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"gopa/internal/core/domain"
	"gopa/internal/core/ports"
)

type CreateTaskInput struct {
	Title, Description string
	Status             domain.TaskStatus
	Priority           domain.TaskPriority
	Category           domain.TaskCategory
	DueDate            *time.Time
}
type UpdateTaskInput struct {
	Title, Description string
	Status             domain.TaskStatus
	Priority           domain.TaskPriority
	Category           domain.TaskCategory
	DueDate            *time.Time
	Position           *int64
}
type TaskService struct {
	tasks ports.TaskRepository
	now   func() time.Time
}

func NewTaskService(tasks ports.TaskRepository) *TaskService {
	return &TaskService{tasks: tasks, now: time.Now}
}
func (s *TaskService) List(ctx context.Context, userID uuid.UUID) ([]domain.Task, error) {
	return s.tasks.List(ctx, userID)
}
func (s *TaskService) Create(ctx context.Context, userID uuid.UUID, input CreateTaskInput) (domain.Task, error) {
	if input.Status == "" {
		input.Status = domain.TaskStatusTodo
	}
	if input.Priority == "" {
		input.Priority = domain.TaskPriorityMedium
	}
	if input.Category == "" {
		input.Category = domain.TaskCategoryWork
	}
	position, err := s.tasks.NextPosition(ctx, userID, input.Status)
	if err != nil {
		return domain.Task{}, err
	}
	now := s.now().UTC()
	task := domain.Task{ID: uuid.New(), UserID: userID, Title: strings.TrimSpace(input.Title), Description: input.Description, Status: input.Status, Priority: input.Priority, Category: input.Category, DueDate: input.DueDate, Position: position, CreatedAt: now, UpdatedAt: now}
	if err := task.Validate(); err != nil {
		return domain.Task{}, err
	}
	return s.tasks.Create(ctx, task)
}
func (s *TaskService) Update(ctx context.Context, userID, id uuid.UUID, input UpdateTaskInput) (domain.Task, error) {
	task, err := s.tasks.Get(ctx, userID, id)
	if err != nil {
		return domain.Task{}, err
	}
	task.Title = strings.TrimSpace(input.Title)
	task.Description = input.Description
	oldStatus := task.Status
	task.Status = input.Status
	task.Priority = input.Priority
	task.Category = input.Category
	task.DueDate = input.DueDate
	if input.Position != nil {
		task.Position = *input.Position
	} else if oldStatus != task.Status {
		position, err := s.tasks.NextPosition(ctx, userID, task.Status)
		if err != nil {
			return domain.Task{}, fmt.Errorf("set moved task position: %w", err)
		}
		task.Position = position
	}
	task.UpdatedAt = s.now().UTC()
	if err := task.Validate(); err != nil {
		return domain.Task{}, err
	}
	return s.tasks.Update(ctx, task)
}
func (s *TaskService) Delete(ctx context.Context, userID, id uuid.UUID) error {
	return s.tasks.Delete(ctx, userID, id)
}
