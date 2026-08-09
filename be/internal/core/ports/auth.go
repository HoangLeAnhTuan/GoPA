package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gopa/internal/core/domain"
)

type UserRepository interface {
	Create(context.Context, domain.User) (domain.User, error)
	FindByID(context.Context, uuid.UUID) (domain.User, error)
	FindByEmail(context.Context, string) (domain.User, error)
}

type RefreshSessionStore interface {
	Save(context.Context, uuid.UUID, uuid.UUID, string, time.Duration) error
	Rotate(context.Context, uuid.UUID, uuid.UUID, string, uuid.UUID, string, time.Duration) (bool, error)
	Delete(context.Context, uuid.UUID, uuid.UUID) error
}

type RateLimiter interface {
	Allow(context.Context, string, int, time.Duration) (bool, error)
}
