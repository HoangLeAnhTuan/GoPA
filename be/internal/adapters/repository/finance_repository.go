package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/shopspring/decimal"
	"gopa/internal/core/domain"
	"gopa/internal/core/ports"
	"gopa/pkg/database"
)

const accountColumns = `id, user_id, name, type, currency, initial_balance, current_balance, color, icon, is_archived, created_at, updated_at`
const categoryColumns = `id, user_id, name, type, icon, color, parent_id, monthly_budget, created_at, updated_at`
const transactionColumns = `id, user_id, account_id, to_account_id, category_id, type, amount, currency, exchange_rate, description, merchant, tags, occurred_at, is_recurring, recurring_rule, created_at, updated_at`

type FinanceRepository struct{ db *sqlx.DB }

func NewFinanceRepository(db *sqlx.DB) *FinanceRepository { return &FinanceRepository{db: db} }

// transactionRow keeps PostgreSQL-specific array handling at the repository
// boundary. The domain model deliberately stays independent of pq types.
type transactionRow struct {
	ID            uuid.UUID              `db:"id"`
	UserID        uuid.UUID              `db:"user_id"`
	AccountID     uuid.UUID              `db:"account_id"`
	ToAccountID   *uuid.UUID             `db:"to_account_id"`
	CategoryID    *uuid.UUID             `db:"category_id"`
	Type          domain.TransactionType `db:"type"`
	Amount        decimal.Decimal        `db:"amount"`
	Currency      string                 `db:"currency"`
	ExchangeRate  decimal.Decimal        `db:"exchange_rate"`
	Description   string                 `db:"description"`
	Merchant      *string                `db:"merchant"`
	Tags          pq.StringArray         `db:"tags"`
	OccurredAt    time.Time              `db:"occurred_at"`
	IsRecurring   bool                   `db:"is_recurring"`
	RecurringRule []byte                 `db:"recurring_rule"`
	CreatedAt     time.Time              `db:"created_at"`
	UpdatedAt     time.Time              `db:"updated_at"`
}

func (row transactionRow) domain() domain.Transaction {
	tags := make([]string, len(row.Tags))
	copy(tags, row.Tags)
	return domain.Transaction{
		ID:            row.ID,
		UserID:        row.UserID,
		AccountID:     row.AccountID,
		ToAccountID:   row.ToAccountID,
		CategoryID:    row.CategoryID,
		Type:          row.Type,
		Amount:        row.Amount,
		Currency:      row.Currency,
		ExchangeRate:  row.ExchangeRate,
		Description:   row.Description,
		Merchant:      row.Merchant,
		Tags:          tags,
		OccurredAt:    row.OccurredAt,
		IsRecurring:   row.IsRecurring,
		RecurringRule: row.RecurringRule,
		CreatedAt:     row.CreatedAt,
		UpdatedAt:     row.UpdatedAt,
	}
}

func (r *FinanceRepository) ListAccounts(ctx context.Context, userID uuid.UUID, includeArchived bool) ([]domain.Account, error) {
	accounts := make([]domain.Account, 0)
	query := `SELECT ` + accountColumns + ` FROM accounts WHERE user_id = $1`
	if !includeArchived {
		query += ` AND is_archived = FALSE`
	}
	query += ` ORDER BY is_archived, created_at DESC`
	if err := r.db.SelectContext(ctx, &accounts, query, userID); err != nil {
		return nil, fmt.Errorf("list accounts: %w", err)
	}
	return accounts, nil
}

func (r *FinanceRepository) GetAccount(ctx context.Context, userID, id uuid.UUID) (domain.Account, error) {
	var account domain.Account
	if err := r.db.GetContext(ctx, &account, `SELECT `+accountColumns+` FROM accounts WHERE id = $1 AND user_id = $2`, id, userID); err != nil {
		return domain.Account{}, mapNotFound(err, "get account")
	}
	return account, nil
}

func (r *FinanceRepository) CreateAccount(ctx context.Context, account domain.Account) (domain.Account, error) {
	const query = `INSERT INTO accounts (` + accountColumns + `) VALUES (:id, :user_id, :name, :type, :currency, :initial_balance, :current_balance, :color, :icon, :is_archived, :created_at, :updated_at)`
	if _, err := r.db.NamedExecContext(ctx, query, account); err != nil {
		return domain.Account{}, fmt.Errorf("create account: %w", err)
	}
	return account, nil
}

func (r *FinanceRepository) UpdateAccount(ctx context.Context, account domain.Account) (domain.Account, error) {
	result, err := r.db.NamedExecContext(ctx, `UPDATE accounts SET name = :name, type = :type, currency = :currency, color = :color, icon = :icon, is_archived = :is_archived, updated_at = :updated_at WHERE id = :id AND user_id = :user_id`, account)
	if err != nil {
		return domain.Account{}, fmt.Errorf("update account: %w", err)
	}
	if err := requireAffected(result, "update account"); err != nil {
		return domain.Account{}, err
	}
	return account, nil
}

func (r *FinanceRepository) DeleteAccount(ctx context.Context, userID, id uuid.UUID) error {
	return database.WithTx(ctx, r.db, func(tx *sqlx.Tx) error {
		var transactionCount int
		if err := tx.GetContext(ctx, &transactionCount, `SELECT COUNT(*) FROM transactions WHERE user_id = $1 AND (account_id = $2 OR to_account_id = $2)`, userID, id); err != nil {
			return fmt.Errorf("count account transactions: %w", err)
		}
		if transactionCount > 0 {
			return domain.ErrConflict
		}
		result, err := tx.ExecContext(ctx, `DELETE FROM accounts WHERE id = $1 AND user_id = $2`, id, userID)
		if err != nil {
			return fmt.Errorf("delete account: %w", err)
		}
		return requireAffected(result, "delete account")
	})
}

func (r *FinanceRepository) ListCategories(ctx context.Context, userID uuid.UUID) ([]domain.Category, error) {
	categories := make([]domain.Category, 0)
	if err := r.db.SelectContext(ctx, &categories, `SELECT `+categoryColumns+` FROM categories WHERE user_id = $1 ORDER BY type, name`, userID); err != nil {
		return nil, fmt.Errorf("list categories: %w", err)
	}
	return categories, nil
}

func (r *FinanceRepository) GetCategory(ctx context.Context, userID, id uuid.UUID) (domain.Category, error) {
	var category domain.Category
	if err := r.db.GetContext(ctx, &category, `SELECT `+categoryColumns+` FROM categories WHERE id = $1 AND user_id = $2`, id, userID); err != nil {
		return domain.Category{}, mapNotFound(err, "get category")
	}
	return category, nil
}

func (r *FinanceRepository) CreateCategory(ctx context.Context, category domain.Category) (domain.Category, error) {
	if category.ParentID != nil {
		if err := r.ensureCategoryParent(ctx, r.db, category.UserID, category.ID, category.ParentID); err != nil {
			return domain.Category{}, err
		}
	}
	const query = `INSERT INTO categories (` + categoryColumns + `) VALUES (:id, :user_id, :name, :type, :icon, :color, :parent_id, :monthly_budget, :created_at, :updated_at)`
	if _, err := r.db.NamedExecContext(ctx, query, category); err != nil {
		return domain.Category{}, fmt.Errorf("create category: %w", err)
	}
	return category, nil
}

func (r *FinanceRepository) UpdateCategory(ctx context.Context, category domain.Category) (domain.Category, error) {
	if category.ParentID != nil {
		if err := r.ensureCategoryParent(ctx, r.db, category.UserID, category.ID, category.ParentID); err != nil {
			return domain.Category{}, err
		}
	}
	result, err := r.db.NamedExecContext(ctx, `UPDATE categories SET name = :name, type = :type, icon = :icon, color = :color, parent_id = :parent_id, monthly_budget = :monthly_budget, updated_at = :updated_at WHERE id = :id AND user_id = :user_id`, category)
	if err != nil {
		return domain.Category{}, fmt.Errorf("update category: %w", err)
	}
	if err := requireAffected(result, "update category"); err != nil {
		return domain.Category{}, err
	}
	return category, nil
}

func (r *FinanceRepository) DeleteCategory(ctx context.Context, userID, id uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM categories WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("delete category: %w", err)
	}
	return requireAffected(result, "delete category")
}

func (r *FinanceRepository) ListTransactions(ctx context.Context, userID uuid.UUID, filter ports.TransactionFilter) ([]domain.Transaction, error) {
	where := []string{"user_id = $1"}
	args := []any{userID}
	if filter.AccountID != nil {
		args = append(args, *filter.AccountID)
		where = append(where, fmt.Sprintf("(account_id = $%d OR to_account_id = $%d)", len(args), len(args)))
	}
	if filter.CategoryID != nil {
		args = append(args, *filter.CategoryID)
		where = append(where, fmt.Sprintf("category_id = $%d", len(args)))
	}
	if filter.Type != nil {
		args = append(args, *filter.Type)
		where = append(where, fmt.Sprintf("type = $%d", len(args)))
	}
	if filter.From != nil {
		args = append(args, filter.From.UTC())
		where = append(where, fmt.Sprintf("occurred_at >= $%d", len(args)))
	}
	if filter.To != nil {
		args = append(args, filter.To.UTC())
		where = append(where, fmt.Sprintf("occurred_at <= $%d", len(args)))
	}
	if strings.TrimSpace(filter.Query) != "" {
		args = append(args, "%"+strings.TrimSpace(filter.Query)+"%")
		where = append(where, fmt.Sprintf("(description ILIKE $%d OR merchant ILIKE $%d)", len(args), len(args)))
	}
	args = append(args, filter.Limit)
	rows := make([]transactionRow, 0)
	query := fmt.Sprintf(`SELECT %s FROM transactions WHERE %s ORDER BY occurred_at DESC, created_at DESC LIMIT $%d`, transactionColumns, strings.Join(where, " AND "), len(args))
	if err := r.db.SelectContext(ctx, &rows, query, args...); err != nil {
		return nil, fmt.Errorf("list transactions: %w", err)
	}
	transactions := make([]domain.Transaction, 0, len(rows))
	for _, row := range rows {
		transactions = append(transactions, row.domain())
	}
	return transactions, nil
}

func (r *FinanceRepository) GetTransaction(ctx context.Context, userID, id uuid.UUID) (domain.Transaction, error) {
	return getTransaction(ctx, r.db, userID, id, false)
}

func (r *FinanceRepository) CreateTransaction(ctx context.Context, transaction domain.Transaction) (domain.Transaction, error) {
	if err := database.WithTx(ctx, r.db, func(tx *sqlx.Tx) error {
		if err := r.validateTransactionReferences(ctx, tx, transaction); err != nil {
			return err
		}
		if _, err := tx.NamedExecContext(ctx, `INSERT INTO transactions (`+transactionColumns+`) VALUES (:id, :user_id, :account_id, :to_account_id, :category_id, :type, :amount, :currency, :exchange_rate, :description, :merchant, :tags, :occurred_at, :is_recurring, :recurring_rule, :created_at, :updated_at)`, transactionParams(transaction)); err != nil {
			return fmt.Errorf("insert transaction: %w", err)
		}
		return applyBalanceEffect(ctx, tx, transaction, false)
	}); err != nil {
		return domain.Transaction{}, err
	}
	return transaction, nil
}

func (r *FinanceRepository) UpdateTransaction(ctx context.Context, transaction domain.Transaction) (domain.Transaction, error) {
	if err := database.WithTx(ctx, r.db, func(tx *sqlx.Tx) error {
		previous, err := getTransaction(ctx, tx, transaction.UserID, transaction.ID, true)
		if err != nil {
			return err
		}
		if err := applyBalanceEffect(ctx, tx, previous, true); err != nil {
			return err
		}
		if err := r.validateTransactionReferences(ctx, tx, transaction); err != nil {
			return err
		}
		result, err := tx.NamedExecContext(ctx, `UPDATE transactions SET account_id = :account_id, to_account_id = :to_account_id, category_id = :category_id, type = :type, amount = :amount, currency = :currency, exchange_rate = :exchange_rate, description = :description, merchant = :merchant, tags = :tags, occurred_at = :occurred_at, is_recurring = :is_recurring, recurring_rule = :recurring_rule, updated_at = :updated_at WHERE id = :id AND user_id = :user_id`, transactionParams(transaction))
		if err != nil {
			return fmt.Errorf("update transaction: %w", err)
		}
		if err := requireAffected(result, "update transaction"); err != nil {
			return err
		}
		return applyBalanceEffect(ctx, tx, transaction, false)
	}); err != nil {
		return domain.Transaction{}, err
	}
	return transaction, nil
}

func (r *FinanceRepository) DeleteTransaction(ctx context.Context, userID, id uuid.UUID, updatedAt time.Time) error {
	return database.WithTx(ctx, r.db, func(tx *sqlx.Tx) error {
		transaction, err := getTransaction(ctx, tx, userID, id, true)
		if err != nil {
			return err
		}
		if err := applyBalanceEffect(ctx, tx, transaction, true); err != nil {
			return err
		}
		result, err := tx.ExecContext(ctx, `DELETE FROM transactions WHERE id = $1 AND user_id = $2`, id, userID)
		if err != nil {
			return fmt.Errorf("delete transaction: %w", err)
		}
		if err := requireAffected(result, "delete transaction"); err != nil {
			return err
		}
		return nil
	})
}

func (r *FinanceRepository) Cashflow(ctx context.Context, userID uuid.UUID, from, to time.Time) (domain.CashflowSummary, error) {
	var summary struct{ Income, Expense decimal.Decimal }
	if err := r.db.GetContext(ctx, &summary, `SELECT COALESCE(SUM(amount * exchange_rate) FILTER (WHERE type = 'INCOME'), 0) AS income, COALESCE(SUM(amount * exchange_rate) FILTER (WHERE type = 'EXPENSE'), 0) AS expense FROM transactions WHERE user_id = $1 AND occurred_at >= $2 AND occurred_at < $3`, userID, from.UTC(), to.UTC()); err != nil {
		return domain.CashflowSummary{}, fmt.Errorf("get cashflow: %w", err)
	}
	net := summary.Income.Sub(summary.Expense)
	rate := decimal.Zero
	if summary.Income.IsPositive() {
		rate = net.Div(summary.Income).Mul(decimal.NewFromInt(100)).Round(2)
	}
	return domain.CashflowSummary{Income: summary.Income, Expense: summary.Expense, Net: net, SavingsRate: rate}, nil
}

func (r *FinanceRepository) SpendingByCategory(ctx context.Context, userID uuid.UUID, from, to time.Time) ([]domain.CategorySpending, error) {
	spending := make([]domain.CategorySpending, 0)
	const query = `SELECT t.category_id, COALESCE(c.name, 'Uncategorized') AS category_name, COALESCE(c.color, '#64748B') AS color, COALESCE(SUM(t.amount * t.exchange_rate), 0) AS amount FROM transactions t LEFT JOIN categories c ON c.id = t.category_id AND c.user_id = t.user_id WHERE t.user_id = $1 AND t.type = 'EXPENSE' AND t.occurred_at >= $2 AND t.occurred_at < $3 GROUP BY t.category_id, c.name, c.color ORDER BY amount DESC`
	if err := r.db.SelectContext(ctx, &spending, query, userID, from.UTC(), to.UTC()); err != nil {
		return nil, fmt.Errorf("get spending by category: %w", err)
	}
	return spending, nil
}

func (r *FinanceRepository) NetWorth(ctx context.Context, userID uuid.UUID) (decimal.Decimal, error) {
	var total decimal.Decimal
	if err := r.db.GetContext(ctx, &total, `SELECT COALESCE(SUM(current_balance), 0) FROM accounts WHERE user_id = $1 AND is_archived = FALSE`, userID); err != nil {
		return decimal.Zero, fmt.Errorf("get net worth: %w", err)
	}
	return total, nil
}

type sqlxGetter interface {
	GetContext(context.Context, any, string, ...any) error
}

func getTransaction(ctx context.Context, getter sqlxGetter, userID, id uuid.UUID, lock bool) (domain.Transaction, error) {
	var transaction transactionRow
	query := `SELECT ` + transactionColumns + ` FROM transactions WHERE id = $1 AND user_id = $2`
	if lock {
		query += ` FOR UPDATE`
	}
	if err := getter.GetContext(ctx, &transaction, query, id, userID); err != nil {
		return domain.Transaction{}, mapNotFound(err, "get transaction")
	}
	return transaction.domain(), nil
}

func (r *FinanceRepository) validateTransactionReferences(ctx context.Context, tx *sqlx.Tx, transaction domain.Transaction) error {
	if _, err := getAccountForUpdate(ctx, tx, transaction.UserID, transaction.AccountID); err != nil {
		return err
	}
	if transaction.ToAccountID != nil {
		if _, err := getAccountForUpdate(ctx, tx, transaction.UserID, *transaction.ToAccountID); err != nil {
			return err
		}
	}
	if transaction.CategoryID != nil {
		category, err := getCategoryForUpdate(ctx, tx, transaction.UserID, *transaction.CategoryID)
		if err != nil {
			return err
		}
		if (transaction.Type == domain.TransactionIncome && category.Type != domain.CategoryIncome) || (transaction.Type == domain.TransactionExpense && category.Type != domain.CategoryExpense) {
			return domain.ErrValidation
		}
	}
	return nil
}

func getAccountForUpdate(ctx context.Context, tx *sqlx.Tx, userID, id uuid.UUID) (domain.Account, error) {
	var account domain.Account
	if err := tx.GetContext(ctx, &account, `SELECT `+accountColumns+` FROM accounts WHERE id = $1 AND user_id = $2 FOR UPDATE`, id, userID); err != nil {
		return domain.Account{}, mapNotFound(err, "get transaction account")
	}
	if account.IsArchived {
		return domain.Account{}, domain.ErrValidation
	}
	return account, nil
}

func getCategoryForUpdate(ctx context.Context, tx *sqlx.Tx, userID, id uuid.UUID) (domain.Category, error) {
	var category domain.Category
	if err := tx.GetContext(ctx, &category, `SELECT `+categoryColumns+` FROM categories WHERE id = $1 AND user_id = $2 FOR UPDATE`, id, userID); err != nil {
		return domain.Category{}, mapNotFound(err, "get transaction category")
	}
	return category, nil
}

func applyBalanceEffect(ctx context.Context, tx *sqlx.Tx, transaction domain.Transaction, reverse bool) error {
	operation := transaction.Type
	if reverse {
		switch operation {
		case domain.TransactionIncome:
			operation = domain.TransactionExpense
		case domain.TransactionExpense:
			operation = domain.TransactionIncome
		}
	}
	updatedAt := transaction.UpdatedAt.UTC()
	switch operation {
	case domain.TransactionIncome:
		return updateBalance(ctx, tx, transaction.UserID, transaction.AccountID, transaction.Amount, transaction.ExchangeRate, true, updatedAt)
	case domain.TransactionExpense:
		return updateBalance(ctx, tx, transaction.UserID, transaction.AccountID, transaction.Amount, transaction.ExchangeRate, false, updatedAt)
	case domain.TransactionTransfer:
		if transaction.ToAccountID == nil {
			return domain.ErrValidation
		}
		if reverse {
			if err := updateBalance(ctx, tx, transaction.UserID, transaction.AccountID, transaction.Amount, decimal.NewFromInt(1), true, updatedAt); err != nil {
				return err
			}
			return updateBalance(ctx, tx, transaction.UserID, *transaction.ToAccountID, transaction.Amount, transaction.ExchangeRate, false, updatedAt)
		}
		if err := updateBalance(ctx, tx, transaction.UserID, transaction.AccountID, transaction.Amount, decimal.NewFromInt(1), false, updatedAt); err != nil {
			return err
		}
		return updateBalance(ctx, tx, transaction.UserID, *transaction.ToAccountID, transaction.Amount, transaction.ExchangeRate, true, updatedAt)
	default:
		return domain.ErrValidation
	}
}

func updateBalance(ctx context.Context, tx *sqlx.Tx, userID, accountID uuid.UUID, amount, exchangeRate decimal.Decimal, add bool, updatedAt time.Time) error {
	operator := "+"
	if !add {
		operator = "-"
	}
	result, err := tx.ExecContext(ctx, `UPDATE accounts SET current_balance = current_balance `+operator+` ($1::numeric * $2::numeric), updated_at = $3 WHERE id = $4 AND user_id = $5`, amount, exchangeRate, updatedAt, accountID, userID)
	if err != nil {
		return fmt.Errorf("update account balance: %w", err)
	}
	return requireAffected(result, "update account balance")
}

func (r *FinanceRepository) ensureCategoryParent(ctx context.Context, getter sqlxGetter, userID, categoryID uuid.UUID, parentID *uuid.UUID) error {
	if parentID == nil {
		return nil
	}
	if *parentID == categoryID {
		return domain.ErrValidation
	}
	var parent struct {
		ID       uuid.UUID  `db:"id"`
		ParentID *uuid.UUID `db:"parent_id"`
	}
	if err := getter.GetContext(ctx, &parent, `SELECT id, parent_id FROM categories WHERE id = $1 AND user_id = $2`, *parentID, userID); err != nil {
		return mapNotFound(err, "get category parent")
	}
	if parent.ParentID != nil {
		return domain.ErrValidation
	}
	return nil
}

func transactionParams(transaction domain.Transaction) map[string]any {
	var recurringRule any
	if len(transaction.RecurringRule) > 0 {
		recurringRule = string(transaction.RecurringRule)
	}
	tags := transaction.Tags
	if tags == nil {
		tags = make([]string, 0)
	}
	return map[string]any{
		"id": transaction.ID, "user_id": transaction.UserID, "account_id": transaction.AccountID, "to_account_id": transaction.ToAccountID, "category_id": transaction.CategoryID,
		"type": transaction.Type, "amount": transaction.Amount, "currency": transaction.Currency, "exchange_rate": transaction.ExchangeRate,
		"description": transaction.Description, "merchant": transaction.Merchant, "tags": pq.StringArray(tags), "occurred_at": transaction.OccurredAt.UTC(),
		"is_recurring": transaction.IsRecurring, "recurring_rule": recurringRule, "created_at": transaction.CreatedAt.UTC(), "updated_at": transaction.UpdatedAt.UTC(),
	}
}

func requireAffected(result sql.Result, operation string) error {
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("count %s: %w", operation, err)
	}
	if rows == 0 {
		return domain.ErrNotFound
	}
	return nil
}

var _ ports.FinanceRepository = (*FinanceRepository)(nil)
