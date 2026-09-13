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

// ConvertedAmount calculates the transaction amount in the account's base currency using the exchange rate.
func (t Transaction) ConvertedAmount() decimal.Decimal {
	return t.Amount.Mul(t.ExchangeRate)
}

type BudgetPeriod string

const (
	BudgetPeriodWeekly    BudgetPeriod = "WEEKLY"
	BudgetPeriodMonthly   BudgetPeriod = "MONTHLY"
	BudgetPeriodQuarterly BudgetPeriod = "QUARTERLY"
	BudgetPeriodYearly    BudgetPeriod = "YEARLY"
	BudgetPeriodCustom    BudgetPeriod = "CUSTOM"
)

type Budget struct {
	ID             uuid.UUID       `json:"id" db:"id"`
	UserID         uuid.UUID       `json:"user_id" db:"user_id"`
	CategoryID     *uuid.UUID      `json:"category_id,omitempty" db:"category_id"`
	Name           string          `json:"name" db:"name"`
	Amount         decimal.Decimal `json:"amount" db:"amount"`
	Period         BudgetPeriod    `json:"period" db:"period"`
	StartDate      time.Time       `json:"start_date" db:"start_date"`
	EndDate        time.Time       `json:"end_date" db:"end_date"`
	AlertThreshold decimal.Decimal `json:"alert_threshold" db:"alert_threshold"`
	IsActive       bool            `json:"is_active" db:"is_active"`
	CreatedAt      time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at" db:"updated_at"`
}

func (b Budget) Validate() error {
	if len(strings.TrimSpace(b.Name)) == 0 || len(b.Name) > 150 || !b.Amount.IsPositive() || !validBudgetPeriod(b.Period) || b.StartDate.IsZero() || b.EndDate.IsZero() || b.EndDate.Before(b.StartDate) || !b.AlertThreshold.IsPositive() || b.AlertThreshold.GreaterThan(decimal.NewFromInt(1)) {
		return ErrValidation
	}
	return nil
}

func validBudgetPeriod(p BudgetPeriod) bool {
	return p == BudgetPeriodWeekly || p == BudgetPeriodMonthly || p == BudgetPeriodQuarterly || p == BudgetPeriodYearly || p == BudgetPeriodCustom
}

type BudgetStatus struct {
	Budget           Budget          `json:"budget"`
	SpentAmount      decimal.Decimal `json:"spent_amount"`
	UtilizationRate  decimal.Decimal `json:"utilization_rate"`
	IsAlertTriggered bool            `json:"is_alert_triggered"`
}

type SavingsGoal struct {
	ID              uuid.UUID       `json:"id" db:"id"`
	UserID          uuid.UUID       `json:"user_id" db:"user_id"`
	Name            string          `json:"name" db:"name"`
	TargetAmount    decimal.Decimal `json:"target_amount" db:"target_amount"`
	CurrentAmount   decimal.Decimal `json:"current_amount" db:"current_amount"`
	LinkedAccountID *uuid.UUID      `json:"linked_account_id,omitempty" db:"linked_account_id"`
	TargetDate      *time.Time      `json:"target_date,omitempty" db:"target_date"`
	Color           string          `json:"color" db:"color"`
	Icon            string          `json:"icon" db:"icon"`
	CreatedAt       time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at" db:"updated_at"`
}

func (g SavingsGoal) Validate() error {
	if len(strings.TrimSpace(g.Name)) == 0 || len(g.Name) > 150 || !g.TargetAmount.IsPositive() || g.CurrentAmount.IsNegative() || !validHexColor(g.Color) || len(strings.TrimSpace(g.Icon)) == 0 || len(g.Icon) > 50 {
		return ErrValidation
	}
	return nil
}

type CashflowSummary struct {
	Income      decimal.Decimal `json:"income"`
	Expense     decimal.Decimal `json:"expense"`
	Net         decimal.Decimal `json:"net"`
	SavingsRate decimal.Decimal `json:"savings_rate"`
}

// CalculateCashflow computes net cashflow and percentage savings rate from income and expense.
// When income is positive, savings rate is (net / income) * 100 rounded to 2 decimal places.
// When income is zero or negative, savings rate defaults to 0.
func CalculateCashflow(income, expense decimal.Decimal) CashflowSummary {
	net := income.Sub(expense)
	rate := decimal.Zero
	if income.IsPositive() {
		rate = net.Div(income).Mul(decimal.NewFromInt(100)).Round(2)
	}
	return CashflowSummary{
		Income:      income,
		Expense:     expense,
		Net:         net,
		SavingsRate: rate,
	}
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
