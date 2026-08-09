import { apiClient } from "../../../lib/api-client";
import { unwrap, unwrapCollection, type ApiEnvelope } from "../../../lib/api-envelope";
import type { AccountDto, AccountInput, CategoryDto, CategoryInput, FinanceDashboardDto, TransactionDto, TransactionFilters, TransactionInput } from "../types";

export async function getDashboard(): Promise<FinanceDashboardDto> { return unwrap(await apiClient.get<ApiEnvelope<FinanceDashboardDto>>("/finance/dashboard")); }
export async function listAccounts(): Promise<AccountDto[]> { return unwrapCollection(await apiClient.get<ApiEnvelope<AccountDto[] | null>>("/accounts")); }
export async function createAccount(input: AccountInput): Promise<AccountDto> { return unwrap(await apiClient.post<ApiEnvelope<AccountDto>>("/accounts", { name:input.name, type:input.type, currency:input.currency, initial_balance:input.initialBalance, color:input.color, icon:input.icon, is_archived:input.isArchived })); }
export async function listCategories(): Promise<CategoryDto[]> { return unwrapCollection(await apiClient.get<ApiEnvelope<CategoryDto[] | null>>("/categories")); }
export async function createCategory(input: CategoryInput): Promise<CategoryDto> { return unwrap(await apiClient.post<ApiEnvelope<CategoryDto>>("/categories", { name:input.name, type:input.type, icon:input.icon, color:input.color, parent_id:input.parentID, monthly_budget:input.monthlyBudget })); }
export async function listTransactions(filters: TransactionFilters = {}): Promise<TransactionDto[]> { return unwrapCollection(await apiClient.get<ApiEnvelope<TransactionDto[] | null>>("/transactions", { params:filters })); }
export async function createTransaction(input: TransactionInput): Promise<TransactionDto> { return unwrap(await apiClient.post<ApiEnvelope<TransactionDto>>("/transactions", { account_id:input.accountID, to_account_id:input.toAccountID, category_id:input.categoryID, type:input.type, amount:input.amount, currency:input.currency, exchange_rate:input.exchangeRate, description:input.description, merchant:input.merchant, tags:input.tags, occurred_at:input.occurredAt, is_recurring:input.isRecurring })); }
