export type AccountType = "CASH" | "BANK" | "CREDIT_CARD" | "SAVINGS" | "INVESTMENT" | "CRYPTO";
export type CategoryType = "INCOME" | "EXPENSE";
export type TransactionType = "INCOME" | "EXPENSE" | "TRANSFER";
export type BudgetPeriod = "WEEKLY" | "MONTHLY" | "QUARTERLY" | "YEARLY" | "CUSTOM";

export interface AccountDto {
  id: string;
  name: string;
  type: AccountType;
  currency: string;
  initial_balance: string;
  current_balance: string;
  color: string;
  icon: string;
  is_archived: boolean;
  created_at: string;
  updated_at: string;
}

export interface CategoryDto {
  id: string;
  name: string;
  type: CategoryType;
  icon: string;
  color: string;
  parent_id: string | null;
  monthly_budget: string | null;
  created_at: string;
  updated_at: string;
}

export interface TransactionDto {
  id: string;
  account_id: string;
  to_account_id: string | null;
  category_id: string | null;
  type: TransactionType;
  amount: string;
  currency: string;
  exchange_rate: string;
  description: string;
  merchant: string | null;
  tags: string[];
  occurred_at: string;
  is_recurring: boolean;
  created_at: string;
  updated_at: string;
}

export interface CashflowDto {
  income: string;
  expense: string;
  net: string;
  savings_rate: string;
}

export interface CategorySpendingDto {
  category_id: string | null;
  category_name: string;
  color: string;
  amount: string;
}

export interface FinanceDashboardDto {
  net_worth: string;
  cashflow: CashflowDto;
  spending: CategorySpendingDto[];
  recent_transactions: TransactionDto[];
}

export interface Budget {
  id: string;
  user_id?: string;
  category_id?: string | null;
  name: string;
  amount: string;
  period: BudgetPeriod;
  start_date: string;
  end_date: string;
  alert_threshold: string;
  is_active: boolean;
  created_at?: string;
  updated_at?: string;
}

export interface BudgetStatus {
  budget: Budget;
  spent_amount: string;
  utilization_rate: string;
  is_alert_triggered: boolean;
}

export interface BudgetInput {
  category_id?: string | null;
  name: string;
  amount: string;
  period: BudgetPeriod;
  start_date: string;
  end_date: string;
  alert_threshold?: string;
  is_active?: boolean;
}

export interface SavingsGoal {
  id: string;
  user_id?: string;
  name: string;
  target_amount: string;
  current_amount: string;
  linked_account_id?: string | null;
  target_date?: string | null;
  color: string;
  icon: string;
  created_at?: string;
  updated_at?: string;
}

export interface SavingsGoalInput {
  name: string;
  target_amount: string;
  current_amount?: string;
  linked_account_id?: string | null;
  target_date?: string | null;
  color: string;
  icon: string;
}

export interface AccountInput {
  name: string;
  type: AccountType;
  currency: string;
  initialBalance: string;
  color: string;
  icon: string;
  isArchived: boolean;
}

export interface CategoryInput {
  name: string;
  type: CategoryType;
  icon: string;
  color: string;
  parentID: string | null;
  monthlyBudget: string | null;
}

export interface TransactionInput {
  accountID: string;
  toAccountID: string | null;
  categoryID: string | null;
  type: TransactionType;
  amount: string;
  currency: string;
  exchangeRate: string;
  description: string;
  merchant: string | null;
  tags: string[];
  occurredAt: string;
  isRecurring: boolean;
}

export interface TransactionFilters {
  type?: TransactionType;
  account?: string;
  category?: string;
  from?: string;
  to?: string;
  q?: string;
}
