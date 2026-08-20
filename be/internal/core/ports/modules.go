package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gopa/internal/core/domain"
)

type TaskRepository interface {
	List(context.Context, uuid.UUID) ([]domain.Task, error)
	Create(context.Context, domain.Task) (domain.Task, error)
	Get(context.Context, uuid.UUID, uuid.UUID) (domain.Task, error)
	Update(context.Context, domain.Task) (domain.Task, error)
	UpdateStatus(ctx context.Context, userID, id uuid.UUID, status domain.TaskStatus, targetPosition *int64) (domain.Task, error)
	Delete(context.Context, uuid.UUID, uuid.UUID) error
	NextPosition(context.Context, uuid.UUID, domain.TaskStatus) (int64, error)
}

type VocabularyRepository interface {
	List(context.Context, uuid.UUID, int) ([]domain.Vocabulary, error)
	ListDue(context.Context, uuid.UUID, int, time.Time) ([]domain.Vocabulary, error)
	Get(context.Context, uuid.UUID, uuid.UUID) (domain.Vocabulary, error)
	Create(context.Context, domain.Vocabulary) (domain.Vocabulary, error)
	BulkCreate(context.Context, uuid.UUID, []domain.Vocabulary) ([]domain.Vocabulary, error)
	Update(context.Context, domain.Vocabulary) (domain.Vocabulary, error)
	Delete(context.Context, uuid.UUID, uuid.UUID) error
	RecordReview(context.Context, domain.VocabularyReview, domain.Vocabulary) error
	CreateLearningSession(context.Context, domain.LearningSession) (domain.LearningSession, error)
	UpdateLearningSession(context.Context, domain.LearningSession) (domain.LearningSession, error)
	GetStats(context.Context, uuid.UUID) (domain.VocabularyStats, error)
}

type JournalRepository interface {
	List(context.Context, uuid.UUID, domain.JournalFilter) ([]domain.Journal, error)
	Search(ctx context.Context, userID uuid.UUID, query string) ([]domain.Journal, error)
	Get(context.Context, uuid.UUID, uuid.UUID) (domain.Journal, error)
	Create(context.Context, domain.Journal) (domain.Journal, error)
	Update(context.Context, domain.Journal) (domain.Journal, error)
	Delete(context.Context, uuid.UUID, uuid.UUID) error
	Link(context.Context, uuid.UUID, uuid.UUID, uuid.UUID) error
	Unlink(context.Context, uuid.UUID, uuid.UUID, uuid.UUID) error
	Stats(context.Context, uuid.UUID) (domain.JournalStats, error)
}

type PomodoroStore interface {
	Get(context.Context, uuid.UUID) (domain.PomodoroState, error)
	Set(context.Context, uuid.UUID, domain.PomodoroState) error
	Delete(context.Context, uuid.UUID) error
}

type PomodoroHistoryRepository interface {
	Create(context.Context, domain.PomodoroHistory) error
	List(ctx context.Context, userID uuid.UUID, limit int, cursor string) ([]domain.PomodoroHistory, error)
}
