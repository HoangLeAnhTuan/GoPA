package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"gopa/internal/core/domain"
)

type vocabRepoStub struct {
	vocabs   map[uuid.UUID]domain.Vocabulary
	reviews  []domain.VocabularyReview
	sessions map[uuid.UUID]domain.LearningSession
}

func newVocabRepoStub() *vocabRepoStub {
	return &vocabRepoStub{
		vocabs:   make(map[uuid.UUID]domain.Vocabulary),
		sessions: make(map[uuid.UUID]domain.LearningSession),
	}
}

func (s *vocabRepoStub) List(_ context.Context, userID uuid.UUID, _ int) ([]domain.Vocabulary, error) {
	var list []domain.Vocabulary
	for _, v := range s.vocabs {
		if v.UserID == userID {
			list = append(list, v)
		}
	}
	return list, nil
}
func (s *vocabRepoStub) ListDue(_ context.Context, userID uuid.UUID, _ int, now time.Time) ([]domain.Vocabulary, error) {
	var list []domain.Vocabulary
	for _, v := range s.vocabs {
		if v.UserID == userID && !v.NextReviewAt.After(now) {
			list = append(list, v)
		}
	}
	return list, nil
}
func (s *vocabRepoStub) Get(_ context.Context, userID, id uuid.UUID) (domain.Vocabulary, error) {
	v, ok := s.vocabs[id]
	if !ok || v.UserID != userID {
		return domain.Vocabulary{}, domain.ErrNotFound
	}
	return v, nil
}
func (s *vocabRepoStub) Create(_ context.Context, v domain.Vocabulary) (domain.Vocabulary, error) {
	s.vocabs[v.ID] = v
	return v, nil
}
func (s *vocabRepoStub) BulkCreate(_ context.Context, userID uuid.UUID, items []domain.Vocabulary) ([]domain.Vocabulary, error) {
	for _, item := range items {
		item.UserID = userID
		s.vocabs[item.ID] = item
	}
	return items, nil
}
func (s *vocabRepoStub) Update(_ context.Context, v domain.Vocabulary) (domain.Vocabulary, error) {
	if _, ok := s.vocabs[v.ID]; !ok {
		return domain.Vocabulary{}, domain.ErrNotFound
	}
	s.vocabs[v.ID] = v
	return v, nil
}
func (s *vocabRepoStub) Delete(_ context.Context, userID, id uuid.UUID) error {
	v, ok := s.vocabs[id]
	if !ok || v.UserID != userID {
		return domain.ErrNotFound
	}
	delete(s.vocabs, id)
	return nil
}
func (s *vocabRepoStub) RecordReview(_ context.Context, review domain.VocabularyReview, updated domain.Vocabulary) error {
	s.reviews = append(s.reviews, review)
	s.vocabs[updated.ID] = updated
	return nil
}
func (s *vocabRepoStub) CreateLearningSession(_ context.Context, session domain.LearningSession) (domain.LearningSession, error) {
	s.sessions[session.ID] = session
	return session, nil
}
func (s *vocabRepoStub) UpdateLearningSession(_ context.Context, session domain.LearningSession) (domain.LearningSession, error) {
	s.sessions[session.ID] = session
	return session, nil
}
func (s *vocabRepoStub) GetStats(_ context.Context, _ uuid.UUID) (domain.VocabularyStats, error) {
	return domain.VocabularyStats{
		TotalWords:      10,
		MasteredWords:   4,
		DueCount:        2,
		BoxDistribution: map[int16]int{1: 3, 2: 3, 3: 2, 4: 2},
		ReviewHeatmap:   map[string]int{"2026-08-20": 5},
	}, nil
}

func TestVocabularyService_ReviewSM2(t *testing.T) {
	userID := uuid.New()
	repo := newVocabRepoStub()
	svc := NewVocabularyService(repo)

	reading := "ねこ"
	v, err := svc.Create(context.Background(), userID, VocabularyInput{
		Language: domain.VocabularyLanguageJapanese,
		Word:     "猫",
		Reading:  &reading,
		Meaning:  "Cat",
	})
	if err != nil {
		t.Fatalf("create vocab: %v", err)
	}

	// 1. Review with quality = 4
	if err := svc.Review(context.Background(), userID, v.ID, 4); err != nil {
		t.Fatalf("review: %v", err)
	}

	updated, err := svc.Get(context.Background(), userID, v.ID)
	if err != nil {
		t.Fatalf("get updated: %v", err)
	}
	if updated.RepetitionCount != 1 || updated.IntervalDays != 1 {
		t.Fatalf("expected 1 rep, 1 interval, got %d rep, %d interval", updated.RepetitionCount, updated.IntervalDays)
	}
	if len(repo.reviews) != 1 {
		t.Fatalf("expected 1 review recorded, got %d", len(repo.reviews))
	}

	// 2. Review with invalid quality
	if err := svc.Review(context.Background(), userID, v.ID, 6); !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("expected ErrValidation for quality 6, got %v", err)
	}
}

func TestVocabularyService_ImportVocabulariesJSON(t *testing.T) {
	userID := uuid.New()
	repo := newVocabRepoStub()
	svc := NewVocabularyService(repo)

	jsonData := []byte(`[
		{"word": "犬", "reading": "いぬ", "meaning": "Dog", "language": "JP", "difficulty_level": "BEGINNER"},
		{"word": "走る", "reading": "はしる", "meaning": "To run", "language": "JP"}
	]`)

	count, err := svc.ImportVocabularies(context.Background(), userID, jsonData, "json")
	if err != nil {
		t.Fatalf("import json: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected 2 imported items, got %d", count)
	}
}

func TestVocabularyService_ImportVocabulariesCSV(t *testing.T) {
	userID := uuid.New()
	repo := newVocabRepoStub()
	svc := NewVocabularyService(repo)

	csvData := []byte("word,meaning,reading,language,example,tags\n水,Water,みず,JP,水を飲む,nature;drink\n火,Fire,ひ,JP,,elements")

	count, err := svc.ImportVocabularies(context.Background(), userID, csvData, "csv")
	if err != nil {
		t.Fatalf("import csv: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected 2 imported items, got %d", count)
	}
}

func TestVocabularyService_LearningSessions(t *testing.T) {
	userID := uuid.New()
	repo := newVocabRepoStub()
	svc := NewVocabularyService(repo)

	// 1. Start Session
	session, err := svc.StartSession(context.Background(), userID, domain.VocabularyLanguageJapanese, "SRS_REVIEW")
	if err != nil {
		t.Fatalf("start session: %v", err)
	}
	if session.SessionType != "SRS_REVIEW" || session.Language != domain.VocabularyLanguageJapanese {
		t.Fatalf("unexpected session: %+v", session)
	}

	// 2. End Session
	ended, err := svc.EndSession(context.Background(), userID, session.ID, 10, 8, 300)
	if err != nil {
		t.Fatalf("end session: %v", err)
	}
	if ended.ItemsReviewed != 10 || ended.ItemsCorrect != 8 || ended.DurationSeconds != 300 || ended.EndedAt == nil {
		t.Fatalf("unexpected ended session: %+v", ended)
	}

	// 3. Validation errors
	_, err = svc.StartSession(context.Background(), userID, domain.VocabularyLanguage("ES"), "SRS_REVIEW")
	if !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("expected ErrValidation for invalid language, got %v", err)
	}

	_, err = svc.EndSession(context.Background(), userID, session.ID, 5, 8, 100)
	if !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("expected ErrValidation when correct > reviewed, got %v", err)
	}
}

func TestVocabularyService_GetStats(t *testing.T) {
	userID := uuid.New()
	repo := newVocabRepoStub()
	svc := NewVocabularyService(repo)

	stats, err := svc.GetStats(context.Background(), userID)
	if err != nil {
		t.Fatalf("get stats: %v", err)
	}
	if stats.TotalWords != 10 || stats.MasteredWords != 4 {
		t.Fatalf("unexpected stats: %+v", stats)
	}
}
