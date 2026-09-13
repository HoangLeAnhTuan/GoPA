import {
  ArrowRight,
  BookOpen,
  Calendar,
  Layers,
  Plus,
  RotateCw,
  Sparkles,
} from "lucide-react";
import React, { useMemo, useState } from "react";
import { Link } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { ButtonInButton } from "../../../components/design-system/button-in-button";
import { DoubleBezelCard } from "../../../components/design-system/double-bezel-card";
import { PageHeader } from "../../../components/design-system/page-header";
import { ErrorState } from "../../../components/feedback/error-state";
import { LoadingState } from "../../../components/feedback/loading-state";
import {
  useCreateVocabulary,
  useReviewQueue,
  useVocabularyStats,
  useVocabularies,
} from "../hooks/use-vocabularies";
import type { VocabularyInput } from "../types";

export function LearnPage() {
  const { t } = useTranslation("learning");
  const allVocabsQuery = useVocabularies();
  const queueQuery = useReviewQueue();
  const statsQuery = useVocabularyStats();
  const createMutation = useCreateVocabulary();

  const [inputWord, setInputWord] = useState("");
  const [inputMeaning, setInputMeaning] = useState("");
  const [inputReading, setInputReading] = useState("");
  const [inputLanguage, setInputLanguage] = useState<"JP" | "EN">("JP");

  const allVocabs = useMemo(() => allVocabsQuery.data ?? [], [allVocabsQuery.data]);
  const dueQueue = useMemo(() => queueQuery.data ?? [], [queueQuery.data]);
  const stats = statsQuery.data;

  // Box counts calculation
  const boxCounts = useMemo(() => {
    const counts: Record<number, number> = { 1: 0, 2: 0, 3: 0, 4: 0, 5: 0 };
    for (const v of allVocabs) {
      const b = v.current_box || 1;
      counts[b] = (counts[b] || 0) + 1;
    }
    return counts;
  }, [allVocabs]);

  // Generate 90-day heatmap data
  const heatmapData = useMemo(() => {
    const days: { date: string; count: number; level: number }[] = [];
    const today = new Date();
    const serverHeatmap = stats?.review_heatmap ?? {};

    for (let i = 89; i >= 0; i--) {
      const d = new Date(today);
      d.setDate(d.getDate() - i);
      const key = d.toISOString().split("T")[0] ?? "";
      const count = (serverHeatmap as Record<string, number>)[key] ?? 0;
      let level = 0;
      if (count >= 15) level = 4;
      else if (count >= 10) level = 3;
      else if (count >= 5) level = 2;
      else if (count > 0) level = 1;
      days.push({ date: key, count, level });
    }
    return days;
  }, [stats?.review_heatmap]);

  if (allVocabsQuery.isPending || queueQuery.isPending) return <LoadingState />;
  if (allVocabsQuery.isError || queueQuery.isError) return <ErrorState />;

  const handleQuickAdd = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!inputWord.trim() || !inputMeaning.trim()) return;

    const payload: VocabularyInput = {
      language: inputLanguage,
      word: inputWord.trim(),
      reading: inputReading.trim() || null,
      meaning: inputMeaning.trim(),
      exampleSentence: null,
    };

    await createMutation.mutateAsync(payload);
    setInputWord("");
    setInputMeaning("");
    setInputReading("");
  };

  return (
    <div className="flex h-full w-full flex-col overflow-y-auto space-y-4">
      <div className="flex flex-wrap items-center justify-between gap-4">
        <PageHeader
          description={t("hub.description")}
          title={t("hub.title")}
        />

        <div className="flex items-center gap-2">
          <Link to="/app/learn/vocabulary">
            <button
              className="flex h-10 items-center gap-2 rounded-xl border border-slate-200/80 bg-white/80 px-4 text-xs font-semibold text-foreground transition hover:bg-slate-100 dark:border-white/10 dark:bg-white/5 dark:text-muted-foreground dark:hover:bg-white/10 dark:hover:text-foreground"
              type="button"
            >
              <BookOpen className="h-4 w-4" />
              {t("hub.manageVocab", { count: allVocabs.length })}
            </button>
          </Link>

          <Link to="/app/learn/review">
            <ButtonInButton icon={ArrowRight} variant="primary">
              {t("hub.reviewSession", { count: dueQueue.length })}
            </ButtonInButton>
          </Link>
        </div>
      </div>

      {/* Top Asymmetric Bento: Quick Review CTA + Leitner Mastery Stack */}
      <div className="grid gap-4 lg:grid-cols-12">
        {/* Left 7 cols: Review Status & Heatmap */}
        <div className="lg:col-span-7">
          <DoubleBezelCard className="flex flex-col justify-between p-6">
            <div className="flex items-start justify-between gap-4">
              <div>
                <span className="inline-flex items-center gap-1.5 rounded-full border border-violet-500/30 bg-violet-500/10 px-3 py-1 text-xs font-semibold text-violet-700 dark:text-violet-400">
                  <Sparkles className="h-3.5 w-3.5" />
                  {t("hub.srActive")}
                </span>
                <h2 className="mt-3 text-2xl font-bold tracking-tight text-foreground">
                  {dueQueue.length > 0 ? (
                    <span>
                      {t("hub.readyToReview", { count: dueQueue.length }).includes("<span") ? (
                        <span
                          dangerouslySetInnerHTML={{
                            __html: t("hub.readyToReview", { count: dueQueue.length }),
                          }}
                        />
                      ) : (
                        <span>
                          {t("hub.readyToReview", { count: dueQueue.length })}
                        </span>
                      )}
                    </span>
                  ) : (
                    <span>{t("hub.allCompleteToday")}</span>
                  )}
                </h2>
                <p className="mt-1 text-xs text-muted-foreground">
                  {t("hub.strengthenNeural")}
                </p>
              </div>

              <Link to="/app/learn/review">
                <ButtonInButton icon={RotateCw} variant="primary">
                  {t("hub.start")}
                </ButtonInButton>
              </Link>
            </div>

            {/* 90-Day Review Matrix (Heatmap) */}
            <div className="mt-6 border-t border-slate-200/80 dark:border-white/10 pt-4">
              <div className="flex items-center justify-between text-xs text-muted-foreground mb-2">
                <span className="flex items-center gap-1 font-semibold">
                  <Calendar className="h-3.5 w-3.5" /> {t("hub.retentionHeatmap")}
                </span>
                <span className="text-[11px]">{t("hub.totalLearned", { count: stats?.total_words ?? allVocabs.length })}</span>
              </div>

              <div className="flex flex-wrap gap-1.5">
                {heatmapData.map((d, i) => {
                  let bg = "bg-slate-100 border border-slate-200/60 dark:border-transparent dark:bg-white/5";
                  if (d.level === 1) bg = "bg-violet-200 dark:bg-violet-900/50";
                  else if (d.level === 2) bg = "bg-violet-400 dark:bg-violet-700/70";
                  else if (d.level === 3) bg = "bg-violet-500 text-white";
                  else if (d.level === 4) bg = "bg-violet-600 dark:bg-violet-400 shadow-[0_0_8px_rgba(167,139,250,0.5)]";

                  return (
                    <div
                      className={`h-3 w-3 rounded-sm transition hover:scale-125 ${bg}`}
                      key={i}
                      title={`${d.date}: ${d.count} cards reviewed`}
                    />
                  );
                })}
              </div>

              <div className="mt-2 flex items-center justify-end gap-2 text-[10px] text-muted-foreground">
                <span>{t("hub.less")}</span>
                <div className="h-2.5 w-2.5 rounded-sm bg-slate-100 border border-slate-200/60 dark:border-transparent dark:bg-white/5" />
                <div className="h-2.5 w-2.5 rounded-sm bg-violet-200 dark:bg-violet-900/50" />
                <div className="h-2.5 w-2.5 rounded-sm bg-violet-400 dark:bg-violet-700/70" />
                <div className="h-2.5 w-2.5 rounded-sm bg-violet-500" />
                <div className="h-2.5 w-2.5 rounded-sm bg-violet-600 dark:bg-violet-400" />
                <span>{t("hub.more")}</span>
              </div>
            </div>
          </DoubleBezelCard>
        </div>

        {/* Right 5 cols: Leitner Boxes Mastery Breakdown */}
        <div className="lg:col-span-5">
          <DoubleBezelCard className="flex flex-col justify-between p-6">
            <div className="flex items-center justify-between">
              <h3 className="text-sm font-bold flex items-center gap-2 text-foreground">
                <Layers className="h-4 w-4 text-violet-500 dark:text-violet-400" />
                {t("hub.boxProgression")}
              </h3>
              <span className="text-xs font-mono text-muted-foreground">{t("hub.sm2Ladder")}</span>
            </div>

            {/* Leitner Box Segments */}
            <div className="mt-4 space-y-2.5">
              {[
                { box: 1, labelKey: "boxes.box1", count: boxCounts[1] || 0, color: "bg-rose-500", barColor: "bg-rose-500/20" },
                { box: 2, labelKey: "boxes.box2", count: boxCounts[2] || 0, color: "bg-amber-500", barColor: "bg-amber-500/20" },
                { box: 3, labelKey: "boxes.box3", count: boxCounts[3] || 0, color: "bg-yellow-500", barColor: "bg-yellow-500/20" },
                { box: 4, labelKey: "boxes.box4", count: boxCounts[4] || 0, color: "bg-blue-500", barColor: "bg-blue-500/20" },
                { box: 5, labelKey: "boxes.box5", count: boxCounts[5] || 0, color: "bg-emerald-500", barColor: "bg-emerald-500/20" },
              ].map(({ box, labelKey, count, color, barColor }) => {
                const total = allVocabs.length || 1;
                const pct = Math.round((count / total) * 100);

                return (
                  <div className="space-y-1" key={box}>
                    <div className="flex items-center justify-between text-xs">
                      <span className="flex items-center gap-1.5 font-medium text-foreground">
                        <span className={`h-2 w-2 rounded-full ${color}`} />
                        {t(labelKey)}
                      </span>
                      <span className="font-mono text-muted-foreground">
                        {count} ({pct}%)
                      </span>
                    </div>
                    <div className={`h-1.5 w-full overflow-hidden rounded-full ${barColor}`}>
                      <div
                        className={`h-full rounded-full ${color}`}
                        style={{ width: `${pct}%` }}
                      />
                    </div>
                  </div>
                );
              })}
            </div>

            <div className="mt-4 border-t border-slate-200/80 dark:border-white/10 pt-3 flex items-center justify-between text-xs text-muted-foreground">
              <span>{t("hub.retentionTarget")}</span>
              <Link className="font-semibold text-violet-600 dark:text-violet-400 hover:underline" to="/app/learn/vocabulary">
                {t("hub.viewAll")}
              </Link>
            </div>
          </DoubleBezelCard>
        </div>
      </div>

      {/* Bottom Grid: Quick Add Card & Recent Words */}
      <div className="grid gap-4 lg:grid-cols-12">
        {/* Quick Add Word Form */}
        <div className="lg:col-span-5">
          <DoubleBezelCard className="p-6">
            <h3 className="text-sm font-bold flex items-center gap-2 text-foreground">
              <Plus className="h-4 w-4 text-violet-500 dark:text-violet-400" />
              {t("hub.quickAddTitle")}
            </h3>
            <p className="mt-1 text-xs text-muted-foreground">
              {t("hub.quickAddSubtitle")}
            </p>

            <form className="mt-4 space-y-3" onSubmit={handleQuickAdd}>
              <div className="grid grid-cols-2 gap-2">
                <div>
                  <label className="text-[11px] font-semibold text-muted-foreground">{t("hub.language")}</label>
                  <select
                    className="mt-1 h-9 w-full rounded-xl border border-slate-200 bg-white px-3 text-xs text-foreground outline-none transition focus:border-violet-500 focus:ring-2 focus:ring-violet-500/20 dark:border-white/10 dark:bg-slate-900/60 dark:text-foreground"
                    onChange={(e) => setInputLanguage(e.target.value as "JP" | "EN")}
                    value={inputLanguage}
                  >
                    <option value="JP">{t("hub.langJp")}</option>
                    <option value="EN">{t("hub.langEn")}</option>
                  </select>
                </div>

                <div>
                  <label className="text-[11px] font-semibold text-muted-foreground">{t("hub.wordKanji")}</label>
                  <input
                    className="mt-1 h-9 w-full rounded-xl border border-slate-200 bg-white px-3 text-xs text-foreground outline-none transition focus:border-violet-500 focus:ring-2 focus:ring-violet-500/20 dark:border-white/10 dark:bg-slate-900/60 dark:text-foreground"
                    onChange={(e) => setInputWord(e.target.value)}
                    placeholder={t("hub.wordPlaceholder")}
                    required
                    value={inputWord}
                  />
                </div>
              </div>

              <div>
                <label className="text-[11px] font-semibold text-muted-foreground">{t("hub.readingFurigana")}</label>
                <input
                  className="mt-1 h-9 w-full rounded-xl border border-slate-200 bg-white px-3 text-xs text-foreground outline-none transition focus:border-violet-500 focus:ring-2 focus:ring-violet-500/20 dark:border-white/10 dark:bg-slate-900/60 dark:text-foreground"
                  onChange={(e) => setInputReading(e.target.value)}
                  placeholder={t("hub.readingPlaceholder")}
                  value={inputReading}
                />
              </div>

              <div>
                <label className="text-[11px] font-semibold text-muted-foreground">{t("hub.meaning")}</label>
                <input
                  className="mt-1 h-9 w-full rounded-xl border border-slate-200 bg-white px-3 text-xs text-foreground outline-none transition focus:border-violet-500 focus:ring-2 focus:ring-violet-500/20 dark:border-white/10 dark:bg-slate-900/60 dark:text-foreground"
                  onChange={(e) => setInputMeaning(e.target.value)}
                  placeholder={t("hub.meaningPlaceholder")}
                  required
                  value={inputMeaning}
                />
              </div>

              <div className="pt-2">
                <ButtonInButton icon={Plus} variant="primary">
                  {t("hub.saveToBox1")}
                </ButtonInButton>
              </div>
            </form>
          </DoubleBezelCard>
        </div>

        {/* Recently Added Vocabulary Preview */}
        <div className="lg:col-span-7">
          <DoubleBezelCard className="p-6">
            <div className="flex items-center justify-between border-b border-slate-200/80 dark:border-white/10 pb-3">
              <h3 className="text-sm font-bold text-foreground">{t("hub.recentTerms")}</h3>
              <Link className="text-xs font-semibold text-violet-600 dark:text-violet-400 hover:underline" to="/app/learn/vocabulary">
                {t("hub.browseAll", { count: allVocabs.length })}
              </Link>
            </div>

            <div className="mt-3 divide-y divide-slate-100 dark:divide-white/5">
              {allVocabs.slice(0, 5).map((vocab) => (
                <div className="flex items-center justify-between py-2.5" key={vocab.id}>
                  <div>
                    <div className="flex items-center gap-2">
                      <span className="font-semibold text-foreground text-sm" lang={vocab.language === "JP" ? "ja" : undefined}>
                        {vocab.word}
                      </span>
                      {vocab.reading && (
                        <span className="font-mono text-xs text-muted-foreground">
                          ({vocab.reading})
                        </span>
                      )}
                      <span className="rounded bg-slate-100 dark:bg-white/10 px-1 py-0.2 text-[9px] font-mono text-muted-foreground">
                        {vocab.language}
                      </span>
                    </div>
                    <p className="text-xs text-muted-foreground mt-0.5">{vocab.meaning}</p>
                  </div>

                  <div className="flex items-center gap-2">
                    <span className="inline-flex h-5 w-5 items-center justify-center rounded-full bg-violet-500/15 font-mono text-[11px] font-bold text-violet-600 dark:text-violet-400">
                      B{vocab.current_box}
                    </span>
                  </div>
                </div>
              ))}

              {allVocabs.length === 0 && (
                <p className="py-8 text-center text-xs text-muted-foreground">
                  {t("hub.emptyVocab")}
                </p>
              )}
            </div>
          </DoubleBezelCard>
        </div>
      </div>
    </div>
  );
}
