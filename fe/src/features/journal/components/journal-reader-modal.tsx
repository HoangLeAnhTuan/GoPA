import { AnimatePresence, motion } from "framer-motion";
import { BookOpen, Calendar, Edit3, FileText, Smile, Tag as TagIcon, X } from "lucide-react";
import { useEffect } from "react";
import { createPortal } from "react-dom";
import { useTranslation } from "react-i18next";
import { getMoodConfig, MOOD_CONFIGS } from "../journal-helpers";
import type { Journal, JournalMood } from "../types";

interface ReaderModalProps {
  entry: Journal | null;
  onClose: () => void;
  onEdit: (entry: Journal) => void;
  onMoodChange?: (entry: Journal, mood: JournalMood | null) => void;
}

const MOOD_EMOJI: Record<string, string> = {
  VERY_POSITIVE: "😄",
  POSITIVE: "😊",
  NEUTRAL: "😐",
  NEGATIVE: "😔",
  VERY_NEGATIVE: "😢",
};

function formatDate(dateStr: string) {
  try {
    return new Date(dateStr).toLocaleDateString(undefined, {
      weekday: "long",
      year: "numeric",
      month: "long",
      day: "numeric",
    });
  } catch {
    return dateStr.slice(0, 10);
  }
}

function wordCount(content: string) {
  return content.trim() ? content.trim().split(/\s+/).length : 0;
}

export function JournalReaderModal({ entry, onClose, onEdit, onMoodChange }: ReaderModalProps) {
  const { t } = useTranslation("journal");

  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (e.key === "Escape") onClose();
    };
    if (entry) window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, [entry, onClose]);

  if (!entry) return null;

  const moodCfg = getMoodConfig(entry.mood);
  const wc = wordCount(entry.content);
  const charCount = entry.content.length;

  return createPortal(
    <AnimatePresence>
      {entry && (
        <div className="fixed inset-0 z-[999] flex items-center justify-center p-3 sm:p-6">
          {/* Backdrop */}
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            onClick={onClose}
            className="fixed inset-0 bg-slate-950/50 backdrop-blur-md"
          />

          {/* Modal */}
          <motion.div
            initial={{ opacity: 0, scale: 0.96, y: 16 }}
            animate={{ opacity: 1, scale: 1, y: 0 }}
            exit={{ opacity: 0, scale: 0.96, y: 12 }}
            transition={{ type: "spring", stiffness: 420, damping: 30 }}
            className="relative z-10 flex flex-col w-full max-w-2xl max-h-[90vh] overflow-hidden rounded-3xl shadow-2xl shadow-slate-950/40"
          >
            {/* ── Decorative top gradient header band ── */}
            <div className="relative shrink-0 bg-gradient-to-br from-amber-50 via-orange-50 to-white dark:from-slate-800 dark:via-slate-850 dark:to-slate-900 border-b border-amber-100/80 dark:border-white/8 px-6 pt-6 pb-5">
              {/* Ambient glow blobs */}
              <div className="absolute top-0 right-0 w-48 h-32 bg-gradient-to-bl from-amber-400/20 to-transparent rounded-br-3xl pointer-events-none" />
              <div className="absolute bottom-0 left-0 w-32 h-20 bg-gradient-to-tr from-orange-400/10 to-transparent pointer-events-none" />

              {/* Close button — top right */}
              <button
                type="button"
                onClick={onClose}
                className="absolute top-4 right-4 rounded-xl p-2 text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 hover:bg-black/8 dark:hover:bg-white/10 transition-all z-10"
              >
                <X size={18} />
              </button>

              {/* Mood badge */}
              {moodCfg && (
                <div className="mb-3">
                  <span className={`inline-flex items-center gap-2 rounded-2xl border px-3 py-1.5 text-xs font-bold ${moodCfg.badgeStyle}`}>
                    <span className="text-base leading-none">{MOOD_EMOJI[entry.mood!]}</span>
                    <span>{t(moodCfg.translationKey)}</span>
                  </span>
                </div>
              )}

              {/* Title */}
              <h2 className="text-2xl font-extrabold tracking-tight text-slate-900 dark:text-white leading-tight pr-10 mb-3">
                {entry.title || <span className="italic font-normal text-slate-400">{t("form.title")}</span>}
              </h2>

              {/* Meta row */}
              <div className="flex flex-wrap items-center gap-x-4 gap-y-1.5">
                {/* Date */}
                <span className="flex items-center gap-1.5 text-sm font-medium text-slate-500 dark:text-slate-400">
                  <Calendar size={13} className="text-amber-500 shrink-0" />
                  {formatDate(entry.published_date || entry.created_at)}
                </span>
              </div>

              {/* Tags row */}
              {entry.tags && entry.tags.length > 0 && (
                <div className="flex flex-wrap items-center gap-1.5 mt-3">
                  <TagIcon size={12} className="text-amber-500 shrink-0" />
                  {entry.tags.map((tag) => (
                    <span
                      key={tag}
                      className="inline-flex rounded-xl border border-amber-400/40 bg-amber-500/12 px-2.5 py-0.5 text-xs font-semibold text-amber-700 dark:text-amber-300"
                    >
                      #{tag}
                    </span>
                  ))}
                </div>
              )}

              {/* Stats chips */}
              <div className="flex items-center gap-2 mt-4 flex-wrap">
                <span className="inline-flex items-center gap-1.5 rounded-xl border border-slate-200 dark:border-white/10 bg-white/70 dark:bg-white/5 px-3 py-1 text-xs font-semibold text-slate-600 dark:text-slate-300">
                  <BookOpen size={11} className="text-violet-500" />
                  {wc} {t("stats.words")}
                </span>
                <span className="inline-flex items-center gap-1.5 rounded-xl border border-slate-200 dark:border-white/10 bg-white/70 dark:bg-white/5 px-3 py-1 text-xs font-semibold text-slate-600 dark:text-slate-300">
                  <FileText size={11} className="text-sky-500" />
                  {charCount} {t("chars")}
                </span>
              </div>
            </div>

            {/* ── Content body ── */}
            <div className="flex-1 min-h-0 overflow-y-auto bg-white dark:bg-slate-900 scrollbar-thin">
              {entry.content ? (
                <div className="p-6">
                  {/* Content box with left accent border */}
                  <div className="relative pl-5 border-l-2 border-amber-400/60">
                    <p className="whitespace-pre-wrap font-mono text-sm leading-7 text-slate-700 dark:text-slate-300">
                      {entry.content}
                    </p>
                  </div>
                </div>
              ) : (
                <div className="flex items-center justify-center h-32 text-sm italic text-muted-foreground/60">
                  {t("emptyContent")}
                </div>
              )}
            </div>

            {/* ── Mood picker section ── */}
            {onMoodChange && (
              <div className="shrink-0 px-6 py-4 bg-white dark:bg-slate-900 border-t border-slate-100 dark:border-white/8">
                <div className="space-y-2.5">
                  <p className="text-[10px] font-bold uppercase tracking-widest text-slate-500 dark:text-slate-400 flex items-center gap-1.5">
                    <Smile size={10} />
                    {t("mood.label")}
                  </p>
                  <div className="flex flex-wrap gap-1.5">
                    {/* Clear mood button */}
                    <motion.button
                      type="button"
                      whileTap={{ scale: 0.93 }}
                      onClick={() => onMoodChange(entry, null)}
                      className={`h-8 px-3 rounded-xl text-xs font-semibold border transition-all duration-200 ${
                        !entry.mood
                          ? "bg-amber-500/20 border-amber-500/50 text-amber-700 dark:text-amber-300 shadow-sm shadow-amber-500/20"
                          : "border-slate-200 dark:border-white/10 bg-slate-100/80 dark:bg-white/5 text-slate-600 dark:text-slate-400 hover:bg-slate-200/80 dark:hover:bg-white/10 hover:text-slate-800 dark:hover:text-slate-200"
                      }`}
                    >
                      ✶ {t("mood.all")}
                    </motion.button>
                    {MOOD_CONFIGS.map((m) => {
                      const isActive = entry.mood === m.value;
                      return (
                        <motion.button
                          key={m.value}
                          type="button"
                          whileTap={{ scale: 0.9 }}
                          onClick={() => onMoodChange(entry, isActive ? null : m.value)}
                          className={`h-8 px-3 rounded-xl text-xs font-semibold border transition-all duration-200 flex items-center gap-1.5 ${
                            isActive
                              ? m.buttonActiveStyle
                              : "border-slate-200 dark:border-white/10 bg-slate-100/80 dark:bg-white/5 text-slate-600 dark:text-slate-400 hover:bg-slate-200/80 dark:hover:bg-white/10 hover:text-slate-800 dark:hover:text-slate-200"
                          }`}
                          title={t(m.translationKey)}
                        >
                          <span>{MOOD_EMOJI[m.value]}</span>
                          <span>{t(m.translationKey)}</span>
                        </motion.button>
                      );
                    })}
                  </div>
                </div>
              </div>
            )}

            {/* ── Footer actions ── */}
            <div className="shrink-0 flex items-center justify-between px-6 py-4 bg-white dark:bg-slate-900 border-t border-slate-100 dark:border-white/8">
              <p className="text-xs text-muted-foreground/60 tabular-nums">
                {formatDate(entry.updated_at || entry.created_at)}
              </p>
              <button
                type="button"
                onClick={() => { onClose(); onEdit(entry); }}
                className="flex items-center gap-2 rounded-xl bg-gradient-to-r from-amber-500 to-orange-500 px-5 py-2.5 text-sm font-bold text-slate-950 hover:from-amber-400 hover:to-orange-400 transition-all shadow-md shadow-amber-500/30 hover:scale-[1.02] hover:shadow-lg hover:shadow-amber-500/35"
              >
                <Edit3 size={14} />
                <span>{t("editTitle")}</span>
              </button>
            </div>
          </motion.div>
        </div>
      )}
    </AnimatePresence>,
    document.body
  );
}
