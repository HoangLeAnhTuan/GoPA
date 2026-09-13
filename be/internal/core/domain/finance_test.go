package domain

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

func TestTransactionValidate_RequiresTargetOnlyForTransfers(t *testing.T) {
	accountID := uuid.New()
	targetID := uuid.New()
	base := Transaction{AccountID: accountID, Type: TransactionTransfer, Amount: decimal.NewFromInt(100), Currency: "VND", ExchangeRate: decimal.NewFromInt(1), OccurredAt: time.Now()}
	if err := base.Validate(); err == nil {
		t.Fatal("transfer without a target must be invalid")
	}
	base.ToAccountID = &targetID
	if err := base.Validate(); err != nil {
		t.Fatalf("valid transfer: %v", err)
	}
	base.CategoryID = &targetID
	if err := base.Validate(); err == nil {
		t.Fatal("transfer with a category must be invalid")
	}
}

func TestTransactionValidate_NormalizesTags(t *testing.T) {
	transaction := Transaction{AccountID: uuid.New(), Type: TransactionExpense, Amount: decimal.NewFromInt(100), Currency: "USD", ExchangeRate: decimal.NewFromInt(1), OccurredAt: time.Now(), Tags: []string{" Food ", "Travel"}}
	if err := transaction.Validate(); err != nil {
		t.Fatalf("validate transaction: %v", err)
	}
	if transaction.Tags[0] != "food" || transaction.Tags[1] != "travel" {
		t.Fatalf("tags were not normalized: %#v", transaction.Tags)
	}
}

func TestTransactionValidate_AllowsNoTags(t *testing.T) {
	transaction := Transaction{AccountID: uuid.New(), Type: TransactionIncome, Amount: decimal.NewFromInt(100), Currency: "VND", ExchangeRate: decimal.NewFromInt(1), OccurredAt: time.Now()}
	if err := transaction.Validate(); err != nil {
		t.Fatalf("a transaction without tags must be valid: %v", err)
	}
}

func TestAccountValidate_AllowsNegativeCurrentBalanceForCreditDebt(t *testing.T) {
	account := Account{Name: "Card", Type: AccountCreditCard, Currency: "VND", InitialBalance: decimal.Zero, CurrentBalance: decimal.NewFromInt(-100), Color: "#10B981", Icon: "credit-card"}
	if err := account.Validate(); err != nil {
		t.Fatalf("credit debt should be representable: %v", err)
	}
}

func TestTransaction_ConvertedAmount(t *testing.T) {
	tests := []struct {
		name         string
		amount       decimal.Decimal
		exchangeRate decimal.Decimal
		expected     decimal.Decimal
	}{
		{
			name:         "USD to VND conversion",
			amount:       decimal.NewFromInt(100),
			exchangeRate: decimal.RequireFromString("25450.50"),
			expected:     decimal.RequireFromString("2545050.00"),
		},
		{
			name:         "JPY to VND conversion",
			amount:       decimal.NewFromInt(10000),
			exchangeRate: decimal.RequireFromString("168.25"),
			expected:     decimal.RequireFromString("1682500.00"),
		},
		{
			name:         "Domestic same-currency (rate 1.0)",
			amount:       decimal.RequireFromString("500000"),
			exchangeRate: decimal.NewFromInt(1),
			expected:     decimal.RequireFromString("500000"),
		},
		{
			name:         "Fractional amount and exchange rate",
			amount:       decimal.RequireFromString("12.50"),
			exchangeRate: decimal.RequireFromString("1.50"),
			expected:     decimal.RequireFromString("18.75"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tx := Transaction{
				Amount:       tc.amount,
				ExchangeRate: tc.exchangeRate,
			}
			actual := tx.ConvertedAmount()
			if !actual.Equal(tc.expected) {
				t.Fatalf("expected converted amount %s, got %s", tc.expected, actual)
			}
		})
	}
}

func TestCalculateCashflow(t *testing.T) {
	tests := []struct {
		name        string
		income      decimal.Decimal
		expense     decimal.Decimal
		expectedNet decimal.Decimal
		expectedSR  decimal.Decimal
	}{
		{
			name:        "Standard surplus",
			income:      decimal.NewFromInt(5000),
			expense:     decimal.NewFromInt(3000),
			expectedNet: decimal.NewFromInt(2000),
			expectedSR:  decimal.NewFromInt(40), // (2000/5000)*100 = 40%
		},
		{
			name:        "Balanced budget",
			income:      decimal.NewFromInt(3000),
			expense:     decimal.NewFromInt(3000),
			expectedNet: decimal.Zero,
			expectedSR:  decimal.Zero,
		},
		{
			name:        "Deficit spending",
			income:      decimal.NewFromInt(2000),
			expense:     decimal.NewFromInt(3000),
			expectedNet: decimal.NewFromInt(-1000),
			expectedSR:  decimal.NewFromInt(-50), // (-1000/2000)*100 = -50%
		},
		{
			name:        "Zero income edge case",
			income:      decimal.Zero,
			expense:     decimal.NewFromInt(1500),
			expectedNet: decimal.NewFromInt(-1500),
			expectedSR:  decimal.Zero, // No division by zero, rate = 0
		},
		{
			name:        "Rounding to two decimals",
			income:      decimal.NewFromInt(3000),
			expense:     decimal.NewFromInt(1000),
			expectedNet: decimal.NewFromInt(2000),
			expectedSR:  decimal.RequireFromString("66.67"), // (2000/3000)*100 = 66.666... -> 66.67
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			summary := CalculateCashflow(tc.income, tc.expense)
			if !summary.Income.Equal(tc.income) {
				t.Errorf("expected income %s, got %s", tc.income, summary.Income)
			}
			if !summary.Expense.Equal(tc.expense) {
				t.Errorf("expected expense %s, got %s", tc.expense, summary.Expense)
			}
			if !summary.Net.Equal(tc.expectedNet) {
				t.Errorf("expected net %s, got %s", tc.expectedNet, summary.Net)
			}
			if !summary.SavingsRate.Equal(tc.expectedSR) {
				t.Errorf("expected savings rate %s, got %s", tc.expectedSR, summary.SavingsRate)
			}
		})
	}
}

func TestBudgetValidate_Constraints(t *testing.T) {
	now := time.Now().UTC()
	validBudget := Budget{
		Name:           "Food",
		Amount:         decimal.NewFromInt(500),
		Period:         BudgetPeriodMonthly,
		StartDate:      now,
		EndDate:        now.AddDate(0, 1, 0),
		AlertThreshold: decimal.RequireFromString("0.80"),
	}
	if err := validBudget.Validate(); err != nil {
		t.Fatalf("expected valid budget: %v", err)
	}

	// Threshold > 1.0 invalid
	invalidThreshold := validBudget
	invalidThreshold.AlertThreshold = decimal.RequireFromString("1.05")
	if err := invalidThreshold.Validate(); err == nil {
		t.Fatal("threshold > 1.0 must be invalid")
	}

	// Threshold <= 0 invalid
	invalidZeroThreshold := validBudget
	invalidZeroThreshold.AlertThreshold = decimal.Zero
	if err := invalidZeroThreshold.Validate(); err == nil {
		t.Fatal("threshold <= 0 must be invalid")
	}

	// EndDate before StartDate invalid
	invalidDates := validBudget
	invalidDates.EndDate = now.AddDate(0, -1, 0)
	if err := invalidDates.Validate(); err == nil {
		t.Fatal("end date before start date must be invalid")
	}

	// Amount <= 0 invalid
	invalidAmount := validBudget
	invalidAmount.Amount = decimal.NewFromInt(-10)
	if err := invalidAmount.Validate(); err == nil {
		t.Fatal("negative budget amount must be invalid")
	}
}

func TestSavingsGoalValidate_Constraints(t *testing.T) {
	validGoal := SavingsGoal{
		Name:          "New Laptop",
		TargetAmount:  decimal.NewFromInt(2000),
		CurrentAmount: decimal.Zero,
		Color:         "#3B82F6",
		Icon:          "laptop",
	}
	if err := validGoal.Validate(); err != nil {
		t.Fatalf("expected valid goal: %v", err)
	}

	// Target amount <= 0 invalid
	zeroTarget := validGoal
	zeroTarget.TargetAmount = decimal.Zero
	if err := zeroTarget.Validate(); err == nil {
		t.Fatal("target amount <= 0 must be invalid")
	}

	// Current amount negative invalid
	negCurrent := validGoal
	negCurrent.CurrentAmount = decimal.NewFromInt(-100)
	if err := negCurrent.Validate(); err == nil {
		t.Fatal("negative current amount must be invalid")
	}

	// Invalid hex color
	badColor := validGoal
	badColor.Color = "blue"
	if err := badColor.Validate(); err == nil {
		t.Fatal("invalid hex color must be rejected")
	}
}
