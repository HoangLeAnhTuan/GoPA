import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { createAccount, createCategory, createTransaction, getDashboard, listAccounts, listCategories, listTransactions } from "../api/finance-api";
import { financeKeys } from "../api/finance-keys";
import type { AccountInput, CategoryInput, TransactionFilters, TransactionInput } from "../types";

export function useFinanceDashboard() { return useQuery({ queryKey: financeKeys.dashboard(), queryFn: getDashboard, staleTime: 30_000 }); }
export function useAccounts() { return useQuery({ queryKey: financeKeys.accounts(), queryFn: listAccounts, staleTime: 60_000 }); }
export function useCategories() { return useQuery({ queryKey: financeKeys.categories(), queryFn: listCategories, staleTime: 60_000 }); }
export function useTransactions(filters:TransactionFilters = {}) { return useQuery({ queryKey: financeKeys.transactions(filters), queryFn: () => listTransactions(filters), staleTime: 10_000 }); }
export function useCreateAccount() { const client=useQueryClient(); return useMutation({ mutationFn:(input:AccountInput)=>createAccount(input), onSuccess:()=>client.invalidateQueries({queryKey:financeKeys.all}) }); }
export function useCreateCategory() { const client=useQueryClient(); return useMutation({ mutationFn:(input:CategoryInput)=>createCategory(input), onSuccess:()=>client.invalidateQueries({queryKey:financeKeys.all}) }); }
export function useCreateTransaction() { const client=useQueryClient(); return useMutation({ mutationFn:(input:TransactionInput)=>createTransaction(input), onSuccess:()=>client.invalidateQueries({queryKey:financeKeys.all}) }); }
