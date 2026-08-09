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
	Title         string
	Content       string
	Tags          []string
	Mood          *domain.JournalMood
	PublishedDate *time.Time
}

type JournalService struct {
	journals ports.JournalRepository
	now      func() time.Time
}

func NewJournalService(journals ports.JournalRepository) *JournalService {
	return &JournalService{journals: journals, now: time.Now}
}

func (s *JournalService) List(ctx context.Context, userID uuid.UUID, filter domain.JournalFilter) ([]domain.Journal, error) {
	filter.Term, filter.Tag = strings.TrimSpace(filter.Term), strings.TrimSpace(strings.ToLower(filter.Tag))
	filter.Limit = boundedLimit(filter.Limit)
	if !isValidJournalMood(filter.Mood) || (filter.From != nil && filter.To != nil && filter.From.After(*filter.To)) {
		return nil, domain.ErrValidation
	}
	return s.journals.List(ctx, userID, filter)
}

func (s *JournalService) Get(ctx context.Context, userID, id uuid.UUID) (domain.Journal, error) {
	return s.journals.Get(ctx, userID, id)
}

func (s *JournalService) Create(ctx context.Context, userID uuid.UUID, input JournalInput) (domain.Journal, error) {
	now := s.now().UTC()
	publishedDate := now
	if input.PublishedDate != nil {
		publishedDate = input.PublishedDate.UTC()
	}
	journal := domain.Journal{ID: uuid.New(), UserID: userID, Title: strings.TrimSpace(input.Title), Content: input.Content, Tags: append([]string(nil), input.Tags...), Mood: input.Mood, PublishedDate: publishedDate, CreatedAt: now, UpdatedAt: now}
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
	journal.Title, journal.Content, journal.Tags, journal.Mood = strings.TrimSpace(input.Title), input.Content, append([]string(nil), input.Tags...), input.Mood
	if input.PublishedDate != nil {
		journal.PublishedDate = input.PublishedDate.UTC()
	}
	journal.UpdatedAt = s.now().UTC()
	if err := journal.Validate(); err != nil {
		return domain.Journal{}, err
	}
	return s.journals.Update(ctx, journal)
}

func (s *JournalService) Delete(ctx context.Context, userID, id uuid.UUID) error {
	return s.journals.Delete(ctx, userID, id)
}

func (s *JournalService) Link(ctx context.Context, userID, journalID, linkedID uuid.UUID) error {
	if journalID == linkedID {
		return domain.ErrValidation
	}
	return s.journals.Link(ctx, userID, journalID, linkedID)
}

func (s *JournalService) Unlink(ctx context.Context, userID, journalID, linkedID uuid.UUID) error {
	if journalID == linkedID {
		return domain.ErrValidation
	}
	return s.journals.Unlink(ctx, userID, journalID, linkedID)
}

func (s *JournalService) Stats(ctx context.Context, userID uuid.UUID) (domain.JournalStats, error) {
	return s.journals.Stats(ctx, userID)
}

func isValidJournalMood(mood *domain.JournalMood) bool {
	return mood == nil || *mood == domain.JournalMoodVeryNegative || *mood == domain.JournalMoodNegative || *mood == domain.JournalMoodNeutral || *mood == domain.JournalMoodPositive || *mood == domain.JournalMoodVeryPositive
}
