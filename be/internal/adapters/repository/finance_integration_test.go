package repository

import (
	"context"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopa/internal/core/domain"
	"gopa/pkg/database"
)

func getTestDSN() string {
	if dsn := os.Getenv("POSTGRES_DSN"); dsn != "" {
		return dsn
	}
	return "postgres://gopa:gopa@localhost:5432/gopa?sslmode=disable"
}

func setupTestDB(t *testing.T) (*sqlx.DB, uuid.UUID) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := database.Open(ctx, getTestDSN())
	if err != nil {
		t.Skipf("Skipping integration test: cannot connect to postgres at %s: %v", getTestDSN(), err)
		return nil, uuid.Nil
	}

	userID := uuid.New()
	email := "test-" + userID.String()[:8] + "@example.com"
	_, err = db.ExecContext(ctx, `INSERT INTO users (id, email, password_hash, display_name, role, created_at, updated_at) VALUES ($1, $2, 'hash', 'Test User', 'USER', now(), now())`, userID, email)
	require.NoError(t, err, "insert test user")

	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cleanupCancel()
		_, _ = db.ExecContext(cleanupCtx, `DELETE FROM outbox_events WHERE payload::text LIKE '%`+userID.String()+`%'`)
		_, _ = db.ExecContext(cleanupCtx, `DELETE FROM transactions WHERE user_id = $1`, userID)
		_, _ = db.ExecContext(cleanupCtx, `DELETE FROM categories WHERE user_id = $1`, userID)
		_, _ = db.ExecContext(cleanupCtx, `DELETE FROM accounts WHERE user_id = $1`, userID)
		_, _ = db.ExecContext(cleanupCtx, `DELETE FROM users WHERE id = $1`, userID)
		db.Close()
	})

	return db, userID
}

func TestFinanceRepo_CreateTransaction_AtomicBalanceAndRollback(t *testing.T) {
	db, userID := setupTestDB(t)
	repo := NewFinanceRepository(db)
	ctx := context.Background()

	// 1. Create account
	account, err := repo.CreateAccount(ctx, domain.Account{
		ID:             uuid.New(),
		UserID:         userID,
		Name:           "Main Checking",
		Type:           domain.AccountBank,
		Currency:       "VND",
		InitialBalance: decimal.NewFromInt(10000),
		CurrentBalance: decimal.NewFromInt(10000),
		Color:          "#10B981",
		Icon:           "wallet",
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	})
	require.NoError(t, err)

	// 2. Deposit 5,000 VND (INCOME)
	incTx, err := repo.CreateTransaction(ctx, domain.Transaction{
		ID:           uuid.New(),
		UserID:       userID,
		AccountID:    account.ID,
		Type:         domain.TransactionIncome,
		Amount:       decimal.NewFromInt(5000),
		Currency:     "VND",
		ExchangeRate: decimal.NewFromInt(1),
		Description:  "Salary bonus",
		OccurredAt:   time.Now().UTC(),
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	})
	require.NoError(t, err)
	assert.NotEmpty(t, incTx.ID)

	accAfterIncome, err := repo.GetAccount(ctx, userID, account.ID)
	require.NoError(t, err)
	assert.True(t, accAfterIncome.CurrentBalance.Equal(decimal.NewFromInt(15000)), "Balance should increase by 5000 to 15000")

	// 3. Spend 3,000 VND (EXPENSE)
	expTx, err := repo.CreateTransaction(ctx, domain.Transaction{
		ID:           uuid.New(),
		UserID:       userID,
		AccountID:    account.ID,
		Type:         domain.TransactionExpense,
		Amount:       decimal.NewFromInt(3000),
		Currency:     "VND",
		ExchangeRate: decimal.NewFromInt(1),
		Description:  "Dinner",
		OccurredAt:   time.Now().UTC(),
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	})
	require.NoError(t, err)
	assert.NotEmpty(t, expTx.ID)

	accAfterExpense, err := repo.GetAccount(ctx, userID, account.ID)
	require.NoError(t, err)
	assert.True(t, accAfterExpense.CurrentBalance.Equal(decimal.NewFromInt(12000)), "Balance should decrease by 3000 to 12000")

	// 4. Test Rollback on invalid category reference
	invalidCategoryID := uuid.New() // Non-existent category
	_, err = repo.CreateTransaction(ctx, domain.Transaction{
		ID:           uuid.New(),
		UserID:       userID,
		AccountID:    account.ID,
		CategoryID:   &invalidCategoryID,
		Type:         domain.TransactionExpense,
		Amount:       decimal.NewFromInt(2000),
		Currency:     "VND",
		ExchangeRate: decimal.NewFromInt(1),
		Description:  "Failed transaction",
		OccurredAt:   time.Now().UTC(),
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	})
	assert.Error(t, err, "Transaction with non-existent category must fail")

	// Verify balance was untouched (still 12,000)
	accAfterRollback, err := repo.GetAccount(ctx, userID, account.ID)
	require.NoError(t, err)
	assert.True(t, accAfterRollback.CurrentBalance.Equal(decimal.NewFromInt(12000)), "Account balance must be preserved after rollback")
}

func TestFinanceRepo_Transfer_AtomicDoubleEntryConsistency(t *testing.T) {
	db, userID := setupTestDB(t)
	repo := NewFinanceRepository(db)
	ctx := context.Background()

	// Create Bank account (20,000,000 VND) and Cash wallet (1,000,000 VND)
	bank, err := repo.CreateAccount(ctx, domain.Account{
		ID:             uuid.New(),
		UserID:         userID,
		Name:           "Bank",
		Type:           domain.AccountBank,
		Currency:       "VND",
		InitialBalance: decimal.NewFromInt(20000000),
		CurrentBalance: decimal.NewFromInt(20000000),
		Color:          "#3B82F6",
		Icon:           "building",
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	})
	require.NoError(t, err)

	cash, err := repo.CreateAccount(ctx, domain.Account{
		ID:             uuid.New(),
		UserID:         userID,
		Name:           "Cash Wallet",
		Type:           domain.AccountCash,
		Currency:       "VND",
		InitialBalance: decimal.NewFromInt(1000000),
		CurrentBalance: decimal.NewFromInt(1000000),
		Color:          "#10B981",
		Icon:           "wallet",
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	})
	require.NoError(t, err)

	// Transfer 5,000,000 VND from Bank to Cash
	transferAmount := decimal.NewFromInt(5000000)
	transferTx, err := repo.CreateTransaction(ctx, domain.Transaction{
		ID:           uuid.New(),
		UserID:       userID,
		AccountID:    bank.ID,
		ToAccountID:  &cash.ID,
		Type:         domain.TransactionTransfer,
		Amount:       transferAmount,
		Currency:     "VND",
		ExchangeRate: decimal.NewFromInt(1),
		Description:  "ATM withdrawal to cash",
		OccurredAt:   time.Now().UTC(),
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	})
	require.NoError(t, err)

	// Verify double-entry consistency:
	// Bank decreased by 5,000,000 -> 15,000,000
	// Cash increased by 5,000,000 -> 6,000,000
	updatedBank, err := repo.GetAccount(ctx, userID, bank.ID)
	require.NoError(t, err)
	assert.True(t, updatedBank.CurrentBalance.Equal(decimal.NewFromInt(15000000)))

	updatedCash, err := repo.GetAccount(ctx, userID, cash.ID)
	require.NoError(t, err)
	assert.True(t, updatedCash.CurrentBalance.Equal(decimal.NewFromInt(6000000)))

	// Delete transfer transaction and verify reverse consistency
	err = repo.DeleteTransaction(ctx, userID, transferTx.ID, time.Now().UTC())
	require.NoError(t, err)

	restoredBank, err := repo.GetAccount(ctx, userID, bank.ID)
	require.NoError(t, err)
	assert.True(t, restoredBank.CurrentBalance.Equal(decimal.NewFromInt(20000000)), "Bank balance restored to 20,000,000")

	restoredCash, err := repo.GetAccount(ctx, userID, cash.ID)
	require.NoError(t, err)
	assert.True(t, restoredCash.CurrentBalance.Equal(decimal.NewFromInt(1000000)), "Cash balance restored to 1,000,000")
}

func TestFinanceRepo_ConcurrentTransactions_ACIDConsistency(t *testing.T) {
	db, userID := setupTestDB(t)
	repo := NewFinanceRepository(db)
	ctx := context.Background()

	// Initial balance: 100,000 VND
	initialBalance := decimal.NewFromInt(100000)
	account, err := repo.CreateAccount(ctx, domain.Account{
		ID:             uuid.New(),
		UserID:         userID,
		Name:           "Concurrency Test Account",
		Type:           domain.AccountBank,
		Currency:       "VND",
		InitialBalance: initialBalance,
		CurrentBalance: initialBalance,
		Color:          "#8B5CF6",
		Icon:           "cpu",
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	})
	require.NoError(t, err)

	// 10 concurrent deposits of 1,000 VND (+10,000)
	// 10 concurrent withdrawals of 500 VND (-5,000)
	// Expected net change: +5,000 VND -> Final balance: 105,000 VND
	const goroutines = 10
	var wg sync.WaitGroup
	errCh := make(chan error, goroutines*2)

	for i := 0; i < goroutines; i++ {
		wg.Add(2)

		// Deposit goroutine
		go func() {
			defer wg.Done()
			_, err := repo.CreateTransaction(ctx, domain.Transaction{
				ID:           uuid.New(),
				UserID:       userID,
				AccountID:    account.ID,
				Type:         domain.TransactionIncome,
				Amount:       decimal.NewFromInt(1000),
				Currency:     "VND",
				ExchangeRate: decimal.NewFromInt(1),
				Description:  "Concurrent Deposit",
				OccurredAt:   time.Now().UTC(),
				CreatedAt:    time.Now().UTC(),
				UpdatedAt:    time.Now().UTC(),
			})
			if err != nil {
				errCh <- err
			}
		}()

		// Withdrawal goroutine
		go func() {
			defer wg.Done()
			_, err := repo.CreateTransaction(ctx, domain.Transaction{
				ID:           uuid.New(),
				UserID:       userID,
				AccountID:    account.ID,
				Type:         domain.TransactionExpense,
				Amount:       decimal.NewFromInt(500),
				Currency:     "VND",
				ExchangeRate: decimal.NewFromInt(1),
				Description:  "Concurrent Withdrawal",
				OccurredAt:   time.Now().UTC(),
				CreatedAt:    time.Now().UTC(),
				UpdatedAt:    time.Now().UTC(),
			})
			if err != nil {
				errCh <- err
			}
		}()
	}

	wg.Wait()
	close(errCh)

	for err := range errCh {
		require.NoError(t, err, "concurrent transaction failed")
	}

	finalAcc, err := repo.GetAccount(ctx, userID, account.ID)
	require.NoError(t, err)

	expectedFinal := decimal.NewFromInt(105000)
	assert.True(t, finalAcc.CurrentBalance.Equal(expectedFinal),
		"Expected final balance %s, got %s (lost updates detected!)", expectedFinal, finalAcc.CurrentBalance)
}
