import { AnimatePresence, motion } from "framer-motion";
import { Calendar, RotateCcw, Search, SlidersHorizontal, Sparkles, Tag, X } from "lucide-react";
import { useEffect, useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { cn } from "../../../lib/cn";
import { getMoodConfig, MOOD_CONFIGS } from "../journal-helpers";
import type { JournalMood } from "../types";

const MOOD_EMOJI: Record<JournalMood, string> = {
  VERY_POSITIVE: "😄",
  POSITIVE: "😊",
  NEUTRAL: "😐",
  NEGATIVE: "😔",
  VERY_NEGATIVE: "😢",
};

// FIX: all classes now use proper light/dark dual variants — no hardcoded dark-only colors
const MOOD_GLOW: Record<JournalMood, string> = {
  VERY_POSITIVE: "shadow-emerald-500/40 border-emerald-500/60 bg-emerald-500/20 text-emerald-700 dark:text-emerald-300",
  POSITIVE: "shadow-teal-500/40 border-teal-500/60 bg-teal-500/20 text-teal-700 dark:text-teal-300",
  NEUTRAL: "shadow-sky-500/40 border-sky-500/60 bg-sky-500/20 text-sky-700 dark:text-sky-300",
  NEGATIVE: "shadow-amber-500/40 border-amber-500/60 bg-amber-500/20 text-amber-700 dark:text-amber-300",
  VERY_NEGATIVE: "shadow-rose-500/40 border-rose-500/60 bg-rose-500/20 text-rose-700 dark:text-rose-300",
};

export interface JournalSearchBarProps {
  q: string;
  onExecuteSearch: (query: string) => void;
  moodFilter: JournalMood | "";
  setMoodFilter: (v: JournalMood | "") => void;
  tagFilter: string;
  setTagFilter: (v: string) => void;
  fromDate: string;
  setFromDate: (v: string) => void;
  toDate: string;
  setToDate: (v: string) => void;
  clearFilters: () => void;
  resultCount?: number;
}

const PILL_COLORS = {
  amber: "bg-amber-500/15 border-amber-500/40 text-amber-700 dark:text-amber-300",
  teal: "bg-teal-500/15 border-teal-500/40 text-teal-700 dark:text-teal-300",
  violet: "bg-violet-500/15 border-violet-500/40 text-violet-700 dark:text-violet-300",
  sky: "bg-sky-500/15 border-sky-500/40 text-sky-700 dark:text-sky-300",
} as const;

function ActivePill({ label, color, onRemove }: { label: string; color: keyof typeof PILL_COLORS; onRemove: () => void }) {
  return (
    <motion.span
      initial={{ scale: 0.8, opacity: 0 }} animate={{ scale: 1, opacity: 1 }} exit={{ scale: 0.8, opacity: 0 }}
      className={cn("inline-flex items-center gap-1 rounded-lg border px-2 py-0.5 text-[11px] font-medium", PILL_COLORS[color])}
    >
      {label}
      <button type="button" onClick={onRemove} className="ml-0.5 rounded-md hover:bg-black/10 dark:hover:bg-white/10 p-0.5 transition-colors">
        <X size={10} />
      </button>
    </motion.span>
  );
}

export function JournalSearchBar({
  q, onExecuteSearch, moodFilter, setMoodFilter,
  tagFilter, setTagFilter, fromDate, setFromDate,
  toDate, setToDate, clearFilters, resultCount,
}: JournalSearchBarProps) {
  const { t } = useTranslation("journal");
  const [localQ, setLocalQ] = useState(q);
  const [isFocused, setIsFocused] = useState(false);
  const [showFilters, setShowFilters] = useState(false);
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => { setLocalQ(q); }, [q]);
  useEffect(() => {
    if (moodFilter || tagFilter || fromDate || toDate) setShowFilters(true);
  }, [moodFilter, tagFilter, fromDate, toDate]);

  const hasActiveFilters = Boolean(q || moodFilter || tagFilter || fromDate || toDate);
  const activeFilterCount = [q, moodFilter, tagFilter, fromDate || toDate].filter(Boolean).length;
  const activeMoodConfig = moodFilter ? getMoodConfig(moodFilter) : undefined;

  const handleSubmit = (e: React.FormEvent) => { e.preventDefault(); onExecuteSearch(localQ); };
  const handleClearSearch = () => { setLocalQ(""); onExecuteSearch(""); inputRef.current?.focus(); };
  const handleClearAll = () => { setLocalQ(""); clearFilters(); setShowFilters(false); };

  return (
    <div className="relative w-full space-y-2">
      {/* ── Search form ── */}
      <form onSubmit={handleSubmit} className="relative group">
        <motion.div
          animate={isFocused ? { opacity: 1, scale: 1.02 } : { opacity: 0, scale: 1 }}
          className="absolute inset-0 rounded-2xl bg-gradient-to-r from-amber-500/20 via-orange-400/15 to-amber-600/20 blur-xl pointer-events-none"
          transition={{ duration: 0.3 }}
        />
        <div className={cn(
          "relative flex items-center rounded-2xl border transition-all duration-300 overflow-hidden shadow-lg backdrop-blur-2xl",
          isFocused
            ? "border-amber-500/60 bg-white/95 dark:bg-slate-900/95 shadow-amber-500/15"
            : "border-slate-200/80 dark:border-white/12 bg-white/80 dark:bg-slate-900/80 hover:border-amber-400/50 dark:hover:border-white/22 hover:bg-white/90 dark:hover:bg-slate-900/90"
        )}>
          {/* Search icon */}
          <div className={cn("flex items-center justify-center pl-4 pr-3 shrink-0 transition-colors duration-300", isFocused ? "text-amber-500" : "text-slate-400")}>
            {isFocused ? (
              <motion.div animate={{ rotate: [0, -10, 10, 0], scale: [1, 1.1, 1] }} transition={{ duration: 0.4 }}>
                <Sparkles size={18} className="text-amber-500" />
              </motion.div>
            ) : <Search size={18} />}
          </div>
          {/* Input */}
          <input ref={inputRef} type="text"
            className="flex-1 h-12 bg-transparent text-sm font-medium text-foreground outline-none placeholder:text-slate-400 dark:placeholder:text-slate-500"
            placeholder={t("search")} value={localQ}
            onChange={(e) => setLocalQ(e.target.value)}
            onFocus={() => setIsFocused(true)} onBlur={() => setIsFocused(false)}
            onKeyDown={(e) => { if (e.key === "Escape") e.currentTarget.blur(); }}
          />
          {/* result count */}
          <AnimatePresence>
            {resultCount !== undefined && (
              <motion.span initial={{ opacity: 0 }} animate={{ opacity: 1 }} exit={{ opacity: 0 }}
                className="hidden sm:block text-[11px] font-semibold text-muted-foreground mr-2 shrink-0 tabular-nums">
                {resultCount} {t("filters.results")}
              </motion.span>
            )}
          </AnimatePresence>

          {/* Clear X */}
          <AnimatePresence>
            {localQ && (
              <motion.button initial={{ opacity: 0, scale: 0.8 }} animate={{ opacity: 1, scale: 1 }} exit={{ opacity: 0, scale: 0.8 }} transition={{ duration: 0.15 }}
                type="button" onClick={handleClearSearch}
                className="flex items-center justify-center size-7 mr-1 rounded-xl bg-slate-200/70 dark:bg-white/10 hover:bg-slate-300/70 dark:hover:bg-white/20 text-slate-500 dark:text-slate-400 hover:text-slate-700 dark:hover:text-slate-200 transition-all shrink-0">
                <X size={13} />
              </motion.button>
            )}
          </AnimatePresence>
          <div className="h-6 w-px bg-slate-200 dark:bg-white/10 mx-1 shrink-0" />
          {/* Filter toggle */}
          <motion.button type="button" whileTap={{ scale: 0.95 }} onClick={() => setShowFilters((v) => !v)}
            className={cn("relative flex items-center gap-1.5 h-full px-3 text-xs font-semibold transition-all shrink-0",
              showFilters
                ? "text-amber-600 dark:text-amber-400"
                : "text-slate-500 dark:text-slate-400 hover:text-slate-700 dark:hover:text-slate-200")}>
            <SlidersHorizontal size={15} />
            <span className="hidden sm:inline text-[11px]">{t("filters.toggle")}</span>
            {activeFilterCount > 0 && (
              <motion.span initial={{ scale: 0 }} animate={{ scale: 1 }}
                className="absolute -top-0.5 -right-0.5 size-4 rounded-full bg-amber-500 text-[9px] font-black text-slate-950 flex items-center justify-center shadow-md shadow-amber-500/40">
                {activeFilterCount}
              </motion.span>
            )}
          </motion.button>
          {/* Submit */}
          <motion.button type="submit" whileHover={{ scale: 1.03 }} whileTap={{ scale: 0.97 }}
            className="m-1.5 flex items-center gap-1.5 rounded-xl px-4 h-9 text-xs font-bold shrink-0 bg-gradient-to-r from-amber-500 to-orange-500 text-slate-950 shadow-md shadow-amber-500/30 hover:from-amber-400 hover:to-orange-400 transition-all">
            <Search size={13} />
            <span className="hidden sm:inline">{t("filters.search")}</span>
          </motion.button>
        </div>
      </form>

      {/* ── Active filter pills ── */}
      <AnimatePresence>
        {hasActiveFilters && (
          <motion.div
            initial={{ opacity: 0, height: 0 }}
            animate={{ opacity: 1, height: "auto" }}
            exit={{ opacity: 0, height: 0 }}
            className="flex flex-wrap items-center gap-1.5 overflow-hidden pl-1"
          >
            {q && (
              <ActivePill label={`"${q}"`} color="amber" onRemove={() => { setLocalQ(""); onExecuteSearch(""); }} />
            )}
            {moodFilter && activeMoodConfig && (
              <ActivePill label={t(activeMoodConfig.translationKey)} color="teal" onRemove={() => setMoodFilter("")} />
            )}
            {tagFilter && (
              <ActivePill label={`#${tagFilter}`} color="violet" onRemove={() => setTagFilter("")} />
            )}
            {(fromDate || toDate) && (
              <ActivePill label={`${fromDate || "…"} → ${toDate || "…"}`} color="sky"
                onRemove={() => { setFromDate(""); setToDate(""); }} />
            )}
            <motion.button
              type="button" whileTap={{ scale: 0.95 }} onClick={handleClearAll}
              className="inline-flex items-center gap-1 rounded-lg border border-rose-300/60 dark:border-rose-500/30 bg-rose-50 dark:bg-rose-500/10 px-2 py-0.5 text-[11px] font-medium text-rose-600 dark:text-rose-400 hover:bg-rose-100 dark:hover:bg-rose-500/20 transition-colors"
            >
              <RotateCcw size={9} />
              {t("filters.clear")}
            </motion.button>
          </motion.div>
        )}
      </AnimatePresence>

      {/* ── Expandable Filters Panel ── */}
      <AnimatePresence>
        {showFilters && (
          <motion.div key="filters"
            initial={{ opacity: 0, y: -8, height: 0 }}
            animate={{ opacity: 1, y: 0, height: "auto" }}
            exit={{ opacity: 0, y: -8, height: 0 }}
            transition={{ type: "spring", stiffness: 380, damping: 30 }}
            className="overflow-hidden">
            {/* FIX: use proper light/dark surface — was hardcoded dark bg-slate-900/90 */}
            <div className="relative rounded-2xl border border-slate-200/80 dark:border-white/10 bg-white/90 dark:bg-slate-900/90 backdrop-blur-2xl p-4 space-y-3.5 shadow-xl shadow-black/10 dark:shadow-black/30 overflow-hidden">
              <div className="absolute inset-0 bg-gradient-to-br from-amber-500/5 via-transparent to-orange-500/5 pointer-events-none rounded-2xl" />

              {/* Mood pills */}
              <div className="space-y-2">
                <p className="text-[10px] font-bold uppercase tracking-widest text-amber-600 dark:text-amber-500/70 flex items-center gap-1.5">
                  <Sparkles size={10} />{t("filters.mood")}
                </p>
                <div className="flex flex-wrap gap-1.5">
                  <motion.button whileTap={{ scale: 0.93 }} type="button" onClick={() => setMoodFilter("")}
                    className={cn("h-8 px-3 rounded-xl text-xs font-semibold border transition-all duration-200",
                      !moodFilter
                        ? "bg-amber-500/20 border-amber-500/50 text-amber-700 dark:text-amber-300 shadow-sm shadow-amber-500/20"
                        : "border-slate-200 dark:border-white/10 bg-slate-100/80 dark:bg-white/5 text-slate-600 dark:text-slate-400 hover:bg-slate-200/80 dark:hover:bg-white/10 hover:text-slate-800 dark:hover:text-slate-200")}>
                    ✶ {t("filters.allMoods")}
                  </motion.button>
                  {MOOD_CONFIGS.map((mood) => {
                    const isActive = moodFilter === mood.value;
                    return (
                      <motion.button key={mood.value} whileTap={{ scale: 0.9 }} type="button"
                        onClick={() => setMoodFilter(isActive ? "" : mood.value)}
                        className={cn("h-8 px-3 rounded-xl text-xs font-semibold border transition-all duration-200 flex items-center gap-1.5",
                          isActive
                            ? cn("shadow-md", MOOD_GLOW[mood.value])
                            : "border-slate-200 dark:border-white/10 bg-slate-100/80 dark:bg-white/5 text-slate-600 dark:text-slate-400 hover:bg-slate-200/80 dark:hover:bg-white/10 hover:text-slate-800 dark:hover:text-slate-200")}
                        title={t(mood.translationKey)}>
                        <span>{MOOD_EMOJI[mood.value]}</span>
                        <span className="hidden sm:inline">{t(mood.translationKey)}</span>
                      </motion.button>
                    );
                  })}
                </div>
              </div>

              {/* Tag + Date row */}
              <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
                {/* Tag filter */}
                <div className="space-y-1.5">
                  <label className="text-[10px] font-bold uppercase tracking-widest text-slate-500 dark:text-slate-400 flex items-center gap-1.5">
                    <Tag size={9} />{t("filters.tag")}
                  </label>
                  <div className="relative">
                    <Tag className="absolute left-3 top-1/2 -translate-y-1/2 size-3.5 text-muted-foreground pointer-events-none" />
                    <input
                      className="h-9 w-full rounded-xl border border-slate-200 dark:border-white/15 bg-white/90 dark:bg-white/5 pl-8 pr-3 text-xs font-medium text-foreground focus:ring-2 focus:ring-amber-500/50 outline-none placeholder:text-muted-foreground/60 transition-all"
                      placeholder={t("filters.tag")}
                      value={tagFilter}
                      onChange={(e) => setTagFilter(e.target.value)}
                    />
                  </div>
                </div>
                {/* From date */}
                <div className="space-y-1.5">
                  <label className="text-[10px] font-bold uppercase tracking-widest text-slate-500 dark:text-slate-400 flex items-center gap-1.5">
                    <Calendar size={9} />{t("filters.from")}
                  </label>
                  <div className="relative">
                    <Calendar className="absolute left-3 top-1/2 -translate-y-1/2 size-3.5 text-muted-foreground pointer-events-none" />
                    <input
                      type="date"
                      className="h-9 w-full rounded-xl border border-slate-200 dark:border-white/15 bg-white/90 dark:bg-white/5 pl-8 pr-3 text-xs font-medium text-foreground focus:ring-2 focus:ring-amber-500/50 outline-none transition-all"
                      value={fromDate}
                      onChange={(e) => setFromDate(e.target.value)}
                    />
                  </div>
                </div>
                {/* To date */}
                <div className="space-y-1.5">
                  <label className="text-[10px] font-bold uppercase tracking-widest text-slate-500 dark:text-slate-400 flex items-center gap-1.5">
                    <Calendar size={9} />{t("filters.to")}
                  </label>
                  <div className="relative">
                    <Calendar className="absolute left-3 top-1/2 -translate-y-1/2 size-3.5 text-muted-foreground pointer-events-none" />
                    <input
                      type="date"
                      className="h-9 w-full rounded-xl border border-slate-200 dark:border-white/15 bg-white/90 dark:bg-white/5 pl-8 pr-3 text-xs font-medium text-foreground focus:ring-2 focus:ring-amber-500/50 outline-none transition-all"
                      value={toDate}
                      onChange={(e) => setToDate(e.target.value)}
                    />
                  </div>
                </div>
              </div>
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
