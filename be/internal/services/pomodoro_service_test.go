package services

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"gopa/internal/core/domain"
)

type pomodoroStoreStub struct {
	states map[uuid.UUID]domain.PomodoroState
}

func newPomodoroStoreStub() *pomodoroStoreStub {
	return &pomodoroStoreStub{states: make(map[uuid.UUID]domain.PomodoroState)}
}
func (s *pomodoroStoreStub) Get(_ context.Context, userID uuid.UUID) (domain.PomodoroState, error) {
	st, ok := s.states[userID]
	if !ok {
		return domain.PomodoroState{}, domain.ErrNotFound
	}
	return st, nil
}
func (s *pomodoroStoreStub) Set(_ context.Context, userID uuid.UUID, state domain.PomodoroState) error {
	s.states[userID] = state
	return nil
}
func (s *pomodoroStoreStub) Delete(_ context.Context, userID uuid.UUID) error {
	delete(s.states, userID)
	return nil
}

type pomodoroHistoryStub struct {
	history []domain.PomodoroHistory
}

func (s *pomodoroHistoryStub) Create(_ context.Context, h domain.PomodoroHistory) error {
	s.history = append(s.history, h)
	return nil
}
func (s *pomodoroHistoryStub) List(_ context.Context, userID uuid.UUID, limit int, _ string) ([]domain.PomodoroHistory, error) {
	var list []domain.PomodoroHistory
	for _, h := range s.history {
		if h.UserID == userID {
			list = append(list, h)
			if len(list) >= limit {
				break
			}
		}
	}
	return list, nil
}

func TestPomodoroService_GetSetStop(t *testing.T) {
	userID := uuid.New()
	taskID := uuid.New()
	store := newPomodoroStoreStub()
	hist := &pomodoroHistoryStub{}
	svc := NewPomodoroService(store, hist)

	// 1. Initial Get is nil
	initial, err := svc.Get(context.Background(), userID)
	if err != nil {
		t.Fatalf("get initial: %v", err)
	}
	if initial != nil {
		t.Fatalf("expected nil initial state, got %+v", initial)
	}

	// 2. Set State
	now := time.Now().UTC()
	start := now.Add(-10 * time.Minute)
	err = svc.Set(context.Background(), userID, domain.PomodoroState{
		Status:          domain.PomodoroRunning,
		StartedAt:       start,
		DurationSeconds: 1500,
		TaskID:          &taskID,
	})
	if err != nil {
		t.Fatalf("set state: %v", err)
	}

	// 3. Get State
	cur, err := svc.Get(context.Background(), userID)
	if err != nil {
		t.Fatalf("get state: %v", err)
	}
	if cur == nil || cur.Status != domain.PomodoroRunning || *cur.TaskID != taskID {
		t.Fatalf("unexpected state: %+v", cur)
	}

	// 4. Stop
	if err := svc.Stop(context.Background(), userID); err != nil {
		t.Fatalf("stop timer: %v", err)
	}

	// 5. Verify deleted from store and inserted into history
	stopped, err := svc.Get(context.Background(), userID)
	if err != nil {
		t.Fatalf("get after stop: %v", err)
	}
	if stopped != nil {
		t.Fatalf("expected state to be cleared after stop, got %+v", stopped)
	}
	if len(hist.history) != 1 {
		t.Fatalf("expected 1 history record, got %d", len(hist.history))
	}
	if *hist.history[0].TaskID != taskID {
		t.Fatalf("expected history task ID %s, got %v", taskID, hist.history[0].TaskID)
	}

	// 6. List History
	list, err := svc.ListHistory(context.Background(), userID, 10, "")
	if err != nil {
		t.Fatalf("list history: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 history entry, got %d", len(list))
	}
}
