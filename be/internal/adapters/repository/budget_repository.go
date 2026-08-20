package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"gopa/internal/core/domain"
	"gopa/internal/core/ports"
)

const budgetColumns = `id, user_id, category_id, name, amount, period, start_date, end_date, alert_threshold, is_active, created_at, updated_at`
const savingsGoalColumns = `id, user_id, name, target_amount, current_amount, linked_account_id, target_date, color, icon, created_at, updated_at`

// BudgetRepository implements ports.BudgetRepository.
type BudgetRepository struct{ db *sqlx.DB }

func NewBudgetRepository(db *sqlx.DB) *BudgetRepository { return &BudgetRepository{db: db} }

func (r *BudgetRepository) ListBudgets(ctx context.Context, userID uuid.UUID, activeOnly bool) ([]domain.Budget, error) {
	budgets := make([]domain.Budget, 0)
	query := `SELECT ` + budgetColumns + ` FROM budgets WHERE user_id = $1`
	if activeOnly {
		query += ` AND is_active = TRUE`
	}
	query += ` ORDER BY created_at DESC`
	if err := r.db.SelectContext(ctx, &budgets, query, userID); err != nil {
		return nil, fmt.Errorf("list budgets: %w", err)
	}
	return budgets, nil
}

func (r *BudgetRepository) GetBudget(ctx context.Context, userID, id uuid.UUID) (domain.Budget, error) {
	var budget domain.Budget
	if err := r.db.GetContext(ctx, &budget, `SELECT `+budgetColumns+` FROM budgets WHERE id = $1 AND user_id = $2`, id, userID); err != nil {
		return domain.Budget{}, mapNotFound(err, "get budget")
	}
	return budget, nil
}

func (r *BudgetRepository) CreateBudget(ctx context.Context, budget domain.Budget) (domain.Budget, error) {
	const query = `INSERT INTO budgets (` + budgetColumns + `) VALUES (:id, :user_id, :category_id, :name, :amount, :period, :start_date, :end_date, :alert_threshold, :is_active, :created_at, :updated_at)`
	if _, err := r.db.NamedExecContext(ctx, query, budget); err != nil {
		return domain.Budget{}, fmt.Errorf("create budget: %w", err)
	}
	return budget, nil
}

func (r *BudgetRepository) UpdateBudget(ctx context.Context, budget domain.Budget) (domain.Budget, error) {
	result, err := r.db.NamedExecContext(ctx,
		`UPDATE budgets SET category_id = :category_id, name = :name, amount = :amount, period = :period, start_date = :start_date, end_date = :end_date, alert_threshold = :alert_threshold, is_active = :is_active, updated_at = :updated_at WHERE id = :id AND user_id = :user_id`,
		budget)
	if err != nil {
		return domain.Budget{}, fmt.Errorf("update budget: %w", err)
	}
	if err := requireAffected(result, "update budget"); err != nil {
		return domain.Budget{}, err
	}
	return budget, nil
}

func (r *BudgetRepository) DeleteBudget(ctx context.Context, userID, id uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM budgets WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("delete budget: %w", err)
	}
	return requireAffected(result, "delete budget")
}

// GetBudgetStatus calculates the budget utilisation by summing expense transactions
// for the budget's linked category within the budget's date range.
func (r *BudgetRepository) GetBudgetStatus(ctx context.Context, userID, id uuid.UUID) (domain.BudgetStatus, error) {
	budget, err := r.GetBudget(ctx, userID, id)
	if err != nil {
		return domain.BudgetStatus{}, err
	}

	var spent decimal.Decimal
	query := `SELECT COALESCE(SUM(amount * exchange_rate), 0) FROM transactions
		WHERE user_id = $1 AND type = 'EXPENSE'
		AND occurred_at >= $2 AND occurred_at < $3`
	args := []any{userID, budget.StartDate.UTC(), budget.EndDate.UTC()}
	if budget.CategoryID != nil {
		query += ` AND category_id = $4`
		args = append(args, *budget.CategoryID)
	}
	if err := r.db.GetContext(ctx, &spent, query, args...); err != nil {
		return domain.BudgetStatus{}, fmt.Errorf("get budget spent: %w", err)
	}

	var utilizationRate decimal.Decimal
	if budget.Amount.IsPositive() {
		utilizationRate = spent.Div(budget.Amount).Round(4)
	}

	return domain.BudgetStatus{
		Budget:           budget,
		SpentAmount:      spent,
		UtilizationRate:  utilizationRate,
		IsAlertTriggered: utilizationRate.GreaterThanOrEqual(budget.AlertThreshold),
	}, nil
}

// ListActiveBudgetsForCategory returns active budgets that cover the current time for a category.
func (r *BudgetRepository) ListActiveBudgetsForCategory(ctx context.Context, userID, categoryID uuid.UUID) ([]domain.Budget, error) {
	now := time.Now().UTC()
	budgets := make([]domain.Budget, 0)
	if err := r.db.SelectContext(ctx, &budgets,
		`SELECT `+budgetColumns+` FROM budgets WHERE user_id = $1 AND is_active = TRUE AND category_id = $2 AND start_date <= $3 AND end_date >= $3`,
		userID, categoryID, now); err != nil {
		return nil, fmt.Errorf("list active budgets for category: %w", err)
	}
	return budgets, nil
}

var _ ports.BudgetRepository = (*BudgetRepository)(nil)

// ─────────────────────────────────────────────────────────────────
// SavingsGoalRepository
// ─────────────────────────────────────────────────────────────────

// SavingsGoalRepository implements ports.SavingsGoalRepository.
type SavingsGoalRepository struct{ db *sqlx.DB }

func NewSavingsGoalRepository(db *sqlx.DB) *SavingsGoalRepository {
	return &SavingsGoalRepository{db: db}
}

func (r *SavingsGoalRepository) ListSavingsGoals(ctx context.Context, userID uuid.UUID) ([]domain.SavingsGoal, error) {
	goals := make([]domain.SavingsGoal, 0)
	if err := r.db.SelectContext(ctx, &goals, `SELECT `+savingsGoalColumns+` FROM savings_goals WHERE user_id = $1 ORDER BY created_at DESC`, userID); err != nil {
		return nil, fmt.Errorf("list savings goals: %w", err)
	}
	return goals, nil
}

func (r *SavingsGoalRepository) GetSavingsGoal(ctx context.Context, userID, id uuid.UUID) (domain.SavingsGoal, error) {
	var goal domain.SavingsGoal
	if err := r.db.GetContext(ctx, &goal, `SELECT `+savingsGoalColumns+` FROM savings_goals WHERE id = $1 AND user_id = $2`, id, userID); err != nil {
		return domain.SavingsGoal{}, mapNotFound(err, "get savings goal")
	}
	return goal, nil
}

func (r *SavingsGoalRepository) CreateSavingsGoal(ctx context.Context, goal domain.SavingsGoal) (domain.SavingsGoal, error) {
	const query = `INSERT INTO savings_goals (` + savingsGoalColumns + `) VALUES (:id, :user_id, :name, :target_amount, :current_amount, :linked_account_id, :target_date, :color, :icon, :created_at, :updated_at)`
	if _, err := r.db.NamedExecContext(ctx, query, goal); err != nil {
		return domain.SavingsGoal{}, fmt.Errorf("create savings goal: %w", err)
	}
	return goal, nil
}

func (r *SavingsGoalRepository) UpdateSavingsGoal(ctx context.Context, goal domain.SavingsGoal) (domain.SavingsGoal, error) {
	result, err := r.db.NamedExecContext(ctx,
		`UPDATE savings_goals SET name = :name, target_amount = :target_amount, linked_account_id = :linked_account_id, target_date = :target_date, color = :color, icon = :icon, updated_at = :updated_at WHERE id = :id AND user_id = :user_id`,
		goal)
	if err != nil {
		return domain.SavingsGoal{}, fmt.Errorf("update savings goal: %w", err)
	}
	if err := requireAffected(result, "update savings goal"); err != nil {
		return domain.SavingsGoal{}, err
	}
	return goal, nil
}

func (r *SavingsGoalRepository) DeleteSavingsGoal(ctx context.Context, userID, id uuid.UUID) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM savings_goals WHERE id = $1 AND user_id = $2`, id, userID)
	if err != nil {
		return fmt.Errorf("delete savings goal: %w", err)
	}
	return requireAffected(result, "delete savings goal")
}

// UpdateGoalProgress atomically increments the current_amount for a savings goal.
func (r *SavingsGoalRepository) UpdateGoalProgress(ctx context.Context, userID, id uuid.UUID, amount decimal.Decimal) (domain.SavingsGoal, error) {
	var goal domain.SavingsGoal
	err := r.db.GetContext(ctx, &goal,
		`UPDATE savings_goals SET current_amount = current_amount + $1, updated_at = now() WHERE id = $2 AND user_id = $3 RETURNING `+savingsGoalColumns,
		amount, id, userID)
	if err != nil {
		return domain.SavingsGoal{}, mapNotFound(err, "update goal progress")
	}
	return goal, nil
}

var _ ports.SavingsGoalRepository = (*SavingsGoalRepository)(nil)
