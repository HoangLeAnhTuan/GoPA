import {
  ArrowLeft,
  Download,
  Link2,
  Pin,
  PinOff,
  Tag,
  Trash2,
  Zap,
} from "lucide-react";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { Link, useNavigate, useParams } from "react-router-dom";
import { DoubleBezelCard } from "../../../components/design-system/double-bezel-card";
import { ErrorState } from "../../../components/feedback/error-state";
import { LoadingState } from "../../../components/feedback/loading-state";
import {
  useDeleteJournal,
  useJournal,
  useJournals,
  useLinkJournal,
  useUpdateJournal,
} from "../hooks/use-journals";
import { getMoodConfig } from "../journal-helpers";

export function JournalDetailPage() {
  const { t } = useTranslation("journal");
  const { journalId } = useParams<{ journalId: string }>();
  const navigate = useNavigate();

  const journalQ = useJournal(journalId);
  const allJournalsQ = useJournals();
  const updateMutation = useUpdateJournal();
  const deleteMutation = useDeleteJournal();
  const linkMutation = useLinkJournal();

  const [isLinkingOpen, setIsLinkingOpen] = useState(false);
  const [selectedLinkTargetId, setSelectedLinkTargetId] = useState("");

  if (journalQ.isPending) return <LoadingState />;
  if (journalQ.isError || !journalQ.data) return <ErrorState />;

  const journal = journalQ.data;
  const allJournals = allJournalsQ.data ?? [];
  const moodCfg = getMoodConfig(journal.mood);

  const availableToLink = allJournals.filter((j) => j.id !== journal.id);

  const handleTogglePin = async () => {
    await updateMutation.mutateAsync({
      id: journal.id,
      input: {
        title: journal.title,
        content: journal.content,
        tags: journal.tags,
        mood: journal.mood,
        energy_level: journal.energy_level,
        pinned: !journal.pinned,
      },
    });
  };

  const handleExportMarkdown = () => {
    const header = `---\ntitle: "${journal.title}"\ndate: ${journal.published_date}\nmood: ${journal.mood || "none"}\nenergy_level: ${journal.energy_level || 0}\ntags: [${(journal.tags || []).join(", ")}]\n---\n\n`;
    const fullText = header + journal.content;
    const blob = new Blob([fullText], { type: "text/markdown;charset=utf-8;" });
    const url = URL.createObjectURL(blob);
    const link = document.createElement("a");
    link.href = url;
    link.download = `${journal.title.toLowerCase().replace(/[^a-z0-9]/g, "-") || "journal"}.md`;
    link.click();
    URL.revokeObjectURL(url);
  };

  const handleDelete = async () => {
    if (confirm(t("detail.confirmDelete", { title: journal.title }))) {
      await deleteMutation.mutateAsync(journal.id);
      navigate("/app/journal");
    }
  };

  const handleLinkSubmit = async () => {
    if (!selectedLinkTargetId) return;
    await linkMutation.mutateAsync({
      journalId: journal.id,
      linkedId: selectedLinkTargetId,
    });
    setIsLinkingOpen(false);
    setSelectedLinkTargetId("");
  };

  // Render markdown helper (processes headers, codeblocks, bullet lists, blockquotes)
  const renderMarkdown = (text: string) => {
    const lines = text.split("\n");
    return lines.map((line, idx) => {
      // Headers
      if (line.startsWith("### ")) {
        return (
          <h3 className="mt-4 mb-2 text-base font-bold tracking-tight text-foreground" key={idx}>
            {line.slice(4)}
          </h3>
        );
      }
      if (line.startsWith("## ")) {
        return (
          <h2 className="mt-6 mb-3 text-lg font-bold tracking-tight text-foreground border-b border-white/10 pb-1" key={idx}>
            {line.slice(3)}
          </h2>
        );
      }
      if (line.startsWith("# ")) {
        return (
          <h1 className="mt-8 mb-4 text-2xl font-extrabold tracking-tight text-foreground" key={idx}>
            {line.slice(2)}
          </h1>
        );
      }
      // Blockquote
      if (line.startsWith("> ")) {
        return (
          <blockquote className="my-2 border-l-2 border-amber-400 pl-4 italic text-muted-foreground text-sm" key={idx}>
            {line.slice(2)}
          </blockquote>
        );
      }
      // List
      if (line.startsWith("- ") || line.startsWith("* ")) {
        return (
          <li className="ml-5 list-disc text-sm text-slate-300 py-0.5" key={idx}>
            {line.slice(2)}
          </li>
        );
      }
      // Empty line
      if (line.trim() === "") {
        return <div className="h-3" key={idx} />;
      }
      // Paragraph
      return (
        <p className="text-sm leading-relaxed text-foreground/90 py-0.5" key={idx}>
          {line}
        </p>
      );
    });
  };

  const wordCount = journal.content.trim() ? journal.content.trim().split(/\s+/).length : 0;

  return (
    <div className="flex h-full w-full flex-col overflow-y-auto space-y-4">
      {/* Top Header & Navigation Bar */}
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div className="flex items-center gap-3">
          <Link
            aria-label={t("backToList")}
            className="flex h-9 w-9 items-center justify-center rounded-xl border border-slate-200/80 bg-white/80 text-foreground transition hover:bg-slate-100 dark:border-white/10 dark:bg-white/5 dark:text-muted-foreground dark:hover:bg-white/10 dark:hover:text-foreground"
            to="/app/journal"
          >
            <ArrowLeft className="h-4 w-4" />
          </Link>
          <div>
            <h1 className="text-xl font-bold tracking-tight">{journal.title}</h1>
            <p className="text-xs text-muted-foreground flex items-center gap-2">
              <span>{new Date(journal.published_date).toLocaleDateString()}</span>
              <span>·</span>
              <span>{wordCount} {t("stats.words")}</span>
              {journal.pinned && (
                <span className="flex items-center gap-1 text-amber-500 dark:text-amber-400 font-semibold">
                  <Pin className="h-3 w-3 fill-amber-500 dark:fill-amber-400" /> {t("detail.pinned")}
                </span>
              )}
            </p>
          </div>
        </div>

        <div className="flex items-center gap-2">
          <button
            aria-label={journal.pinned ? t("detail.unpin") : t("detail.pin")}
            className={`flex h-9 items-center gap-1.5 rounded-xl border px-3 text-xs font-semibold transition ${
              journal.pinned
                ? "border-amber-500/40 bg-amber-500/15 text-amber-600 dark:text-amber-300"
                : "border-slate-200 bg-white text-muted-foreground hover:bg-slate-100 hover:text-foreground dark:border-white/10 dark:bg-white/5 dark:hover:bg-white/10"
            }`}
            onClick={handleTogglePin}
            type="button"
          >
            {journal.pinned ? <PinOff className="h-3.5 w-3.5" /> : <Pin className="h-3.5 w-3.5" />}
            {journal.pinned ? t("detail.unpin") : t("detail.pin")}
          </button>

          <button
            aria-label={t("detail.exportMd")}
            className="flex h-9 items-center gap-1.5 rounded-xl border border-slate-200 bg-white px-3 text-xs font-semibold text-muted-foreground transition hover:bg-slate-100 hover:text-foreground dark:border-white/10 dark:bg-white/5 dark:hover:bg-white/10 dark:hover:text-foreground"
            onClick={handleExportMarkdown}
            type="button"
          >
            <Download className="h-3.5 w-3.5" />
            {t("detail.exportMd")}
          </button>

          <button
            aria-label={t("detail.delete")}
            className="flex h-9 items-center gap-1.5 rounded-xl border border-rose-500/20 bg-rose-500/10 px-3 text-xs font-semibold text-rose-400 transition hover:bg-rose-500/20"
            onClick={handleDelete}
            type="button"
          >
            <Trash2 className="h-3.5 w-3.5" />
            {t("detail.delete")}
          </button>
        </div>
      </div>

      {/* Metadata Badges Strip */}
      <div className="flex flex-wrap items-center gap-2 rounded-2xl border border-slate-200/80 bg-white/80 p-3 shadow-sm backdrop-blur-xl dark:border-white/10 dark:bg-white/5">
        {/* Mood Badge */}
        {moodCfg && (
          <span className={`inline-flex items-center gap-1.5 rounded-xl border px-3 py-1 text-xs font-bold ${moodCfg.badgeStyle}`}>
            <span className={`h-2 w-2 rounded-full ${moodCfg.dotColor}`} />
            <span>{t("detail.moodLabel", { mood: journal.mood ? (t(`mood.${journal.mood}`, { defaultValue: journal.mood })) : "" })}</span>
          </span>
        )}

        {/* Energy Level */}
        {journal.energy_level && (
          <span className="inline-flex items-center gap-1.5 rounded-xl border border-amber-500/30 bg-amber-500/10 px-3 py-1 text-xs font-bold text-amber-500 dark:text-amber-400">
            <Zap className="h-3.5 w-3.5 fill-amber-500 dark:fill-amber-400" />
            <span>{t("detail.energyLabel", { level: journal.energy_level })}</span>
          </span>
        )}

        {/* Tags */}
        {(journal.tags || []).map((tag) => (
          <span
            className="inline-flex items-center gap-1 rounded-xl border border-slate-200 bg-slate-100/70 dark:border-white/10 dark:bg-white/5 px-2.5 py-1 text-xs font-mono text-muted-foreground"
            key={tag}
          >
            <Tag className="h-3 w-3" />
            #{tag}
          </span>
        ))}
      </div>

      {/* Main Journal Article Content */}
      <DoubleBezelCard className="p-8 sm:p-10">
        <article className="prose dark:prose-invert max-w-none text-foreground">
          {renderMarkdown(journal.content)}
        </article>
      </DoubleBezelCard>

      {/* Connected Knowledge & Backlinks Section */}
      <DoubleBezelCard className="p-6">
        <div className="flex items-center justify-between border-b border-white/10 pb-3">
          <div className="flex items-center gap-2">
            <Link2 className="h-4 w-4 text-violet-400" />
            <h3 className="text-sm font-bold">{t("detail.backlinksTitle")}</h3>
          </div>
          <button
            className="rounded-lg border border-white/10 bg-white/5 px-2.5 py-1 text-xs font-semibold text-muted-foreground hover:bg-white/10 hover:text-foreground"
            onClick={() => setIsLinkingOpen(true)}
            type="button"
          >
            {t("detail.linkNote")}
          </button>
        </div>

        {isLinkingOpen && (
          <div className="mt-3 flex items-center gap-2">
            <select
              aria-label={t("detail.selectEntry")}
              className="h-9 flex-1 rounded-xl border border-slate-200 bg-white text-foreground dark:border-white/10 dark:bg-black/30 px-3 text-xs outline-none focus:border-violet-500"
              onChange={(e) => setSelectedLinkTargetId(e.target.value)}
              value={selectedLinkTargetId}
            >
              <option value="">{t("detail.selectEntry")}</option>
              {availableToLink.map((j) => (
                <option key={j.id} value={j.id}>
                  {j.title} ({new Date(j.published_date).toLocaleDateString()})
                </option>
              ))}
            </select>
            <button
              className="h-9 rounded-xl bg-violet-600 px-4 text-xs font-semibold text-white shadow hover:bg-violet-500"
              disabled={!selectedLinkTargetId || linkMutation.isPending}
              onClick={handleLinkSubmit}
              type="button"
            >
              {t("detail.confirm")}
            </button>
            <button
              className="h-9 rounded-xl border border-slate-200 bg-slate-100 text-foreground hover:bg-slate-200 dark:border-white/10 dark:bg-transparent dark:text-muted-foreground px-3 text-xs font-semibold dark:hover:bg-white/5"
              onClick={() => setIsLinkingOpen(false)}
              type="button"
            >
              {t("detail.cancel")}
            </button>
          </div>
        )}

        <div className="mt-3">
          <p className="text-xs text-muted-foreground">
            {t("detail.backlinksDesc")}
          </p>
        </div>
      </DoubleBezelCard>
    </div>
  );
}
