package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"gopa/internal/constants"
	"gopa/internal/core/domain"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type PomodoroStore struct{ client *redis.Client }

func NewPomodoroStore(client *redis.Client) *PomodoroStore { return &PomodoroStore{client: client} }
func (s *PomodoroStore) Get(ctx context.Context, userID uuid.UUID) (domain.PomodoroState, error) {
	raw, err := s.client.Get(ctx, pomodoroKey(userID)).Bytes()
	if errors.Is(err, redis.Nil) {
		return domain.PomodoroState{}, domain.ErrNotFound
	}
	if err != nil {
		return domain.PomodoroState{}, fmt.Errorf("get pomodoro state: %w", err)
	}
	var state domain.PomodoroState
	if err := json.Unmarshal(raw, &state); err != nil {
		return domain.PomodoroState{}, fmt.Errorf("unmarshal pomodoro state: %w", err)
	}
	return state, nil
}
func (s *PomodoroStore) Set(ctx context.Context, userID uuid.UUID, state domain.PomodoroState) error {
	raw, err := json.Marshal(state)
	if err != nil {
		return fmt.Errorf("marshal pomodoro state: %w", err)
	}
	if err := s.client.Set(ctx, pomodoroKey(userID), raw, 24*time.Hour).Err(); err != nil {
		return fmt.Errorf("set pomodoro state: %w", err)
	}
	_ = s.client.Publish(ctx, pomodoroChannel(userID), raw).Err()
	return nil
}
func (s *PomodoroStore) Delete(ctx context.Context, userID uuid.UUID) error {
	if err := s.client.Del(ctx, pomodoroKey(userID)).Err(); err != nil {
		return fmt.Errorf("delete pomodoro state: %w", err)
	}
	_ = s.client.Publish(ctx, pomodoroChannel(userID), []byte("null")).Err()
	return nil
}
func (s *PomodoroStore) Subscribe(ctx context.Context, userID uuid.UUID) *redis.PubSub {
	return s.client.Subscribe(ctx, pomodoroChannel(userID))
}
func pomodoroKey(userID uuid.UUID) string { return constants.CachePomodoroStateKey + userID.String() }
func pomodoroChannel(userID uuid.UUID) string {
	return constants.CachePomodoroChannelPrefix + userID.String()
}
