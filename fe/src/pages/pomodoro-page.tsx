import { AnimatePresence, motion } from "framer-motion";
import { Coffee, Moon, Pause, Play, RotateCcw, Target, Timer, TrendingUp, Zap } from "lucide-react";
import { useEffect, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { GlassPanel } from "../components/design-system/glass-panel";
import { PageHeader } from "../components/design-system/page-header";
import { usePomodoro, useSetPomodoro, useStopPomodoro } from "../features/pomodoro/hooks/use-pomodoro";
import { cn } from "../lib/cn";

type SessionMode = "focus" | "short" | "long";

const SESSION_CONFIG: Record<SessionMode, { durationSecs: number; color: string; ringColor: string; glowColor: string }> = {
  focus: { durationSecs: 1500, color: "from-blue-500 to-indigo-600",   ringColor: "stroke-blue-500",   glowColor: "shadow-blue-500/25" },
  short: { durationSecs: 300,  color: "from-emerald-500 to-teal-500",  ringColor: "stroke-emerald-500", glowColor: "shadow-emerald-500/25" },
  long:  { durationSecs: 900,  color: "from-violet-500 to-purple-600", ringColor: "stroke-violet-500",  glowColor: "shadow-violet-500/25" },
};

const CIRCUMFERENCE = 2 * Math.PI * 42; // r=42

export function PomodoroPage() {
  const { t } = useTranslation();
  const pomodoro = usePomodoro();
  const [now, setNow] = useState(Date.now());
  const [mode, setMode] = useState<SessionMode>("focus");
  const [sessionCount, setSessionCount] = useState(0);
  const setPomodoro = useSetPomodoro();
  const stopPomodoro = useStopPomodoro();

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
  const isRunning = pomodoro.data !== null && pomodoro.data !== undefined;

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
    long:  t("pomodoro.modesLong"),
  };

  return (
    <div className="flex h-[calc(100vh-8rem)] flex-col gap-3 overflow-hidden">
      <PageHeader
        description={t("pomodoro.pageDescription")}
        title={t("pomodoro.title")}
      />

      {/* Main layout: ring left, stats/tip right — everything fits on one screen */}
      <div className="grid min-h-0 flex-1 gap-3 lg:grid-cols-[1fr_280px]">

        {/* Left: Mode tabs + Ring + Buttons */}
        <GlassPanel className="flex flex-col items-center justify-center gap-4 p-5">
          {/* Mode selector tabs */}
          <div className="flex w-full max-w-sm gap-1 rounded-2xl border border-white/30 bg-white/30 p-1 dark:border-white/10 dark:bg-slate-900/40">
            {(["focus", "short", "long"] as SessionMode[]).map((m) => (
              <button
                className={cn(
                  "flex flex-1 items-center justify-center gap-1.5 rounded-xl py-1.5 text-xs font-semibold transition-all duration-200",
                  mode === m
                    ? `bg-gradient-to-r ${SESSION_CONFIG[m].color} text-white shadow-md`
                    : "text-muted-foreground hover:bg-white/50 dark:hover:bg-slate-800/60"
                )}
                disabled={isRunning}
                key={m}
                onClick={() => setMode(m)}
                type="button"
              >
                {modeIcons[m]}
                <span className="hidden sm:inline">{modeLabels[m]}</span>
                <span className="sm:hidden">{modeLabels[m].split(" ·")[0]}</span>
              </button>
            ))}
          </div>

          {/* SVG Ring Timer — compact size */}
          <div className={cn("relative flex size-52 items-center justify-center rounded-full shadow-2xl sm:size-60", config.glowColor)}>
            <div className={cn("absolute inset-5 rounded-full opacity-10 blur-2xl bg-gradient-to-br", config.color)} />

            <svg className="absolute inset-0 size-full -rotate-90" viewBox="0 0 100 100">
              <circle
                className="text-slate-200/60 dark:text-slate-800"
                cx="50" cy="50" fill="none"
                r="42" stroke="currentColor" strokeWidth="4"
              />
              <motion.circle
                animate={{ strokeDashoffset }}
                className={config.ringColor}
                cx="50" cy="50" fill="none" r="42"
                stroke="currentColor"
                strokeDasharray={CIRCUMFERENCE}
                strokeDashoffset={strokeDashoffset}
                strokeLinecap="round" strokeWidth="4"
                transition={{ duration: 0.5, ease: "linear" }}
              />
            </svg>

            <div className="relative z-10 flex flex-col items-center">
              <div className={cn("mb-1 flex size-7 items-center justify-center rounded-full bg-gradient-to-br", config.color)}>
                <Timer className="text-white" size={14} />
              </div>
              <span className="font-mono text-5xl font-bold tracking-tight tabular-nums sm:text-[3.25rem]">
                {formattedTime}
              </span>
              <AnimatePresence mode="wait">
                <motion.p
                  animate={{ opacity: 1, y: 0 }}
                  className="mt-0.5 text-[10px] font-semibold uppercase tracking-widest text-muted-foreground"
                  exit={{ opacity: 0, y: -3 }}
                  initial={{ opacity: 0, y: 3 }}
                  key={isRunning ? "running" : "idle"}
                  transition={{ duration: 0.2 }}
                >
                  {isRunning ? modeStatusLabel[mode] : t("pomodoro.idle")}
                </motion.p>
              </AnimatePresence>
            </div>
          </div>

          {/* Action buttons */}
          <div className="flex items-center gap-2.5">
            {!isRunning ? (
              <motion.button
                className={cn(
                  "inline-flex items-center gap-2 rounded-2xl bg-gradient-to-r px-7 py-2.5 text-sm font-bold text-white shadow-lg transition-all",
                  config.color, config.glowColor
                )}
                onClick={start}
                type="button"
                whileHover={{ scale: 1.04 }}
                whileTap={{ scale: 0.96 }}
              >
                <Play fill="white" size={16} />
                {t("pomodoro.start")}
              </motion.button>
            ) : (
              <>
                <motion.button
                  className="inline-flex items-center gap-1.5 rounded-2xl bg-rose-500 px-5 py-2.5 text-sm font-bold text-white shadow-lg shadow-rose-500/25"
                  onClick={reset}
                  type="button"
                  whileHover={{ scale: 1.03 }}
                  whileTap={{ scale: 0.96 }}
                >
                  <Pause fill="white" size={15} />
                  {t("pomodoro.pause")}
                </motion.button>
                <motion.button
                  className="inline-flex items-center gap-1.5 rounded-2xl border border-white/40 bg-white/60 px-4 py-2.5 text-sm font-semibold text-foreground shadow-sm hover:bg-white/80 dark:border-white/10 dark:bg-slate-800"
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
          <div className="flex items-center gap-2">
            {Array.from({ length: 4 }).map((_, i) => (
              <div
                className={cn(
                  "size-2 rounded-full transition-all duration-300",
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

        {/* Right: Stats + Tip stacked */}
        <div className="flex flex-col gap-3">
          {/* Stats */}
          <GlassPanel className="flex flex-col gap-3 p-4">
            <p className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              {t("pomodoro.stats.sessions")}
            </p>
            <div className="grid grid-cols-3 gap-2 lg:grid-cols-1">
              <StatRow
                color="text-blue-500"
                icon={<Target size={13} />}
                label={t("pomodoro.stats.sessions")}
                value={String(sessionCount)}
              />
              <StatRow
                color="text-emerald-500"
                icon={<TrendingUp size={13} />}
                label={t("pomodoro.stats.focus")}
                value={`${Math.floor((sessionCount * 25) / 60)}h ${(sessionCount * 25) % 60}m`}
              />
              <StatRow
                color="text-violet-500"
                icon={<Zap size={13} />}
                label={t("pomodoro.stats.streak")}
                value={sessionCount > 0 ? `🔥 ${sessionCount}` : "—"}
              />
            </div>
          </GlassPanel>

          {/* Pomodoro tip */}
          <GlassPanel className="flex flex-1 flex-col justify-between gap-2 p-4">
            <div className="flex items-center gap-2">
              <div className={cn("flex size-7 shrink-0 items-center justify-center rounded-xl bg-gradient-to-br", SESSION_CONFIG.focus.color)}>
                <Timer className="text-white" size={13} />
              </div>
              <p className="text-xs font-semibold text-foreground">{t("pomodoro.tipTitle")}</p>
            </div>
            <p className="text-xs leading-relaxed text-muted-foreground">{t("pomodoro.tipBody")}</p>
          </GlassPanel>
        </div>
      </div>
    </div>
  );
}

function StatRow({ icon, label, value, color }: { icon: React.ReactNode; label: string; value: string; color: string }) {
  return (
    <div className="flex items-center justify-between lg:flex-row">
      <div className={cn("flex items-center gap-1.5", color)}>
        {icon}
        <span className="text-xs text-muted-foreground">{label}</span>
      </div>
      <span className="font-mono text-sm font-bold tabular-nums text-foreground">{value}</span>
    </div>
  );
}
