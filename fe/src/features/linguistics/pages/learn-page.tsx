import { Plus } from "lucide-react";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { ErrorState } from "../../../components/feedback/error-state";
import { LoadingState } from "../../../components/feedback/loading-state";
import { GlassPanel } from "../../../components/design-system/glass-panel";
import { PageHeader } from "../../../components/design-system/page-header";
import { MacSelect } from "../../../components/design-system/mac-select";
import { useCreateVocabulary, useReviewQueue, useReviewVocabulary, useVocabularies } from "../hooks/use-vocabularies";
import type { VocabularyInput } from "../types";

const empty: VocabularyInput = { language: "JP", word: "", reading: null, meaning: "", exampleSentence: null };

export function LearnPage() {
  const { t } = useTranslation("learning");
  const all = useVocabularies();
  const due = useReviewQueue();
  const create = useCreateVocabulary();
  const review = useReviewVocabulary();
  const [input, setInput] = useState(empty);

  if (all.isPending || due.isPending) return <LoadingState />;
  if (all.isError || due.isError) return <ErrorState />;

  const next = due.data[0];

  const submit = async (event: React.FormEvent) => {
    event.preventDefault();
    if (!input.word.trim() || !input.meaning.trim()) return;
    await create.mutateAsync(input);
    setInput(empty);
  };

  const languageOptions = [
    { value: "JP", label: "Japanese (日本語)" },
    { value: "EN", label: "English" },
  ];

  return (
    <div>
      <PageHeader description={t("description")} title={t("title")} />
      <section className="mt-6 grid gap-5 lg:grid-cols-[1.1fr_1.9fr]">
        <GlassPanel className="p-5">
          <h2 className="font-semibold">{t("review.title")}</h2>
          {next === undefined ? (
            <p className="mt-3 text-sm text-muted-foreground">{t("review.empty")}</p>
          ) : (
            <div className="mt-5">
              <p className="text-3xl font-semibold" lang={next.language === "JP" ? "ja" : undefined}>
                {next.word}
              </p>
              {next.reading === null ? null : <p className="mt-1 text-muted-foreground">{next.reading}</p>}
              <p className="mt-5">{next.meaning}</p>
              <div className="mt-5 grid grid-cols-3 gap-2">
                {[1, 3, 5].map((quality) => (
                  <button
                    className="h-10 rounded-xl border border-white/40 bg-white/60 text-sm font-medium transition hover:bg-white/80 dark:border-white/10 dark:bg-slate-900/60 dark:hover:bg-slate-800"
                    key={quality}
                    onClick={() => review.mutate({ id: next.id, quality })}
                    type="button"
                  >
                    {t(`review.quality.${quality}`)}
                  </button>
                ))}
              </div>
            </div>
          )}
        </GlassPanel>

        <GlassPanel className="p-5">
          <h2 className="font-semibold">{t("form.title")}</h2>
          <form className="mt-4 grid gap-3 sm:grid-cols-2" onSubmit={submit}>
            <input
              className="h-11 rounded-xl border border-white/40 bg-white/70 px-3.5 text-sm outline-none transition focus:border-primary dark:border-white/10 dark:bg-slate-900/70"
              onChange={(e) => setInput({ ...input, word: e.target.value })}
              placeholder={t("form.word")}
              value={input.word}
            />
            <input
              className="h-11 rounded-xl border border-white/40 bg-white/70 px-3.5 text-sm outline-none transition focus:border-primary dark:border-white/10 dark:bg-slate-900/70"
              onChange={(e) => setInput({ ...input, meaning: e.target.value })}
              placeholder={t("form.meaning")}
              value={input.meaning}
            />
            <input
              className="h-11 rounded-xl border border-white/40 bg-white/70 px-3.5 text-sm outline-none transition focus:border-primary dark:border-white/10 dark:bg-slate-900/70"
              onChange={(e) => setInput({ ...input, reading: e.target.value || null })}
              placeholder={t("form.reading")}
              value={input.reading ?? ""}
            />
            <MacSelect
              onChange={(e) => setInput({ ...input, language: e as VocabularyInput["language"] })}
              options={languageOptions}
              value={input.language}
            />
            <button
              className="inline-flex h-11 items-center justify-center gap-2 rounded-xl bg-primary px-4 font-semibold text-primary-foreground shadow-lg shadow-primary/25 transition hover:opacity-90 sm:col-span-2"
              type="submit"
            >
              <Plus size={17} />
              {t("form.add")}
            </button>
          </form>
          <div className="mt-6 divide-y divide-white/20 dark:divide-white/10">
            {all.data.map((v) => (
              <div className="py-3" key={v.id}>
                <p className="font-medium">
                  {v.word} <span className="text-sm text-muted-foreground">{v.reading}</span>
                </p>
                <p className="text-sm text-muted-foreground">
                  {v.meaning} · {t("box", { count: v.current_box })}
                </p>
              </div>
            ))}
          </div>
        </GlassPanel>
      </section>
    </div>
  );
}

