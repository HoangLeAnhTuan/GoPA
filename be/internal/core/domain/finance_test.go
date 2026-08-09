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
	if err := transaction.Validate(); err != nil { t.Fatalf("a transaction without tags must be valid: %v", err) }
}

func TestAccountValidate_AllowsNegativeCurrentBalanceForCreditDebt(t *testing.T) {
	account := Account{Name: "Card", Type: AccountCreditCard, Currency: "VND", InitialBalance: decimal.Zero, CurrentBalance: decimal.NewFromInt(-100), Color: "#10B981", Icon: "credit-card"}
	if err := account.Validate(); err != nil {
		t.Fatalf("credit debt should be representable: %v", err)
	}
}
