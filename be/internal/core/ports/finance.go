package ports

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gopa/internal/core/domain"
)

type TransactionFilter struct {
	AccountID  *uuid.UUID
	CategoryID *uuid.UUID
	Type       *domain.TransactionType
	From       *time.Time
	To         *time.Time
	Query      string
	Limit      int
}

type FinanceRepository interface {
	ListAccounts(context.Context, uuid.UUID, bool) ([]domain.Account, error)
	GetAccount(context.Context, uuid.UUID, uuid.UUID) (domain.Account, error)
	CreateAccount(context.Context, domain.Account) (domain.Account, error)
	UpdateAccount(context.Context, domain.Account) (domain.Account, error)
	DeleteAccount(context.Context, uuid.UUID, uuid.UUID) error
	ListCategories(context.Context, uuid.UUID) ([]domain.Category, error)
	GetCategory(context.Context, uuid.UUID, uuid.UUID) (domain.Category, error)
	CreateCategory(context.Context, domain.Category) (domain.Category, error)
	UpdateCategory(context.Context, domain.Category) (domain.Category, error)
	DeleteCategory(context.Context, uuid.UUID, uuid.UUID) error
	ListTransactions(context.Context, uuid.UUID, TransactionFilter) ([]domain.Transaction, error)
	GetTransaction(context.Context, uuid.UUID, uuid.UUID) (domain.Transaction, error)
	CreateTransaction(context.Context, domain.Transaction) (domain.Transaction, error)
	UpdateTransaction(context.Context, domain.Transaction) (domain.Transaction, error)
	DeleteTransaction(context.Context, uuid.UUID, uuid.UUID, time.Time) error
	Cashflow(context.Context, uuid.UUID, time.Time, time.Time) (domain.CashflowSummary, error)
	SpendingByCategory(context.Context, uuid.UUID, time.Time, time.Time) ([]domain.CategorySpending, error)
	NetWorth(context.Context, uuid.UUID) (decimal.Decimal, error)
}
