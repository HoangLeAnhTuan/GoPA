package services

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"gopa/internal/core/domain"
	"gopa/internal/core/ports"
)

type VocabularyInput struct {
	Language        domain.VocabularyLanguage
	Word            string
	Reading         *string
	Meaning         string
	ExampleSentence *string
}
type VocabularyService struct {
	vocabularies ports.VocabularyRepository
	now          func() time.Time
}

func NewVocabularyService(vocabularies ports.VocabularyRepository) *VocabularyService {
	return &VocabularyService{vocabularies: vocabularies, now: time.Now}
}
func (s *VocabularyService) List(ctx context.Context, userID uuid.UUID, limit int) ([]domain.Vocabulary, error) {
	return s.vocabularies.List(ctx, userID, boundedLimit(limit))
}
func (s *VocabularyService) Due(ctx context.Context, userID uuid.UUID, limit int) ([]domain.Vocabulary, error) {
	return s.vocabularies.ListDue(ctx, userID, boundedLimit(limit), s.now().UTC())
}
func (s *VocabularyService) Get(ctx context.Context, userID, id uuid.UUID) (domain.Vocabulary, error) {
	return s.vocabularies.Get(ctx, userID, id)
}
func (s *VocabularyService) Create(ctx context.Context, userID uuid.UUID, input VocabularyInput) (domain.Vocabulary, error) {
	now := s.now().UTC()
	vocabulary := domain.Vocabulary{ID: uuid.New(), UserID: userID, Language: input.Language, Word: strings.TrimSpace(input.Word), Reading: input.Reading, Meaning: strings.TrimSpace(input.Meaning), ExampleSentence: input.ExampleSentence, CurrentBox: 1, NextReviewAt: now, CreatedAt: now, UpdatedAt: now}
	if err := vocabulary.Validate(); err != nil {
		return domain.Vocabulary{}, err
	}
	return s.vocabularies.Create(ctx, vocabulary)
}
func (s *VocabularyService) Update(ctx context.Context, userID, id uuid.UUID, input VocabularyInput) (domain.Vocabulary, error) {
	vocabulary, err := s.vocabularies.Get(ctx, userID, id)
	if err != nil {
		return domain.Vocabulary{}, err
	}
	vocabulary.Language = input.Language
	vocabulary.Word = strings.TrimSpace(input.Word)
	vocabulary.Reading = input.Reading
	vocabulary.Meaning = strings.TrimSpace(input.Meaning)
	vocabulary.ExampleSentence = input.ExampleSentence
	vocabulary.UpdatedAt = s.now().UTC()
	if err := vocabulary.Validate(); err != nil {
		return domain.Vocabulary{}, err
	}
	return s.vocabularies.Update(ctx, vocabulary)
}
func (s *VocabularyService) Delete(ctx context.Context, userID, id uuid.UUID) error {
	return s.vocabularies.Delete(ctx, userID, id)
}
func (s *VocabularyService) Review(ctx context.Context, userID, id uuid.UUID, quality int16) error {
	vocabulary, err := s.vocabularies.Get(ctx, userID, id)
	if err != nil {
		return err
	}
	now := s.now().UTC()
	schedule, err := domain.ScheduleReview(vocabulary.CurrentBox, quality, now)
	if err != nil {
		return err
	}
	return s.vocabularies.RecordReview(ctx, userID, id, quality, now, schedule)
}
func boundedLimit(value int) int {
	if value <= 0 {
		return 20
	}
	if value > 100 {
		return 100
	}
	return value
}
