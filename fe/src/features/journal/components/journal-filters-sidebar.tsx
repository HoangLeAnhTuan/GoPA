import { AnimatePresence, motion } from "framer-motion";
import { BookOpen, Calendar, NotebookPen, Tag } from "lucide-react";
import { useTranslation } from "react-i18next";
import { getMoodConfig } from "../journal-helpers";
import type { Journal, JournalMood } from "../types";
import { JournalSearchBar } from "./journal-search-bar";

const MOOD_EMOJI: Record<JournalMood, string> = {
  VERY_POSITIVE: "😄",
  POSITIVE: "😊",
  NEUTRAL: "😐",
  NEGATIVE: "😔",
  VERY_NEGATIVE: "😢",
};

interface SidebarProps {
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
  journals: Journal[];
  selectedId?: string;
  onOpenReader: (entry: Journal) => void;
}

/** Format a date string as a short readable date */
function formatDate(dateStr: string) {
  try {
    return new Date(dateStr).toLocaleDateString(undefined, {
      month: "short",
      day: "numeric",
      year: "numeric",
    });
  } catch {
    return dateStr.slice(0, 10);
  }
}

/** Journal entry card for the grid */
function JournalCard({
  entry,
  isSelected,
  onClick,
}: {
  entry: Journal;
  isSelected: boolean;
  onClick: () => void;
}) {
  const { t } = useTranslation("journal");
  const moodCfg = getMoodConfig(entry.mood);
  const wordCount = entry.content.trim() ? entry.content.trim().split(/\s+/).length : 0;
  const preview = entry.content.replace(/[#*`_[\]>]/g, "").trim().slice(0, 120);

  return (
    <motion.button
      type="button"
      onClick={onClick}
      layout
      initial={{ opacity: 0, y: 8, scale: 0.98 }}
      animate={{ opacity: 1, y: 0, scale: 1 }}
      exit={{ opacity: 0, y: -4, scale: 0.97 }}
      whileHover={{ y: -2, scale: 1.005 }}
      whileTap={{ scale: 0.99 }}
      transition={{ type: "spring", stiffness: 400, damping: 28 }}
      className={`group relative w-full text-left rounded-2xl border p-4 transition-all duration-300 overflow-hidden backdrop-blur-xl
        ${isSelected
          ? "border-amber-500/50 bg-white/90 dark:bg-slate-800/80 shadow-lg shadow-amber-500/12"
          : "border-slate-200/70 dark:border-white/10 bg-white/70 dark:bg-slate-900/60 hover:border-amber-400/50 dark:hover:border-amber-500/30 hover:bg-white/90 dark:hover:bg-slate-900/80 hover:shadow-md hover:shadow-amber-500/10"
        }`}
    >
      {/* Ambient glow on hover */}
      <div className="absolute inset-0 opacity-0 group-hover:opacity-100 transition-opacity duration-500 bg-gradient-to-br from-amber-500/5 via-transparent to-transparent pointer-events-none rounded-2xl" />

      {/* Header: title + mood badge */}
      <div className="flex items-start justify-between gap-3 mb-2">
        <h3 className="font-bold text-sm text-foreground leading-snug line-clamp-1 flex-1">
          {entry.title || <span className="text-muted-foreground italic font-normal">{t("form.title")}</span>}
        </h3>
        {moodCfg && (
          <span className={`inline-flex items-center gap-1 rounded-lg border px-2 py-0.5 text-[10px] font-semibold shrink-0 ${moodCfg.badgeStyle}`}>
            <span>{MOOD_EMOJI[entry.mood!]}</span>
            <span className="hidden sm:inline">{t(moodCfg.translationKey)}</span>
          </span>
        )}
      </div>

      {/* Content preview */}
      {preview && (
        <p className="text-xs text-muted-foreground/80 line-clamp-2 leading-relaxed mb-3 font-mono">
          {preview}{entry.content.length > 120 ? "…" : ""}
        </p>
      )}

      {/* Footer: date + tags + word count */}
      <div className="flex items-center justify-between gap-2 flex-wrap">
        <div className="flex items-center gap-2.5 flex-wrap">
          {/* Date */}
          <span className="flex items-center gap-1 text-[10px] font-medium text-muted-foreground/70">
            <Calendar size={10} className="text-amber-500/70" />
            {formatDate(entry.published_date || entry.created_at)}
          </span>

          {/* Tags */}
          {entry.tags && entry.tags.length > 0 && (
            <span className="flex items-center gap-1 text-[10px] font-mono text-amber-600 dark:text-amber-400 font-semibold">
              <Tag size={9} />
              {entry.tags.slice(0, 2).map(t => `#${t}`).join(" ")}
              {entry.tags.length > 2 && <span className="text-muted-foreground/60">+{entry.tags.length - 2}</span>}
            </span>
          )}
        </div>

        {/* Word count */}
        <span className="flex items-center gap-1 text-[10px] font-medium text-muted-foreground/60 ml-auto shrink-0">
          <BookOpen size={9} />
          {wordCount} {t("stats.words")}
        </span>
      </div>

      {/* Selected indicator bar */}
      <AnimatePresence>
        {isSelected && (
          <motion.div
            initial={{ scaleY: 0, opacity: 0 }}
            animate={{ scaleY: 1, opacity: 1 }}
            exit={{ scaleY: 0, opacity: 0 }}
            className="absolute left-0 top-3 bottom-3 w-0.5 rounded-full bg-gradient-to-b from-amber-400 to-orange-500"
          />
        )}
      </AnimatePresence>
    </motion.button>
  );
}

/** Full-featured filters sidebar: search bar + optional filter panel + scrollable journal card grid */
export function JournalFiltersSidebar({
  q,
  onExecuteSearch,
  moodFilter, setMoodFilter,
  tagFilter, setTagFilter,
  fromDate, setFromDate,
  toDate, setToDate,
  clearFilters,
  journals,
  selectedId,
  onOpenReader,
}: SidebarProps) {
  const { t } = useTranslation("journal");

  return (
    <div className="flex flex-col h-full min-h-0 gap-3">
      {/* Search + collapsible filters — shrinks to fit content */}
      <div className="shrink-0">
        <JournalSearchBar
          q={q}
          onExecuteSearch={onExecuteSearch}
          moodFilter={moodFilter}
          setMoodFilter={setMoodFilter}
          tagFilter={tagFilter}
          setTagFilter={setTagFilter}
          fromDate={fromDate}
          setFromDate={setFromDate}
          toDate={toDate}
          setToDate={setToDate}
          clearFilters={clearFilters}
          resultCount={journals.length}
        />
      </div>

      {/* Journal card grid — flex-1 with internal scroll only */}
      <div className="flex-1 min-h-0 overflow-y-auto pr-1 scrollbar-thin">
        <AnimatePresence mode="popLayout">
          {journals.length === 0 ? (
            <motion.div
              key="empty"
              initial={{ opacity: 0, y: 8 }}
              animate={{ opacity: 1, y: 0 }}
              exit={{ opacity: 0 }}
              className="flex flex-col items-center justify-center h-48 gap-3 text-center"
            >
              <div className="size-14 rounded-2xl border border-dashed border-amber-400/40 bg-amber-500/5 flex items-center justify-center">
                <NotebookPen size={24} className="text-amber-500/50" />
              </div>
              <div>
                <p className="text-sm font-semibold text-muted-foreground">{t("noEntries")}</p>
                <p className="text-xs text-muted-foreground/60 mt-0.5">{t("empty")}</p>
              </div>
            </motion.div>
          ) : (
            <motion.div
              key="grid"
              className="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-3 2xl:grid-cols-4 gap-3 pb-2"
            >
              <AnimatePresence mode="popLayout">
                {journals.map((entry, i) => (
                  <motion.div key={entry.id} layout
                    initial={{ opacity: 0, y: 10 }}
                    animate={{ opacity: 1, y: 0, transition: { delay: Math.min(i * 0.03, 0.18) } }}
                    exit={{ opacity: 0, scale: 0.96 }}
                  >
                    <JournalCard
                      entry={entry}
                      isSelected={entry.id === selectedId}
                      onClick={() => onOpenReader(entry)}
                    />
                  </motion.div>
                ))}
              </AnimatePresence>
            </motion.div>
          )}
        </AnimatePresence>
      </div>
    </div>
  );
}
