package services

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/google/uuid"
	"gopa/internal/core/domain"
	"gopa/internal/core/domain/linguistics"
	"gopa/internal/core/ports"
)

type VocabularyInput struct {
	Language           domain.VocabularyLanguage
	Word               string
	Reading            *string
	Meaning            string
	ExampleSentence    *string
	ExampleTranslation *string
	Tags               []string
	DifficultyLevel    string
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
	diff := input.DifficultyLevel
	if diff == "" {
		diff = "BEGINNER"
	}
	vocabulary := domain.Vocabulary{
		ID:                 uuid.New(),
		UserID:             userID,
		Language:           input.Language,
		Word:               strings.TrimSpace(input.Word),
		Reading:            input.Reading,
		Meaning:            strings.TrimSpace(input.Meaning),
		ExampleSentence:    input.ExampleSentence,
		ExampleTranslation: input.ExampleTranslation,
		Tags:               append([]string(nil), input.Tags...),
		DifficultyLevel:    diff,
		EaseFactor:         2.50,
		IntervalDays:       0,
		RepetitionCount:    0,
		CurrentBox:         1,
		MasteryScore:       0,
		NextReviewAt:       now,
		CreatedAt:          now,
		UpdatedAt:          now,
	}
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
	vocabulary.ExampleTranslation = input.ExampleTranslation
	if input.DifficultyLevel != "" {
		vocabulary.DifficultyLevel = input.DifficultyLevel
	}
	if input.Tags != nil {
		vocabulary.Tags = append([]string(nil), input.Tags...)
	}
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
	if quality < 0 || quality > 5 {
		return domain.ErrValidation
	}
	now := s.now().UTC()
	result := linguistics.CalculateNextReview(linguistics.VocabularyState{
		CurrentBox:      vocabulary.CurrentBox,
		EaseFactor:      vocabulary.EaseFactor,
		IntervalDays:    vocabulary.IntervalDays,
		RepetitionCount: vocabulary.RepetitionCount,
		MasteryScore:    vocabulary.MasteryScore,
	}, quality, now)

	vocabulary.CurrentBox = result.NextBox
	vocabulary.EaseFactor = result.EaseFactor
	vocabulary.IntervalDays = result.IntervalDays
	vocabulary.RepetitionCount = result.RepetitionCount
	vocabulary.MasteryScore = result.MasteryScore
	vocabulary.NextReviewAt = result.NextReviewAt

	review := domain.VocabularyReview{
		ID:             uuid.New(),
		EventID:        uuid.New(),
		VocabularyID:   id,
		UserID:         userID,
		Quality:        quality,
		StudyMode:      "FLASHCARD",
		ResponseTimeMs: 0,
		WasCorrect:     quality >= 3,
		ReviewedAt:     now,
		CreatedAt:      now,
	}
	return s.vocabularies.RecordReview(ctx, review, vocabulary)
}

type vocabularyImportDTO struct {
	Language           domain.VocabularyLanguage `json:"language"`
	Word               string                    `json:"word"`
	Reading            *string                   `json:"reading"`
	Meaning            string                    `json:"meaning"`
	ExampleSentence    *string                   `json:"example_sentence"`
	ExampleTranslation *string                   `json:"example_translation"`
	Tags               []string                  `json:"tags"`
	DifficultyLevel    string                    `json:"difficulty_level"`
}

func (s *VocabularyService) ImportVocabularies(ctx context.Context, userID uuid.UUID, data []byte, format string) (int, error) {
	if len(data) == 0 {
		return 0, domain.ErrValidation
	}
	now := s.now().UTC()
	var dtos []vocabularyImportDTO

	switch strings.ToLower(strings.TrimSpace(format)) {
	case "csv":
		reader := csv.NewReader(bytes.NewReader(data))
		// Optional header skip
		firstRow := true
		for {
			record, err := reader.Read()
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				return 0, fmt.Errorf("%w: parse csv: %v", domain.ErrValidation, err)
			}
			if len(record) == 0 {
				continue
			}
			if firstRow {
				firstRow = false
				// Check if first row is header
				lower := strings.ToLower(record[0])
				if lower == "word" || lower == "vocabulary" || lower == "language" {
					continue
				}
			}
			dto := vocabularyImportDTO{}
			// Expected columns: word, meaning, [reading], [language], [example], [tags]
			if len(record) >= 1 {
				dto.Word = record[0]
			}
			if len(record) >= 2 {
				dto.Meaning = record[1]
			}
			if len(record) >= 3 && strings.TrimSpace(record[2]) != "" {
				r := strings.TrimSpace(record[2])
				dto.Reading = &r
			}
			if len(record) >= 4 && strings.TrimSpace(record[3]) != "" {
				dto.Language = domain.VocabularyLanguage(strings.ToUpper(strings.TrimSpace(record[3])))
			}
			if len(record) >= 5 && strings.TrimSpace(record[4]) != "" {
				ex := strings.TrimSpace(record[4])
				dto.ExampleSentence = &ex
			}
			if len(record) >= 6 && strings.TrimSpace(record[5]) != "" {
				tags := strings.Split(record[5], ";")
				for _, t := range tags {
					if st := strings.TrimSpace(t); st != "" {
						dto.Tags = append(dto.Tags, st)
					}
				}
			}
			dtos = append(dtos, dto)
		}
	default: // JSON by default
		if err := json.Unmarshal(data, &dtos); err != nil {
			return 0, fmt.Errorf("%w: parse json: %v", domain.ErrValidation, err)
		}
	}

	if len(dtos) == 0 {
		return 0, domain.ErrValidation
	}

	items := make([]domain.Vocabulary, 0, len(dtos))
	for _, dto := range dtos {
		lang := dto.Language
		if lang == "" {
			lang = domain.VocabularyLanguageJapanese
		}
		diff := dto.DifficultyLevel
		if diff == "" {
			diff = "BEGINNER"
		}
		vocab := domain.Vocabulary{
			ID:                 uuid.New(),
			UserID:             userID,
			Language:           lang,
			Word:               strings.TrimSpace(dto.Word),
			Reading:            dto.Reading,
			Meaning:            strings.TrimSpace(dto.Meaning),
			ExampleSentence:    dto.ExampleSentence,
			ExampleTranslation: dto.ExampleTranslation,
			Tags:               append([]string(nil), dto.Tags...),
			DifficultyLevel:    diff,
			EaseFactor:         2.50,
			IntervalDays:       0,
			RepetitionCount:    0,
			CurrentBox:         1,
			MasteryScore:       0,
			NextReviewAt:       now,
			CreatedAt:          now,
			UpdatedAt:          now,
		}
		if err := vocab.Validate(); err != nil {
			return 0, err
		}
		items = append(items, vocab)
	}

	created, err := s.vocabularies.BulkCreate(ctx, userID, items)
	if err != nil {
		return 0, err
	}
	return len(created), nil
}

func (s *VocabularyService) StartSession(ctx context.Context, userID uuid.UUID, language domain.VocabularyLanguage, sessionType string) (domain.LearningSession, error) {
	if language != domain.VocabularyLanguageJapanese && language != domain.VocabularyLanguageEnglish {
		return domain.LearningSession{}, domain.ErrValidation
	}
	st := strings.TrimSpace(sessionType)
	if st == "" {
		st = "SRS_REVIEW"
	}
	now := s.now().UTC()
	session := domain.LearningSession{
		ID:              uuid.New(),
		UserID:          userID,
		Language:        language,
		SessionType:     st,
		ItemsReviewed:   0,
		ItemsCorrect:    0,
		DurationSeconds: 0,
		StartedAt:       now,
	}
	return s.vocabularies.CreateLearningSession(ctx, session)
}

func (s *VocabularyService) EndSession(ctx context.Context, userID, sessionID uuid.UUID, itemsReviewed, itemsCorrect, durationSeconds int) (domain.LearningSession, error) {
	if itemsReviewed < 0 || itemsCorrect < 0 || itemsCorrect > itemsReviewed || durationSeconds < 0 {
		return domain.LearningSession{}, domain.ErrValidation
	}
	endedAt := s.now().UTC()
	session := domain.LearningSession{
		ID:              sessionID,
		UserID:          userID,
		ItemsReviewed:   itemsReviewed,
		ItemsCorrect:    itemsCorrect,
		DurationSeconds: durationSeconds,
		EndedAt:         &endedAt,
	}
	return s.vocabularies.UpdateLearningSession(ctx, session)
}

func (s *VocabularyService) GetStats(ctx context.Context, userID uuid.UUID) (domain.VocabularyStats, error) {
	return s.vocabularies.GetStats(ctx, userID)
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
