package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gopa/internal/core/domain"
	"gopa/internal/core/ports"
)

type financeRepoStub struct {
	accounts     map[uuid.UUID]domain.Account
	categories   map[uuid.UUID]domain.Category
	transactions map[uuid.UUID]domain.Transaction
	seeded       bool
}

func newFinanceRepoStub() *financeRepoStub {
	return &financeRepoStub{
		accounts:     make(map[uuid.UUID]domain.Account),
		categories:   make(map[uuid.UUID]domain.Category),
		transactions: make(map[uuid.UUID]domain.Transaction),
	}
}

func (s *financeRepoStub) ListAccounts(_ context.Context, userID uuid.UUID, includeArchived bool) ([]domain.Account, error) {
	var list []domain.Account
	for _, a := range s.accounts {
		if a.UserID == userID && (includeArchived || !a.IsArchived) {
			list = append(list, a)
		}
	}
	return list, nil
}
func (s *financeRepoStub) GetAccount(_ context.Context, userID, id uuid.UUID) (domain.Account, error) {
	a, ok := s.accounts[id]
	if !ok || a.UserID != userID {
		return domain.Account{}, domain.ErrNotFound
	}
	return a, nil
}
func (s *financeRepoStub) CreateAccount(_ context.Context, a domain.Account) (domain.Account, error) {
	s.accounts[a.ID] = a
	return a, nil
}
func (s *financeRepoStub) UpdateAccount(_ context.Context, a domain.Account) (domain.Account, error) {
	if _, ok := s.accounts[a.ID]; !ok {
		return domain.Account{}, domain.ErrNotFound
	}
	s.accounts[a.ID] = a
	return a, nil
}
func (s *financeRepoStub) DeleteAccount(_ context.Context, userID, id uuid.UUID) error {
	a, ok := s.accounts[id]
	if !ok || a.UserID != userID {
		return domain.ErrNotFound
	}
	delete(s.accounts, id)
	return nil
}
func (s *financeRepoStub) ListCategories(_ context.Context, userID uuid.UUID) ([]domain.Category, error) {
	var list []domain.Category
	for _, c := range s.categories {
		if c.UserID == userID {
			list = append(list, c)
		}
	}
	return list, nil
}
func (s *financeRepoStub) GetCategory(_ context.Context, userID, id uuid.UUID) (domain.Category, error) {
	c, ok := s.categories[id]
	if !ok || c.UserID != userID {
		return domain.Category{}, domain.ErrNotFound
	}
	return c, nil
}
func (s *financeRepoStub) CreateCategory(_ context.Context, c domain.Category) (domain.Category, error) {
	s.categories[c.ID] = c
	return c, nil
}
func (s *financeRepoStub) UpdateCategory(_ context.Context, c domain.Category) (domain.Category, error) {
	if _, ok := s.categories[c.ID]; !ok {
		return domain.Category{}, domain.ErrNotFound
	}
	s.categories[c.ID] = c
	return c, nil
}
func (s *financeRepoStub) DeleteCategory(_ context.Context, userID, id uuid.UUID) error {
	c, ok := s.categories[id]
	if !ok || c.UserID != userID {
		return domain.ErrNotFound
	}
	delete(s.categories, id)
	return nil
}
func (s *financeRepoStub) ListTransactions(_ context.Context, userID uuid.UUID, _ ports.TransactionFilter) ([]domain.Transaction, error) {
	var list []domain.Transaction
	for _, t := range s.transactions {
		if t.UserID == userID {
			list = append(list, t)
		}
	}
	return list, nil
}
func (s *financeRepoStub) GetTransaction(_ context.Context, userID, id uuid.UUID) (domain.Transaction, error) {
	t, ok := s.transactions[id]
	if !ok || t.UserID != userID {
		return domain.Transaction{}, domain.ErrNotFound
	}
	return t, nil
}
func (s *financeRepoStub) CreateTransaction(_ context.Context, t domain.Transaction) (domain.Transaction, error) {
	s.transactions[t.ID] = t
	return t, nil
}
func (s *financeRepoStub) UpdateTransaction(_ context.Context, t domain.Transaction) (domain.Transaction, error) {
	if _, ok := s.transactions[t.ID]; !ok {
		return domain.Transaction{}, domain.ErrNotFound
	}
	s.transactions[t.ID] = t
	return t, nil
}
func (s *financeRepoStub) DeleteTransaction(_ context.Context, userID, id uuid.UUID, _ time.Time) error {
	t, ok := s.transactions[id]
	if !ok || t.UserID != userID {
		return domain.ErrNotFound
	}
	delete(s.transactions, id)
	return nil
}
func (s *financeRepoStub) Cashflow(_ context.Context, _ uuid.UUID, _, _ time.Time) (domain.CashflowSummary, error) {
	return domain.CashflowSummary{
		Income:      decimal.NewFromInt(5000),
		Expense:     decimal.NewFromInt(3000),
		Net:         decimal.NewFromInt(2000),
		SavingsRate: decimal.NewFromInt(40),
	}, nil
}
func (s *financeRepoStub) SpendingByCategory(_ context.Context, _ uuid.UUID, _, _ time.Time) ([]domain.CategorySpending, error) {
	return []domain.CategorySpending{
		{CategoryName: "Food", Color: "#F59E0B", Amount: decimal.NewFromInt(1500)},
	}, nil
}
func (s *financeRepoStub) NetWorth(_ context.Context, _ uuid.UUID) (decimal.Decimal, error) {
	return decimal.NewFromInt(25000), nil
}
func (s *financeRepoStub) SeedDefaultCategories(_ context.Context, _ uuid.UUID) error {
	s.seeded = true
	return nil
}

type budgetRepoStub struct {
	budgets map[uuid.UUID]domain.Budget
}

func newBudgetRepoStub() *budgetRepoStub {
	return &budgetRepoStub{budgets: make(map[uuid.UUID]domain.Budget)}
}
func (s *budgetRepoStub) ListBudgets(_ context.Context, userID uuid.UUID, activeOnly bool) ([]domain.Budget, error) {
	var list []domain.Budget
	for _, b := range s.budgets {
		if b.UserID == userID && (!activeOnly || b.IsActive) {
			list = append(list, b)
		}
	}
	return list, nil
}
func (s *budgetRepoStub) GetBudget(_ context.Context, userID, id uuid.UUID) (domain.Budget, error) {
	b, ok := s.budgets[id]
	if !ok || b.UserID != userID {
		return domain.Budget{}, domain.ErrNotFound
	}
	return b, nil
}
func (s *budgetRepoStub) CreateBudget(_ context.Context, b domain.Budget) (domain.Budget, error) {
	s.budgets[b.ID] = b
	return b, nil
}
func (s *budgetRepoStub) UpdateBudget(_ context.Context, b domain.Budget) (domain.Budget, error) {
	if _, ok := s.budgets[b.ID]; !ok {
		return domain.Budget{}, domain.ErrNotFound
	}
	s.budgets[b.ID] = b
	return b, nil
}
func (s *budgetRepoStub) DeleteBudget(_ context.Context, userID, id uuid.UUID) error {
	b, ok := s.budgets[id]
	if !ok || b.UserID != userID {
		return domain.ErrNotFound
	}
	delete(s.budgets, id)
	return nil
}
func (s *budgetRepoStub) GetBudgetStatus(_ context.Context, userID, id uuid.UUID) (domain.BudgetStatus, error) {
	b, ok := s.budgets[id]
	if !ok || b.UserID != userID {
		return domain.BudgetStatus{}, domain.ErrNotFound
	}
	spent := decimal.NewFromInt(850)
	rate := spent.Div(b.Amount).Round(4)
	return domain.BudgetStatus{
		Budget:           b,
		SpentAmount:      spent,
		UtilizationRate:  rate,
		IsAlertTriggered: rate.GreaterThanOrEqual(b.AlertThreshold),
	}, nil
}
func (s *budgetRepoStub) ListActiveBudgetsForCategory(_ context.Context, _ uuid.UUID, _ uuid.UUID) ([]domain.Budget, error) {
	return nil, nil
}

type savingsGoalRepoStub struct {
	goals map[uuid.UUID]domain.SavingsGoal
}

func newSavingsGoalRepoStub() *savingsGoalRepoStub {
	return &savingsGoalRepoStub{goals: make(map[uuid.UUID]domain.SavingsGoal)}
}
func (s *savingsGoalRepoStub) ListSavingsGoals(_ context.Context, userID uuid.UUID) ([]domain.SavingsGoal, error) {
	var list []domain.SavingsGoal
	for _, g := range s.goals {
		if g.UserID == userID {
			list = append(list, g)
		}
	}
	return list, nil
}
func (s *savingsGoalRepoStub) GetSavingsGoal(_ context.Context, userID, id uuid.UUID) (domain.SavingsGoal, error) {
	g, ok := s.goals[id]
	if !ok || g.UserID != userID {
		return domain.SavingsGoal{}, domain.ErrNotFound
	}
	return g, nil
}
func (s *savingsGoalRepoStub) CreateSavingsGoal(_ context.Context, g domain.SavingsGoal) (domain.SavingsGoal, error) {
	s.goals[g.ID] = g
	return g, nil
}
func (s *savingsGoalRepoStub) UpdateSavingsGoal(_ context.Context, g domain.SavingsGoal) (domain.SavingsGoal, error) {
	if _, ok := s.goals[g.ID]; !ok {
		return domain.SavingsGoal{}, domain.ErrNotFound
	}
	s.goals[g.ID] = g
	return g, nil
}
func (s *savingsGoalRepoStub) DeleteSavingsGoal(_ context.Context, userID, id uuid.UUID) error {
	g, ok := s.goals[id]
	if !ok || g.UserID != userID {
		return domain.ErrNotFound
	}
	delete(s.goals, id)
	return nil
}
func (s *savingsGoalRepoStub) UpdateGoalProgress(_ context.Context, userID, id uuid.UUID, amount decimal.Decimal) (domain.SavingsGoal, error) {
	g, ok := s.goals[id]
	if !ok || g.UserID != userID {
		return domain.SavingsGoal{}, domain.ErrNotFound
	}
	g.CurrentAmount = g.CurrentAmount.Add(amount)
	s.goals[id] = g
	return g, nil
}

func TestFinanceService_BudgetCRUDAndStatus(t *testing.T) {
	userID := uuid.New()
	financeRepo := newFinanceRepoStub()
	budgetRepo := newBudgetRepoStub()
	savingsRepo := newSavingsGoalRepoStub()
	svc := NewFinanceService(financeRepo, budgetRepo, savingsRepo)

	now := time.Now().UTC()
	start := now.AddDate(0, 0, -15)
	end := now.AddDate(0, 0, 15)

	// 1. Create Budget
	budget, err := svc.CreateBudget(context.Background(), userID, BudgetInput{
		Name:           "Monthly Groceries",
		Amount:         decimal.NewFromInt(1000),
		Period:         domain.BudgetPeriodMonthly,
		StartDate:      start,
		EndDate:        end,
		AlertThreshold: decimal.NewFromFloat(0.80),
		IsActive:       true,
	})
	if err != nil {
		t.Fatalf("create budget: %v", err)
	}
	if budget.Name != "Monthly Groceries" || !budget.IsActive {
		t.Fatalf("unexpected budget: %+v", budget)
	}

	// 2. Get Budget
	got, err := svc.GetBudget(context.Background(), userID, budget.ID)
	if err != nil {
		t.Fatalf("get budget: %v", err)
	}
	if got.ID != budget.ID {
		t.Fatalf("expected ID %s, got %s", budget.ID, got.ID)
	}

	// 3. Update Budget
	updated, err := svc.UpdateBudget(context.Background(), userID, budget.ID, BudgetInput{
		Name:           "Groceries & Essentials",
		Amount:         decimal.NewFromInt(1200),
		Period:         domain.BudgetPeriodMonthly,
		StartDate:      start,
		EndDate:        end,
		AlertThreshold: decimal.NewFromFloat(0.75),
		IsActive:       true,
	})
	if err != nil {
		t.Fatalf("update budget: %v", err)
	}
	if updated.Name != "Groceries & Essentials" || !updated.Amount.Equal(decimal.NewFromInt(1200)) {
		t.Fatalf("unexpected updated budget: %+v", updated)
	}

	// 4. Get Budget Status
	status, err := svc.GetBudgetStatus(context.Background(), userID, budget.ID)
	if err != nil {
		t.Fatalf("get budget status: %v", err)
	}
	if !status.SpentAmount.Equal(decimal.NewFromInt(850)) {
		t.Fatalf("expected spent 850, got %s", status.SpentAmount)
	}

	// 5. List Budgets
	list, err := svc.ListBudgets(context.Background(), userID, true)
	if err != nil {
		t.Fatalf("list budgets: %v", err)
	}
	if len(list) != 1 {
		t.Fatalf("expected 1 budget, got %d", len(list))
	}

	// 6. Delete Budget
	if err := svc.DeleteBudget(context.Background(), userID, budget.ID); err != nil {
		t.Fatalf("delete budget: %v", err)
	}
	_, err = svc.GetBudget(context.Background(), userID, budget.ID)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestFinanceService_SavingsGoalCRUDAndProgress(t *testing.T) {
	userID := uuid.New()
	financeRepo := newFinanceRepoStub()
	budgetRepo := newBudgetRepoStub()
	savingsRepo := newSavingsGoalRepoStub()
	svc := NewFinanceService(financeRepo, budgetRepo, savingsRepo)

	// 1. Create Savings Goal
	goal, err := svc.CreateSavingsGoal(context.Background(), userID, SavingsGoalInput{
		Name:          "Emergency Fund",
		TargetAmount:  decimal.NewFromInt(10000),
		CurrentAmount: decimal.NewFromInt(2000),
		Color:         "#10B981",
		Icon:          "shield",
	})
	if err != nil {
		t.Fatalf("create savings goal: %v", err)
	}
	if goal.Name != "Emergency Fund" || !goal.CurrentAmount.Equal(decimal.NewFromInt(2000)) {
		t.Fatalf("unexpected goal: %+v", goal)
	}

	// 2. Update Goal Progress
	progressed, err := svc.UpdateGoalProgress(context.Background(), userID, goal.ID, decimal.NewFromInt(1500))
	if err != nil {
		t.Fatalf("update goal progress: %v", err)
	}
	if !progressed.CurrentAmount.Equal(decimal.NewFromInt(3500)) {
		t.Fatalf("expected current amount 3500, got %s", progressed.CurrentAmount)
	}

	// 3. Update Progress with zero (invalid)
	_, err = svc.UpdateGoalProgress(context.Background(), userID, goal.ID, decimal.Zero)
	if !errors.Is(err, domain.ErrValidation) {
		t.Fatalf("expected ErrValidation for 0 progress, got %v", err)
	}

	// 4. Update Savings Goal
	updated, err := svc.UpdateSavingsGoal(context.Background(), userID, goal.ID, SavingsGoalInput{
		Name:          "Emergency Fund Tier 1",
		TargetAmount:  decimal.NewFromInt(15000),
		CurrentAmount: decimal.NewFromInt(3500),
		Color:         "#22C55E",
		Icon:          "shield-check",
	})
	if err != nil {
		t.Fatalf("update goal: %v", err)
	}
	if updated.Name != "Emergency Fund Tier 1" {
		t.Fatalf("unexpected name %q", updated.Name)
	}

	// 5. Delete Savings Goal
	if err := svc.DeleteSavingsGoal(context.Background(), userID, goal.ID); err != nil {
		t.Fatalf("delete goal: %v", err)
	}
	_, err = svc.GetSavingsGoal(context.Background(), userID, goal.ID)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestFinanceService_SeedDefaultCategories(t *testing.T) {
	userID := uuid.New()
	financeRepo := newFinanceRepoStub()
	budgetRepo := newBudgetRepoStub()
	savingsRepo := newSavingsGoalRepoStub()
	svc := NewFinanceService(financeRepo, budgetRepo, savingsRepo)

	if err := svc.SeedDefaultCategories(context.Background(), userID); err != nil {
		t.Fatalf("seed categories: %v", err)
	}
	if !financeRepo.seeded {
		t.Fatalf("expected seeded to be true")
	}
}

func TestFinanceService_DashboardAndCashflow(t *testing.T) {
	userID := uuid.New()
	financeRepo := newFinanceRepoStub()
	budgetRepo := newBudgetRepoStub()
	savingsRepo := newSavingsGoalRepoStub()
	svc := NewFinanceService(financeRepo, budgetRepo, savingsRepo)

	dashboard, err := svc.Dashboard(context.Background(), userID)
	if err != nil {
		t.Fatalf("dashboard: %v", err)
	}
	if !dashboard.NetWorth.Equal(decimal.NewFromInt(25000)) {
		t.Fatalf("expected net worth 25000, got %s", dashboard.NetWorth)
	}
}
