import { ArrowLeft, Calendar, Save, Sparkles, Tag as TagIcon, Trash2 } from "lucide-react";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { GlassPanel } from "../../../components/design-system/glass-panel";
import { MOOD_CONFIGS, PRESET_TAGS } from "../journal-helpers";
import type { Journal, JournalMood } from "../types";
import { JournalLinkingSection } from "./journal-linking-section";

interface EditorProps {
  selectedId?: string;
  title: string; setTitle: (v: string) => void;
  content: string; setContent: (v: string) => void;
  tags: string; setTags: (v: string) => void;
  mood: JournalMood | null; setMood: (v: JournalMood | null) => void;
  pubDate: string; setPubDate: (v: string) => void;
  handleSave: (e: React.FormEvent) => void;
  handleDelete: () => void;
  isPending: boolean;
  availableToLink: Journal[];
  handleLink: (linkedId: string) => void;
  onBack?: () => void;
}

export function JournalEntryEditor(p: EditorProps) {
  const { t } = useTranslation("journal");
  const [confirmDel, setConfirmDel] = useState(false);

  const wordCount = p.content.trim() ? p.content.trim().split(/\s+/).length : 0;
  const charCount = p.content.length;

  const currentTagList = p.tags
    .split(",")
    .map((t) => t.trim())
    .filter(Boolean);

  const togglePresetTag = (tag: string) => {
    if (currentTagList.includes(tag)) {
      const updated = currentTagList.filter((t) => t !== tag);
      p.setTags(updated.join(", "));
    } else {
      const updated = [...currentTagList, tag];
      p.setTags(updated.join(", "));
    }
  };

  return (
    <GlassPanel className="flex flex-col h-full min-h-0 overflow-hidden p-4 sm:p-5 space-y-4 border-white/25 dark:border-white/10 bg-white/85 dark:bg-slate-950/85 shadow-xl backdrop-blur-2xl rounded-xl">
      {/* ── Header: Back button + word/char counter ── */}
      <div className="flex items-center justify-between pb-3 border-b border-slate-200/80 dark:border-white/10 shrink-0">
        {p.onBack ? (
          <button
            type="button"
            onClick={p.onBack}
            className="flex items-center gap-2 rounded-xl border border-slate-200 dark:border-white/15 bg-white/80 dark:bg-white/5 px-4 py-2 text-xs font-bold text-foreground hover:bg-slate-100 dark:hover:bg-white/10 hover:border-amber-400/50 transition-all shadow-sm hover:scale-[1.01] hover:shadow-md"
          >
            <ArrowLeft size={14} />
            <span>{t("backToList")}</span>
          </button>
        ) : (
          <div className="flex items-center gap-2">
            <Sparkles size={15} className="text-amber-500" />
            <span className="text-sm font-bold text-foreground tracking-tight">
              {t(p.selectedId ? "editTitle" : "newTitle")}
            </span>
          </div>
        )}

        <div className="flex items-center gap-3">
          <span className="text-xs font-mono text-muted-foreground/80 font-medium tabular-nums">
            {wordCount} {t("stats.words")} · {charCount} {t("chars")}
          </span>
        </div>
      </div>

      <form className="flex flex-col flex-1 min-h-0 space-y-4" onSubmit={p.handleSave}>
        {/* ── Title & Published Date Row ── */}
        <div className="grid gap-3 sm:grid-cols-[1fr_13rem] shrink-0">
          <div className="space-y-1.5">
            <label className="text-[10px] font-bold text-muted-foreground uppercase tracking-widest block">
              {t("form.title")}
            </label>
            <input
              className="h-10 w-full rounded-xl border border-slate-200 dark:border-white/15 bg-white/90 dark:bg-white/5 px-4 text-sm font-bold text-foreground focus:ring-2 focus:ring-amber-500/50 focus:border-amber-400/50 outline-none transition-all placeholder:text-muted-foreground/50 shadow-sm hover:border-slate-300 dark:hover:border-white/25"
              onChange={(e) => p.setTitle(e.target.value)}
              placeholder={t("form.title")}
              value={p.title}
            />
          </div>
          <div className="space-y-1.5">
            <label className="text-[10px] font-bold text-muted-foreground uppercase tracking-widest block">
              {t("publishedDate")}
            </label>
            <div className="relative flex items-center">
              <Calendar className="absolute left-3 size-4 text-muted-foreground pointer-events-none" />
              <input
                type="date"
                className="h-10 w-full rounded-xl border border-slate-200 dark:border-white/15 bg-white/90 dark:bg-white/5 pl-9 pr-3 text-xs font-medium text-foreground focus:ring-2 focus:ring-amber-500/50 focus:border-amber-400/50 outline-none shadow-sm transition-all hover:border-slate-300 dark:hover:border-white/25"
                onChange={(e) => p.setPubDate(e.target.value)}
                value={p.pubDate}
              />
            </div>
          </div>
        </div>

        {/* ── Mood Selector ── */}
        <div className="shrink-0 space-y-2">
          <label className="text-[10px] font-bold text-muted-foreground uppercase tracking-widest block">
            {t("mood.label")}
          </label>
          <div className="flex flex-wrap items-center gap-2">
            {MOOD_CONFIGS.map((m) => {
              const active = p.mood === m.value;
              return (
                <button
                  key={m.value}
                  type="button"
                  onClick={() => p.setMood(active ? null : m.value)}
                  className={`rounded-xl border px-3.5 py-2 text-xs font-bold transition-all duration-200 flex items-center gap-2 shadow-sm hover:scale-[1.02] ${
                    active
                      ? m.buttonActiveStyle
                      : "border-slate-200 dark:border-white/10 bg-white/70 dark:bg-white/5 text-muted-foreground hover:bg-white/95 dark:hover:bg-white/10 hover:text-foreground hover:border-slate-300 dark:hover:border-white/20"
                  }`}
                >
                  <span className={`size-2 rounded-full ${m.dotColor}`} />
                  <span>{t(m.translationKey)}</span>
                </button>
              );
            })}
          </div>
        </div>

        {/* ── Preset Quick Tags + Custom Input ── */}
        <div className="shrink-0 space-y-2">
          <div className="flex items-center justify-between">
            <label className="text-[10px] font-bold text-muted-foreground uppercase tracking-widest">
              {t("filters.tag")}
            </label>
            <span className="text-[10px] text-muted-foreground/60">{t("presetTagNotice")}</span>
          </div>

          {/* Preset chips */}
          <div className="flex flex-wrap items-center gap-1.5">
            {PRESET_TAGS.map((tag) => {
              const isSelected = currentTagList.includes(tag);
              return (
                <button
                  key={tag}
                  type="button"
                  onClick={() => togglePresetTag(tag)}
                  className={`rounded-xl border px-3 py-1.5 text-xs font-semibold transition-all hover:scale-[1.02] shadow-sm ${
                    isSelected
                      ? "border-amber-500/70 bg-amber-500/20 text-amber-800 dark:text-amber-200 font-bold"
                      : "border-slate-200 dark:border-white/10 bg-white/70 dark:bg-white/5 text-muted-foreground hover:bg-slate-100 dark:hover:bg-white/10 hover:border-slate-300 dark:hover:border-white/20"
                  }`}
                >
                  #{tag}
                </button>
              );
            })}
          </div>

          {/* Custom tag input */}
          <div className="relative flex items-center">
            <TagIcon className="absolute left-3 size-4 text-muted-foreground pointer-events-none" />
            <input
              className="h-9 w-full rounded-xl border border-slate-200 dark:border-white/15 bg-white/90 dark:bg-white/5 pl-9 pr-3 text-xs font-medium text-foreground focus:ring-2 focus:ring-amber-500/50 outline-none placeholder:text-muted-foreground/60 font-mono shadow-sm transition-all hover:border-slate-300 dark:hover:border-white/25"
              onChange={(e) => p.setTags(e.target.value)}
              placeholder={t("form.tags")}
              value={p.tags}
            />
          </div>
        </div>

        {/* ── Content Textarea (fills remaining space) ── */}
        <div className="flex flex-col flex-1 min-h-0 space-y-1.5">
          <label className="text-[10px] font-bold text-muted-foreground uppercase tracking-widest shrink-0">
            {t("form.content")}
          </label>
          <textarea
            className="flex-1 w-full min-h-0 rounded-xl border border-slate-200/80 dark:border-white/10 bg-white/90 dark:bg-slate-900/70 p-4 font-mono text-sm leading-relaxed text-foreground focus:ring-2 focus:ring-amber-500/50 outline-none resize-none overflow-y-auto placeholder:text-muted-foreground/50 shadow-inner scrollbar-thin transition-all hover:border-slate-300 dark:hover:border-white/20"
            onChange={(e) => p.setContent(e.target.value)}
            placeholder={t("form.content")}
            value={p.content}
          />
        </div>

        {/* ── Bottom Action Row ── */}
        <div className="flex items-center justify-between pt-2.5 border-t border-slate-200/80 dark:border-white/10 shrink-0 gap-3">
          <button
            className="h-10 rounded-xl bg-gradient-to-r from-amber-500 to-orange-500 px-6 text-xs font-bold text-slate-950 hover:from-amber-400 hover:to-orange-400 disabled:opacity-50 transition-all shadow-md shadow-amber-500/30 flex items-center gap-2 hover:scale-[1.01] hover:shadow-lg hover:shadow-amber-500/35"
            type="submit"
            disabled={p.isPending}
          >
            <Save size={15} />
            <span>{t("save")}</span>
          </button>

          {p.selectedId && (
            <div>
              {!confirmDel ? (
                <button
                  type="button"
                  onClick={() => setConfirmDel(true)}
                  className="h-10 rounded-xl border border-rose-400/40 dark:border-rose-500/30 bg-rose-50 dark:bg-rose-500/10 px-4 text-xs font-bold text-rose-600 dark:text-rose-400 hover:bg-rose-100 dark:hover:bg-rose-500/20 hover:border-rose-400/60 transition-all flex items-center gap-2 shadow-sm"
                >
                  <Trash2 size={14} />
                  <span>{t("delete")}</span>
                </button>
              ) : (
                <div className="flex items-center gap-2 rounded-xl border border-rose-400/40 dark:border-rose-500/30 bg-rose-50/80 dark:bg-rose-500/10 px-3.5 py-2">
                  <span className="text-xs font-semibold text-rose-600 dark:text-rose-400">{t("confirmDelete")}</span>
                  <button
                    type="button"
                    onClick={() => { setConfirmDel(false); p.handleDelete(); }}
                    className="h-7 rounded-lg bg-rose-600 px-3 text-xs font-bold text-white shadow-sm hover:bg-rose-500 transition-colors"
                  >
                    {t("confirm")}
                  </button>
                  <button
                    type="button"
                    onClick={() => setConfirmDel(false)}
                    className="h-7 rounded-lg border border-slate-300 dark:border-white/20 px-2.5 text-xs font-medium text-muted-foreground hover:bg-slate-100 dark:hover:bg-white/10 transition-colors"
                  >
                    {t("cancel")}
                  </button>
                </div>
              )}
            </div>
          )}
        </div>
      </form>

      {/* ── Linking Section ── */}
      {p.selectedId && (
        <div className="shrink-0 pt-1">
          <JournalLinkingSection availableToLink={p.availableToLink} onLink={p.handleLink} />
        </div>
      )}
    </GlassPanel>
  );
}
