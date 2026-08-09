import { BookOpen, Flame, TrendingUp } from "lucide-react";
import { useTranslation } from "react-i18next";
import { GlassPanel } from "../../../components/design-system/glass-panel";
import { MOOD_CONFIGS } from "../journal-helpers";
import type { JournalStats } from "../types";

interface JournalStatsBannerProps {
  stats: JournalStats;
}

const MOOD_EMOJI: Record<string, string> = {
  VERY_POSITIVE: "😄",
  POSITIVE: "😊",
  NEUTRAL: "😐",
  NEGATIVE: "😔",
  VERY_NEGATIVE: "😢",
};

export function JournalStatsBanner({ stats }: JournalStatsBannerProps) {
  const { t } = useTranslation("journal");

  return (
    <GlassPanel className="px-3 py-2.5 border-amber-500/20 bg-gradient-to-r from-amber-500/10 via-amber-500/5 to-orange-500/8 shadow-sm backdrop-blur-xl shrink-0 rounded-xl">
      <div className="flex flex-col gap-2.5 sm:flex-row sm:items-center sm:justify-between">

        {/* KPI Metrics Row */}
        <div className="flex flex-wrap items-center gap-2 text-xs">
          {/* Current Streak */}
          <div className="flex items-center gap-2 rounded-xl border border-amber-500/30 bg-amber-500/15 px-3 py-1.5 shadow-sm">
            <Flame size={13} className="text-amber-500 shrink-0" />
            <div className="flex items-center gap-1.5">
              <span className="text-[10px] font-semibold text-amber-700/80 dark:text-amber-300/80 uppercase tracking-wide">{t("stats.currentStreak")}:</span>
              <span className="font-bold tabular-nums text-amber-700 dark:text-amber-300">
                {stats.current_streak} <span className="font-normal opacity-70">{t("stats.days")}</span>
              </span>
            </div>
            <span className="size-1.5 rounded-full bg-amber-500 animate-pulse" />
          </div>

          {/* Longest Streak */}
          <div className="flex items-center gap-2 rounded-xl border border-slate-200/80 dark:border-white/10 bg-white/70 dark:bg-slate-900/60 px-3 py-1.5 shadow-sm">
            <TrendingUp size={13} className="text-blue-500 shrink-0" />
            <div className="flex items-center gap-1.5">
              <span className="text-[10px] font-semibold text-muted-foreground uppercase tracking-wide">{t("stats.longestStreak")}:</span>
              <span className="font-bold text-blue-600 dark:text-blue-400 tabular-nums">
                {stats.longest_streak} <span className="font-normal opacity-70">{t("stats.days")}</span>
              </span>
            </div>
          </div>

          {/* Total Entries */}
          <div className="flex items-center gap-2 rounded-xl border border-slate-200/80 dark:border-white/10 bg-white/70 dark:bg-slate-900/60 px-3 py-1.5 shadow-sm">
            <BookOpen size={13} className="text-violet-500 shrink-0" />
            <div className="flex items-center gap-1.5">
              <span className="text-[10px] font-semibold text-muted-foreground uppercase tracking-wide">{t("stats.entries")}:</span>
              <span className="font-bold text-violet-600 dark:text-violet-400 tabular-nums">{stats.entries}</span>
            </div>
          </div>

          {/* Total Words */}
          <div className="flex items-center gap-2 rounded-xl border border-slate-200/80 dark:border-white/10 bg-white/70 dark:bg-slate-900/60 px-3 py-1.5 shadow-sm">
            <div className="flex items-center gap-1.5">
              <span className="text-[10px] font-semibold text-muted-foreground uppercase tracking-wide">{t("stats.wordCount")}:</span>
              <span className="font-bold text-emerald-600 dark:text-emerald-400 tabular-nums">{stats.word_count.toLocaleString()}</span>
            </div>
          </div>
        </div>

        {/* Mood Distribution Chips */}
        <div className="flex items-center gap-1.5 overflow-x-auto scrollbar-none pb-0.5">
          <span className="shrink-0 mr-0.5 text-[9px] font-bold text-amber-700 dark:text-amber-300 uppercase tracking-widest">
            {t("stats.moodBreakdown")}:
          </span>
          {MOOD_CONFIGS.map((m) => {
            const count = stats.mood_counts?.[m.value] ?? 0;
            return (
              <span
                key={m.value}
                // FIX: badgeStyle from MOOD_CONFIGS now has proper light/dark text classes
                className={`inline-flex items-center gap-1.5 rounded-xl border px-2.5 py-1 text-[10px] font-semibold transition-all shrink-0 ${m.badgeStyle}`}
                title={t(m.translationKey)}
              >
                <span>{MOOD_EMOJI[m.value]}</span>
                <span className="hidden sm:inline font-medium">{t(m.translationKey)}</span>
                <span className="font-bold border-l border-current/25 pl-1.5 tabular-nums">{count}</span>
              </span>
            );
          })}
        </div>

      </div>
    </GlassPanel>
  );
}
