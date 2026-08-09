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
	Delete(context.Context, uuid.UUID, uuid.UUID) error
	NextPosition(context.Context, uuid.UUID, domain.TaskStatus) (int64, error)
}

type VocabularyRepository interface {
	List(context.Context, uuid.UUID, int) ([]domain.Vocabulary, error)
	ListDue(context.Context, uuid.UUID, int, time.Time) ([]domain.Vocabulary, error)
	Get(context.Context, uuid.UUID, uuid.UUID) (domain.Vocabulary, error)
	Create(context.Context, domain.Vocabulary) (domain.Vocabulary, error)
	Update(context.Context, domain.Vocabulary) (domain.Vocabulary, error)
	Delete(context.Context, uuid.UUID, uuid.UUID) error
	RecordReview(context.Context, uuid.UUID, uuid.UUID, int16, time.Time, domain.ReviewSchedule) error
}

type VehicleRepository interface {
	List(context.Context, uuid.UUID) ([]domain.Vehicle, error)
	Get(context.Context, uuid.UUID, uuid.UUID) (domain.Vehicle, error)
	Create(context.Context, domain.Vehicle) (domain.Vehicle, error)
	Update(context.Context, domain.Vehicle) (domain.Vehicle, error)
	ListLogs(context.Context, uuid.UUID, uuid.UUID) ([]domain.VehicleLog, error)
	CreateLog(context.Context, domain.VehicleLog) (domain.VehicleLog, error)
	LatestMaintenanceMileage(context.Context, uuid.UUID, uuid.UUID) (*int, error)
}

type NetworkRepository interface {
	List(context.Context, uuid.UUID) ([]domain.NetworkNode, error)
	Get(context.Context, uuid.UUID, uuid.UUID) (domain.NetworkNode, error)
	Create(context.Context, domain.NetworkNode) (domain.NetworkNode, error)
	Update(context.Context, domain.NetworkNode) (domain.NetworkNode, error)
	Delete(context.Context, uuid.UUID, uuid.UUID) error
}

type JournalRepository interface {
	List(context.Context, uuid.UUID, string, string, int) ([]domain.Journal, error)
	Get(context.Context, uuid.UUID, uuid.UUID) (domain.Journal, error)
	Create(context.Context, domain.Journal) (domain.Journal, error)
	Update(context.Context, domain.Journal) (domain.Journal, error)
	Delete(context.Context, uuid.UUID, uuid.UUID) error
}

type PomodoroStore interface {
	Get(context.Context, uuid.UUID) (domain.PomodoroState, error)
	Set(context.Context, uuid.UUID, domain.PomodoroState) error
	Delete(context.Context, uuid.UUID) error
}

type PomodoroHistoryRepository interface {
	Create(context.Context, domain.PomodoroHistory) error
}
