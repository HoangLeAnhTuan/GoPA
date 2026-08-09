package http

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gopa/internal/core/domain"
	"gopa/internal/core/ports"
	"gopa/internal/services"
)

func (h *ModuleHandler) registerFinanceRoutes(protected *gin.RouterGroup) {
	protected.GET("/finance/dashboard", h.financeDashboard)
	protected.GET("/finance/net-worth", h.financeNetWorth)
	protected.GET("/accounts", h.listAccounts)
	protected.POST("/accounts", h.createAccount)
	protected.PATCH("/accounts/:id", h.updateAccount)
	protected.DELETE("/accounts/:id", h.deleteAccount)
	protected.GET("/categories", h.listCategories)
	protected.POST("/categories", h.createCategory)
	protected.PATCH("/categories/:id", h.updateCategory)
	protected.DELETE("/categories/:id", h.deleteCategory)
	protected.GET("/transactions/summary", h.transactionSummary)
	protected.GET("/transactions/by-category", h.transactionsByCategory)
	protected.GET("/transactions", h.listTransactions)
	protected.POST("/transactions", h.createTransaction)
	protected.GET("/transactions/:id", h.getTransaction)
	protected.PATCH("/transactions/:id", h.updateTransaction)
	protected.DELETE("/transactions/:id", h.deleteTransaction)
}

func (h *ModuleHandler) financeDashboard(c *gin.Context) {
	value, err := h.finance.Dashboard(c, userID(c))
	respond(c, value, err, http.StatusOK)
}
func (h *ModuleHandler) financeNetWorth(c *gin.Context) {
	value, err := h.finance.NetWorth(c, userID(c))
	respond(c, gin.H{"net_worth": value}, err, http.StatusOK)
}

func (h *ModuleHandler) listAccounts(c *gin.Context) {
	values, err := h.finance.ListAccounts(c, userID(c), c.Query("include_archived") == "true")
	respond(c, values, err, http.StatusOK)
}
func (h *ModuleHandler) createAccount(c *gin.Context) {
	var request accountRequest
	if !bind(c, &request) {
		return
	}
	input, err := request.input()
	if err != nil {
		respond(c, nil, err, http.StatusBadRequest)
		return
	}
	value, err := h.finance.CreateAccount(c, userID(c), input)
	respond(c, value, err, http.StatusCreated)
}
func (h *ModuleHandler) updateAccount(c *gin.Context) {
	id, ok := parameterUUID(c)
	if !ok {
		return
	}
	var request accountRequest
	if !bind(c, &request) {
		return
	}
	input, err := request.input()
	if err != nil {
		respond(c, nil, err, http.StatusBadRequest)
		return
	}
	value, err := h.finance.UpdateAccount(c, userID(c), id, input)
	respond(c, value, err, http.StatusOK)
}
func (h *ModuleHandler) deleteAccount(c *gin.Context) {
	id, ok := parameterUUID(c)
	if !ok {
		return
	}
	respond(c, gin.H{}, h.finance.DeleteAccount(c, userID(c), id), http.StatusOK)
}

func (h *ModuleHandler) listCategories(c *gin.Context) {
	values, err := h.finance.ListCategories(c, userID(c))
	respond(c, values, err, http.StatusOK)
}
func (h *ModuleHandler) createCategory(c *gin.Context) {
	var request categoryRequest
	if !bind(c, &request) {
		return
	}
	input, err := request.input()
	if err != nil {
		respond(c, nil, err, http.StatusBadRequest)
		return
	}
	value, err := h.finance.CreateCategory(c, userID(c), input)
	respond(c, value, err, http.StatusCreated)
}
func (h *ModuleHandler) updateCategory(c *gin.Context) {
	id, ok := parameterUUID(c)
	if !ok {
		return
	}
	var request categoryRequest
	if !bind(c, &request) {
		return
	}
	input, err := request.input()
	if err != nil {
		respond(c, nil, err, http.StatusBadRequest)
		return
	}
	value, err := h.finance.UpdateCategory(c, userID(c), id, input)
	respond(c, value, err, http.StatusOK)
}
func (h *ModuleHandler) deleteCategory(c *gin.Context) {
	id, ok := parameterUUID(c)
	if !ok {
		return
	}
	respond(c, gin.H{}, h.finance.DeleteCategory(c, userID(c), id), http.StatusOK)
}

func (h *ModuleHandler) listTransactions(c *gin.Context) {
	filter, err := transactionFilter(c)
	if err != nil {
		respond(c, nil, err, http.StatusOK)
		return
	}
	values, err := h.finance.ListTransactions(c, userID(c), filter)
	respond(c, values, err, http.StatusOK)
}
func (h *ModuleHandler) getTransaction(c *gin.Context) {
	id, ok := parameterUUID(c)
	if !ok {
		return
	}
	value, err := h.finance.GetTransaction(c, userID(c), id)
	respond(c, value, err, http.StatusOK)
}
func (h *ModuleHandler) createTransaction(c *gin.Context) {
	var request transactionRequest
	if !bind(c, &request) {
		return
	}
	input, err := request.input()
	if err != nil {
		respond(c, nil, err, http.StatusBadRequest)
		return
	}
	value, err := h.finance.CreateTransaction(c, userID(c), input)
	respond(c, value, err, http.StatusCreated)
}
func (h *ModuleHandler) updateTransaction(c *gin.Context) {
	id, ok := parameterUUID(c)
	if !ok {
		return
	}
	var request transactionRequest
	if !bind(c, &request) {
		return
	}
	input, err := request.input()
	if err != nil {
		respond(c, nil, err, http.StatusBadRequest)
		return
	}
	value, err := h.finance.UpdateTransaction(c, userID(c), id, input)
	respond(c, value, err, http.StatusOK)
}
func (h *ModuleHandler) deleteTransaction(c *gin.Context) {
	id, ok := parameterUUID(c)
	if !ok {
		return
	}
	respond(c, gin.H{}, h.finance.DeleteTransaction(c, userID(c), id), http.StatusOK)
}

func (h *ModuleHandler) transactionSummary(c *gin.Context) {
	from, to, err := queryDateRange(c)
	if err != nil {
		respond(c, nil, err, http.StatusOK)
		return
	}
	value, err := h.finance.Cashflow(c, userID(c), from, to)
	respond(c, value, err, http.StatusOK)
}
func (h *ModuleHandler) transactionsByCategory(c *gin.Context) {
	from, to, err := queryDateRange(c)
	if err != nil {
		respond(c, nil, err, http.StatusOK)
		return
	}
	values, err := h.finance.SpendingByCategory(c, userID(c), from, to)
	respond(c, values, err, http.StatusOK)
}

type accountRequest struct {
	Name           string             `json:"name" binding:"required,max=120"`
	Type           domain.AccountType `json:"type" binding:"required"`
	Currency       string             `json:"currency" binding:"required,len=3"`
	InitialBalance string             `json:"initial_balance" binding:"required"`
	Color          string             `json:"color"`
	Icon           string             `json:"icon"`
	IsArchived     bool               `json:"is_archived"`
}

func (r accountRequest) input() (services.AccountInput, error) {
	value, err := decimal.NewFromString(r.InitialBalance)
	if err != nil {
		return services.AccountInput{}, domain.ErrValidation
	}
	return services.AccountInput{Name: r.Name, Type: r.Type, Currency: r.Currency, InitialBalance: value, Color: r.Color, Icon: r.Icon, IsArchived: r.IsArchived}, nil
}

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
	if r.MonthlyBudget != nil {
		value, err := decimal.NewFromString(*r.MonthlyBudget)
		if err != nil {
			return services.CategoryInput{}, domain.ErrValidation
		}
		budget = &value
	}
	return services.CategoryInput{Name: r.Name, Type: r.Type, Icon: r.Icon, Color: r.Color, ParentID: r.ParentID, MonthlyBudget: budget}, nil
}

type transactionRequest struct {
	AccountID    uuid.UUID              `json:"account_id" binding:"required"`
	ToAccountID  *uuid.UUID             `json:"to_account_id"`
	CategoryID   *uuid.UUID             `json:"category_id"`
	Type         domain.TransactionType `json:"type" binding:"required"`
	Amount       string                 `json:"amount" binding:"required"`
	Currency     string                 `json:"currency" binding:"required,len=3"`
	ExchangeRate string                 `json:"exchange_rate" binding:"required"`
	Description  string                 `json:"description"`
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
	exchangeRate, err := decimal.NewFromString(r.ExchangeRate)
	if err != nil {
		return services.TransactionInput{}, domain.ErrValidation
	}
	return services.TransactionInput{AccountID: r.AccountID, ToAccountID: r.ToAccountID, CategoryID: r.CategoryID, Type: r.Type, Amount: amount, Currency: r.Currency, ExchangeRate: exchangeRate, Description: r.Description, Merchant: r.Merchant, Tags: r.Tags, OccurredAt: r.OccurredAt, IsRecurring: r.IsRecurring}, nil
}

func transactionFilter(c *gin.Context) (ports.TransactionFilter, error) {
	filter := ports.TransactionFilter{Query: strings.TrimSpace(c.Query("q")), Limit: limit(c)}
	for key, target := range map[string]**uuid.UUID{"account": &filter.AccountID, "category": &filter.CategoryID} {
		if raw := strings.TrimSpace(c.Query(key)); raw != "" {
			value, err := uuid.Parse(raw)
			if err != nil {
				return ports.TransactionFilter{}, domain.ErrValidation
			}
			*target = &value
		}
	}
	if raw := strings.TrimSpace(c.Query("type")); raw != "" {
		value := domain.TransactionType(raw)
		filter.Type = &value
	}
	for key, target := range map[string]**time.Time{"from": &filter.From, "to": &filter.To} {
		if raw := strings.TrimSpace(c.Query(key)); raw != "" {
			value, err := time.Parse(time.RFC3339, raw)
			if err != nil {
				return ports.TransactionFilter{}, domain.ErrValidation
			}
			*target = &value
		}
	}
	return filter, nil
}

func queryDateRange(c *gin.Context) (time.Time, time.Time, error) {
	now := time.Now().UTC()
	from := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	to := from.AddDate(0, 1, 0)
	if raw := strings.TrimSpace(c.Query("from")); raw != "" {
		value, err := time.Parse("2006-01-02", raw)
		if err != nil {
			return time.Time{}, time.Time{}, domain.ErrValidation
		}
		from = value.UTC()
	}
	if raw := strings.TrimSpace(c.Query("to")); raw != "" {
		value, err := time.Parse("2006-01-02", raw)
		if err != nil {
			return time.Time{}, time.Time{}, domain.ErrValidation
		}
		to = value.AddDate(0, 0, 1).UTC()
	}
	if !from.Before(to) {
		return time.Time{}, time.Time{}, domain.ErrValidation
	}
	return from, to, nil
}
