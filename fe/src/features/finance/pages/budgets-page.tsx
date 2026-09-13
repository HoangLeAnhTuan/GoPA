import {
  AlertCircle,
  AlertTriangle,
  ArrowLeft,
  Edit2,
  PieChart,
  Plus,
  Sparkles,
  Tag,
  Trash2,
  X,
} from "lucide-react";
import React, { useState } from "react";
import { Link } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { AmountDisplay } from "../../../components/design-system/amount-display";
import { ButtonInButton } from "../../../components/design-system/button-in-button";
import { DoubleBezelCard } from "../../../components/design-system/double-bezel-card";
import { ErrorState } from "../../../components/feedback/error-state";
import { LoadingState } from "../../../components/feedback/loading-state";
import {
  useBudgets,
  useCategories,
  useCreateBudget,
  useDeleteBudget,
  useSeedCategories,
  useTransactions,
  useUpdateBudget,
} from "../hooks/use-finance";
import type { Budget, BudgetInput, BudgetPeriod } from "../types";

export function BudgetsPage() {
  const { t } = useTranslation("finance");
  const budgetsQuery = useBudgets();
  const categoriesQuery = useCategories();
  const transactionsQuery = useTransactions();
  const createMutation = useCreateBudget();
  const updateMutation = useUpdateBudget();
  const deleteMutation = useDeleteBudget();
  const seedMutation = useSeedCategories();

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingBudget, setEditingBudget] = useState<Budget | null>(null);

  // Form states
  const [formName, setFormName] = useState("");
  const [formAmount, setFormAmount] = useState("");
  const [formCategoryId, setFormCategoryId] = useState<string>("");
  const [formPeriod, setFormPeriod] = useState<BudgetPeriod>("MONTHLY");
  const [formThreshold, setFormThreshold] = useState<string>("0.80");

  if (budgetsQuery.isPending || categoriesQuery.isPending || transactionsQuery.isPending) {
    return <LoadingState />;
  }
  if (budgetsQuery.isError || categoriesQuery.isError || transactionsQuery.isError) {
    return <ErrorState />;
  }

  const budgets = budgetsQuery.data ?? [];
  const categories = categoriesQuery.data ?? [];
  const transactions = transactionsQuery.data ?? [];

  // Calculate actual spending for each budget based on its category
  const budgetsWithSpent = budgets.map((b) => {
    const limitNum = parseFloat(b.amount) || 0;
    const thresholdNum = parseFloat(b.alert_threshold) || 0.8;

    // Filter transactions in this budget's category and date range
    const spent = transactions
      .filter((t) => {
        if (t.type !== "EXPENSE") return false;
        if (b.category_id && t.category_id !== b.category_id) return false;
        const txDate = new Date(t.occurred_at).getTime();
        const start = new Date(b.start_date).getTime();
        const end = new Date(b.end_date).getTime();
        return txDate >= start && txDate <= end;
      })
      .reduce((sum, t) => sum + (parseFloat(t.amount) || 0), 0);

    const remaining = Math.max(0, limitNum - spent);
    const percentage = limitNum > 0 ? Math.round((spent / limitNum) * 100) : 0;
    const isExceeded = spent > limitNum;
    const isNearThreshold = spent / limitNum >= thresholdNum;

    return {
      ...b,
      spentNum: spent,
      remainingNum: remaining,
      percentage,
      isExceeded,
      isNearThreshold,
    };
  });

  const openCreateModal = () => {
    setEditingBudget(null);
    setFormName("");
    setFormAmount("");
    setFormCategoryId(categories[0]?.id ?? "");
    setFormPeriod("MONTHLY");
    setFormThreshold("0.80");
    setIsModalOpen(true);
  };

  const openEditModal = (budget: Budget) => {
    setEditingBudget(budget);
    setFormName(budget.name);
    setFormAmount(budget.amount);
    setFormCategoryId(budget.category_id ?? "");
    setFormPeriod(budget.period);
    setFormThreshold(budget.alert_threshold ?? "0.80");
    setIsModalOpen(true);
  };

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formName.trim() || !formAmount.trim()) return;

    // Default dates for monthly: 1st of month to last of month
    const now = new Date();
    const start = new Date(Date.UTC(now.getUTCFullYear(), now.getUTCMonth(), 1)).toISOString();
    const end = new Date(Date.UTC(now.getUTCFullYear(), now.getUTCMonth() + 1, 0, 23, 59, 59)).toISOString();

    const payload: BudgetInput = {
      name: formName.trim(),
      amount: formAmount.trim(),
      category_id: formCategoryId || null,
      period: formPeriod,
      start_date: start,
      end_date: end,
      alert_threshold: formThreshold,
      is_active: true,
    };

    if (editingBudget) {
      await updateMutation.mutateAsync({ id: editingBudget.id, input: payload });
    } else {
      await createMutation.mutateAsync(payload);
    }
    setIsModalOpen(false);
  };

  return (
    <div className="flex h-full w-full flex-col overflow-y-auto space-y-4">
      {/* Navigation Sub-header */}
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div className="flex items-center gap-3">
          <Link
            aria-label="Back to Finance"
            className="flex h-9 w-9 items-center justify-center rounded-xl border border-white/10 bg-white/5 text-muted-foreground transition hover:bg-white/10 hover:text-foreground"
            to="/app/finance"
          >
            <ArrowLeft className="h-4 w-4" />
          </Link>
          <div>
            <h1 className="text-xl font-bold tracking-tight">{t("budgets.title")}</h1>
            <p className="text-xs text-muted-foreground">
              {t("budgets.subtitle")}
            </p>
          </div>
        </div>

        <div className="flex items-center gap-2">
          {categories.length === 0 && (
            <button
              className="flex h-9 items-center gap-1.5 rounded-xl border border-emerald-500/30 bg-emerald-500/10 px-3 text-xs font-semibold text-emerald-400 transition hover:bg-emerald-500/20"
              disabled={seedMutation.isPending}
              onClick={() => seedMutation.mutate()}
              type="button"
            >
              <Sparkles className="h-3.5 w-3.5" />
              {seedMutation.isPending ? t("budgets.seeding") : t("budgets.seedDefault")}
            </button>
          )}

          <ButtonInButton icon={Plus} onClick={openCreateModal} variant="primary">
            {t("budgets.newBudget")}
          </ButtonInButton>
        </div>
      </div>

      {/* Tabs Row */}
      <div className="flex items-center gap-2 border-b border-white/10 pb-2 text-xs">
        <Link
          className="rounded-lg px-3 py-1.5 font-medium text-muted-foreground transition hover:text-foreground"
          to="/app/finance"
        >
          {t("tabs.overview")}
        </Link>
        <Link
          className="rounded-lg bg-white/10 px-3 py-1.5 font-semibold text-foreground shadow-sm"
          to="/app/finance/budgets"
        >
          {t("tabs.budgetsWithCount", { count: budgets.length })}
        </Link>
        <Link
          className="rounded-lg px-3 py-1.5 font-medium text-muted-foreground transition hover:text-foreground"
          to="/app/finance/goals"
        >
          {t("tabs.goals")}
        </Link>
      </div>

      {/* Budgets Grid */}
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {budgetsWithSpent.map((budget) => {
          const cat = categories.find((c) => c.id === budget.category_id);
          let meterColor = "bg-emerald-500";
          let badgeBorder = "border-white/10";
          if (budget.isExceeded) {
            meterColor = "bg-rose-500";
            badgeBorder = "border-rose-500/40 bg-rose-500/5";
          } else if (budget.isNearThreshold) {
            meterColor = "bg-amber-500";
            badgeBorder = "border-amber-500/40 bg-amber-500/5";
          }

          return (
            <DoubleBezelCard className={`p-5 transition ${badgeBorder}`} key={budget.id}>
              {/* Card Top */}
              <div className="flex items-start justify-between">
                <div className="flex items-center gap-2.5">
                  <div
                    className="flex h-9 w-9 items-center justify-center rounded-xl font-bold shadow-inner"
                    style={{
                      backgroundColor: cat?.color ? `${cat.color}25` : "rgba(255,255,255,0.1)",
                      color: cat?.color || "#ffffff",
                    }}
                  >
                    <Tag className="h-4 w-4" />
                  </div>
                  <div>
                    <h3 className="font-bold text-sm">{budget.name}</h3>
                    <p className="text-[11px] text-muted-foreground capitalize">
                      {cat?.name || "General"} · {budget.period.toLowerCase()}
                    </p>
                  </div>
                </div>

                <div className="flex items-center gap-1">
                  <button
                    aria-label="Edit budget"
                    className="flex h-7 w-7 items-center justify-center rounded-lg border border-white/10 bg-white/5 text-muted-foreground transition hover:bg-white/10 hover:text-foreground"
                    onClick={() => openEditModal(budget)}
                    type="button"
                  >
                    <Edit2 className="h-3 w-3" />
                  </button>
                  <button
                    aria-label="Delete budget"
                    className="flex h-7 w-7 items-center justify-center rounded-lg border border-white/10 bg-white/5 text-muted-foreground transition hover:bg-rose-500/20 hover:text-rose-400"
                    onClick={() => {
                      if (confirm(t("budgets.deleteConfirm", { name: budget.name }))) {
                        deleteMutation.mutate(budget.id);
                      }
                    }}
                    type="button"
                  >
                    <Trash2 className="h-3 w-3" />
                  </button>
                </div>
              </div>

              {/* Amount Figures */}
              <div className="mt-4 flex items-baseline justify-between">
                <div>
                  <span className="text-[10px] uppercase font-mono tracking-wider text-muted-foreground">
                    {t("budgets.spent")}
                  </span>
                  <div className="font-mono text-base font-bold text-foreground">
                    <AmountDisplay amount={budget.spentNum} />
                  </div>
                </div>

                <div className="text-right">
                  <span className="text-[10px] uppercase font-mono tracking-wider text-muted-foreground">
                    {t("budgets.limit")}
                  </span>
                  <div className="font-mono text-base font-bold text-muted-foreground">
                    <AmountDisplay amount={budget.amount} />
                  </div>
                </div>
              </div>

              {/* Progress Meter */}
              <div className="mt-3">
                <div className="h-2 w-full overflow-hidden rounded-full bg-slate-800">
                  <div
                    className={`h-full rounded-full transition-all duration-500 ${meterColor}`}
                    style={{ width: `${Math.min(100, budget.percentage)}%` }}
                  />
                </div>
                <div className="mt-1.5 flex items-center justify-between text-[11px]">
                  <span className="font-mono font-medium text-muted-foreground">
                    {t("budgets.used", { percent: budget.percentage })}
                  </span>
                  <span className="font-mono font-semibold">
                    {budget.isExceeded ? (
                      <span className="text-rose-400">
                        {t("budgets.overBy", { amount: Math.abs(budget.remainingNum).toLocaleString() })}
                      </span>
                    ) : (
                      <span className="text-muted-foreground">
                        {t("budgets.remaining", { amount: budget.remainingNum.toLocaleString() })}
                      </span>
                    )}
                  </span>
                </div>
              </div>

              {/* Threshold Warning Banner */}
              {budget.isExceeded ? (
                <div className="mt-3 flex items-center gap-1.5 rounded-lg border border-rose-500/30 bg-rose-500/10 px-2.5 py-1.5 text-[11px] font-semibold text-rose-400">
                  <AlertCircle className="h-3.5 w-3.5 shrink-0" />
                  <span>{t("budgets.limitExceeded")}</span>
                </div>
              ) : budget.isNearThreshold ? (
                <div className="mt-3 flex items-center gap-1.5 rounded-lg border border-amber-500/30 bg-amber-500/10 px-2.5 py-1.5 text-[11px] font-semibold text-amber-400">
                  <AlertTriangle className="h-3.5 w-3.5 shrink-0" />
                  <span>
                    {t("budgets.thresholdAlert", {
                      threshold: Math.round(parseFloat(budget.alert_threshold) * 100),
                    })}
                  </span>
                </div>
              ) : null}
            </DoubleBezelCard>
          );
        })}

        {budgets.length === 0 && (
          <div className="col-span-full py-16 text-center">
            <div className="mx-auto flex h-14 w-14 items-center justify-center rounded-2xl bg-white/5 text-muted-foreground">
              <PieChart className="h-7 w-7" />
            </div>
            <h2 className="mt-3 text-base font-bold">{t("budgets.emptyTitle")}</h2>
            <p className="mt-1 text-xs text-muted-foreground">
              {t("budgets.emptySubtitle")}
            </p>
            <div className="mt-4 flex justify-center">
              <ButtonInButton icon={Plus} onClick={openCreateModal} variant="primary">
                {t("budgets.createFirst")}
              </ButtonInButton>
            </div>
          </div>
        )}
      </div>

      {/* CREATE / EDIT MODAL */}
      {isModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4 backdrop-blur-sm">
          <div className="w-full max-w-md">
            <DoubleBezelCard className="p-6">
              <div className="flex items-center justify-between border-b border-white/10 pb-3">
                <h2 className="text-base font-bold">
                  {editingBudget ? t("budgets.editModalTitle") : t("budgets.createModalTitle")}
                </h2>
                <button
                  aria-label="Close modal"
                  className="rounded-lg p-1 text-muted-foreground hover:bg-white/10 hover:text-foreground"
                  onClick={() => setIsModalOpen(false)}
                  type="button"
                >
                  <X className="h-4 w-4" />
                </button>
              </div>

              <form className="mt-4 space-y-3" onSubmit={handleSave}>
                <div>
                  <label className="text-[11px] font-semibold text-muted-foreground">{t("budgets.nameLabel")}</label>
                  <input
                    className="mt-1 h-9 w-full rounded-xl border border-white/10 bg-black/30 px-3 text-xs outline-none focus:border-violet-500"
                    onChange={(e) => setFormName(e.target.value)}
                    placeholder={t("budgets.namePlaceholder")}
                    required
                    value={formName}
                  />
                </div>

                <div>
                  <label className="text-[11px] font-semibold text-muted-foreground">{t("budgets.amountLabel")}</label>
                  <input
                    className="mt-1 h-9 w-full rounded-xl border border-white/10 bg-black/30 px-3 text-xs font-mono outline-none focus:border-violet-500"
                    onChange={(e) => setFormAmount(e.target.value)}
                    placeholder={t("budgets.amountPlaceholder")}
                    required
                    type="number"
                    value={formAmount}
                  />
                </div>

                <div className="grid grid-cols-2 gap-2">
                  <div>
                    <label className="text-[11px] font-semibold text-muted-foreground">{t("budgets.categoryLabel")}</label>
                    <select
                      className="mt-1 h-9 w-full rounded-xl border border-white/10 bg-black/30 px-3 text-xs outline-none focus:border-violet-500"
                      onChange={(e) => setFormCategoryId(e.target.value)}
                      value={formCategoryId}
                    >
                      <option value="">{t("budgets.categoryAll")}</option>
                      {categories.map((c) => (
                        <option key={c.id} value={c.id}>
                          {c.name}
                        </option>
                      ))}
                    </select>
                  </div>

                  <div>
                    <label className="text-[11px] font-semibold text-muted-foreground">{t("budgets.periodLabel")}</label>
                    <select
                      className="mt-1 h-9 w-full rounded-xl border border-white/10 bg-black/30 px-3 text-xs outline-none focus:border-violet-500"
                      onChange={(e) => setFormPeriod(e.target.value as BudgetPeriod)}
                      value={formPeriod}
                    >
                      <option value="MONTHLY">{t("budgets.periodMonthly")}</option>
                      <option value="WEEKLY">{t("budgets.periodWeekly")}</option>
                      <option value="YEARLY">{t("budgets.periodYearly")}</option>
                    </select>
                  </div>
                </div>

                <div>
                  <div className="flex items-center justify-between text-[11px]">
                    <label className="font-semibold text-muted-foreground">{t("budgets.alertThresholdLabel")}</label>
                    <span className="font-mono font-bold text-violet-400">
                      {Math.round(parseFloat(formThreshold) * 100)}%
                    </span>
                  </div>
                  <input
                    className="mt-2 w-full accent-violet-500"
                    max="1.00"
                    min="0.50"
                    onChange={(e) => setFormThreshold(e.target.value)}
                    step="0.05"
                    type="range"
                    value={formThreshold}
                  />
                  <p className="mt-1 text-[10px] text-muted-foreground">
                    {t("budgets.thresholdHint")}
                  </p>
                </div>

                <div className="flex justify-end gap-2 pt-2">
                  <button
                    className="rounded-xl border border-white/10 px-4 py-2 text-xs font-semibold hover:bg-white/5"
                    onClick={() => setIsModalOpen(false)}
                    type="button"
                  >
                    {t("budgets.cancel")}
                  </button>
                  <button
                    className="rounded-xl bg-violet-600 px-4 py-2 text-xs font-semibold text-white shadow-lg shadow-violet-600/20 hover:bg-violet-500"
                    type="submit"
                  >
                    {editingBudget ? t("budgets.update") : t("budgets.save")}
                  </button>
                </div>
              </form>
            </DoubleBezelCard>
          </div>
        </div>
      )}
    </div>
  );
}
