import type { TransactionFilters } from "../types";

export const financeKeys = {
  all: ["finance"] as const,
  dashboard: () => [...financeKeys.all, "dashboard"] as const,
  accounts: () => [...financeKeys.all, "accounts"] as const,
  categories: () => [...financeKeys.all, "categories"] as const,
  transactions: (filters: TransactionFilters = {}) => [...financeKeys.all, "transactions", filters] as const,
};
