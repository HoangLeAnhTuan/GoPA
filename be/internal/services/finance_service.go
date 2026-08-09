package services

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gopa/internal/core/domain"
	"gopa/internal/core/ports"
)

type AccountInput struct {
	Name           string
	Type           domain.AccountType
	Currency       string
	InitialBalance decimal.Decimal
	Color          string
	Icon           string
	IsArchived     bool
}

type CategoryInput struct {
	Name          string
	Type          domain.CategoryType
	Icon          string
	Color         string
	ParentID      *uuid.UUID
	MonthlyBudget *decimal.Decimal
}

type TransactionInput struct {
	AccountID    uuid.UUID
	ToAccountID  *uuid.UUID
	CategoryID   *uuid.UUID
	Type         domain.TransactionType
	Amount       decimal.Decimal
	Currency     string
	ExchangeRate decimal.Decimal
	Description  string
	Merchant     *string
	Tags         []string
	OccurredAt   time.Time
	IsRecurring  bool
}

type FinanceService struct {
	repository ports.FinanceRepository
	now        func() time.Time
}

func NewFinanceService(repository ports.FinanceRepository) *FinanceService {
	return &FinanceService{repository: repository, now: time.Now}
}

func (s *FinanceService) ListAccounts(ctx context.Context, userID uuid.UUID, includeArchived bool) ([]domain.Account, error) {
	return s.repository.ListAccounts(ctx, userID, includeArchived)
}

func (s *FinanceService) CreateAccount(ctx context.Context, userID uuid.UUID, input AccountInput) (domain.Account, error) {
	now := s.now().UTC()
	account := domain.Account{ID: uuid.New(), UserID: userID, Name: strings.TrimSpace(input.Name), Type: input.Type, Currency: strings.ToUpper(strings.TrimSpace(input.Currency)), InitialBalance: input.InitialBalance, CurrentBalance: input.InitialBalance, Color: defaultString(input.Color, "#10B981"), Icon: defaultString(input.Icon, "wallet"), IsArchived: input.IsArchived, CreatedAt: now, UpdatedAt: now}
	if err := account.Validate(); err != nil {
		return domain.Account{}, err
	}
	return s.repository.CreateAccount(ctx, account)
}

func (s *FinanceService) UpdateAccount(ctx context.Context, userID, id uuid.UUID, input AccountInput) (domain.Account, error) {
	account, err := s.repository.GetAccount(ctx, userID, id)
	if err != nil {
		return domain.Account{}, err
	}
	account.Name = strings.TrimSpace(input.Name)
	account.Type = input.Type
	account.Currency = strings.ToUpper(strings.TrimSpace(input.Currency))
	account.Color = defaultString(input.Color, "#10B981")
	account.Icon = defaultString(input.Icon, "wallet")
	account.IsArchived = input.IsArchived
	account.UpdatedAt = s.now().UTC()
	if err := account.Validate(); err != nil {
		return domain.Account{}, err
	}
	return s.repository.UpdateAccount(ctx, account)
}

func (s *FinanceService) DeleteAccount(ctx context.Context, userID, id uuid.UUID) error {
	return s.repository.DeleteAccount(ctx, userID, id)
}

func (s *FinanceService) ListCategories(ctx context.Context, userID uuid.UUID) ([]domain.Category, error) {
	return s.repository.ListCategories(ctx, userID)
}

func (s *FinanceService) CreateCategory(ctx context.Context, userID uuid.UUID, input CategoryInput) (domain.Category, error) {
	now := s.now().UTC()
	category := domain.Category{ID: uuid.New(), UserID: userID, Name: strings.TrimSpace(input.Name), Type: input.Type, Icon: defaultString(input.Icon, "tag"), Color: defaultString(input.Color, "#64748B"), ParentID: input.ParentID, MonthlyBudget: input.MonthlyBudget, CreatedAt: now, UpdatedAt: now}
	if err := category.Validate(); err != nil {
		return domain.Category{}, err
	}
	return s.repository.CreateCategory(ctx, category)
}

func (s *FinanceService) UpdateCategory(ctx context.Context, userID, id uuid.UUID, input CategoryInput) (domain.Category, error) {
	category, err := s.repository.GetCategory(ctx, userID, id)
	if err != nil {
		return domain.Category{}, err
	}
	category.Name = strings.TrimSpace(input.Name)
	category.Type = input.Type
	category.Icon = defaultString(input.Icon, "tag")
	category.Color = defaultString(input.Color, "#64748B")
	category.ParentID = input.ParentID
	category.MonthlyBudget = input.MonthlyBudget
	category.UpdatedAt = s.now().UTC()
	if err := category.Validate(); err != nil {
		return domain.Category{}, err
	}
	return s.repository.UpdateCategory(ctx, category)
}

func (s *FinanceService) DeleteCategory(ctx context.Context, userID, id uuid.UUID) error {
	return s.repository.DeleteCategory(ctx, userID, id)
}

func (s *FinanceService) ListTransactions(ctx context.Context, userID uuid.UUID, filter ports.TransactionFilter) ([]domain.Transaction, error) {
	return s.repository.ListTransactions(ctx, userID, filter)
}

func (s *FinanceService) CreateTransaction(ctx context.Context, userID uuid.UUID, input TransactionInput) (domain.Transaction, error) {
	now := s.now().UTC()
	transaction, err := s.newTransaction(userID, uuid.New(), input, now, now)
	if err != nil {
		return domain.Transaction{}, err
	}
	return s.repository.CreateTransaction(ctx, transaction)
}

func (s *FinanceService) UpdateTransaction(ctx context.Context, userID, id uuid.UUID, input TransactionInput) (domain.Transaction, error) {
	current, err := s.repository.GetTransaction(ctx, userID, id)
	if err != nil {
		return domain.Transaction{}, err
	}
	transaction, err := s.newTransaction(userID, id, input, current.CreatedAt, s.now().UTC())
	if err != nil {
		return domain.Transaction{}, err
	}
	return s.repository.UpdateTransaction(ctx, transaction)
}

func (s *FinanceService) GetTransaction(ctx context.Context, userID, id uuid.UUID) (domain.Transaction, error) {
	return s.repository.GetTransaction(ctx, userID, id)
}

func (s *FinanceService) newTransaction(userID, id uuid.UUID, input TransactionInput, createdAt, updatedAt time.Time) (domain.Transaction, error) {
	transaction := domain.Transaction{ID: id, UserID: userID, AccountID: input.AccountID, ToAccountID: input.ToAccountID, CategoryID: input.CategoryID, Type: input.Type, Amount: input.Amount, Currency: strings.ToUpper(strings.TrimSpace(input.Currency)), ExchangeRate: input.ExchangeRate, Description: strings.TrimSpace(input.Description), Merchant: trimOptional(input.Merchant), Tags: append([]string(nil), input.Tags...), OccurredAt: input.OccurredAt.UTC(), IsRecurring: input.IsRecurring, CreatedAt: createdAt.UTC(), UpdatedAt: updatedAt.UTC()}
	if err := transaction.Validate(); err != nil {
		return domain.Transaction{}, err
	}
	return transaction, nil
}

func (s *FinanceService) DeleteTransaction(ctx context.Context, userID, id uuid.UUID) error {
	return s.repository.DeleteTransaction(ctx, userID, id, s.now().UTC())
}

func (s *FinanceService) Dashboard(ctx context.Context, userID uuid.UUID) (domain.FinanceDashboard, error) {
	now := s.now().UTC()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	nextMonth := monthStart.AddDate(0, 1, 0)
	netWorth, err := s.repository.NetWorth(ctx, userID)
	if err != nil {
		return domain.FinanceDashboard{}, err
	}
	cashflow, err := s.repository.Cashflow(ctx, userID, monthStart, nextMonth)
	if err != nil {
		return domain.FinanceDashboard{}, err
	}
	spending, err := s.repository.SpendingByCategory(ctx, userID, monthStart, nextMonth)
	if err != nil {
		return domain.FinanceDashboard{}, err
	}
	recent, err := s.repository.ListTransactions(ctx, userID, ports.TransactionFilter{Limit: 10})
	if err != nil {
		return domain.FinanceDashboard{}, err
	}
	return domain.FinanceDashboard{NetWorth: netWorth, Cashflow: cashflow, Spending: spending, RecentTransactions: recent}, nil
}

func (s *FinanceService) Cashflow(ctx context.Context, userID uuid.UUID, from, to time.Time) (domain.CashflowSummary, error) {
	if !from.Before(to) {
		return domain.CashflowSummary{}, domain.ErrValidation
	}
	return s.repository.Cashflow(ctx, userID, from, to)
}

func (s *FinanceService) SpendingByCategory(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]domain.CategorySpending, error) {
	if !from.Before(to) {
		return nil, domain.ErrValidation
	}
	return s.repository.SpendingByCategory(ctx, userID, from, to)
}

func (s *FinanceService) NetWorth(ctx context.Context, userID uuid.UUID) (decimal.Decimal, error) {
	return s.repository.NetWorth(ctx, userID)
}

func defaultString(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return strings.TrimSpace(value)
}

func trimOptional(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
