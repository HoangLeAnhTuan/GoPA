import { motion } from "framer-motion";
import { Coffee, Moon, Pause, Play, RotateCcw, Target, TrendingUp, Zap, History, CheckCircle2 } from "lucide-react";
import React, { useEffect, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { GlassPanel } from "../components/design-system/glass-panel";
import { PageHeader } from "../components/design-system/page-header";
import { usePomodoro, useSetPomodoro, useStopPomodoro } from "../features/pomodoro/hooks/use-pomodoro";
import { usePomodoroWs, usePomodoroHistory } from "../features/pomodoro/hooks/use-pomodoro-ws";
import { cn } from "../lib/cn";

type SessionMode = "focus" | "short" | "long";

const SESSION_CONFIG: Record<SessionMode, { durationSecs: number; color: string; ringColor: string; glowColor: string }> = {
  focus: { durationSecs: 1500, color: "from-blue-500 to-indigo-600", ringColor: "stroke-blue-500", glowColor: "shadow-blue-500/25" },
  short: { durationSecs: 300, color: "from-emerald-500 to-teal-500", ringColor: "stroke-emerald-500", glowColor: "shadow-emerald-500/25" },
  long: { durationSecs: 900, color: "from-violet-500 to-purple-600", ringColor: "stroke-violet-500", glowColor: "shadow-violet-500/25" },
};

const CIRCUMFERENCE = 2 * Math.PI * 42; // r=42

export function PomodoroPage() {
  const { t } = useTranslation();
  const pomodoro = usePomodoro();
  const { status: wsStatus } = usePomodoroWs();
  const { data: historyItems } = usePomodoroHistory(5);

  const normalizedHistory = useMemo(() => historyItems ?? [], [historyItems]);

  const todayHistory = useMemo(() => {
    const today = new Date().toDateString();
    return normalizedHistory.filter((item) => {
      const d = new Date(item.ended_at);
      return !isNaN(d.getTime()) && d.toDateString() === today;
    });
  }, [normalizedHistory]);

  const [now, setNow] = useState(Date.now());
  const [mode, setMode] = useState<SessionMode>("focus");
  const [sessionCount, setSessionCount] = useState(0);
  const setPomodoro = useSetPomodoro();
  const stopPomodoro = useStopPomodoro();

  const totalSessionsToday = todayHistory.length + sessionCount;
  const totalFocusSecondsToday = useMemo(() => {
    return todayHistory.reduce((acc, it) => acc + it.duration_seconds, sessionCount * 1500);
  }, [todayHistory, sessionCount]);

  const config = SESSION_CONFIG[mode];

  const remainingSeconds = useMemo(() => {
    if (pomodoro.data === null || pomodoro.data === undefined) return config.durationSecs;
    const elapsed = Math.max(0, Math.floor((now - new Date(pomodoro.data.started_at).getTime()) / 1000));
    return Math.max(0, pomodoro.data.duration_seconds - elapsed);
  }, [now, pomodoro.data, config.durationSecs]);

  const formattedTime = useMemo(() => {
    const m = Math.floor(remainingSeconds / 60).toString().padStart(2, "0");
    const s = (remainingSeconds % 60).toString().padStart(2, "0");
    return `${m}:${s}`;
  }, [remainingSeconds]);

  const progressPercent = useMemo(() => {
    const total = pomodoro.data?.duration_seconds ?? config.durationSecs;
    return Math.min(100, Math.max(0, ((total - remainingSeconds) / total) * 100));
  }, [remainingSeconds, pomodoro.data, config.durationSecs]);

  const strokeDashoffset = CIRCUMFERENCE - (CIRCUMFERENCE * progressPercent) / 100;
  const isRunning = pomodoro.data !== null && pomodoro.data !== undefined && pomodoro.data.status === "running";

  useEffect(() => {
    if (!isRunning) return undefined;
    const timer = window.setInterval(() => setNow(Date.now()), 500);
    return () => window.clearInterval(timer);
  }, [isRunning]);

  useEffect(() => {
    if (isRunning && remainingSeconds === 0) {
      stopPomodoro.mutate();
      setSessionCount((n) => n + 1);
    }
  }, [remainingSeconds, isRunning, stopPomodoro]);

  const start = () =>
    setPomodoro.mutate({
      status: "running",
      started_at: new Date().toISOString(),
      duration_seconds: config.durationSecs,
      task_id: null,
      updated_at: new Date().toISOString(),
    });

  const reset = () => stopPomodoro.mutate();

  const modeStatusLabel: Record<SessionMode, string> = {
    focus: t("pomodoro.focusing"),
    short: t("pomodoro.idle"),
    long: t("pomodoro.idle"),
  };

  const modeIcons: Record<SessionMode, React.ReactNode> = {
    focus: <Zap size={12} />,
    short: <Coffee size={12} />,
    long: <Moon size={12} />,
  };

  const modeLabels: Record<SessionMode, string> = {
    focus: t("pomodoro.modesFocus"),
    short: t("pomodoro.modesShort"),
    long: t("pomodoro.modesLong"),
  };

  return (
    <div className="flex flex-col gap-5 pb-12 animate-in fade-in duration-500">
      <div className="flex items-center justify-between">
        <PageHeader
          description={t("pomodoro.pageDescription")}
          title={t("pomodoro.title")}
        />
        <div className="flex items-center gap-2">
          {wsStatus === "connected" ? (
            <span className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-semibold bg-emerald-500/10 text-emerald-400 border border-emerald-500/20">
              <span className="size-1.5 rounded-full bg-emerald-400 animate-ping" />
              {t("pomodoro.wsSync")}
            </span>
          ) : (
            <span className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-xs font-semibold bg-amber-500/10 text-amber-400 border border-amber-500/20">
              {t("pomodoro.wsReconnecting")}
            </span>
          )}
        </div>
      </div>

      {/* Main layout: ring left, stats/tip right */}
      <div className="grid gap-5 lg:grid-cols-[1.5fr_1fr]">
        {/* Left: Timer + Ring + Controls */}
        <GlassPanel className="flex flex-col items-center justify-between p-6">
          {/* Mode selector pills */}
          <div className="flex gap-1.5 rounded-2xl border border-white/30 bg-white/40 p-1.5 dark:border-white/10 dark:bg-slate-900/60">
            {(["focus", "short", "long"] as SessionMode[]).map((m) => (
              <button
                className={cn(
                  "flex items-center gap-1.5 rounded-xl px-4 py-2 text-xs font-semibold transition-all cursor-pointer",
                  mode === m
                    ? "bg-primary text-primary-foreground shadow-md shadow-primary/25"
                    : "text-muted-foreground hover:bg-white/60 hover:text-foreground dark:hover:bg-slate-800/60"
                )}
                disabled={isRunning}
                key={m}
                onClick={() => setMode(m)}
                type="button"
              >
                {modeIcons[m]}
                <span>{modeLabels[m]}</span>
              </button>
            ))}
          </div>

          {/* Circular progress with clock display */}
          <div className="relative my-6 flex size-64 items-center justify-center">
            <svg className="size-full -rotate-90" viewBox="0 0 100 100">
              <circle
                className="stroke-slate-200/50 dark:stroke-slate-800/80"
                cx="50"
                cy="50"
                fill="none"
                r="42"
                strokeWidth="6"
              />
              <circle
                className={cn("transition-all duration-500", config.ringColor)}
                cx="50"
                cy="50"
                fill="none"
                r="42"
                strokeDasharray={CIRCUMFERENCE}
                strokeDashoffset={strokeDashoffset}
                strokeLinecap="round"
                strokeWidth="6"
              />
            </svg>
            <div className="absolute flex flex-col items-center gap-1 text-center">
              <span className="font-mono text-4xl font-bold tracking-tight text-foreground tabular-nums">
                {formattedTime}
              </span>
              <span className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
                {isRunning ? modeStatusLabel[mode] : t("pomodoro.idle")}
              </span>
            </div>
          </div>

          {/* Action buttons */}
          <div className="flex items-center gap-3">
            {!isRunning ? (
              <motion.button
                className={cn(
                  "inline-flex items-center gap-2 rounded-2xl bg-gradient-to-r px-7 py-3 text-sm font-semibold text-white shadow-lg cursor-pointer",
                  config.color,
                  config.glowColor
                )}
                onClick={start}
                type="button"
                whileHover={{ scale: 1.03 }}
                whileTap={{ scale: 0.96 }}
              >
                <Play size={16} />
                {t("pomodoro.start")}
              </motion.button>
            ) : (
              <>
                <motion.button
                  className="inline-flex items-center gap-2 rounded-2xl bg-amber-500 px-6 py-3 text-sm font-semibold text-white shadow-lg shadow-amber-500/25 hover:bg-amber-600 cursor-pointer"
                  onClick={reset}
                  type="button"
                  whileHover={{ scale: 1.02 }}
                  whileTap={{ scale: 0.96 }}
                >
                  <Pause size={16} />
                  {t("pomodoro.pause")}
                </motion.button>
                <motion.button
                  className="inline-flex items-center gap-1.5 rounded-2xl border border-white/40 bg-white/60 px-4 py-3 text-sm font-semibold text-foreground shadow-sm hover:bg-white/80 dark:border-white/10 dark:bg-slate-800 cursor-pointer"
                  onClick={reset}
                  type="button"
                  whileHover={{ scale: 1.02 }}
                  whileTap={{ scale: 0.96 }}
                >
                  <RotateCcw size={14} />
                  {t("pomodoro.reset")}
                </motion.button>
              </>
            )}
          </div>

          {/* Session progress dots */}
          <div className="mt-6 flex items-center gap-2">
            {Array.from({ length: 4 }).map((_, i) => (
              <div
                className={cn(
                  "size-2.5 rounded-full transition-all duration-300",
                  i < (sessionCount % 4)
                    ? `bg-gradient-to-br ${config.color}`
                    : "bg-slate-300/60 dark:bg-slate-700"
                )}
                key={i}
              />
            ))}
            <span className="ml-1 text-xs text-muted-foreground">
              {sessionCount} {t("pomodoro.sessionsCompleted")}
            </span>
          </div>
        </GlassPanel>

        {/* Right: Stats + History stacked */}
        <div className="flex flex-col gap-4">
          {/* Stats card */}
          <GlassPanel className="flex flex-col gap-3 p-5">
            <p className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              {t("pomodoro.stats.sessions")}
            </p>
            <div className="grid grid-cols-1 gap-2.5">
              <StatRow
                color="text-blue-500"
                icon={<Target size={14} />}
                label={t("pomodoro.statsSessions")}
                value={String(totalSessionsToday)}
              />
              <StatRow
                color="text-emerald-500"
                icon={<TrendingUp size={14} />}
                label={t("pomodoro.statsTotalFocus")}
                value={`${Math.floor(totalFocusSecondsToday / 3600)}h ${Math.floor((totalFocusSecondsToday % 3600) / 60)}m`}
              />
              <StatRow
                color="text-violet-500"
                icon={<Zap size={14} />}
                label={t("pomodoro.statsStreak")}
                value={totalSessionsToday > 0 ? `🔥 ${totalSessionsToday}` : "—"}
              />
            </div>
          </GlassPanel>

          {/* Recent History */}
          <GlassPanel className="flex flex-1 flex-col p-5">
            <div className="flex items-center justify-between border-b border-white/10 pb-3 mb-3">
              <div className="flex items-center gap-2">
                <History size={16} className="text-rose-400" />
                <h3 className="text-xs font-semibold uppercase tracking-wider text-foreground">
                  {t("pomodoro.recentHistory")}
                </h3>
              </div>
            </div>

            <div className="space-y-2 overflow-y-auto max-h-56">
              {normalizedHistory.length === 0 ? (
                <p className="text-xs text-muted-foreground text-center py-6">
                  {t("pomodoro.noHistory")}
                </p>
              ) : (
                normalizedHistory.map((item) => {
                  const minutes = Math.max(1, Math.floor(item.duration_seconds / 60));
                  const parsedDate = new Date(item.ended_at);
                  const timeLabel = !isNaN(parsedDate.getTime())
                    ? parsedDate.toLocaleTimeString([], { hour: "2-digit", minute: "2-digit" })
                    : "";
                  return (
                    <div
                      key={item.id}
                      className="flex items-center justify-between p-2.5 rounded-xl bg-white/5 border border-white/5 text-xs"
                    >
                      <div className="flex items-center gap-2">
                        <CheckCircle2 size={13} className="text-emerald-400 shrink-0" />
                        <span className="font-medium text-foreground">
                          {t("pomodoro.minutesFocus", { minutes })}
                        </span>
                      </div>
                      {timeLabel && (
                        <span className="text-[11px] text-muted-foreground">
                          {timeLabel}
                        </span>
                      )}
                    </div>
                  );
                })
              )}
            </div>
          </GlassPanel>
        </div>
      </div>
    </div>
  );
}

function StatRow({ icon, label, value, color }: { icon: React.ReactNode; label: string; value: string; color: string }) {
  return (
    <div className="flex items-center justify-between py-1">
      <div className={cn("flex items-center gap-2", color)}>
        {icon}
        <span className="text-xs text-muted-foreground">{label}</span>
      </div>
      <span className="font-mono text-sm font-bold tabular-nums text-foreground">{value}</span>
    </div>
  );
}
