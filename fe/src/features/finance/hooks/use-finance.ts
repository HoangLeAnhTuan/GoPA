import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  createAccount,
  createBudget,
  createCategory,
  createSavingsGoal,
  createTransaction,
  deleteBudget,
  deleteSavingsGoal,
  getBudgetStatus,
  getDashboard,
  getNetWorth,
  listAccounts,
  listBudgets,
  listCategories,
  listSavingsGoals,
  listTransactions,
  seedCategories,
  updateBudget,
  updateGoalProgress,
  updateSavingsGoal,
} from "../api/finance-api";
import { financeKeys } from "../api/finance-keys";
import type {
  AccountInput,
  BudgetInput,
  CategoryInput,
  SavingsGoalInput,
  TransactionFilters,
  TransactionInput,
} from "../types";

export function useFinanceDashboard() {
  return useQuery({
    queryKey: financeKeys.dashboard(),
    queryFn: getDashboard,
    staleTime: 30_000,
  });
}

export function useNetWorth() {
  return useQuery({
    queryKey: [...financeKeys.all, "net-worth"],
    queryFn: getNetWorth,
    staleTime: 30_000,
  });
}

export function useAccounts() {
  return useQuery({
    queryKey: financeKeys.accounts(),
    queryFn: listAccounts,
    staleTime: 60_000,
  });
}

export function useCategories() {
  return useQuery({
    queryKey: financeKeys.categories(),
    queryFn: listCategories,
    staleTime: 60_000,
  });
}

export function useTransactions(filters: TransactionFilters = {}) {
  return useQuery({
    queryKey: financeKeys.transactions(filters),
    queryFn: () => listTransactions(filters),
    staleTime: 10_000,
  });
}

export function useCreateAccount() {
  const client = useQueryClient();
  return useMutation({
    mutationFn: (input: AccountInput) => createAccount(input),
    onSuccess: () => client.invalidateQueries({ queryKey: financeKeys.all }),
  });
}

export function useCreateCategory() {
  const client = useQueryClient();
  return useMutation({
    mutationFn: (input: CategoryInput) => createCategory(input),
    onSuccess: () => client.invalidateQueries({ queryKey: financeKeys.all }),
  });
}

export function useSeedCategories() {
  const client = useQueryClient();
  return useMutation({
    mutationFn: seedCategories,
    onSuccess: () => client.invalidateQueries({ queryKey: financeKeys.all }),
  });
}

export function useCreateTransaction() {
  const client = useQueryClient();
  return useMutation({
    mutationFn: (input: TransactionInput) => createTransaction(input),
    onSuccess: () => client.invalidateQueries({ queryKey: financeKeys.all }),
  });
}

// --- Budgets Hooks ---

export function useBudgets() {
  return useQuery({
    queryKey: financeKeys.budgets(),
    queryFn: listBudgets,
    staleTime: 30_000,
  });
}

export function useBudgetStatus(id: string) {
  return useQuery({
    queryKey: financeKeys.budgetStatus(id),
    queryFn: () => getBudgetStatus(id),
    enabled: Boolean(id),
  });
}

export function useCreateBudget() {
  const client = useQueryClient();
  return useMutation({
    mutationFn: (input: BudgetInput) => createBudget(input),
    onSuccess: () => client.invalidateQueries({ queryKey: financeKeys.all }),
  });
}

export function useUpdateBudget() {
  const client = useQueryClient();
  return useMutation({
    mutationFn: ({ id, input }: { id: string; input: BudgetInput }) => updateBudget(id, input),
    onSuccess: () => client.invalidateQueries({ queryKey: financeKeys.all }),
  });
}

export function useDeleteBudget() {
  const client = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => deleteBudget(id),
    onSuccess: () => client.invalidateQueries({ queryKey: financeKeys.all }),
  });
}

// --- Savings Goals Hooks ---

export function useSavingsGoals() {
  return useQuery({
    queryKey: financeKeys.savingsGoals(),
    queryFn: listSavingsGoals,
    staleTime: 30_000,
  });
}

export function useCreateSavingsGoal() {
  const client = useQueryClient();
  return useMutation({
    mutationFn: (input: SavingsGoalInput) => createSavingsGoal(input),
    onSuccess: () => client.invalidateQueries({ queryKey: financeKeys.all }),
  });
}

export function useUpdateSavingsGoal() {
  const client = useQueryClient();
  return useMutation({
    mutationFn: ({ id, input }: { id: string; input: SavingsGoalInput }) =>
      updateSavingsGoal(id, input),
    onSuccess: () => client.invalidateQueries({ queryKey: financeKeys.all }),
  });
}

export function useDeleteSavingsGoal() {
  const client = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => deleteSavingsGoal(id),
    onSuccess: () => client.invalidateQueries({ queryKey: financeKeys.all }),
  });
}

export function useUpdateGoalProgress() {
  const client = useQueryClient();
  return useMutation({
    mutationFn: ({ id, amount }: { id: string; amount: string }) =>
      updateGoalProgress(id, amount),
    onSuccess: () => client.invalidateQueries({ queryKey: financeKeys.all }),
  });
}
