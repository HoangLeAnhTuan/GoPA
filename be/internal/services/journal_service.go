package services

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"gopa/internal/core/domain"
	"gopa/internal/core/ports"
)

type JournalInput struct {
	Title, Content string
	Tags           []string
}
type JournalService struct {
	journals ports.JournalRepository
	now      func() time.Time
}

func NewJournalService(journals ports.JournalRepository) *JournalService {
	return &JournalService{journals: journals, now: time.Now}
}
func (s *JournalService) List(ctx context.Context, userID uuid.UUID, term, tag string, limit int) ([]domain.Journal, error) {
	return s.journals.List(ctx, userID, strings.TrimSpace(term), strings.TrimSpace(tag), boundedLimit(limit))
}
func (s *JournalService) Get(ctx context.Context, userID, id uuid.UUID) (domain.Journal, error) {
	return s.journals.Get(ctx, userID, id)
}
func (s *JournalService) Create(ctx context.Context, userID uuid.UUID, input JournalInput) (domain.Journal, error) {
	now := s.now().UTC()
	journal := domain.Journal{ID: uuid.New(), UserID: userID, Title: strings.TrimSpace(input.Title), Content: input.Content, Tags: append([]string(nil), input.Tags...), CreatedAt: now, UpdatedAt: now}
	if err := journal.Validate(); err != nil {
		return domain.Journal{}, err
	}
	return s.journals.Create(ctx, journal)
}
func (s *JournalService) Update(ctx context.Context, userID, id uuid.UUID, input JournalInput) (domain.Journal, error) {
	journal, err := s.journals.Get(ctx, userID, id)
	if err != nil {
		return domain.Journal{}, err
	}
	journal.Title = strings.TrimSpace(input.Title)
	journal.Content = input.Content
	journal.Tags = append([]string(nil), input.Tags...)
	journal.UpdatedAt = s.now().UTC()
	if err := journal.Validate(); err != nil {
		return domain.Journal{}, err
	}
	return s.journals.Update(ctx, journal)
}
func (s *JournalService) Delete(ctx context.Context, userID, id uuid.UUID) error {
	return s.journals.Delete(ctx, userID, id)
}
