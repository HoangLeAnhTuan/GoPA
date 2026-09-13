import type { TransactionFilters } from "../types";

export const financeKeys = {
  all: ["finance"] as const,
  dashboard: () => [...financeKeys.all, "dashboard"] as const,
  accounts: () => [...financeKeys.all, "accounts"] as const,
  categories: () => [...financeKeys.all, "categories"] as const,
  transactions: (filters: TransactionFilters = {}) =>
    [...financeKeys.all, "transactions", filters] as const,
  budgets: () => [...financeKeys.all, "budgets"] as const,
  budgetStatus: (id: string) => [...financeKeys.all, "budgets", id, "status"] as const,
  savingsGoals: () => [...financeKeys.all, "savings-goals"] as const,
};
