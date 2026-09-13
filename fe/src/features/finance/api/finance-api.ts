import { apiClient } from "../../../lib/api-client";
import { unwrap, unwrapCollection, type ApiEnvelope } from "../../../lib/api-envelope";
import type {
  AccountDto,
  AccountInput,
  Budget,
  BudgetInput,
  BudgetStatus,
  CategoryDto,
  CategoryInput,
  FinanceDashboardDto,
  SavingsGoal,
  SavingsGoalInput,
  TransactionDto,
  TransactionFilters,
  TransactionInput,
} from "../types";

export async function getDashboard(): Promise<FinanceDashboardDto> {
  return unwrap(await apiClient.get<ApiEnvelope<FinanceDashboardDto>>("/finance/dashboard"));
}

export async function getNetWorth(): Promise<{ net_worth: string }> {
  return unwrap(await apiClient.get<ApiEnvelope<{ net_worth: string }>>("/finance/net-worth"));
}

export async function listAccounts(): Promise<AccountDto[]> {
  return unwrapCollection(await apiClient.get<ApiEnvelope<AccountDto[] | null>>("/accounts"));
}

export async function createAccount(input: AccountInput): Promise<AccountDto> {
  return unwrap(
    await apiClient.post<ApiEnvelope<AccountDto>>("/accounts", {
      name: input.name,
      type: input.type,
      currency: input.currency,
      initial_balance: input.initialBalance,
      color: input.color,
      icon: input.icon,
      is_archived: input.isArchived,
    })
  );
}

export async function listCategories(): Promise<CategoryDto[]> {
  return unwrapCollection(await apiClient.get<ApiEnvelope<CategoryDto[] | null>>("/categories"));
}

export async function createCategory(input: CategoryInput): Promise<CategoryDto> {
  return unwrap(
    await apiClient.post<ApiEnvelope<CategoryDto>>("/categories", {
      name: input.name,
      type: input.type,
      icon: input.icon,
      color: input.color,
      parent_id: input.parentID,
      monthly_budget: input.monthlyBudget,
    })
  );
}

export async function seedCategories(): Promise<{ seeded_count: number }> {
  return unwrap(await apiClient.post<ApiEnvelope<{ seeded_count: number }>>("/categories/seed"));
}

export async function listTransactions(
  filters: TransactionFilters = {}
): Promise<TransactionDto[]> {
  return unwrapCollection(
    await apiClient.get<ApiEnvelope<TransactionDto[] | null>>("/transactions", { params: filters })
  );
}

export async function createTransaction(input: TransactionInput): Promise<TransactionDto> {
  return unwrap(
    await apiClient.post<ApiEnvelope<TransactionDto>>("/transactions", {
      account_id: input.accountID,
      to_account_id: input.toAccountID,
      category_id: input.categoryID,
      type: input.type,
      amount: input.amount,
      currency: input.currency,
      exchange_rate: input.exchangeRate,
      description: input.description,
      merchant: input.merchant,
      tags: input.tags,
      occurred_at: input.occurredAt,
      is_recurring: input.isRecurring,
    })
  );
}

// --- Budgets APIs ---

export async function listBudgets(): Promise<Budget[]> {
  return unwrapCollection(await apiClient.get<ApiEnvelope<Budget[] | null>>("/budgets"));
}

export async function getBudget(id: string): Promise<Budget> {
  return unwrap(await apiClient.get<ApiEnvelope<Budget>>(`/budgets/${id}`));
}

export async function createBudget(input: BudgetInput): Promise<Budget> {
  return unwrap(
    await apiClient.post<ApiEnvelope<Budget>>("/budgets", {
      category_id: input.category_id || null,
      name: input.name,
      amount: input.amount,
      period: input.period,
      start_date: input.start_date,
      end_date: input.end_date,
      alert_threshold: input.alert_threshold ?? "0.80",
      is_active: input.is_active ?? true,
    })
  );
}

export async function updateBudget(id: string, input: BudgetInput): Promise<Budget> {
  return unwrap(
    await apiClient.patch<ApiEnvelope<Budget>>(`/budgets/${id}`, {
      category_id: input.category_id || null,
      name: input.name,
      amount: input.amount,
      period: input.period,
      start_date: input.start_date,
      end_date: input.end_date,
      alert_threshold: input.alert_threshold ?? "0.80",
      is_active: input.is_active ?? true,
    })
  );
}

export async function deleteBudget(id: string): Promise<void> {
  await apiClient.delete(`/budgets/${id}`);
}

export async function getBudgetStatus(id: string): Promise<BudgetStatus> {
  return unwrap(await apiClient.get<ApiEnvelope<BudgetStatus>>(`/budgets/${id}/status`));
}

// --- Savings Goals APIs ---

export async function listSavingsGoals(): Promise<SavingsGoal[]> {
  return unwrapCollection(await apiClient.get<ApiEnvelope<SavingsGoal[] | null>>("/savings-goals"));
}

export async function getSavingsGoal(id: string): Promise<SavingsGoal> {
  return unwrap(await apiClient.get<ApiEnvelope<SavingsGoal>>(`/savings-goals/${id}`));
}

export async function createSavingsGoal(input: SavingsGoalInput): Promise<SavingsGoal> {
  return unwrap(
    await apiClient.post<ApiEnvelope<SavingsGoal>>("/savings-goals", {
      name: input.name,
      target_amount: input.target_amount,
      current_amount: input.current_amount ?? "0",
      linked_account_id: input.linked_account_id || null,
      target_date: input.target_date || null,
      color: input.color,
      icon: input.icon,
    })
  );
}

export async function updateSavingsGoal(
  id: string,
  input: SavingsGoalInput
): Promise<SavingsGoal> {
  return unwrap(
    await apiClient.patch<ApiEnvelope<SavingsGoal>>(`/savings-goals/${id}`, {
      name: input.name,
      target_amount: input.target_amount,
      current_amount: input.current_amount ?? "0",
      linked_account_id: input.linked_account_id || null,
      target_date: input.target_date || null,
      color: input.color,
      icon: input.icon,
    })
  );
}

export async function deleteSavingsGoal(id: string): Promise<void> {
  await apiClient.delete(`/savings-goals/${id}`);
}

export async function updateGoalProgress(
  id: string,
  amount: string
): Promise<SavingsGoal> {
  return unwrap(
    await apiClient.post<ApiEnvelope<SavingsGoal>>(`/savings-goals/${id}/progress`, {
      amount,
    })
  );
}
