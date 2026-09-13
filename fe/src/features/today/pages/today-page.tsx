import { useEffect, useMemo, useState } from "react";
import { Link } from "react-router-dom";
import { useTranslation } from "react-i18next";
import {
  Sparkles,
  Play,
  RotateCcw,
  ArrowRight,
  GraduationCap,
  CheckCircle2,
  Wallet,
  Flame,
  Plus,
  ListTodo,
  BookOpen,
  ArrowUpRight,
  Clock,
  Check,
  Pause,
} from "lucide-react";
import { DoubleBezelCard } from "../../../components/design-system/double-bezel-card";
import { ButtonInButton } from "../../../components/design-system/button-in-button";
import { AmountDisplay } from "../../../components/design-system/amount-display";
import { ProgressRing } from "../../../components/design-system/progress-ring";
import { useAuth } from "../../auth/hooks/use-auth";
import { useTasks } from "../../tasks/hooks/use-tasks";
import { usePomodoro, useSetPomodoro, useStopPomodoro } from "../../pomodoro/hooks/use-pomodoro";
import { useReviewQueue } from "../../linguistics/hooks/use-vocabularies";
import { useFinanceDashboard } from "../../finance/hooks/use-finance";
import { useJournalStats } from "../../journal/hooks/use-journals";
import { cn } from "../../../lib/cn";

export function TodayPage() {
  const { t, i18n } = useTranslation();
  const { user } = useAuth();

  // Queries
  const { data: tasks, isLoading: tasksLoading } = useTasks();
  const { data: pomodoro } = usePomodoro();
  const { data: reviewQueueData } = useReviewQueue();
  const { data: financeData } = useFinanceDashboard();
  const { data: journalStats } = useJournalStats();

  // Mutations
  const setPomodoro = useSetPomodoro();
  const stopPomodoro = useStopPomodoro();

  // Current timestamp for live timer interpolation
  const [now, setNow] = useState(Date.now());
  const isPomodoroRunning = pomodoro !== null && pomodoro !== undefined && pomodoro.status === "running";

  useEffect(() => {
    if (!isPomodoroRunning) return undefined;
    const interval = window.setInterval(() => setNow(Date.now()), 500);
    return () => window.clearInterval(interval);
  }, [isPomodoroRunning]);

  // Greeting based on time of day
  const hour = new Date().getHours();
  const greetingKey =
    hour < 12
      ? t("today.greetingMorning")
      : hour < 18
      ? t("today.greetingAfternoon")
      : t("today.greetingEvening");
  const userGreeting = user?.displayName ? `${greetingKey}, ${user.displayName}` : greetingKey;

  // Filter next 3 actionable tasks (TODO or IN_PROGRESS)
  const pendingTasks = (tasks ?? [])
    .filter((task) => task.status === "TODO" || task.status === "IN_PROGRESS")
    .slice(0, 3);

  // Pomodoro countdown calculations
  const totalDuration = pomodoro?.duration_seconds ?? 1500;
  const remainingSeconds = useMemo(() => {
    if (!pomodoro) return 1500;
    const elapsed = Math.max(0, Math.floor((now - new Date(pomodoro.started_at).getTime()) / 1000));
    return Math.max(0, pomodoro.duration_seconds - elapsed);
  }, [now, pomodoro]);

  const pomodoroMinutes = Math.floor(remainingSeconds / 60);
  const pomodoroSecs = remainingSeconds % 60;
  const pomodoroTimeStr = `${String(pomodoroMinutes).padStart(2, "0")}:${String(pomodoroSecs).padStart(2, "0")}`;

  const handleTogglePomodoro = () => {
    if (isPomodoroRunning) {
      stopPomodoro.mutate();
    } else {
      setPomodoro.mutate({
        status: "running",
        started_at: new Date().toISOString(),
        duration_seconds: 1500,
        task_id: null,
        updated_at: new Date().toISOString(),
      });
    }
  };

  const handleResetPomodoro = () => {
    stopPomodoro.mutate();
  };

  const queueCount = reviewQueueData?.length ?? 0;
  const streakDays = journalStats?.current_streak ?? 0;
  const netWorth = financeData?.net_worth ?? "0";

  return (
    <div className="flex flex-col gap-6 pb-12 animate-in fade-in duration-500">
      {/* Top Welcome Strip */}
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <div className="inline-flex items-center gap-2 rounded-full border border-indigo-500/20 bg-indigo-500/10 px-3 py-1 text-xs font-semibold text-indigo-400">
            <Sparkles size={13} className="animate-pulse" />
            <span>{t("today.workspaceTag")}</span>
          </div>
          <h1 className="mt-2 text-2xl sm:text-3xl font-bold tracking-tight text-foreground">
            {userGreeting}
          </h1>
          <p className="mt-1 text-sm text-muted-foreground">
            {new Date().toLocaleDateString(
              i18n.language === "vi" ? "vi-VN" : i18n.language === "ja" ? "ja-JP" : "en-US",
              {
                weekday: "long",
                year: "numeric",
                month: "long",
                day: "numeric",
              }
            )}
          </p>
        </div>

        <div className="flex items-center gap-2.5">
          <Link to="/app/tasks">
            <ButtonInButton icon={Plus} variant="primary" size="md">
              {t("today.addTask")}
            </ButtonInButton>
          </Link>
          <Link to="/app/journal">
            <ButtonInButton icon={BookOpen} variant="secondary" size="md">
              {t("today.writeJournal")}
            </ButtonInButton>
          </Link>
        </div>
      </div>

      {/* Primary Asymmetric Bento Grid */}
      <div className="grid grid-cols-12 gap-5">
        {/* Zone 1: Next Focus Tasks (Hero Left - col-span-8) */}
        <DoubleBezelCard
          outerClassName="col-span-12 lg:col-span-8"
          innerClassName="flex flex-col justify-between"
          glowColor="rgba(99, 102, 241, 0.12)"
        >
          <div>
            <div className="flex items-center justify-between border-b border-slate-200/80 dark:border-white/10 pb-4">
              <div className="flex items-center gap-2.5">
                <div className="size-9 rounded-xl bg-blue-500/10 border border-blue-500/20 flex items-center justify-center text-blue-400">
                  <ListTodo size={18} />
                </div>
                <div>
                  <h2 className="text-base font-semibold text-foreground">{t("today.focusTitle")}</h2>
                  <p className="text-xs text-muted-foreground">{t("today.focusSubtitle")}</p>
                </div>
              </div>
              <Link
                to="/app/tasks"
                className="group text-xs font-medium text-blue-400 hover:text-blue-300 inline-flex items-center gap-1 transition-colors"
              >
                <span>{t("today.viewAll")}</span>
                <ArrowRight size={13} className="transition-transform group-hover:translate-x-0.5" />
              </Link>
            </div>

            {/* Task list */}
            <div className="mt-4 space-y-3">
              {tasksLoading ? (
                <div className="space-y-2 py-4">
                  <div className="h-12 w-full rounded-xl bg-slate-100 dark:bg-white/5 animate-pulse" />
                  <div className="h-12 w-full rounded-xl bg-slate-100 dark:bg-white/5 animate-pulse" />
                </div>
              ) : pendingTasks.length === 0 ? (
                <div className="flex flex-col items-center justify-center py-8 text-center">
                  <div className="size-12 rounded-full bg-emerald-500/10 border border-emerald-500/20 flex items-center justify-center text-emerald-400 mb-2">
                    <CheckCircle2 size={24} />
                  </div>
                  <p className="text-sm font-medium text-foreground">{t("today.allDone")}</p>
                  <p className="text-xs text-muted-foreground mt-0.5">{t("today.allDoneSubtitle")}</p>
                </div>
              ) : (
                pendingTasks.map((task) => (
                  <div
                    key={task.id}
                    className="flex items-center justify-between p-3.5 rounded-xl border border-slate-200/80 bg-slate-50/70 hover:bg-slate-100/80 dark:border-white/5 dark:bg-white/[0.02] dark:hover:bg-white/[0.05] transition-colors"
                  >
                    <div className="flex items-center gap-3 min-w-0">
                      <div
                        className={cn(
                          "size-2.5 rounded-full shrink-0",
                          task.priority === "HIGH" || task.priority === "URGENT"
                            ? "bg-rose-500 ring-4 ring-rose-500/20"
                            : task.priority === "MEDIUM"
                            ? "bg-amber-500 ring-4 ring-amber-500/20"
                            : "bg-blue-500 ring-4 ring-blue-500/20"
                        )}
                      />
                      <div className="min-w-0">
                        <p className="text-sm font-medium text-foreground truncate">{task.title}</p>
                        <div className="flex items-center gap-2 mt-0.5">
                          <span className="text-[10px] font-semibold tracking-wider uppercase px-2 py-0.5 rounded-md bg-slate-100 text-slate-600 dark:bg-white/10 dark:text-muted-foreground">
                            {task.category}
                          </span>
                          {task.dueDate && (
                            <span className="text-xs text-muted-foreground flex items-center gap-1">
                              <Clock size={11} />
                              {new Date(task.dueDate).toLocaleDateString()}
                            </span>
                          )}
                        </div>
                      </div>
                    </div>

                    <Link to="/app/tasks">
                      <button
                        type="button"
                        className="size-8 rounded-lg bg-slate-100 hover:bg-slate-200 text-muted-foreground hover:text-foreground dark:bg-white/5 dark:hover:bg-white/10 flex items-center justify-center transition-colors"
                        title={t("today.taskDetail")}
                      >
                        <ArrowRight size={14} />
                      </button>
                    </Link>
                  </div>
                ))
              )}
            </div>
          </div>

          <div className="mt-5 pt-3 border-t border-slate-200/60 dark:border-white/5 flex items-center justify-between text-xs text-muted-foreground">
            <span>
              {t("today.totalTasks")} <strong className="text-foreground">{tasks?.length ?? 0}</strong>
            </span>
            <span className="flex items-center gap-1.5 text-emerald-400 font-medium">
              <Check size={13} />
              {t("today.completed")} {tasks?.filter((t) => t.status === "DONE").length ?? 0}
            </span>
          </div>
        </DoubleBezelCard>

        {/* Zone 2: Live Pomodoro Focus Ring (Hero Right - col-span-4) */}
        <DoubleBezelCard
          outerClassName="col-span-12 lg:col-span-4"
          innerClassName="flex flex-col items-center justify-between text-center"
          glowColor="rgba(244, 63, 94, 0.12)"
        >
          <div className="w-full flex items-center justify-between border-b border-slate-200/80 dark:border-white/10 pb-3">
            <div className="flex items-center gap-2">
              <div
                className={cn(
                  "size-2 rounded-full",
                  isPomodoroRunning ? "bg-rose-500 animate-ping" : "bg-muted-foreground"
                )}
              />
              <span className="text-xs font-semibold uppercase tracking-wider text-rose-400">
                {isPomodoroRunning ? t("today.running") : t("pomodoro.idle")}
              </span>
            </div>
            <Link
              to="/app/pomodoro"
              className="text-xs text-muted-foreground hover:text-foreground transition-colors"
            >
              {t("today.expand")}
            </Link>
          </div>

          {/* Circular Progress Display */}
          <div className="py-4 my-auto">
            <ProgressRing
              value={totalDuration - remainingSeconds}
              max={totalDuration}
              size={140}
              strokeWidth={10}
              color="#f43f5e"
            >
              <div className="flex flex-col items-center justify-center">
                <span className="font-mono text-2xl font-bold tracking-tight text-foreground tabular-nums">
                  {pomodoroTimeStr}
                </span>
                <span className="text-[10px] font-semibold text-muted-foreground uppercase mt-0.5">
                  {isPomodoroRunning ? t("today.running") : t("today.paused")}
                </span>
              </div>
            </ProgressRing>
          </div>

          {/* Controls */}
          <div className="w-full flex items-center justify-center gap-3 pt-2">
            <button
              type="button"
              onClick={handleTogglePomodoro}
              className={cn(
                "flex-1 inline-flex items-center justify-center gap-2 rounded-full py-2 px-4 text-xs font-semibold transition-all shadow-md active:scale-95 cursor-pointer",
                isPomodoroRunning
                  ? "bg-amber-600 hover:bg-amber-500 text-white shadow-amber-950/40"
                  : "bg-rose-600 hover:bg-rose-500 text-white shadow-rose-950/40"
              )}
            >
              {isPomodoroRunning ? <Pause size={14} /> : <Play size={14} />}
              <span>{isPomodoroRunning ? t("today.stop") : t("today.start")}</span>
            </button>
            <button
              type="button"
              onClick={handleResetPomodoro}
              title={t("today.resetTip")}
              className="size-9 rounded-full bg-slate-100 hover:bg-slate-200 border border-slate-200 text-muted-foreground hover:text-foreground dark:bg-white/10 dark:hover:bg-white/15 dark:border-white/10 flex items-center justify-center transition-colors cursor-pointer"
            >
              <RotateCcw size={14} />
            </button>
          </div>
        </DoubleBezelCard>

        {/* Zone 3: Linguistics SRS Queue (col-span-12 md:col-span-4) */}
        <DoubleBezelCard
          outerClassName="col-span-12 md:col-span-4"
          innerClassName="flex flex-col justify-between"
          glowColor="rgba(168, 85, 247, 0.12)"
        >
          <div>
            <div className="flex items-center justify-between">
              <div className="size-9 rounded-xl bg-violet-500/10 border border-violet-500/20 flex items-center justify-center text-violet-400">
                <GraduationCap size={18} />
              </div>
              <span className="text-xs font-semibold px-2 py-0.5 rounded-full bg-violet-500/10 text-violet-400 border border-violet-500/20">
                SM-2 SRS
              </span>
            </div>

            <div className="mt-4">
              <h3 className="text-sm font-semibold text-foreground">{t("today.reviewVocab")}</h3>
              <p className="text-xs text-muted-foreground mt-0.5">{t("today.reviewVocabSubtitle")}</p>
            </div>

            <div className="mt-4 flex items-baseline gap-2">
              <span className="font-mono text-3xl font-bold tracking-tight text-violet-400 tabular-nums">
                {queueCount}
              </span>
              <span className="text-xs text-muted-foreground">{t("today.wordsDueToday")}</span>
            </div>
          </div>

          <div className="mt-6">
            <Link to="/app/learn" className="w-full block">
              <ButtonInButton
                icon={ArrowRight}
                variant="violet"
                size="sm"
                className="w-full justify-between"
              >
                {t("today.startLearning")}
              </ButtonInButton>
            </Link>
          </div>
        </DoubleBezelCard>

        {/* Zone 4: Finance Snapshot (col-span-12 md:col-span-4) */}
        <DoubleBezelCard
          outerClassName="col-span-12 md:col-span-4"
          innerClassName="flex flex-col justify-between"
          glowColor="rgba(168, 85, 247, 0.12)"
        >
          <div>
            <div className="flex items-center justify-between">
              <div className="size-9 rounded-xl bg-emerald-500/10 border border-emerald-500/20 flex items-center justify-center text-emerald-400">
                <Wallet size={18} />
              </div>
              <span className="text-xs font-semibold px-2 py-0.5 rounded-full bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
                {t("navigation.assets")}
              </span>
            </div>

            <div className="mt-4">
              <h3 className="text-sm font-semibold text-foreground">{t("today.netWorth")}</h3>
              <p className="text-xs text-muted-foreground mt-0.5">{t("today.netWorthSubtitle")}</p>
            </div>

            <div className="mt-4">
              <AmountDisplay
                amount={netWorth}
                currency="VND"
                className="text-2xl font-bold text-foreground"
              />
              <div className="flex items-center gap-3 mt-2 text-xs text-muted-foreground">
                <span>{t("today.income")} +{financeData?.cashflow?.income ?? "0"} ₫</span>
                <span>{t("today.expense")} -{financeData?.cashflow?.expense ?? "0"} ₫</span>
              </div>
            </div>
          </div>

          <div className="mt-6">
            <Link to="/app/finance" className="w-full block">
              <ButtonInButton
                icon={ArrowUpRight}
                variant="emerald"
                size="sm"
                className="w-full justify-between"
              >
                {t("today.cashflowDetail")}
              </ButtonInButton>
            </Link>
          </div>
        </DoubleBezelCard>

        {/* Zone 5: Journal Streak (col-span-12 md:col-span-4) */}
        <DoubleBezelCard
          outerClassName="col-span-12 md:col-span-4"
          innerClassName="flex flex-col justify-between"
          glowColor="rgba(245, 158, 11, 0.12)"
        >
          <div>
            <div className="flex items-center justify-between">
              <div className="size-9 rounded-xl bg-amber-500/10 border border-amber-500/20 flex items-center justify-center text-amber-400">
                <Flame size={18} />
              </div>
              <span className="text-xs font-semibold px-2 py-0.5 rounded-full bg-amber-500/10 text-amber-400 border border-amber-500/20">
                {t("navigation.journal")}
              </span>
            </div>

            <div className="mt-4">
              <h3 className="text-sm font-semibold text-foreground">{t("today.journalStreak")}</h3>
              <p className="text-xs text-muted-foreground mt-0.5">{t("today.journalStreakSubtitle")}</p>
            </div>

            <div className="mt-4 flex items-baseline gap-2">
              <span className="font-mono text-3xl font-bold tracking-tight text-amber-400 tabular-nums">
                {streakDays}
              </span>
              <span className="text-xs text-muted-foreground">{t("today.consecutiveDays")}</span>
            </div>
          </div>

          <div className="mt-6">
            <Link to="/app/journal" className="w-full block">
              <ButtonInButton
                icon={Plus}
                variant="amber"
                size="sm"
                className="w-full justify-between"
              >
                {t("today.writeEntry")}
              </ButtonInButton>
            </Link>
          </div>
        </DoubleBezelCard>
      </div>
    </div>
  );
}
