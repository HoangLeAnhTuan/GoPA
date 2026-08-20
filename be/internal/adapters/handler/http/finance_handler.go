package http

import (
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gopa/internal/core/domain"
	"gopa/internal/core/ports"
	"gopa/internal/services"
	"gopa/pkg/utils/response"
)

type FinanceHandler struct {
	finance *services.FinanceService
}

func NewFinanceHandler(finance *services.FinanceService) *FinanceHandler {
	return &FinanceHandler{finance: finance}
}

func (h *FinanceHandler) RegisterRoutes(api *gin.RouterGroup, authenticate gin.HandlerFunc) {
	fin := api.Group("", authenticate)

	// Dashboard & Analytics
	fin.GET("/finance/dashboard", h.dashboard)
	fin.GET("/finance/net-worth", h.netWorth)
	fin.GET("/finance/cashflow", h.cashflow)
	fin.GET("/finance/spending", h.spending)

	// Accounts
	fin.GET("/accounts", h.listAccounts)
	fin.POST("/accounts", h.createAccount)
	fin.PATCH("/accounts/:id", h.updateAccount)
	fin.DELETE("/accounts/:id", h.deleteAccount)

	// Categories
	fin.GET("/categories", h.listCategories)
	fin.POST("/categories", h.createCategory)
	fin.PATCH("/categories/:id", h.updateCategory)
	fin.DELETE("/categories/:id", h.deleteCategory)
	fin.POST("/categories/seed", h.seedCategories)

	// Transactions
	fin.GET("/transactions", h.listTransactions)
	fin.POST("/transactions", h.createTransaction)
	fin.GET("/transactions/:id", h.getTransaction)
	fin.PATCH("/transactions/:id", h.updateTransaction)
	fin.DELETE("/transactions/:id", h.deleteTransaction)

	// Budgets
	fin.GET("/budgets", h.listBudgets)
	fin.POST("/budgets", h.createBudget)
	fin.GET("/budgets/:id", h.getBudget)
	fin.PATCH("/budgets/:id", h.updateBudget)
	fin.DELETE("/budgets/:id", h.deleteBudget)
	fin.GET("/budgets/:id/status", h.getBudgetStatus)

	// Savings Goals
	fin.GET("/savings-goals", h.listSavingsGoals)
	fin.POST("/savings-goals", h.createSavingsGoal)
	fin.GET("/savings-goals/:id", h.getSavingsGoal)
	fin.PATCH("/savings-goals/:id", h.updateSavingsGoal)
	fin.DELETE("/savings-goals/:id", h.deleteSavingsGoal)
	fin.POST("/savings-goals/:id/progress", h.updateGoalProgress)
}

// --- Dashboard & Analytics ---

func (h *FinanceHandler) dashboard(c *gin.Context) {
	val, err := h.finance.Dashboard(c.Request.Context(), userID(c))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, val)
}

func (h *FinanceHandler) netWorth(c *gin.Context) {
	val, err := h.finance.NetWorth(c.Request.Context(), userID(c))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"net_worth": val})
}

func (h *FinanceHandler) cashflow(c *gin.Context) {
	from, to, ok := parseDateRange(c)
	if !ok {
		return
	}
	val, err := h.finance.Cashflow(c.Request.Context(), userID(c), from, to)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, val)
}

func (h *FinanceHandler) spending(c *gin.Context) {
	from, to, ok := parseDateRange(c)
	if !ok {
		return
	}
	val, err := h.finance.SpendingByCategory(c.Request.Context(), userID(c), from, to)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, val)
}

// --- Accounts ---

type accountRequest struct {
	Name           string             `json:"name" binding:"required,max=120"`
	Type           domain.AccountType `json:"type" binding:"required"`
	Currency       string             `json:"currency" binding:"required,len=3"`
	InitialBalance string             `json:"initial_balance"`
	Color          string             `json:"color"`
	Icon           string             `json:"icon"`
	IsArchived     bool               `json:"is_archived"`
}

func (r accountRequest) input() (services.AccountInput, error) {
	balance := decimal.Zero
	if r.InitialBalance != "" {
		parsed, err := decimal.NewFromString(r.InitialBalance)
		if err != nil {
			return services.AccountInput{}, domain.ErrValidation
		}
		balance = parsed
	}
	return services.AccountInput{
		Name:           r.Name,
		Type:           r.Type,
		Currency:       r.Currency,
		InitialBalance: balance,
		Color:          r.Color,
		Icon:           r.Icon,
		IsArchived:     r.IsArchived,
	}, nil
}

func (h *FinanceHandler) listAccounts(c *gin.Context) {
	archived := c.Query("include_archived") == "true"
	items, err := h.finance.ListAccounts(c.Request.Context(), userID(c), archived)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, items)
}

func (h *FinanceHandler) createAccount(c *gin.Context) {
	var req accountRequest
	if !bindJSON(c, &req) {
		return
	}
	inp, err := req.input()
	if err != nil {
		response.Fail(c, err)
		return
	}
	val, err := h.finance.CreateAccount(c.Request.Context(), userID(c), inp)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, val)
}

func (h *FinanceHandler) updateAccount(c *gin.Context) {
	id, ok := parseParamUUID(c, "id")
	if !ok {
		return
	}
	var req accountRequest
	if !bindJSON(c, &req) {
		return
	}
	inp, err := req.input()
	if err != nil {
		response.Fail(c, err)
		return
	}
	val, err := h.finance.UpdateAccount(c.Request.Context(), userID(c), id, inp)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, val)
}

func (h *FinanceHandler) deleteAccount(c *gin.Context) {
	id, ok := parseParamUUID(c, "id")
	if !ok {
		return
	}
	if err := h.finance.DeleteAccount(c.Request.Context(), userID(c), id); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"id": id})
}

// --- Categories ---

type categoryRequest struct {
	Name          string              `json:"name" binding:"required,max=100"`
	Type          domain.CategoryType `json:"type" binding:"required"`
	Icon          string              `json:"icon"`
	Color         string              `json:"color"`
	ParentID      *uuid.UUID          `json:"parent_id"`
	MonthlyBudget *string             `json:"monthly_budget"`
}

func (r categoryRequest) input() (services.CategoryInput, error) {
	var budget *decimal.Decimal
	if r.MonthlyBudget != nil && *r.MonthlyBudget != "" {
		parsed, err := decimal.NewFromString(*r.MonthlyBudget)
		if err != nil {
			return services.CategoryInput{}, domain.ErrValidation
		}
		budget = &parsed
	}
	return services.CategoryInput{
		Name:          r.Name,
		Type:          r.Type,
		Icon:          r.Icon,
		Color:         r.Color,
		ParentID:      r.ParentID,
		MonthlyBudget: budget,
	}, nil
}

func (h *FinanceHandler) listCategories(c *gin.Context) {
	items, err := h.finance.ListCategories(c.Request.Context(), userID(c))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, items)
}

func (h *FinanceHandler) createCategory(c *gin.Context) {
	var req categoryRequest
	if !bindJSON(c, &req) {
		return
	}
	inp, err := req.input()
	if err != nil {
		response.Fail(c, err)
		return
	}
	val, err := h.finance.CreateCategory(c.Request.Context(), userID(c), inp)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, val)
}

func (h *FinanceHandler) updateCategory(c *gin.Context) {
	id, ok := parseParamUUID(c, "id")
	if !ok {
		return
	}
	var req categoryRequest
	if !bindJSON(c, &req) {
		return
	}
	inp, err := req.input()
	if err != nil {
		response.Fail(c, err)
		return
	}
	val, err := h.finance.UpdateCategory(c.Request.Context(), userID(c), id, inp)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, val)
}

func (h *FinanceHandler) deleteCategory(c *gin.Context) {
	id, ok := parseParamUUID(c, "id")
	if !ok {
		return
	}
	if err := h.finance.DeleteCategory(c.Request.Context(), userID(c), id); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"id": id})
}

func (h *FinanceHandler) seedCategories(c *gin.Context) {
	if err := h.finance.SeedDefaultCategories(c.Request.Context(), userID(c)); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"seeded": true})
}

// --- Transactions ---

type transactionRequest struct {
	AccountID    uuid.UUID              `json:"account_id" binding:"required"`
	ToAccountID  *uuid.UUID             `json:"to_account_id"`
	CategoryID   *uuid.UUID             `json:"category_id"`
	Type         domain.TransactionType `json:"type" binding:"required"`
	Amount       string                 `json:"amount" binding:"required"`
	Currency     string                 `json:"currency" binding:"required,len=3"`
	ExchangeRate string                 `json:"exchange_rate"`
	Description  string                 `json:"description" binding:"max=500"`
	Merchant     *string                `json:"merchant"`
	Tags         []string               `json:"tags"`
	OccurredAt   time.Time              `json:"occurred_at" binding:"required"`
	IsRecurring  bool                   `json:"is_recurring"`
}

func (r transactionRequest) input() (services.TransactionInput, error) {
	amount, err := decimal.NewFromString(r.Amount)
	if err != nil {
		return services.TransactionInput{}, domain.ErrValidation
	}
	rate := decimal.NewFromInt(1)
	if r.ExchangeRate != "" {
		parsed, err := decimal.NewFromString(r.ExchangeRate)
		if err != nil {
			return services.TransactionInput{}, domain.ErrValidation
		}
		rate = parsed
	}
	return services.TransactionInput{
		AccountID:    r.AccountID,
		ToAccountID:  r.ToAccountID,
		CategoryID:   r.CategoryID,
		Type:         r.Type,
		Amount:       amount,
		Currency:     r.Currency,
		ExchangeRate: rate,
		Description:  r.Description,
		Merchant:     r.Merchant,
		Tags:         r.Tags,
		OccurredAt:   r.OccurredAt,
		IsRecurring:  r.IsRecurring,
	}, nil
}

func (h *FinanceHandler) listTransactions(c *gin.Context) {
	filter := ports.TransactionFilter{
		Limit: queryLimit(c, 20),
		Query: strings.TrimSpace(c.Query("q")),
	}
	if val := c.Query("account_id"); val != "" {
		if id, err := uuid.Parse(val); err == nil {
			filter.AccountID = &id
		}
	}
	if val := c.Query("category_id"); val != "" {
		if id, err := uuid.Parse(val); err == nil {
			filter.CategoryID = &id
		}
	}
	if val := c.Query("type"); val != "" {
		t := domain.TransactionType(val)
		filter.Type = &t
	}
	if val := c.Query("from"); val != "" {
		if parsed, err := time.Parse("2006-01-02", val); err == nil {
			filter.From = &parsed
		}
	}
	if val := c.Query("to"); val != "" {
		if parsed, err := time.Parse("2006-01-02", val); err == nil {
			filter.To = &parsed
		}
	}
	items, err := h.finance.ListTransactions(c.Request.Context(), userID(c), filter)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, items)
}

func (h *FinanceHandler) getTransaction(c *gin.Context) {
	id, ok := parseParamUUID(c, "id")
	if !ok {
		return
	}
	item, err := h.finance.GetTransaction(c.Request.Context(), userID(c), id)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, item)
}

func (h *FinanceHandler) createTransaction(c *gin.Context) {
	var req transactionRequest
	if !bindJSON(c, &req) {
		return
	}
	inp, err := req.input()
	if err != nil {
		response.Fail(c, err)
		return
	}
	val, err := h.finance.CreateTransaction(c.Request.Context(), userID(c), inp)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, val)
}

func (h *FinanceHandler) updateTransaction(c *gin.Context) {
	id, ok := parseParamUUID(c, "id")
	if !ok {
		return
	}
	var req transactionRequest
	if !bindJSON(c, &req) {
		return
	}
	inp, err := req.input()
	if err != nil {
		response.Fail(c, err)
		return
	}
	val, err := h.finance.UpdateTransaction(c.Request.Context(), userID(c), id, inp)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, val)
}

func (h *FinanceHandler) deleteTransaction(c *gin.Context) {
	id, ok := parseParamUUID(c, "id")
	if !ok {
		return
	}
	if err := h.finance.DeleteTransaction(c.Request.Context(), userID(c), id); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"id": id})
}

// --- Budgets ---

type budgetRequest struct {
	CategoryID     *uuid.UUID          `json:"category_id"`
	Name           string              `json:"name" binding:"required,max=150"`
	Amount         string              `json:"amount" binding:"required"`
	Period         domain.BudgetPeriod `json:"period" binding:"required"`
	StartDate      time.Time           `json:"start_date" binding:"required"`
	EndDate        time.Time           `json:"end_date" binding:"required"`
	AlertThreshold *string             `json:"alert_threshold"`
	IsActive       bool                `json:"is_active"`
}

func (r budgetRequest) input() (services.BudgetInput, error) {
	amount, err := decimal.NewFromString(r.Amount)
	if err != nil {
		return services.BudgetInput{}, domain.ErrValidation
	}
	threshold := decimal.NewFromFloat(0.80)
	if r.AlertThreshold != nil && *r.AlertThreshold != "" {
		parsed, err := decimal.NewFromString(*r.AlertThreshold)
		if err != nil {
			return services.BudgetInput{}, domain.ErrValidation
		}
		threshold = parsed
	}
	return services.BudgetInput{
		CategoryID:     r.CategoryID,
		Name:           r.Name,
		Amount:         amount,
		Period:         r.Period,
		StartDate:      r.StartDate,
		EndDate:        r.EndDate,
		AlertThreshold: threshold,
		IsActive:       r.IsActive,
	}, nil
}

func (h *FinanceHandler) listBudgets(c *gin.Context) {
	activeOnly := c.Query("active_only") == "true"
	items, err := h.finance.ListBudgets(c.Request.Context(), userID(c), activeOnly)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, items)
}

func (h *FinanceHandler) getBudget(c *gin.Context) {
	id, ok := parseParamUUID(c, "id")
	if !ok {
		return
	}
	item, err := h.finance.GetBudget(c.Request.Context(), userID(c), id)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, item)
}

func (h *FinanceHandler) createBudget(c *gin.Context) {
	var req budgetRequest
	if !bindJSON(c, &req) {
		return
	}
	inp, err := req.input()
	if err != nil {
		response.Fail(c, err)
		return
	}
	val, err := h.finance.CreateBudget(c.Request.Context(), userID(c), inp)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, val)
}

func (h *FinanceHandler) updateBudget(c *gin.Context) {
	id, ok := parseParamUUID(c, "id")
	if !ok {
		return
	}
	var req budgetRequest
	if !bindJSON(c, &req) {
		return
	}
	inp, err := req.input()
	if err != nil {
		response.Fail(c, err)
		return
	}
	val, err := h.finance.UpdateBudget(c.Request.Context(), userID(c), id, inp)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, val)
}

func (h *FinanceHandler) deleteBudget(c *gin.Context) {
	id, ok := parseParamUUID(c, "id")
	if !ok {
		return
	}
	if err := h.finance.DeleteBudget(c.Request.Context(), userID(c), id); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"id": id})
}

func (h *FinanceHandler) getBudgetStatus(c *gin.Context) {
	id, ok := parseParamUUID(c, "id")
	if !ok {
		return
	}
	status, err := h.finance.GetBudgetStatus(c.Request.Context(), userID(c), id)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, status)
}

// --- Savings Goals ---

type savingsGoalRequest struct {
	Name            string     `json:"name" binding:"required,max=150"`
	TargetAmount    string     `json:"target_amount" binding:"required"`
	CurrentAmount   string     `json:"current_amount"`
	LinkedAccountID *uuid.UUID `json:"linked_account_id"`
	TargetDate      *time.Time `json:"target_date"`
	Color           string     `json:"color"`
	Icon            string     `json:"icon"`
}

func (r savingsGoalRequest) input() (services.SavingsGoalInput, error) {
	target, err := decimal.NewFromString(r.TargetAmount)
	if err != nil {
		return services.SavingsGoalInput{}, domain.ErrValidation
	}
	current := decimal.Zero
	if r.CurrentAmount != "" {
		parsed, err := decimal.NewFromString(r.CurrentAmount)
		if err != nil {
			return services.SavingsGoalInput{}, domain.ErrValidation
		}
		current = parsed
	}
	return services.SavingsGoalInput{
		Name:            r.Name,
		TargetAmount:    target,
		CurrentAmount:   current,
		LinkedAccountID: r.LinkedAccountID,
		TargetDate:      r.TargetDate,
		Color:           r.Color,
		Icon:            r.Icon,
	}, nil
}

func (h *FinanceHandler) listSavingsGoals(c *gin.Context) {
	items, err := h.finance.ListSavingsGoals(c.Request.Context(), userID(c))
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, items)
}

func (h *FinanceHandler) getSavingsGoal(c *gin.Context) {
	id, ok := parseParamUUID(c, "id")
	if !ok {
		return
	}
	item, err := h.finance.GetSavingsGoal(c.Request.Context(), userID(c), id)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, item)
}

func (h *FinanceHandler) createSavingsGoal(c *gin.Context) {
	var req savingsGoalRequest
	if !bindJSON(c, &req) {
		return
	}
	inp, err := req.input()
	if err != nil {
		response.Fail(c, err)
		return
	}
	val, err := h.finance.CreateSavingsGoal(c.Request.Context(), userID(c), inp)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.Created(c, val)
}

func (h *FinanceHandler) updateSavingsGoal(c *gin.Context) {
	id, ok := parseParamUUID(c, "id")
	if !ok {
		return
	}
	var req savingsGoalRequest
	if !bindJSON(c, &req) {
		return
	}
	inp, err := req.input()
	if err != nil {
		response.Fail(c, err)
		return
	}
	val, err := h.finance.UpdateSavingsGoal(c.Request.Context(), userID(c), id, inp)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, val)
}

func (h *FinanceHandler) deleteSavingsGoal(c *gin.Context) {
	id, ok := parseParamUUID(c, "id")
	if !ok {
		return
	}
	if err := h.finance.DeleteSavingsGoal(c.Request.Context(), userID(c), id); err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, gin.H{"id": id})
}

type progressRequest struct {
	Amount string `json:"amount" binding:"required"`
}

func (h *FinanceHandler) updateGoalProgress(c *gin.Context) {
	id, ok := parseParamUUID(c, "id")
	if !ok {
		return
	}
	var req progressRequest
	if !bindJSON(c, &req) {
		return
	}
	amt, err := decimal.NewFromString(req.Amount)
	if err != nil {
		response.Fail(c, domain.ErrValidation)
		return
	}
	val, err := h.finance.UpdateGoalProgress(c.Request.Context(), userID(c), id, amt)
	if err != nil {
		response.Fail(c, err)
		return
	}
	response.OK(c, val)
}

func parseDateRange(c *gin.Context) (time.Time, time.Time, bool) {
	now := time.Now().UTC()
	fromStr := c.Query("from")
	toStr := c.Query("to")

	from := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	to := from.AddDate(0, 1, 0)

	if fromStr != "" {
		parsed, err := time.Parse("2006-01-02", fromStr)
		if err != nil {
			response.Fail(c, domain.ErrValidation)
			return time.Time{}, time.Time{}, false
		}
		from = parsed.UTC()
	}
	if toStr != "" {
		parsed, err := time.Parse("2006-01-02", toStr)
		if err != nil {
			response.Fail(c, domain.ErrValidation)
			return time.Time{}, time.Time{}, false
		}
		to = parsed.UTC()
	}
	return from, to, true
}
