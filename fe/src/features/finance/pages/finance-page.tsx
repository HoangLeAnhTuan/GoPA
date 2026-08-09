import { zodResolver } from "@hookform/resolvers/zod";
import { ArrowDownRight, ArrowUpRight, Check, Landmark, Plus, ReceiptText, SlidersHorizontal, Sparkles, Tags, WalletCards } from "lucide-react";
import React, { useState } from "react";
import { useForm } from "react-hook-form";
import { useTranslation } from "react-i18next";
import { z } from "zod";
import { AmountDisplay } from "../../../components/design-system/amount-display";
import { GlassPanel } from "../../../components/design-system/glass-panel";
import { PageHeader } from "../../../components/design-system/page-header";
import { MacSelect } from "../../../components/design-system/mac-select";
import { ErrorState } from "../../../components/feedback/error-state";
import { LoadingState } from "../../../components/feedback/loading-state";
import { cn } from "../../../lib/cn";
import { useAccounts, useCategories, useCreateAccount, useCreateCategory, useCreateTransaction, useFinanceDashboard } from "../hooks/use-finance";
import type { AccountInput, AccountType, CategoryInput, CategoryType, TransactionInput, TransactionType } from "../types";

const accountSchema = z.object({ name:z.string().trim().min(1).max(120), type:z.enum(["CASH","BANK","CREDIT_CARD","SAVINGS","INVESTMENT","CRYPTO"]), currency:z.string().trim().regex(/^[A-Za-z]{3}$/), initialBalance:z.string().regex(/^\d+(\.\d{1,2})?$/), color:z.string().regex(/^#[0-9A-Fa-f]{6}$/), icon:z.string().trim().min(1).max(50) });
const categorySchema = z.object({ name:z.string().trim().min(1).max(100), type:z.enum(["INCOME","EXPENSE"]), color:z.string().regex(/^#[0-9A-Fa-f]{6}$/) });
const transactionSchema = z.object({ accountID:z.string().uuid(), type:z.enum(["INCOME","EXPENSE","TRANSFER"]), amount:z.string().regex(/^\d+(\.\d{1,2})?$/), currency:z.string().regex(/^[A-Za-z]{3}$/), description:z.string().max(500), categoryID:z.string(), toAccountID:z.string() });

type AccountValues = z.infer<typeof accountSchema>;
type CategoryValues = z.infer<typeof categorySchema>;
type TransactionValues = z.infer<typeof transactionSchema>;

// Preset 6 Colors for Categories & Accounts
const PRESET_COLORS = [
  { hex: "#10B981", key: "emerald" },
  { hex: "#F59E0B", key: "amber" },
  { hex: "#3B82F6", key: "blue" },
  { hex: "#F43F5E", key: "rose" },
  { hex: "#8B5CF6", key: "purple" },
  { hex: "#06B6D4", key: "cyan" },
] as const;

// Preset 6 Popular Categories
const POPULAR_CATEGORIES: Array<{ key: string; type: CategoryType; color: string; icon: string }> = [
  { key: "food", type: "EXPENSE", color: "#F59E0B", icon: "utensils" },
  { key: "transport", type: "EXPENSE", color: "#3B82F6", icon: "car" },
  { key: "shopping", type: "EXPENSE", color: "#F43F5E", icon: "shopping-bag" },
  { key: "income", type: "INCOME", color: "#10B981", icon: "wallet" },
  { key: "entertainment", type: "EXPENSE", color: "#8B5CF6", icon: "film" },
  { key: "savings", type: "INCOME", color: "#06B6D4", icon: "piggy-bank" },
];

type ManageTab = "transaction" | "account" | "category";

const POPULAR_CURRENCIES = ["VND", "USD", "EUR", "JPY", "GBP", "KRW", "SGD"] as const;

export function FinancePage() {
  const { t } = useTranslation("finance");
  const [tab, setTab] = useState<ManageTab>("transaction");
  const dashboard = useFinanceDashboard(); const accounts = useAccounts(); const categories = useCategories();
  if (dashboard.isPending || accounts.isPending || categories.isPending) return <LoadingState />;
  if (dashboard.isError || accounts.isError || categories.isError) return <ErrorState onRetry={() => { void dashboard.refetch(); void accounts.refetch(); void categories.refetch(); }} />;

  const tabItems: Array<{ key: ManageTab; label: string; icon: React.ReactNode }> = [
    { key: "transaction", label: t("quickAdd.title"), icon: <ReceiptText size={13} /> },
    { key: "account",     label: t("accounts.title"), icon: <Landmark size={13} /> },
    { key: "category",   label: t("categories.title"), icon: <Sparkles size={13} /> },
  ];

  return (
    <div className="flex flex-col gap-4">
      <PageHeader title={t("title")} description={t("description")} />

      {/* 4 overview metric cards */}
      <FinanceOverview />

      {/* Main area: Quick-add tabs (left) + Recent transactions (right) */}
      <div className="grid gap-3 lg:grid-cols-[1.1fr_0.9fr]" style={{ minHeight: 0 }}>
        {/* Left: tab switcher + form */}
        <GlassPanel className="flex flex-col overflow-hidden p-4">
          {/* Tab header */}
          <div className="mb-3 flex gap-1 rounded-xl border border-white/30 bg-white/30 p-1 dark:border-white/10 dark:bg-slate-900/40">
            {tabItems.map(({ key, label, icon }) => (
              <button
                className={cn(
                  "flex flex-1 items-center justify-center gap-1.5 rounded-lg py-1.5 text-xs font-semibold transition-all",
                  tab === key
                    ? "bg-primary text-primary-foreground shadow-sm"
                    : "text-muted-foreground hover:bg-white/50 dark:hover:bg-slate-800/60"
                )}
                key={key}
                onClick={() => setTab(key)}
                type="button"
              >
                {icon}
                <span className="hidden sm:inline">{label}</span>
              </button>
            ))}
          </div>
          {/* Tab content */}
          <div className="overflow-y-auto">
            {tab === "transaction" && <TransactionQuickAdd />}
            {tab === "account" && <AccountQuickAdd />}
            {tab === "category" && <CategoryQuickAdd />}
          </div>
        </GlassPanel>

        {/* Right: Recent transactions */}
        <GlassPanel className="flex flex-col overflow-hidden p-4">
          <h2 className="mb-2 shrink-0 text-sm font-semibold">{t("recent.title")}</h2>
          <div className="overflow-y-auto">
            <RecentTransactions />
          </div>
        </GlassPanel>
      </div>
    </div>
  );
}

function FinanceOverview() {
  const { t } = useTranslation("finance"); const dashboard = useFinanceDashboard();
  if (dashboard.data === undefined) return null;
  const currency = dashboard.data.recent_transactions[0]?.currency ?? "VND";
  const net = Number(dashboard.data.cashflow.net);
  const savingsRate = dashboard.data.cashflow.savings_rate;
  return (
    <div className="grid grid-cols-2 gap-3 xl:grid-cols-4">
      {/* Net Worth */}
      <GlassPanel className="relative overflow-hidden p-4">
        <div className="absolute inset-0 bg-gradient-to-br from-emerald-500/15 via-teal-400/8 to-transparent" />
        <div className="relative">
          <div className="flex items-center gap-2 text-emerald-600 dark:text-emerald-400">
            <div className="flex size-7 items-center justify-center rounded-lg bg-emerald-500/15">
              <WalletCards size={14} />
            </div>
            <span className="text-[10px] font-bold uppercase tracking-widest">{t("netWorth")}</span>
          </div>
          <AmountDisplay amount={dashboard.data.net_worth} className="mt-2 block text-xl font-bold tracking-tight" currency={currency} />
          <p className="mt-0.5 text-[10px] text-muted-foreground">{t("netWorthDescription")}</p>
        </div>
      </GlassPanel>

      {/* Income */}
      <GlassPanel className="relative overflow-hidden p-4">
        <div className="absolute inset-0 bg-gradient-to-br from-blue-500/12 via-sky-400/6 to-transparent" />
        <div className="relative">
          <div className="flex items-center gap-2 text-blue-500">
            <div className="flex size-7 items-center justify-center rounded-lg bg-blue-500/15">
              <ArrowUpRight size={14} />
            </div>
            <span className="text-[10px] font-bold uppercase tracking-widest">{t("income")}</span>
          </div>
          <AmountDisplay amount={dashboard.data.cashflow.income} className="mt-2 block text-xl font-bold tracking-tight" currency={currency} signed="positive" />
          <p className="mt-0.5 text-[10px] text-muted-foreground">{t("thisMonth")}</p>
        </div>
      </GlassPanel>

      {/* Expenses */}
      <GlassPanel className="relative overflow-hidden p-4">
        <div className="absolute inset-0 bg-gradient-to-br from-rose-500/12 via-pink-400/6 to-transparent" />
        <div className="relative">
          <div className="flex items-center gap-2 text-rose-500">
            <div className="flex size-7 items-center justify-center rounded-lg bg-rose-500/15">
              <ArrowDownRight size={14} />
            </div>
            <span className="text-[10px] font-bold uppercase tracking-widest">{t("expenses")}</span>
          </div>
          <AmountDisplay amount={dashboard.data.cashflow.expense} className="mt-2 block text-xl font-bold tracking-tight" currency={currency} signed="negative" />
          <p className="mt-0.5 text-[10px] text-muted-foreground">{t("thisMonth")}</p>
        </div>
      </GlassPanel>

      {/* Monthly Cashflow */}
      <GlassPanel className="relative overflow-hidden p-4">
        <div className={cn("absolute inset-0 bg-gradient-to-br to-transparent", net >= 0 ? "from-emerald-500/12 via-green-400/6" : "from-rose-500/12 via-red-400/6")} />
        <div className="relative">
          <div className={cn("flex items-center gap-2", net >= 0 ? "text-emerald-500" : "text-rose-500")}>
            <div className={cn("flex size-7 items-center justify-center rounded-lg", net >= 0 ? "bg-emerald-500/15" : "bg-rose-500/15")}>
              <Tags size={14} />
            </div>
            <span className="text-[10px] font-bold uppercase tracking-widest">{t("cashflow")}</span>
          </div>
          <AmountDisplay amount={dashboard.data.cashflow.net} className="mt-2 block text-xl font-bold tracking-tight" currency={currency} signed={net >= 0 ? "positive" : "negative"} />
          <p className="mt-0.5 text-[10px] text-muted-foreground">{t("savingsRate", { value: savingsRate })}</p>
        </div>
      </GlassPanel>
    </div>
  );
}

function TransactionQuickAdd() {
  const { t } = useTranslation("finance"); const accounts = useAccounts(); const categories = useCategories(); const create = useCreateTransaction();
  const form = useForm<TransactionValues>({ resolver:zodResolver(transactionSchema), defaultValues:{ accountID:"", type:"EXPENSE", amount:"", currency:"VND", description:"", categoryID:"", toAccountID:"" } });
  const type = form.watch("type"); const accountID = form.watch("accountID");
  const eligibleCategories = (categories.data ?? []).filter((category) => category.type === type);

  const submit = async (values:TransactionValues) => {
    const input:TransactionInput={ accountID:values.accountID, toAccountID:values.type === "TRANSFER" && values.toAccountID !== "" ? values.toAccountID : null, categoryID:values.type === "TRANSFER" || values.categoryID === "" ? null : values.categoryID, type:values.type as TransactionType, amount:values.amount, currency:values.currency.toUpperCase(), exchangeRate:"1", description:values.description, merchant:null, tags:[], occurredAt:new Date().toISOString(), isRecurring:false };
    await create.mutateAsync(input);
    form.reset({ ...form.getValues(), amount:"", description:"", categoryID:"", toAccountID:"" });
  };

  const accountOptions = [
    { value: "", label: t("quickAdd.selectAccount") },
    ...(accounts.data?.map((acc) => ({ value: acc.id, label: `${acc.name} (${acc.currency})` })) ?? []),
  ];

  const toAccountOptions = [
    { value: "", label: t("quickAdd.selectAccount") },
    ...(accounts.data?.filter((acc) => acc.id !== accountID).map((acc) => ({ value: acc.id, label: `${acc.name} (${acc.currency})` })) ?? []),
  ];

  const categoryOptions = [
    { value: "", label: t("quickAdd.uncategorized") },
    ...eligibleCategories.map((cat) => ({ value: cat.id, label: cat.name })),
  ];
  const currencyOptions = POPULAR_CURRENCIES.map((currency) => ({ value: currency, label: t(`currencies.${currency}`) }));

  return (
    <div>
      {accounts.data?.length === 0 ? (
        <p className="rounded-xl border border-dashed p-3 text-xs text-muted-foreground">{t("quickAdd.accountRequired")}</p>
      ) : (
        <form className="grid gap-2.5 sm:grid-cols-2" onSubmit={form.handleSubmit(submit)}>
          <fieldset className="sm:col-span-2">
            <div className="grid grid-cols-3 gap-1 rounded-xl border border-white/40 bg-white/40 p-1 dark:border-white/10 dark:bg-slate-900/40">
              {(["INCOME", "EXPENSE", "TRANSFER"] as TransactionType[]).map((value) => (
                <label className="cursor-pointer" key={value}>
                  <input className="peer sr-only" type="radio" value={value} {...form.register("type")} />
                  <span className="flex h-8 items-center justify-center rounded-lg text-xs font-semibold transition peer-checked:bg-emerald-500 peer-checked:text-white peer-checked:shadow-sm">
                    {t(`types.${value}`)}
                  </span>
                </label>
              ))}
            </div>
          </fieldset>

          <Field label={t("quickAdd.amount")} error={form.formState.errors.amount?.message}>
            <input inputMode="decimal" {...form.register("amount")} className="field h-9 text-sm" placeholder="0.00" />
          </Field>

          <Field label={t("quickAdd.currency")} error={form.formState.errors.currency?.message}>
            <MacSelect
              onChange={(val) => form.setValue("currency", val, { shouldValidate: true })}
              options={currencyOptions}
              value={form.watch("currency")}
            />
          </Field>

          <Field label={t("quickAdd.account")} error={form.formState.errors.accountID?.message}>
            <MacSelect onChange={(val) => form.setValue("accountID", val, { shouldValidate: true })} options={accountOptions} value={form.watch("accountID")} />
          </Field>

          {type === "TRANSFER" ? (
            <Field label={t("quickAdd.toAccount")} error={form.formState.errors.toAccountID?.message}>
              <MacSelect onChange={(val) => form.setValue("toAccountID", val, { shouldValidate: true })} options={toAccountOptions} value={form.watch("toAccountID")} />
            </Field>
          ) : (
            <Field label={t("quickAdd.category")} error={form.formState.errors.categoryID?.message}>
              <MacSelect onChange={(val) => form.setValue("categoryID", val, { shouldValidate: true })} options={categoryOptions} value={form.watch("categoryID")} />
            </Field>
          )}

          <Field label={t("quickAdd.description")} error={form.formState.errors.description?.message} className="sm:col-span-2">
            <input {...form.register("description")} className="field h-9 text-sm" placeholder={t("quickAdd.descriptionPlaceholder")} />
          </Field>

          <button className="inline-flex h-9 items-center justify-center gap-2 rounded-xl bg-emerald-500 px-4 text-sm font-semibold text-white shadow-lg shadow-emerald-500/20 transition hover:bg-emerald-600 disabled:opacity-60 sm:col-span-2" disabled={create.isPending} type="submit">
            <Plus size={15} />
            {create.isPending ? t("quickAdd.saving") : t("quickAdd.save")}
          </button>
        </form>
      )}
    </div>
  );
}

function AccountQuickAdd() {
  const { t } = useTranslation("finance"); const create = useCreateAccount();
  const form = useForm<AccountValues>({ resolver:zodResolver(accountSchema), defaultValues:{ name:"", type:"CASH", currency:"VND", initialBalance:"0.00", color:"#10B981", icon:"wallet" } });
  const submit = async (values:AccountValues) => { const input:AccountInput={ ...values, type:values.type as AccountType, currency:values.currency.toUpperCase(), isArchived:false }; await create.mutateAsync(input); form.reset(); };

  const accountTypeOptions = (["CASH","BANK","CREDIT_CARD","SAVINGS","INVESTMENT","CRYPTO"] as AccountType[]).map((type) => ({
    value: type,
    label: t(`accountTypes.${type}`),
  }));
  const currencyOptions = POPULAR_CURRENCIES.map((currency) => ({ value: currency, label: t(`currencies.${currency}`) }));

  return (
    <div>
      <form className="grid gap-2.5 sm:grid-cols-2" onSubmit={form.handleSubmit(submit)}>
        <Field label={t("accounts.name")} error={form.formState.errors.name?.message}>
          <input {...form.register("name")} className="field h-9 text-sm" placeholder={t("accounts.namePlaceholder")} />
        </Field>
        <Field label={t("accounts.type")} error={form.formState.errors.type?.message}>
          <MacSelect onChange={(val) => form.setValue("type", val as AccountType, { shouldValidate: true })} options={accountTypeOptions} value={form.watch("type")} />
        </Field>
        <Field label={t("accounts.balance")} error={form.formState.errors.initialBalance?.message}>
          <input inputMode="decimal" {...form.register("initialBalance")} className="field h-9 text-sm" placeholder="0.00" />
        </Field>
        <Field label={t("quickAdd.currency")} error={form.formState.errors.currency?.message}>
          <MacSelect
            onChange={(val) => form.setValue("currency", val, { shouldValidate: true })}
            options={currencyOptions}
            value={form.watch("currency")}
          />
        </Field>
        <div className="flex justify-end sm:col-span-2">
          <button
            className="inline-flex items-center gap-2 rounded-2xl bg-gradient-to-r from-emerald-500 to-teal-500 px-5 py-2 text-sm font-semibold text-white shadow-lg shadow-emerald-500/25 transition-all hover:from-emerald-400 hover:to-teal-400 active:scale-[0.97] disabled:opacity-60"
            disabled={create.isPending}
            type="submit"
          >
            <Plus size={15} />
            {create.isPending ? t("quickAdd.saving") : t("accounts.add")}
          </button>
        </div>
      </form>
    </div>
  );
}

function CategoryQuickAdd() {
  const { t } = useTranslation("finance");
  const create = useCreateCategory();
  const [showCustomForm, setShowCustomForm] = useState(false);
  const form = useForm<CategoryValues>({
    resolver: zodResolver(categorySchema),
    defaultValues: { name: "", type: "EXPENSE", color: "#F59E0B" },
  });

  const selectedColor = form.watch("color");

  const submit = async (values: CategoryValues) => {
    const input: CategoryInput = { name: values.name, type: values.type as CategoryType, icon: "tag", color: values.color, parentID: null, monthlyBudget: null };
    await create.mutateAsync(input);
    form.reset();
  };

  const handleQuickAddPopular = async (cat: typeof POPULAR_CATEGORIES[number]) => {
    const input: CategoryInput = { name: t(`popularCategories.${cat.key}`), type: cat.type, icon: cat.icon, color: cat.color, parentID: null, monthlyBudget: null };
    await create.mutateAsync(input);
  };

  return (
    <div>
      <div className="flex items-center justify-between">
        <p className="text-xs text-muted-foreground">{t("categories.popularHint")}</p>
        <button
          className="flex items-center gap-1.5 rounded-xl px-2.5 py-1 text-xs font-medium text-amber-500 transition hover:bg-amber-500/10"
          onClick={() => setShowCustomForm(!showCustomForm)}
          type="button"
        >
          <SlidersHorizontal size={13} />
          {showCustomForm ? t("categories.customToggleClose") : t("categories.customToggleOpen")}
        </button>
      </div>

      <div className="mt-2 grid grid-cols-2 gap-1.5 sm:grid-cols-3">
        {POPULAR_CATEGORIES.map((cat) => (
          <button
            className="flex items-center gap-2 rounded-xl border border-white/40 bg-white/50 p-2 text-left text-xs font-medium transition hover:border-amber-500/50 hover:bg-amber-500/10 dark:border-white/10 dark:bg-slate-900/50"
            disabled={create.isPending}
            key={cat.key}
            onClick={() => void handleQuickAddPopular(cat)}
            type="button"
          >
            <span className="size-2 shrink-0 rounded-full" style={{ backgroundColor: cat.color }} />
            <span className="truncate">{t(`popularCategories.${cat.key}`)}</span>
          </button>
        ))}
      </div>

      {showCustomForm && (
        <form className="mt-3 border-t border-white/20 pt-3 dark:border-white/10 grid gap-2.5" onSubmit={form.handleSubmit(submit)}>
          <div className="grid gap-2 sm:grid-cols-2">
            <Field label={t("categories.name")} error={form.formState.errors.name?.message}>
              <input {...form.register("name")} className="field h-9 text-sm" placeholder={t("categories.namePlaceholder")} />
            </Field>
            <Field label={t("categories.type")} error={form.formState.errors.type?.message}>
              <MacSelect
                onChange={(val) => form.setValue("type", val as CategoryType, { shouldValidate: true })}
                options={[
                  { value: "EXPENSE", label: t("types.EXPENSE") },
                  { value: "INCOME", label: t("types.INCOME") },
                ]}
                value={form.watch("type")}
              />
            </Field>
          </div>
          <div>
            <label className="mb-1.5 block text-xs font-medium text-muted-foreground">{t("categories.color")}:</label>
            <div className="flex items-center gap-2">
              {PRESET_COLORS.map((c) => (
                <button
                  aria-label={t(`colors.${c.key}`)}
                  className="relative flex size-6 items-center justify-center rounded-full transition-transform hover:scale-110"
                  key={c.hex}
                  onClick={() => form.setValue("color", c.hex)}
                  style={{ backgroundColor: c.hex }}
                  type="button"
                >
                  {selectedColor === c.hex && <Check className="size-3.5 text-white drop-shadow" />}
                </button>
              ))}
            </div>
          </div>
          <div className="flex justify-end">
            <button
              className="inline-flex items-center gap-2 rounded-2xl bg-gradient-to-r from-amber-500 to-orange-500 px-5 py-2 text-sm font-semibold text-white shadow-lg shadow-amber-500/25 transition-all active:scale-[0.97] disabled:opacity-60"
              disabled={create.isPending}
              type="submit"
            >
              <Plus size={15} />
              {create.isPending ? t("quickAdd.saving") : t("categories.add")}
            </button>
          </div>
        </form>
      )}
    </div>
  );
}

function RecentTransactions() {
  const { t, i18n } = useTranslation("finance"); const dashboard = useFinanceDashboard();
  if (dashboard.data === undefined) return null;
  if (dashboard.data.recent_transactions.length === 0) {
    return <p className="text-xs text-muted-foreground">{t("recent.empty")}</p>;
  }
  return (
    <div className="divide-y divide-white/10">
      {dashboard.data.recent_transactions.map((transaction) => (
        <div className="flex items-center justify-between gap-3 py-2 text-xs" key={transaction.id}>
          <div className="min-w-0">
            <p className="truncate font-medium">{transaction.description || t(`types.${transaction.type}`)}</p>
            <p className="text-[10px] text-muted-foreground">
              {new Intl.DateTimeFormat(i18n.resolvedLanguage === "vi" ? "vi-VN" : i18n.resolvedLanguage === "ja" ? "ja-JP" : "en-US", { dateStyle: "medium" }).format(new Date(transaction.occurred_at))}
            </p>
          </div>
          <AmountDisplay amount={transaction.amount} currency={transaction.currency} signed={transaction.type === "INCOME" ? "positive" : transaction.type === "EXPENSE" ? "negative" : "neutral"} />
        </div>
      ))}
    </div>
  );
}

function Field({ label, error, children, className }:{label:string;error?:string;children:React.ReactNode;className?:string}) {
  return (
    <label className={cn("block text-xs font-medium", className)}>
      <span className="mb-1 block text-muted-foreground">{label}</span>
      {children}
      {error === undefined ? null : <span className="mt-1 block text-xs text-rose-500">{error}</span>}
    </label>
  );
}
