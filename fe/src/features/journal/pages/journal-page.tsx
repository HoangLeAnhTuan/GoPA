import { AnimatePresence, motion } from "framer-motion";
import { Plus } from "lucide-react";
import { useEffect, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { PageHeader } from "../../../components/design-system/page-header";
import { ErrorState } from "../../../components/feedback/error-state";
import { LoadingState } from "../../../components/feedback/loading-state";
import { JournalEntryEditor } from "../components/journal-entry-editor";
import { JournalFiltersSidebar } from "../components/journal-filters-sidebar";
import { JournalReaderModal } from "../components/journal-reader-modal";
import { JournalStatsBanner } from "../components/journal-stats-banner";
import {
  useCreateJournal,
  useDeleteJournal,
  useJournals,
  useJournalStats,
  useLinkJournal,
  useUpdateJournal,
} from "../hooks/use-journals";
import type { Journal, JournalMood } from "../types";

const todayStr = () => new Date().toISOString().slice(0, 10);

export function JournalPage() {
  const { t } = useTranslation("journal");

  // ── View mode ──
  const [viewMode, setViewMode] = useState<"list" | "editor">("list");
  const [readerTarget, setReaderTarget] = useState<Journal | null>(null);

  // ── Search & filter state ──
  const [q, setQ] = useState("");
  const [debouncedQ, setDebouncedQ] = useState("");
  const [moodFilter, setMoodFilter] = useState<JournalMood | "">("");
  const [tagFilter, setTagFilter] = useState("");
  const [fromDate, setFromDate] = useState("");
  const [toDate, setToDate] = useState("");

  // ── Editor state ──
  const [selectedId, setSelectedId] = useState<string | undefined>();
  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [tags, setTags] = useState("");
  const [mood, setMood] = useState<JournalMood | null>(null);
  const [pubDate, setPubDate] = useState(todayStr());

  // Debounce search query (300 ms)
  useEffect(() => {
    const timer = setTimeout(() => setDebouncedQ(q), 300);
    return () => clearTimeout(timer);
  }, [q]);

  // Stable filter object for TanStack Query key
  const filter = useMemo(
    () => ({
      q: debouncedQ.trim() || undefined,
      mood: moodFilter || undefined,
      tag: tagFilter.trim() || undefined,
      from: fromDate || undefined,
      to: toDate || undefined,
    }),
    [debouncedQ, moodFilter, tagFilter, fromDate, toDate]
  );

  const journalsQ = useJournals(filter);
  const allQ = useJournals();
  const statsQ = useJournalStats();

  const create = useCreateJournal();
  const update = useUpdateJournal();
  const del = useDeleteJournal();
  const link = useLinkJournal();

  // Client-side instant filter fallback — applies both text query AND mood filter
  // so results are snappy even before a new server response arrives
  const journals = useMemo(() => {
    let result = journalsQ.data ?? [];

    // Apply mood filter client-side (as a fast fallback while server re-fetches)
    if (moodFilter) {
      result = result.filter((j) => j.mood === moodFilter);
    }

    // Apply text filter client-side
    if (debouncedQ.trim()) {
      const term = debouncedQ.trim().toLowerCase();
      result = result.filter(
        (j) =>
          j.title.toLowerCase().includes(term) ||
          j.content.toLowerCase().includes(term) ||
          (j.tags && j.tags.some((tag) => tag.toLowerCase().includes(term)))
      );
    }

    return result;
  }, [journalsQ.data, debouncedQ, moodFilter]);

  const selected = (allQ.data ?? []).find((j) => j.id === selectedId);

  // Populate editor when selected journal changes
  useEffect(() => {
    if (selected) {
      setTitle(selected.title);
      setContent(selected.content);
      setTags(selected.tags ? selected.tags.join(", ") : "");
      setMood(selected.mood ?? null);
      setPubDate(
        selected.published_date ? selected.published_date.slice(0, 10) : todayStr()
      );
    }
  }, [selected]);

  const resetForNew = () => {
    setSelectedId(undefined);
    setTitle("");
    setContent("");
    setTags("");
    setMood(null);
    setPubDate(todayStr());
    setViewMode("editor");
  };

  // FIX: this now correctly sets both q and debouncedQ for immediate effect on explicit submit
  const handleExecuteSearch = (query: string) => {
    setQ(query);
    setDebouncedQ(query);
  };

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title.trim()) return;
    const input = {
      title,
      content,
      tags: tags.split(",").map((x) => x.trim()).filter(Boolean),
      mood: mood || null,
      published_date: pubDate || undefined,
    };
    if (!selectedId) {
      const res = await create.mutateAsync(input);
      setSelectedId(res.id);
    } else {
      await update.mutateAsync({ id: selectedId, input });
    }
  };

  const handleDelete = async () => {
    if (!selectedId) return;
    await del.mutateAsync(selectedId);
    setSelectedId(undefined);
    setViewMode("list");
  };

  const handleLink = async (linkedId: string) => {
    if (!selectedId) return;
    await link.mutateAsync({ journalId: selectedId, linkedId });
  };

  const handleReaderMoodChange = async (entry: Journal, mood: JournalMood | null) => {
    // Optimistically update the reader target so the modal reflects the change immediately
    setReaderTarget((prev) => (prev?.id === entry.id ? { ...prev, mood } : prev));
    await update.mutateAsync({
      id: entry.id,
      input: {
        title: entry.title,
        content: entry.content,
        tags: entry.tags,
        mood,
        published_date: entry.published_date ? entry.published_date.slice(0, 10) : undefined,
      },
    });
  };

  const clearFilters = () => {
    setQ("");
    setDebouncedQ("");
    setMoodFilter("");
    setTagFilter("");
    setFromDate("");
    setToDate("");
  };

  if (journalsQ.isPending && !journalsQ.data) return <LoadingState />;
  if (journalsQ.isError) return <ErrorState />;

  const availableToLink = (allQ.data ?? []).filter((j) => j.id !== selectedId);

  return (
    // Root: full height, no overflow — the only scrolling happens inside child panels
    <div className="flex flex-col h-full w-full min-h-0 overflow-hidden">
      <AnimatePresence mode="wait">

        {/* ── LIST VIEW ── */}
        {viewMode === "list" && (
          <motion.div
            key="list-view"
            initial={{ opacity: 0, scale: 0.99, y: 4 }}
            animate={{ opacity: 1, scale: 1, y: 0 }}
            exit={{ opacity: 0, scale: 0.99, y: -4 }}
            transition={{ type: "spring", stiffness: 400, damping: 30 }}
            className="flex flex-col h-full min-h-0 gap-3"
          >
            {/* Page Header — fixed height */}
            <div className="shrink-0">
              <PageHeader
                title={t("title")}
                description={t("description")}
                actions={
                  <button
                    id="journal-new-entry-btn"
                    className="h-10 items-center justify-center rounded-xl bg-gradient-to-r from-amber-500 to-orange-500 px-4 text-xs font-bold text-slate-950 shadow-md shadow-amber-500/30 transition-all hover:from-amber-400 hover:to-orange-400 hover:scale-[1.02] hover:shadow-lg hover:shadow-amber-500/35 flex items-center gap-2"
                    onClick={resetForNew}
                    type="button"
                  >
                    <Plus size={16} />
                    <span>{t("newEntry")}</span>
                  </button>
                }
              />
            </div>

            {/* Stats Banner — fixed height, only shown when data ready */}
            {statsQ.data && (
              <div className="shrink-0">
                <JournalStatsBanner stats={statsQ.data} />
              </div>
            )}

            {/* Search + Filter + Card Grid — takes all remaining space, no outer scroll */}
            <div className="flex-1 min-h-0 overflow-hidden">
              <JournalFiltersSidebar
                q={q}
                onExecuteSearch={handleExecuteSearch}
                moodFilter={moodFilter} setMoodFilter={setMoodFilter}
                tagFilter={tagFilter} setTagFilter={setTagFilter}
                fromDate={fromDate} setFromDate={setFromDate}
                toDate={toDate} setToDate={setToDate}
                clearFilters={clearFilters}
                journals={journals}
                selectedId={selectedId}
                onOpenReader={(entry) => setReaderTarget(entry)}
              />
            </div>
          </motion.div>
        )}

        {/* ── EDITOR VIEW ── */}
        {viewMode === "editor" && (
          <motion.div
            key="editor-view"
            initial={{ opacity: 0, scale: 0.99, y: 6 }}
            animate={{ opacity: 1, scale: 1, y: 0 }}
            exit={{ opacity: 0, scale: 0.99, y: -6 }}
            transition={{ type: "spring", stiffness: 400, damping: 30 }}
            className="flex-1 min-h-0 h-full overflow-hidden"
          >
            <JournalEntryEditor
              selectedId={selectedId}
              title={title} setTitle={setTitle}
              content={content} setContent={setContent}
              tags={tags} setTags={setTags}
              mood={mood} setMood={setMood}
              pubDate={pubDate} setPubDate={setPubDate}
              handleSave={handleSave}
              handleDelete={handleDelete}
              isPending={create.isPending || update.isPending}
              availableToLink={availableToLink}
              handleLink={handleLink}
              onBack={() => setViewMode("list")}
            />
          </motion.div>
        )}
      </AnimatePresence>

      {/* Reader Modal — portal, always mounted */}
      <JournalReaderModal
        entry={readerTarget}
        onClose={() => setReaderTarget(null)}
        onEdit={(entry) => {
          setSelectedId(entry.id);
          setReaderTarget(null);
          setViewMode("editor");
        }}
        onMoodChange={handleReaderMoodChange}
      />
    </div>
  );
}
