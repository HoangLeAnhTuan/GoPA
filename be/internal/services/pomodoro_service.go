package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gopa/internal/core/domain"
	"gopa/internal/core/ports"
)

type PomodoroService struct {
	states  ports.PomodoroStore
	history ports.PomodoroHistoryRepository
	now     func() time.Time
}

func NewPomodoroService(states ports.PomodoroStore, history ports.PomodoroHistoryRepository) *PomodoroService {
	return &PomodoroService{states: states, history: history, now: time.Now}
}
func (s *PomodoroService) Get(ctx context.Context, userID uuid.UUID) (*domain.PomodoroState, error) {
	state, err := s.states.Get(ctx, userID)
	if errors.Is(err, domain.ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &state, nil
}
func (s *PomodoroService) Set(ctx context.Context, userID uuid.UUID, state domain.PomodoroState) error {
	state.UpdatedAt = s.now().UTC()
	if state.StartedAt.IsZero() {
		state.StartedAt = state.UpdatedAt
	}
	if err := state.Validate(); err != nil {
		return err
	}
	return s.states.Set(ctx, userID, state)
}
func (s *PomodoroService) Stop(ctx context.Context, userID uuid.UUID) error {
	state, err := s.states.Get(ctx, userID)
	if errors.Is(err, domain.ErrNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	if err := s.states.Delete(ctx, userID); err != nil {
		return err
	}
	endedAt := s.now().UTC()
	elapsed := int(endedAt.Sub(state.StartedAt).Seconds())
	if elapsed <= 0 {
		return nil
	}
	if elapsed > state.DurationSeconds {
		elapsed = state.DurationSeconds
	}
	history := domain.PomodoroHistory{ID: uuid.New(), EventID: uuid.New(), UserID: userID, TaskID: state.TaskID, StartedAt: state.StartedAt, EndedAt: endedAt, DurationSeconds: elapsed, CreatedAt: endedAt}
	return s.history.Create(ctx, history)
}

func (s *PomodoroService) ListHistory(ctx context.Context, userID uuid.UUID, limit int, cursor string) ([]domain.PomodoroHistory, error) {
	return s.history.List(ctx, userID, boundedLimit(limit), strings.TrimSpace(cursor))
}
