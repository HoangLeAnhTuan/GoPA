import {
  ArrowLeft,
  Calendar,
  Coins,
  Edit2,
  PiggyBank,
  Plus,
  Target,
  Trash2,
  Wallet,
  X,
} from "lucide-react";
import React, { useState } from "react";
import { Link } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { AmountDisplay } from "../../../components/design-system/amount-display";
import { ButtonInButton } from "../../../components/design-system/button-in-button";
import { DoubleBezelCard } from "../../../components/design-system/double-bezel-card";
import { ProgressRing } from "../../../components/design-system/progress-ring";
import { ErrorState } from "../../../components/feedback/error-state";
import { LoadingState } from "../../../components/feedback/loading-state";
import {
  useAccounts,
  useCreateSavingsGoal,
  useDeleteSavingsGoal,
  useSavingsGoals,
  useUpdateGoalProgress,
  useUpdateSavingsGoal,
} from "../hooks/use-finance";
import type { SavingsGoal, SavingsGoalInput } from "../types";

export function GoalsPage() {
  const { t } = useTranslation("finance");
  const goalsQuery = useSavingsGoals();
  const accountsQuery = useAccounts();
  const createMutation = useCreateSavingsGoal();
  const updateMutation = useUpdateSavingsGoal();
  const deleteMutation = useDeleteSavingsGoal();
  const progressMutation = useUpdateGoalProgress();

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingGoal, setEditingGoal] = useState<SavingsGoal | null>(null);

  // Deposit progress dialog
  const [depositGoal, setDepositGoal] = useState<SavingsGoal | null>(null);
  const [depositAmount, setDepositAmount] = useState("");

  // Form states
  const [formName, setFormName] = useState("");
  const [formTarget, setFormTarget] = useState("");
  const [formCurrent, setFormCurrent] = useState("0");
  const [formAccountId, setFormAccountId] = useState<string>("");
  const [formDate, setFormDate] = useState<string>("");
  const [formColor, setFormColor] = useState("#8b5cf6");
  const [formIcon, setFormIcon] = useState("target");

  if (goalsQuery.isPending || accountsQuery.isPending) return <LoadingState />;
  if (goalsQuery.isError || accountsQuery.isError) return <ErrorState />;

  const goals = goalsQuery.data ?? [];
  const accounts = accountsQuery.data ?? [];

  const openCreateModal = () => {
    setEditingGoal(null);
    setFormName("");
    setFormTarget("");
    setFormCurrent("0");
    setFormAccountId(accounts[0]?.id ?? "");
    setFormDate("");
    setFormColor("#8b5cf6");
    setFormIcon("target");
    setIsModalOpen(true);
  };

  const openEditModal = (goal: SavingsGoal) => {
    setEditingGoal(goal);
    setFormName(goal.name);
    setFormTarget(goal.target_amount);
    setFormCurrent(goal.current_amount);
    setFormAccountId(goal.linked_account_id ?? "");
    setFormDate(goal.target_date ? goal.target_date.split("T")[0] ?? "" : "");
    setFormColor(goal.color || "#8b5cf6");
    setFormIcon(goal.icon || "target");
    setIsModalOpen(true);
  };

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formName.trim() || !formTarget.trim()) return;

    const payload: SavingsGoalInput = {
      name: formName.trim(),
      target_amount: formTarget.trim(),
      current_amount: formCurrent.trim() || "0",
      linked_account_id: formAccountId || null,
      target_date: formDate ? new Date(formDate).toISOString() : null,
      color: formColor,
      icon: formIcon,
    };

    if (editingGoal) {
      await updateMutation.mutateAsync({ id: editingGoal.id, input: payload });
    } else {
      await createMutation.mutateAsync(payload);
    }
    setIsModalOpen(false);
  };

  const handleDepositSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!depositGoal || !depositAmount.trim()) return;

    await progressMutation.mutateAsync({
      id: depositGoal.id,
      amount: depositAmount.trim(),
    });
    setDepositGoal(null);
    setDepositAmount("");
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
            <h1 className="text-xl font-bold tracking-tight">{t("goals.title")}</h1>
            <p className="text-xs text-muted-foreground">
              {t("goals.subtitle")}
            </p>
          </div>
        </div>

        <ButtonInButton icon={Plus} onClick={openCreateModal} variant="primary">
          {t("goals.newGoal")}
        </ButtonInButton>
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
          className="rounded-lg px-3 py-1.5 font-medium text-muted-foreground transition hover:text-foreground"
          to="/app/finance/budgets"
        >
          {t("tabs.budgets")}
        </Link>
        <Link
          className="rounded-lg bg-white/10 px-3 py-1.5 font-semibold text-foreground shadow-sm"
          to="/app/finance/goals"
        >
          {t("tabs.goalsWithCount", { count: goals.length })}
        </Link>
      </div>

      {/* Savings Goals Grid */}
      <div className="grid gap-4 sm:grid-cols-2 lg:grid-cols-3">
        {goals.map((goal) => {
          const targetNum = parseFloat(goal.target_amount) || 1;
          const currentNum = parseFloat(goal.current_amount) || 0;
          const percentage = Math.min(100, Math.round((currentNum / targetNum) * 100));
          const linkedAccount = accounts.find((a) => a.id === goal.linked_account_id);

          // Days remaining calculation
          let daysText = t("goals.noDeadline");
          if (goal.target_date) {
            const diffMs = new Date(goal.target_date).getTime() - Date.now();
            const days = Math.ceil(diffMs / (1000 * 60 * 60 * 24));
            if (days < 0) daysText = t("goals.targetDatePassed");
            else if (days === 0) daysText = t("goals.dueToday");
            else daysText = t("goals.daysLeft", { days });
          }

          return (
            <DoubleBezelCard className="p-5" key={goal.id}>
              <div className="flex items-start justify-between">
                <div className="flex items-center gap-2.5">
                  <div
                    className="flex h-9 w-9 items-center justify-center rounded-xl shadow-inner"
                    style={{
                      backgroundColor: `${goal.color || "#8b5cf6"}20`,
                      color: goal.color || "#8b5cf6",
                    }}
                  >
                    <Target className="h-4 w-4" />
                  </div>
                  <div>
                    <h3 className="font-bold text-sm">{goal.name}</h3>
                    <p className="text-[11px] text-muted-foreground flex items-center gap-1">
                      <Calendar className="h-3 w-3" />
                      {daysText}
                    </p>
                  </div>
                </div>

                <div className="flex items-center gap-1">
                  <button
                    aria-label="Edit goal"
                    className="flex h-7 w-7 items-center justify-center rounded-lg border border-white/10 bg-white/5 text-muted-foreground transition hover:bg-white/10 hover:text-foreground"
                    onClick={() => openEditModal(goal)}
                    type="button"
                  >
                    <Edit2 className="h-3 w-3" />
                  </button>
                  <button
                    aria-label="Delete goal"
                    className="flex h-7 w-7 items-center justify-center rounded-lg border border-white/10 bg-white/5 text-muted-foreground transition hover:bg-rose-500/20 hover:text-rose-400"
                    onClick={() => {
                      if (confirm(t("goals.deleteConfirm", { name: goal.name }))) {
                        deleteMutation.mutate(goal.id);
                      }
                    }}
                    type="button"
                  >
                    <Trash2 className="h-3 w-3" />
                  </button>
                </div>
              </div>

              {/* Progress Ring & Balance Layout */}
              <div className="mt-5 flex items-center justify-between gap-4">
                <div>
                  <span className="text-[10px] uppercase font-mono tracking-wider text-muted-foreground">
                    {t("goals.saved")}
                  </span>
                  <div className="font-mono text-lg font-bold text-foreground">
                    <AmountDisplay amount={currentNum} />
                  </div>
                  <div className="text-[11px] text-muted-foreground mt-0.5">
                    {t("goals.target", { amount: goal.target_amount })}
                  </div>
                  {linkedAccount && (
                    <p className="mt-2 text-[10px] text-muted-foreground flex items-center gap-1">
                      <Wallet className="h-3 w-3" /> {linkedAccount.name}
                    </p>
                  )}
                </div>

                <div className="flex flex-col items-center">
                  <ProgressRing
                    className="text-violet-500"
                    color={goal.color || "#8b5cf6"}
                    size={76}
                    strokeWidth={7}
                    value={percentage}
                  >
                    <span className="font-mono text-xs font-bold text-foreground">
                      {percentage}%
                    </span>
                  </ProgressRing>
                </div>
              </div>

              {/* Deposit Quick Action */}
              <div className="mt-4 border-t border-white/10 pt-3">
                <button
                  className="flex w-full items-center justify-center gap-1.5 rounded-xl border border-white/10 bg-white/5 py-2 text-xs font-semibold transition hover:bg-white/10 hover:border-violet-500/40"
                  onClick={() => {
                    setDepositGoal(goal);
                    setDepositAmount("");
                  }}
                  type="button"
                >
                  <Coins className="h-3.5 w-3.5 text-violet-400" />
                  {t("goals.depositFunds")}
                </button>
              </div>
            </DoubleBezelCard>
          );
        })}

        {goals.length === 0 && (
          <div className="col-span-full py-16 text-center">
            <div className="mx-auto flex h-14 w-14 items-center justify-center rounded-2xl bg-white/5 text-muted-foreground">
              <PiggyBank className="h-7 w-7" />
            </div>
            <h2 className="mt-3 text-base font-bold">{t("goals.emptyTitle")}</h2>
            <p className="mt-1 text-xs text-muted-foreground">
              {t("goals.emptySubtitle")}
            </p>
            <div className="mt-4 flex justify-center">
              <ButtonInButton icon={Plus} onClick={openCreateModal} variant="primary">
                {t("goals.createFirst")}
              </ButtonInButton>
            </div>
          </div>
        )}
      </div>

      {/* QUICK DEPOSIT DIALOG */}
      {depositGoal && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4 backdrop-blur-sm">
          <div className="w-full max-w-sm">
            <DoubleBezelCard className="p-6">
              <div className="flex items-center justify-between border-b border-white/10 pb-3">
                <div>
                  <h3 className="text-base font-bold">{t("goals.depositModalTitle")}</h3>
                  <p className="text-xs text-muted-foreground">{depositGoal.name}</p>
                </div>
                <button
                  aria-label="Close modal"
                  className="rounded-lg p-1 text-muted-foreground hover:bg-white/10 hover:text-foreground"
                  onClick={() => setDepositGoal(null)}
                  type="button"
                >
                  <X className="h-4 w-4" />
                </button>
              </div>

              <form className="mt-4 space-y-3" onSubmit={handleDepositSubmit}>
                <div>
                  <label className="text-[11px] font-semibold text-muted-foreground">{t("goals.depositAmountLabel")}</label>
                  <input
                    autoFocus
                    className="mt-1 h-9 w-full rounded-xl border border-white/10 bg-black/30 px-3 text-xs font-mono outline-none focus:border-violet-500"
                    onChange={(e) => setDepositAmount(e.target.value)}
                    placeholder={t("goals.depositAmountPlaceholder")}
                    required
                    type="number"
                    value={depositAmount}
                  />
                </div>

                <div className="grid grid-cols-3 gap-1.5 pt-1">
                  {[500000, 1000000, 5000000].map((quick) => (
                    <button
                      className="rounded-lg border border-white/10 bg-white/5 py-1 text-[11px] font-mono text-muted-foreground hover:bg-white/10 hover:text-foreground"
                      key={quick}
                      onClick={() => setDepositAmount(String(quick))}
                      type="button"
                    >
                      +{quick.toLocaleString()}
                    </button>
                  ))}
                </div>

                <div className="flex justify-end gap-2 pt-3">
                  <button
                    className="rounded-xl border border-white/10 px-4 py-2 text-xs font-semibold hover:bg-white/5"
                    onClick={() => setDepositGoal(null)}
                    type="button"
                  >
                    {t("goals.cancel")}
                  </button>
                  <button
                    className="rounded-xl bg-violet-600 px-4 py-2 text-xs font-semibold text-white shadow-lg shadow-violet-600/20 hover:bg-violet-500"
                    disabled={progressMutation.isPending}
                    type="submit"
                  >
                    {progressMutation.isPending ? t("goals.depositing") : t("goals.confirmDeposit")}
                  </button>
                </div>
              </form>
            </DoubleBezelCard>
          </div>
        </div>
      )}

      {/* CREATE / EDIT MODAL */}
      {isModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4 backdrop-blur-sm">
          <div className="w-full max-w-md">
            <DoubleBezelCard className="p-6">
              <div className="flex items-center justify-between border-b border-white/10 pb-3">
                <h2 className="text-base font-bold">
                  {editingGoal ? t("goals.editModalTitle") : t("goals.createModalTitle")}
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
                  <label className="text-[11px] font-semibold text-muted-foreground">{t("goals.nameLabel")}</label>
                  <input
                    className="mt-1 h-9 w-full rounded-xl border border-white/10 bg-black/30 px-3 text-xs outline-none focus:border-violet-500"
                    onChange={(e) => setFormName(e.target.value)}
                    placeholder={t("goals.namePlaceholder")}
                    required
                    value={formName}
                  />
                </div>

                <div className="grid grid-cols-2 gap-2">
                  <div>
                    <label className="text-[11px] font-semibold text-muted-foreground">{t("goals.targetAmountLabel")}</label>
                    <input
                      className="mt-1 h-9 w-full rounded-xl border border-white/10 bg-black/30 px-3 text-xs font-mono outline-none focus:border-violet-500"
                      onChange={(e) => setFormTarget(e.target.value)}
                      placeholder={t("goals.targetAmountPlaceholder")}
                      required
                      type="number"
                      value={formTarget}
                    />
                  </div>

                  <div>
                    <label className="text-[11px] font-semibold text-muted-foreground">{t("goals.currentAmountLabel")}</label>
                    <input
                      className="mt-1 h-9 w-full rounded-xl border border-white/10 bg-black/30 px-3 text-xs font-mono outline-none focus:border-violet-500"
                      onChange={(e) => setFormCurrent(e.target.value)}
                      placeholder={t("goals.currentAmountPlaceholder")}
                      type="number"
                      value={formCurrent}
                    />
                  </div>
                </div>

                <div className="grid grid-cols-2 gap-2">
                  <div>
                    <label className="text-[11px] font-semibold text-muted-foreground">{t("goals.linkedAccountLabel")}</label>
                    <select
                      className="mt-1 h-9 w-full rounded-xl border border-white/10 bg-black/30 px-3 text-xs outline-none focus:border-violet-500"
                      onChange={(e) => setFormAccountId(e.target.value)}
                      value={formAccountId}
                    >
                      <option value="">{t("goals.linkedAccountNone")}</option>
                      {accounts.map((a) => (
                        <option key={a.id} value={a.id}>
                          {a.name} ({a.currency})
                        </option>
                      ))}
                    </select>
                  </div>

                  <div>
                    <label className="text-[11px] font-semibold text-muted-foreground">{t("goals.targetDateLabel")}</label>
                    <input
                      className="mt-1 h-9 w-full rounded-xl border border-white/10 bg-black/30 px-3 text-xs outline-none focus:border-violet-500"
                      onChange={(e) => setFormDate(e.target.value)}
                      type="date"
                      value={formDate}
                    />
                  </div>
                </div>

                <div>
                  <label className="text-[11px] font-semibold text-muted-foreground">{t("goals.colorThemeLabel")}</label>
                  <div className="mt-1.5 flex items-center gap-2">
                    {["#8b5cf6", "#10b981", "#3b82f6", "#f59e0b", "#ec4899", "#06b6d4"].map((c) => (
                      <button
                        className={`h-6 w-6 rounded-full transition ${formColor === c ? "ring-2 ring-white ring-offset-2 ring-offset-black scale-110" : ""}`}
                        key={c}
                        onClick={() => setFormColor(c)}
                        style={{ backgroundColor: c }}
                        type="button"
                      />
                    ))}
                  </div>
                </div>

                <div className="flex justify-end gap-2 pt-2">
                  <button
                    className="rounded-xl border border-white/10 px-4 py-2 text-xs font-semibold hover:bg-white/5"
                    onClick={() => setIsModalOpen(false)}
                    type="button"
                  >
                    {t("goals.cancel")}
                  </button>
                  <button
                    className="rounded-xl bg-violet-600 px-4 py-2 text-xs font-semibold text-white shadow-lg shadow-violet-600/20 hover:bg-violet-500"
                    type="submit"
                  >
                    {editingGoal ? t("goals.update") : t("goals.save")}
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
