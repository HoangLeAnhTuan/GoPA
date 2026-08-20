package services

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"gopa/internal/core/domain"
)

type taskRepoStub struct {
	tasks           map[uuid.UUID]domain.Task
	nextPos         map[domain.TaskStatus]int64
	createErr       error
	getErr          error
	updateErr       error
	updateStatusErr error
	deleteErr       error
}

func newTaskRepoStub() *taskRepoStub {
	return &taskRepoStub{
		tasks:   make(map[uuid.UUID]domain.Task),
		nextPos: make(map[domain.TaskStatus]int64),
	}
}

func (s *taskRepoStub) List(_ context.Context, userID uuid.UUID) ([]domain.Task, error) {
	list := make([]domain.Task, 0)
	for _, t := range s.tasks {
		if t.UserID == userID {
			list = append(list, t)
		}
	}
	return list, nil
}

func (s *taskRepoStub) Get(_ context.Context, userID, id uuid.UUID) (domain.Task, error) {
	if s.getErr != nil {
		return domain.Task{}, s.getErr
	}
	t, ok := s.tasks[id]
	if !ok || t.UserID != userID {
		return domain.Task{}, domain.ErrNotFound
	}
	return t, nil
}

func (s *taskRepoStub) Create(_ context.Context, task domain.Task) (domain.Task, error) {
	if s.createErr != nil {
		return domain.Task{}, s.createErr
	}
	s.tasks[task.ID] = task
	return task, nil
}

func (s *taskRepoStub) Update(_ context.Context, task domain.Task) (domain.Task, error) {
	if s.updateErr != nil {
		return domain.Task{}, s.updateErr
	}
	if _, ok := s.tasks[task.ID]; !ok {
		return domain.Task{}, domain.ErrNotFound
	}
	s.tasks[task.ID] = task
	return task, nil
}

func (s *taskRepoStub) UpdateStatus(_ context.Context, userID, id uuid.UUID, status domain.TaskStatus, targetPosition *int64) (domain.Task, error) {
	if s.updateStatusErr != nil {
		return domain.Task{}, s.updateStatusErr
	}
	t, ok := s.tasks[id]
	if !ok || t.UserID != userID {
		return domain.Task{}, domain.ErrNotFound
	}
	t.Status = status
	if targetPosition != nil {
		t.Position = *targetPosition
	} else {
		t.Position = s.nextPos[status]
	}
	s.tasks[id] = t
	return t, nil
}

func (s *taskRepoStub) Delete(_ context.Context, userID, id uuid.UUID) error {
	if s.deleteErr != nil {
		return s.deleteErr
	}
	t, ok := s.tasks[id]
	if !ok || t.UserID != userID {
		return domain.ErrNotFound
	}
	delete(s.tasks, id)
	return nil
}

func (s *taskRepoStub) NextPosition(_ context.Context, _ uuid.UUID, status domain.TaskStatus) (int64, error) {
	return s.nextPos[status], nil
}

func TestTaskService_Create(t *testing.T) {
	userID := uuid.New()
	repo := newTaskRepoStub()
	repo.nextPos[domain.TaskStatusTodo] = 3
	svc := NewTaskService(repo)

	// 1. Create with defaults
	task, err := svc.Create(context.Background(), userID, CreateTaskInput{
		Title:              "  Implement Task Service  ",
		Description:        "Clean architecture Kanban service",
		EstimatedPomodoros: 4,
	})
	if err != nil {
		t.Fatalf("create task failed: %v", err)
	}
	if task.Title != "Implement Task Service" {
		t.Fatalf("expected trimmed title, got %q", task.Title)
	}
	if task.Status != domain.TaskStatusTodo {
		t.Fatalf("expected default status TODO, got %s", task.Status)
	}
	if task.Priority != domain.TaskPriorityMedium {
		t.Fatalf("expected default priority MEDIUM, got %s", task.Priority)
	}
	if task.Category != domain.TaskCategoryWork {
		t.Fatalf("expected default category WORK, got %s", task.Category)
	}
	if task.Position != 3 {
		t.Fatalf("expected position 3, got %d", task.Position)
	}
	if task.EstimatedPomodoros != 4 || task.CompletedPomodoros != 0 {
		t.Fatalf("expected 4 estimated, 0 completed pomodoros, got %d/%d", task.EstimatedPomodoros, task.CompletedPomodoros)
	}

	// 2. Validation error on empty title
	_, err = svc.Create(context.Background(), userID, CreateTaskInput{
		Title: "   ",
	})
	if !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("expected ErrValidation for empty title, got %v", err)
	}

	// 3. Validation error on negative estimated pomodoros
	_, err = svc.Create(context.Background(), userID, CreateTaskInput{
		Title:              "Task",
		EstimatedPomodoros: -1,
	})
	if !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("expected ErrValidation for negative pomodoros, got %v", err)
	}
}

func TestTaskService_GetAndList(t *testing.T) {
	userID := uuid.New()
	repo := newTaskRepoStub()
	svc := NewTaskService(repo)

	created, err := svc.Create(context.Background(), userID, CreateTaskInput{
		Title: "Sample Task",
	})
	if err != nil {
		t.Fatalf("create task: %v", err)
	}

	// Get existing
	got, err := svc.Get(context.Background(), userID, created.ID)
	if err != nil {
		t.Fatalf("get task: %v", err)
	}
	if got.ID != created.ID {
		t.Fatalf("expected task ID %s, got %s", created.ID, got.ID)
	}

	// Get not found
	_, err = svc.Get(context.Background(), userID, uuid.New())
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound for non-existing task, got %v", err)
	}

	// List
	list, err := svc.List(context.Background(), userID)
	if err != nil {
		t.Fatalf("list tasks: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 task in list, got %d", len(list))
	}
}

func TestTaskService_Update(t *testing.T) {
	userID := uuid.New()
	repo := newTaskRepoStub()
	repo.nextPos[domain.TaskStatusInProgress] = 5
	svc := NewTaskService(repo)

	task, err := svc.Create(context.Background(), userID, CreateTaskInput{
		Title: "Original Title",
	})
	if err != nil {
		t.Fatalf("create task: %v", err)
	}

	// 1. Update fields and change status (without explicit position -> auto calculates next position)
	est := 6
	comp := 2
	updated, err := svc.Update(context.Background(), userID, task.ID, UpdateTaskInput{
		Title:              "Updated Title",
		Description:        "Updated description",
		Status:             domain.TaskStatusInProgress,
		Priority:           domain.TaskPriorityUrgent,
		Category:           domain.TaskCategoryLeetcode,
		EstimatedPomodoros: &est,
		CompletedPomodoros: &comp,
	})
	if err != nil {
		t.Fatalf("update task: %v", err)
	}
	if updated.Title != "Updated Title" || updated.Status != domain.TaskStatusInProgress || updated.Priority != domain.TaskPriorityUrgent || updated.Category != domain.TaskCategoryLeetcode {
		t.Fatalf("unexpected updated fields: %+v", updated)
	}
	if updated.Position != 5 {
		t.Fatalf("expected position 5 after status shift, got %d", updated.Position)
	}
	if updated.EstimatedPomodoros != 6 || updated.CompletedPomodoros != 2 {
		t.Fatalf("expected 6 estimated, 2 completed, got %d/%d", updated.EstimatedPomodoros, updated.CompletedPomodoros)
	}

	// 2. Update with explicit position
	pos := int64(10)
	updated2, err := svc.Update(context.Background(), userID, task.ID, UpdateTaskInput{
		Title:    "Updated Title",
		Status:   domain.TaskStatusInProgress,
		Priority: domain.TaskPriorityUrgent,
		Category: domain.TaskCategoryLeetcode,
		Position: &pos,
	})
	if err != nil {
		t.Fatalf("update with explicit position: %v", err)
	}
	if updated2.Position != 10 {
		t.Fatalf("expected explicit position 10, got %d", updated2.Position)
	}

	// 3. Update not found
	_, err = svc.Update(context.Background(), userID, uuid.New(), UpdateTaskInput{
		Title:    "Nonexistent",
		Status:   domain.TaskStatusTodo,
		Priority: domain.TaskPriorityLow,
		Category: domain.TaskCategoryWork,
	})
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound for non-existing task, got %v", err)
	}
}

func TestTaskService_UpdateStatus(t *testing.T) {
	userID := uuid.New()
	repo := newTaskRepoStub()
	repo.nextPos[domain.TaskStatusDone] = 7
	svc := NewTaskService(repo)

	task, err := svc.Create(context.Background(), userID, CreateTaskInput{
		Title: "Task to move",
	})
	if err != nil {
		t.Fatalf("create task: %v", err)
	}

	// 1. Valid status update without targetPosition
	moved, err := svc.UpdateStatus(context.Background(), userID, task.ID, domain.TaskStatusDone, nil)
	if err != nil {
		t.Fatalf("update status: %v", err)
	}
	if moved.Status != domain.TaskStatusDone {
		t.Fatalf("expected status DONE, got %s", moved.Status)
	}
	if moved.Position != 7 {
		t.Fatalf("expected position 7, got %d", moved.Position)
	}

	// 2. Valid status update with targetPosition
	targetPos := int64(2)
	movedWithPos, err := svc.UpdateStatus(context.Background(), userID, task.ID, domain.TaskStatusInProgress, &targetPos)
	if err != nil {
		t.Fatalf("update status with position: %v", err)
	}
	if movedWithPos.Status != domain.TaskStatusInProgress || movedWithPos.Position != 2 {
		t.Fatalf("expected IN_PROGRESS at position 2, got %s at %d", movedWithPos.Status, movedWithPos.Position)
	}

	// 3. Invalid status validation error
	_, err = svc.UpdateStatus(context.Background(), userID, task.ID, domain.TaskStatus("INVALID"), nil)
	if !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("expected ErrValidation for invalid status, got %v", err)
	}

	// 4. Negative target position validation error
	negPos := int64(-1)
	_, err = svc.UpdateStatus(context.Background(), userID, task.ID, domain.TaskStatusDone, &negPos)
	if !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("expected ErrValidation for negative targetPosition, got %v", err)
	}

	// 5. Not found error
	_, err = svc.UpdateStatus(context.Background(), userID, uuid.New(), domain.TaskStatusDone, nil)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound for nonexistent task, got %v", err)
	}
}

func TestTaskService_Delete(t *testing.T) {
	userID := uuid.New()
	repo := newTaskRepoStub()
	svc := NewTaskService(repo)

	task, err := svc.Create(context.Background(), userID, CreateTaskInput{
		Title: "Task to delete",
	})
	if err != nil {
		t.Fatalf("create task: %v", err)
	}

	// Delete success
	if err := svc.Delete(context.Background(), userID, task.ID); err != nil {
		t.Fatalf("delete task: %v", err)
	}

	// Delete not found
	if err := svc.Delete(context.Background(), userID, task.ID); !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound on subsequent delete, got %v", err)
	}
}
