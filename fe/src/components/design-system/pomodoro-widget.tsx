import { Minimize2, Pause, Play, Timer } from "lucide-react";
import { useEffect, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { GlassPanel } from "./glass-panel";
import { usePomodoro, useSetPomodoro, useStopPomodoro } from "../../features/pomodoro/hooks/use-pomodoro";
import { useUIPreferencesStore } from "../../stores/ui-preferences-store";

export function PomodoroWidget() {
  const { t } = useTranslation();
  const pomodoro = usePomodoro();
  const [now, setNow] = useState(Date.now());
  const setPomodoro = useSetPomodoro();
  const stopPomodoro = useStopPomodoro();

  const showFloating = useUIPreferencesStore((state) => state.showFloatingPomodoro);
  const toggleFloating = useUIPreferencesStore((state) => state.toggleFloatingPomodoro);

  const remaining = useMemo(() => {
    if (pomodoro.data === null || pomodoro.data === undefined) return "25:00";
    const elapsed = Math.max(0, Math.floor((now - new Date(pomodoro.data.started_at).getTime()) / 1000));
    const seconds = Math.max(0, pomodoro.data.duration_seconds - elapsed);
    return `${Math.floor(seconds / 60).toString().padStart(2, "0")}:${(seconds % 60).toString().padStart(2, "0")}`;
  }, [now, pomodoro.data]);

  useEffect(() => {
    if (pomodoro.data === null || pomodoro.data === undefined) return undefined;
    const timer = window.setInterval(() => setNow(Date.now()), 1_000);
    return () => window.clearInterval(timer);
  }, [pomodoro.data]);

  const start = () =>
    setPomodoro.mutate({
      status: "running",
      started_at: new Date().toISOString(),
      duration_seconds: 1500,
      task_id: null,
      updated_at: new Date().toISOString(),
    });

  // If minimized (default), show a tiny floating button at the bottom-right
  if (!showFloating) {
    return (
      <button
        aria-label={t("pomodoro.title")}
        className="fixed bottom-5 right-5 z-40 flex size-11 items-center justify-center rounded-full border border-white/50 bg-white/80 text-primary shadow-xl shadow-slate-950/15 backdrop-blur-2xl transition-transform hover:scale-110 active:scale-95 dark:border-white/10 dark:bg-slate-900/85"
        onClick={toggleFloating}
        title={t("pomodoro.title")}
        type="button"
      >
        <Timer size={20} />
        {pomodoro.data !== null && pomodoro.data !== undefined && (
          <span className="absolute -top-1 -right-1 flex size-3">
            <span className="absolute inline-flex size-full animate-ping rounded-full bg-emerald-400 opacity-75" />
            <span className="relative inline-flex size-3 rounded-full bg-emerald-500" />
          </span>
        )}
      </button>
    );
  }

  return (
    <GlassPanel className="fixed bottom-5 right-5 z-40 flex items-center gap-3 border-white/60 p-3.5 shadow-2xl backdrop-blur-2xl dark:border-white/15">
      <Timer aria-hidden="true" className="text-primary shrink-0" size={18} />
      <div>
        <p className="text-xs font-semibold text-foreground">{t("pomodoro.title")}</p>
        <p className="font-mono text-xs font-bold text-muted-foreground">
          {pomodoro.data === null || pomodoro.data === undefined ? t("pomodoro.idle") : remaining}
        </p>
      </div>

      <div className="flex items-center gap-1">
        {pomodoro.data === null || pomodoro.data === undefined ? (
          <button aria-label={t("pomodoro.title")} className="rounded-lg p-1.5 text-primary hover:bg-slate-200/60 dark:hover:bg-slate-800" onClick={start} type="button">
            <Play size={16} />
          </button>
        ) : (
          <button aria-label={t("pomodoro.title")} className="rounded-lg p-1.5 text-primary hover:bg-slate-200/60 dark:hover:bg-slate-800" onClick={() => stopPomodoro.mutate()} type="button">
            <Pause size={16} />
          </button>
        )}

        <button
          aria-label={t("pomodoro.minimize")}
          className="rounded-lg p-1.5 text-muted-foreground transition hover:bg-slate-200/60 hover:text-foreground dark:hover:bg-slate-800"
          onClick={toggleFloating}
          title={t("pomodoro.minimize")}
          type="button"
        >
          <Minimize2 size={15} />
        </button>
      </div>
    </GlassPanel>
  );
}
