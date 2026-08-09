package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AccountType string

const (
	AccountCash       AccountType = "CASH"
	AccountBank       AccountType = "BANK"
	AccountCreditCard AccountType = "CREDIT_CARD"
	AccountSavings    AccountType = "SAVINGS"
	AccountInvestment AccountType = "INVESTMENT"
	AccountCrypto     AccountType = "CRYPTO"
)

type Account struct {
	ID             uuid.UUID       `json:"id" db:"id"`
	UserID         uuid.UUID       `json:"user_id" db:"user_id"`
	Name           string          `json:"name" db:"name"`
	Type           AccountType     `json:"type" db:"type"`
	Currency       string          `json:"currency" db:"currency"`
	InitialBalance decimal.Decimal `json:"initial_balance" db:"initial_balance"`
	CurrentBalance decimal.Decimal `json:"current_balance" db:"current_balance"`
	Color          string          `json:"color" db:"color"`
	Icon           string          `json:"icon" db:"icon"`
	IsArchived     bool            `json:"is_archived" db:"is_archived"`
	CreatedAt      time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at" db:"updated_at"`
}

func (a Account) Validate() error {
	if len(strings.TrimSpace(a.Name)) == 0 || len(a.Name) > 120 || !validAccountType(a.Type) || !validCurrency(a.Currency) || a.InitialBalance.IsNegative() || !validHexColor(a.Color) || len(strings.TrimSpace(a.Icon)) == 0 || len(a.Icon) > 50 {
		return ErrValidation
	}
	return nil
}

func validAccountType(value AccountType) bool {
	return value == AccountCash || value == AccountBank || value == AccountCreditCard || value == AccountSavings || value == AccountInvestment || value == AccountCrypto
}

type CategoryType string

const (
	CategoryIncome  CategoryType = "INCOME"
	CategoryExpense CategoryType = "EXPENSE"
)

type Category struct {
	ID            uuid.UUID        `json:"id" db:"id"`
	UserID        uuid.UUID        `json:"user_id" db:"user_id"`
	Name          string           `json:"name" db:"name"`
	Type          CategoryType     `json:"type" db:"type"`
	Icon          string           `json:"icon" db:"icon"`
	Color         string           `json:"color" db:"color"`
	ParentID      *uuid.UUID       `json:"parent_id" db:"parent_id"`
	MonthlyBudget *decimal.Decimal `json:"monthly_budget" db:"monthly_budget"`
	CreatedAt     time.Time        `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at" db:"updated_at"`
}

func (c Category) Validate() error {
	if len(strings.TrimSpace(c.Name)) == 0 || len(c.Name) > 100 || (c.Type != CategoryIncome && c.Type != CategoryExpense) || len(strings.TrimSpace(c.Icon)) == 0 || len(c.Icon) > 50 || !validHexColor(c.Color) || (c.ParentID != nil && *c.ParentID == c.ID) || (c.MonthlyBudget != nil && !c.MonthlyBudget.IsPositive()) {
		return ErrValidation
	}
	return nil
}

type TransactionType string

const (
	TransactionIncome   TransactionType = "INCOME"
	TransactionExpense  TransactionType = "EXPENSE"
	TransactionTransfer TransactionType = "TRANSFER"
)

type Transaction struct {
	ID            uuid.UUID       `json:"id" db:"id"`
	UserID        uuid.UUID       `json:"user_id" db:"user_id"`
	AccountID     uuid.UUID       `json:"account_id" db:"account_id"`
	ToAccountID   *uuid.UUID      `json:"to_account_id" db:"to_account_id"`
	CategoryID    *uuid.UUID      `json:"category_id" db:"category_id"`
	Type          TransactionType `json:"type" db:"type"`
	Amount        decimal.Decimal `json:"amount" db:"amount"`
	Currency      string          `json:"currency" db:"currency"`
	ExchangeRate  decimal.Decimal `json:"exchange_rate" db:"exchange_rate"`
	Description   string          `json:"description" db:"description"`
	Merchant      *string         `json:"merchant" db:"merchant"`
	Tags          []string        `json:"tags" db:"tags"`
	OccurredAt    time.Time       `json:"occurred_at" db:"occurred_at"`
	IsRecurring   bool            `json:"is_recurring" db:"is_recurring"`
	RecurringRule []byte          `json:"-" db:"recurring_rule"`
	CreatedAt     time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at" db:"updated_at"`
}

func (t *Transaction) Validate() error {
	if !t.Amount.IsPositive() || t.Amount.GreaterThan(decimal.NewFromInt(999999999999)) || !t.ExchangeRate.IsPositive() || !validCurrency(t.Currency) || len(t.Description) > 500 || (t.Merchant != nil && len(*t.Merchant) > 150) || t.OccurredAt.IsZero() {
		return ErrValidation
	}
	if t.Type != TransactionIncome && t.Type != TransactionExpense && t.Type != TransactionTransfer {
		return ErrValidation
	}
	if t.Type == TransactionTransfer {
		if t.ToAccountID == nil || *t.ToAccountID == t.AccountID || t.CategoryID != nil {
			return ErrValidation
		}
	} else if t.ToAccountID != nil {
		return ErrValidation
	}
	return normalizeTags(&t.Tags)
}

type CashflowSummary struct {
	Income      decimal.Decimal `json:"income"`
	Expense     decimal.Decimal `json:"expense"`
	Net         decimal.Decimal `json:"net"`
	SavingsRate decimal.Decimal `json:"savings_rate"`
}

type CategorySpending struct {
	CategoryID   *uuid.UUID      `json:"category_id" db:"category_id"`
	CategoryName string          `json:"category_name" db:"category_name"`
	Color        string          `json:"color" db:"color"`
	Amount       decimal.Decimal `json:"amount" db:"amount"`
}

type FinanceDashboard struct {
	NetWorth           decimal.Decimal    `json:"net_worth"`
	Cashflow           CashflowSummary    `json:"cashflow"`
	Spending           []CategorySpending `json:"spending"`
	RecentTransactions []Transaction      `json:"recent_transactions"`
}

func validCurrency(value string) bool {
	if len(value) != 3 {
		return false
	}
	for _, character := range value {
		if character < 'A' || character > 'Z' {
			return false
		}
	}
	return true
}

func validHexColor(value string) bool {
	if len(value) != 7 || value[0] != '#' {
		return false
	}
	for _, character := range value[1:] {
		if !(character >= '0' && character <= '9') && !(character >= 'a' && character <= 'f') && !(character >= 'A' && character <= 'F') {
			return false
		}
	}
	return true
}

func normalizeTags(tags *[]string) error {
	if len(*tags) > 20 {
		return ErrValidation
	}
	seen := make(map[string]struct{}, len(*tags))
	for index, raw := range *tags {
		value := strings.TrimSpace(strings.ToLower(raw))
		if value == "" || len(value) > 50 {
			return ErrValidation
		}
		if _, exists := seen[value]; exists {
			return ErrValidation
		}
		seen[value] = struct{}{}
		(*tags)[index] = value
	}
	return nil
}
