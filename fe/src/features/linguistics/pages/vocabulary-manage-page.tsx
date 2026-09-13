import {
  ArrowLeft,
  Check,
  Edit2,
  FileSpreadsheet,
  Filter,
  Plus,
  Search,
  Trash2,
  Upload,
  Volume2,
  X,
} from "lucide-react";
import React, { useMemo, useRef, useState } from "react";
import { Link } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { ButtonInButton } from "../../../components/design-system/button-in-button";
import { DoubleBezelCard } from "../../../components/design-system/double-bezel-card";
import { ErrorState } from "../../../components/feedback/error-state";
import { LoadingState } from "../../../components/feedback/loading-state";
import {
  useCreateVocabulary,
  useDeleteVocabulary,
  useGenerateAudio,
  useImportVocabularies,
  useUpdateVocabulary,
  useVocabularies,
} from "../hooks/use-vocabularies";
import type { Vocabulary, VocabularyInput, VocabularyLanguage } from "../types";

export function VocabularyManagePage() {
  const { t } = useTranslation("learning");
  const vocabsQuery = useVocabularies();
  const createMutation = useCreateVocabulary();
  const updateMutation = useUpdateVocabulary();
  const deleteMutation = useDeleteVocabulary();
  const audioMutation = useGenerateAudio();
  const importMutation = useImportVocabularies();

  const [searchQuery, setSearchQuery] = useState("");
  const [selectedBox, setSelectedBox] = useState<number | "all">("all");
  const [selectedLanguage, setSelectedLanguage] = useState<VocabularyLanguage | "all">("all");
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [isImportOpen, setIsImportOpen] = useState(false);
  const [editingVocab, setEditingVocab] = useState<Vocabulary | null>(null);
  const [importText, setImportText] = useState("");
  const [importFormat, setImportFormat] = useState<"csv" | "json">("csv");
  const [activeAudioId, setActiveAudioId] = useState<string | null>(null);

  // Form State
  const [formWord, setFormWord] = useState("");
  const [formReading, setFormReading] = useState("");
  const [formMeaning, setFormMeaning] = useState("");
  const [formSentence, setFormSentence] = useState("");
  const [formTranslation, setFormTranslation] = useState("");
  const [formLanguage, setFormLanguage] = useState<VocabularyLanguage>("JP");
  const [formTags, setFormTags] = useState("");

  const audioPlayerRef = useRef<HTMLAudioElement | null>(null);

  const openCreateModal = () => {
    setEditingVocab(null);
    setFormWord("");
    setFormReading("");
    setFormMeaning("");
    setFormSentence("");
    setFormTranslation("");
    setFormLanguage("JP");
    setFormTags("");
    setIsModalOpen(true);
  };

  const openEditModal = (vocab: Vocabulary) => {
    setEditingVocab(vocab);
    setFormWord(vocab.word);
    setFormReading(vocab.reading ?? "");
    setFormMeaning(vocab.meaning);
    setFormSentence(vocab.example_sentence ?? "");
    setFormTranslation(vocab.example_translation ?? "");
    setFormLanguage(vocab.language);
    setFormTags((vocab.tags ?? []).join(", "));
    setIsModalOpen(true);
  };

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!formWord.trim() || !formMeaning.trim()) return;

    const tagsArray = formTags
      .split(",")
      .map((tag) => tag.trim())
      .filter(Boolean);

    const payload: VocabularyInput = {
      language: formLanguage,
      word: formWord.trim(),
      reading: formReading.trim() || null,
      meaning: formMeaning.trim(),
      exampleSentence: formSentence.trim() || null,
      exampleTranslation: formTranslation.trim() || null,
      tags: tagsArray,
    };

    if (editingVocab) {
      await updateMutation.mutateAsync({
        id: editingVocab.id,
        input: payload,
      });
    } else {
      await createMutation.mutateAsync(payload);
    }
    setIsModalOpen(false);
  };

  const handlePlayAudio = async (vocab: Vocabulary) => {
    setActiveAudioId(vocab.id);
    try {
      if (vocab.audio_url) {
        if (!audioPlayerRef.current) audioPlayerRef.current = new Audio(vocab.audio_url);
        else audioPlayerRef.current.src = vocab.audio_url;
        await audioPlayerRef.current.play();
      } else {
        const res = await audioMutation.mutateAsync(vocab.id).catch(() => null);
        if (res?.audio_url) {
          if (!audioPlayerRef.current) audioPlayerRef.current = new Audio(res.audio_url);
          else audioPlayerRef.current.src = res.audio_url;
          await audioPlayerRef.current.play();
        } else if ("speechSynthesis" in window) {
          const utterance = new SpeechSynthesisUtterance(vocab.word);
          utterance.lang = vocab.language === "JP" ? "ja-JP" : "en-US";
          window.speechSynthesis.speak(utterance);
        }
      }
    } catch {
      // Ignored
    } finally {
      setTimeout(() => setActiveAudioId(null), 1000);
    }
  };

  const handleImportSubmit = async () => {
    if (!importText.trim()) return;
    try {
      await importMutation.mutateAsync({
        data: importText,
        format: importFormat,
      });
      setIsImportOpen(false);
      setImportText("");
    } catch {
      // Handle error
    }
  };

  const filteredVocabs = useMemo(() => {
    const list = vocabsQuery.data ?? [];
    return list.filter((v) => {
      const matchQuery =
        searchQuery.trim() === "" ||
        v.word.toLowerCase().includes(searchQuery.toLowerCase()) ||
        v.meaning.toLowerCase().includes(searchQuery.toLowerCase()) ||
        (v.reading && v.reading.toLowerCase().includes(searchQuery.toLowerCase()));

      const matchBox = selectedBox === "all" || v.current_box === selectedBox;
      const matchLang = selectedLanguage === "all" || v.language === selectedLanguage;

      return matchQuery && matchBox && matchLang;
    });
  }, [vocabsQuery.data, searchQuery, selectedBox, selectedLanguage]);

  if (vocabsQuery.isPending) return <LoadingState />;
  if (vocabsQuery.isError) return <ErrorState />;

  const allVocabs = vocabsQuery.data ?? [];

  return (
    <div className="flex h-full w-full flex-col overflow-y-auto space-y-4">
      {/* Top Header & Breadcrumb */}
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div className="flex items-center gap-3">
          <Link
            aria-label="Back to Learn"
            className="flex h-9 w-9 items-center justify-center rounded-xl border border-white/10 bg-white/5 text-muted-foreground transition hover:bg-white/10 hover:text-foreground"
            to="/app/learn"
          >
            <ArrowLeft className="h-4 w-4" />
          </Link>
          <div>
            <h1 className="text-xl font-bold tracking-tight">{t("manage.title")}</h1>
            <p className="text-xs text-muted-foreground">
              {t("manage.entriesCount", { count: allVocabs.length })}
            </p>
          </div>
        </div>

        <div className="flex items-center gap-2">
          <button
            className="flex h-9 items-center gap-1.5 rounded-xl border border-white/10 bg-white/5 px-3 text-xs font-semibold text-muted-foreground transition hover:bg-white/10 hover:text-foreground"
            onClick={() => setIsImportOpen(true)}
            type="button"
          >
            <Upload className="h-3.5 w-3.5" />
            {t("manage.importBtn")}
          </button>

          <ButtonInButton icon={Plus} onClick={openCreateModal} variant="primary">
            {t("manage.addWordBtn")}
          </ButtonInButton>
        </div>
      </div>

      {/* Filter & Search Bar */}
      <div className="flex flex-wrap items-center justify-between gap-3 rounded-2xl border border-slate-200/80 bg-white/80 p-3 shadow-sm backdrop-blur-xl dark:border-white/10 dark:bg-white/5">
        <div className="relative min-w-[240px] flex-1">
          <Search className="absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
          <input
            className="h-9 w-full rounded-xl border border-slate-200 bg-white pl-9 pr-3 text-xs text-foreground outline-none transition focus:border-violet-500 dark:border-white/10 dark:bg-black/20"
            onChange={(e) => setSearchQuery(e.target.value)}
            placeholder={t("manage.searchPlaceholder")}
            value={searchQuery}
          />
        </div>

        {/* Leitner Box filter */}
        <div className="flex items-center gap-1">
          <span className="text-xs text-muted-foreground mr-1 flex items-center gap-1">
            <Filter className="h-3 w-3" /> {t("manage.boxFilter")}
          </span>
          {(["all", 1, 2, 3, 4, 5] as const).map((box) => (
            <button
              className={`rounded-lg px-2.5 py-1 text-xs font-mono font-medium transition ${
                selectedBox === box
                  ? "bg-violet-600 text-white shadow-sm"
                  : "bg-slate-100 text-muted-foreground hover:bg-slate-200 hover:text-foreground dark:bg-white/5 dark:text-muted-foreground dark:hover:bg-white/10 dark:hover:text-foreground"
              }`}
              key={String(box)}
              onClick={() => setSelectedBox(box)}
              type="button"
            >
              {box === "all" ? t("manage.all") : `B${box}`}
            </button>
          ))}
        </div>

        {/* Language filter */}
        <div className="flex items-center gap-1 rounded-xl border border-slate-200 bg-slate-100/70 p-0.5 dark:border-white/10 dark:bg-white/5">
          {(["all", "JP", "EN"] as const).map((lang) => (
            <button
              className={`rounded-lg px-2.5 py-1 text-xs font-semibold transition ${
                selectedLanguage === lang
                  ? "bg-white text-foreground shadow-sm dark:bg-violet-600 dark:text-white"
                  : "text-muted-foreground hover:text-foreground"
              }`}
              key={lang}
              onClick={() => setSelectedLanguage(lang)}
              type="button"
            >
              {lang === "all" ? t("manage.all") : lang}
            </button>
          ))}
        </div>
      </div>

      {/* Vocabulary Table Container */}
      <DoubleBezelCard className="overflow-hidden p-0">
        <div className="w-full overflow-x-auto">
          <table className="w-full border-collapse text-left text-xs">
            <thead>
              <tr className="border-b border-slate-200/80 bg-slate-50/70 text-[11px] font-bold uppercase tracking-wider text-slate-700 dark:border-white/10 dark:bg-white/5 dark:text-muted-foreground">
                <th className="py-3 pl-4 pr-2">{t("manage.colWord")}</th>
                <th className="px-3 py-3">{t("manage.colReading")}</th>
                <th className="px-3 py-3">{t("manage.colMeaning")}</th>
                <th className="px-3 py-3">{t("manage.colSentence")}</th>
                <th className="px-2 py-3 text-center">{t("manage.colBox")}</th>
                <th className="px-2 py-3 text-center">{t("manage.colEase")}</th>
                <th className="py-3 pl-2 pr-4 text-right">{t("manage.colActions")}</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-slate-100 dark:divide-white/5">
              {filteredVocabs.length === 0 ? (
                <tr>
                  <td className="py-12 text-center text-muted-foreground" colSpan={7}>
                    {t("manage.noItems")}
                  </td>
                </tr>
              ) : (
                filteredVocabs.map((vocab) => (
                  <tr
                    className="transition hover:bg-slate-50/80 dark:hover:bg-white/5"
                    key={vocab.id}
                  >
                    {/* Word Column with Furigana Ruby */}
                    <td className="py-3 pl-4 pr-2 font-semibold">
                      <div className="flex items-center gap-2">
                        <button
                          aria-label="Play audio"
                          className="flex h-7 w-7 shrink-0 items-center justify-center rounded-lg border border-slate-200 bg-white text-muted-foreground transition hover:bg-slate-100 hover:text-foreground dark:border-white/10 dark:bg-white/5 dark:hover:bg-white/15 dark:hover:text-foreground"
                          onClick={() => handlePlayAudio(vocab)}
                          type="button"
                        >
                          <Volume2
                            className={`h-3.5 w-3.5 ${activeAudioId === vocab.id ? "animate-pulse text-violet-400" : ""}`}
                          />
                        </button>
                        <span className="font-medium text-sm text-foreground" lang={vocab.language === "JP" ? "ja" : undefined}>
                          {vocab.word}
                        </span>
                        <span className="rounded bg-slate-100 text-slate-600 dark:bg-white/10 dark:text-muted-foreground px-1 py-0.2 text-[9px] font-mono">
                          {vocab.language}
                        </span>
                      </div>
                    </td>

                    {/* Reading */}
                    <td className="px-3 py-3 font-mono text-muted-foreground">
                      {vocab.reading || "—"}
                    </td>

                    {/* Meaning */}
                    <td className="px-3 py-3 font-medium text-foreground">
                      {vocab.meaning}
                    </td>

                    {/* Example Sentence */}
                    <td className="max-w-xs truncate px-3 py-3 text-muted-foreground" title={vocab.example_sentence ?? ""}>
                      {vocab.example_sentence ? (
                        <span>
                          {vocab.example_sentence}
                          {vocab.example_translation && (
                            <span className="block text-[10px] text-muted-foreground/70">
                              {vocab.example_translation}
                            </span>
                          )}
                        </span>
                      ) : (
                        "—"
                      )}
                    </td>

                    {/* Leitner Box */}
                    <td className="px-2 py-3 text-center">
                      <span className="inline-flex h-5 w-5 items-center justify-center rounded-full bg-violet-500/15 font-mono text-[11px] font-bold text-violet-400">
                        {vocab.current_box}
                      </span>
                    </td>

                    {/* Ease factor */}
                    <td className="px-2 py-3 text-center font-mono text-muted-foreground">
                      {vocab.ease_factor ? vocab.ease_factor.toFixed(1) : "2.5"}
                    </td>

                    {/* Actions */}
                    <td className="py-3 pl-2 pr-4 text-right">
                      <div className="flex items-center justify-end gap-1">
                        <button
                          aria-label="Edit word"
                          className="flex h-7 w-7 items-center justify-center rounded-lg border border-white/10 bg-white/5 text-muted-foreground transition hover:bg-white/10 hover:text-foreground"
                          onClick={() => openEditModal(vocab)}
                          type="button"
                        >
                          <Edit2 className="h-3 w-3" />
                        </button>
                        <button
                          aria-label="Delete word"
                          className="flex h-7 w-7 items-center justify-center rounded-lg border border-white/10 bg-white/5 text-muted-foreground transition hover:bg-rose-500/20 hover:text-rose-400"
                          onClick={() => {
                            if (confirm(t("manage.deleteConfirm", { word: vocab.word }))) {
                              deleteMutation.mutate(vocab.id);
                            }
                          }}
                          type="button"
                        >
                          <Trash2 className="h-3 w-3" />
                        </button>
                      </div>
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </DoubleBezelCard>

      {/* CREATE / EDIT MODAL */}
      {isModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4 backdrop-blur-sm">
          <div className="w-full max-w-md">
            <DoubleBezelCard className="p-6">
              <div className="flex items-center justify-between border-b border-white/10 pb-3">
                <h2 className="text-base font-bold">
                  {editingVocab ? t("manage.editModalTitle") : t("manage.createModalTitle")}
                </h2>
                <button
                  aria-label="Close modal"
                  className="rounded-lg p-1 text-muted-foreground hover:bg-white/10 hover:text-foreground"
                  onClick={() => setIsModalOpen(false)}
                  type="button"
                >
                  <X className="h-4 w-4" />
                </button>
              </div>

              <form className="mt-4 space-y-3" onSubmit={handleSave}>
                <div className="grid grid-cols-2 gap-2">
                  <div>
                    <label className="text-[11px] font-semibold text-muted-foreground">{t("hub.language")}</label>
                    <select
                      className="mt-1 h-9 w-full rounded-xl border border-slate-200 bg-white px-3 text-xs text-foreground outline-none transition focus:border-violet-500 dark:border-white/10 dark:bg-black/30"
                      onChange={(e) => setFormLanguage(e.target.value as VocabularyLanguage)}
                      value={formLanguage}
                    >
                      <option value="JP">{t("hub.langJp")}</option>
                      <option value="EN">{t("hub.langEn")}</option>
                    </select>
                  </div>
                  <div>
                    <label className="text-[11px] font-semibold text-muted-foreground">{t("hub.wordKanji")}</label>
                    <input
                      className="mt-1 h-9 w-full rounded-xl border border-slate-200 bg-white px-3 text-xs text-foreground outline-none transition focus:border-violet-500 dark:border-white/10 dark:bg-black/30"
                      onChange={(e) => setFormWord(e.target.value)}
                      placeholder={t("hub.wordPlaceholder")}
                      required
                      value={formWord}
                    />
                  </div>
                </div>

                <div>
                  <label className="text-[11px] font-semibold text-muted-foreground">{t("hub.readingFurigana")}</label>
                  <input
                    className="mt-1 h-9 w-full rounded-xl border border-slate-200 bg-white px-3 text-xs text-foreground outline-none transition focus:border-violet-500 dark:border-white/10 dark:bg-black/30"
                    onChange={(e) => setFormReading(e.target.value)}
                    placeholder={t("hub.readingPlaceholder")}
                    value={formReading}
                  />
                </div>

                <div>
                  <label className="text-[11px] font-semibold text-muted-foreground">{t("hub.meaning")}</label>
                  <input
                    className="mt-1 h-9 w-full rounded-xl border border-slate-200 bg-white px-3 text-xs text-foreground outline-none transition focus:border-violet-500 dark:border-white/10 dark:bg-black/30"
                    onChange={(e) => setFormMeaning(e.target.value)}
                    placeholder={t("hub.meaningPlaceholder")}
                    required
                    value={formMeaning}
                  />
                </div>

                <div>
                  <label className="text-[11px] font-semibold text-muted-foreground">{t("manage.colSentence")}</label>
                  <input
                    className="mt-1 h-9 w-full rounded-xl border border-slate-200 bg-white px-3 text-xs text-foreground outline-none transition focus:border-violet-500 dark:border-white/10 dark:bg-black/30"
                    onChange={(e) => setFormSentence(e.target.value)}
                    placeholder={t("manage.sentencePlaceholder")}
                    value={formSentence}
                  />
                </div>

                <div>
                  <label className="text-[11px] font-semibold text-muted-foreground">{t("manage.translationPlaceholder")}</label>
                  <input
                    className="mt-1 h-9 w-full rounded-xl border border-slate-200 bg-white px-3 text-xs text-foreground outline-none transition focus:border-violet-500 dark:border-white/10 dark:bg-black/30"
                    onChange={(e) => setFormTranslation(e.target.value)}
                    placeholder={t("manage.translationPlaceholder")}
                    value={formTranslation}
                  />
                </div>

                <div>
                  <label className="text-[11px] font-semibold text-muted-foreground">{t("manage.tagsLabel")}</label>
                  <input
                    className="mt-1 h-9 w-full rounded-xl border border-slate-200 bg-white px-3 text-xs text-foreground outline-none transition focus:border-violet-500 dark:border-white/10 dark:bg-black/30"
                    onChange={(e) => setFormTags(e.target.value)}
                    placeholder={t("manage.tagsPlaceholder")}
                    value={formTags}
                  />
                </div>

                <div className="flex justify-end gap-2 pt-2">
                  <button
                    className="rounded-xl border border-white/10 px-4 py-2 text-xs font-semibold hover:bg-white/5"
                    onClick={() => setIsModalOpen(false)}
                    type="button"
                  >
                    {t("manage.cancel")}
                  </button>
                  <button
                    className="rounded-xl bg-violet-600 px-4 py-2 text-xs font-semibold text-white shadow-lg shadow-violet-600/20 hover:bg-violet-500"
                    type="submit"
                  >
                    {editingVocab ? t("manage.updateWord") : t("manage.saveWord")}
                  </button>
                </div>
              </form>
            </DoubleBezelCard>
          </div>
        </div>
      )}

      {/* IMPORT MODAL */}
      {isImportOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 p-4 backdrop-blur-sm">
          <div className="w-full max-w-lg">
            <DoubleBezelCard className="p-6">
              <div className="flex items-center justify-between border-b border-white/10 pb-3">
                <div className="flex items-center gap-2">
                  <FileSpreadsheet className="h-5 w-5 text-violet-400" />
                  <h2 className="text-base font-bold">{t("manage.importModalTitle")}</h2>
                </div>
                <button
                  aria-label="Close modal"
                  className="rounded-lg p-1 text-muted-foreground hover:bg-white/10 hover:text-foreground"
                  onClick={() => setIsImportOpen(false)}
                  type="button"
                >
                  <X className="h-4 w-4" />
                </button>
              </div>

              <div className="mt-4 space-y-3">
                <div className="flex items-center gap-2">
                  <span className="text-xs text-muted-foreground">{t("manage.format")}</span>
                  <button
                    className={`rounded-lg px-2.5 py-1 text-xs font-mono transition ${importFormat === "csv" ? "bg-violet-600 text-white" : "bg-white/5 text-muted-foreground"}`}
                    onClick={() => setImportFormat("csv")}
                    type="button"
                  >
                    CSV
                  </button>
                  <button
                    className={`rounded-lg px-2.5 py-1 text-xs font-mono transition ${importFormat === "json" ? "bg-violet-600 text-white" : "bg-white/5 text-muted-foreground"}`}
                    onClick={() => setImportFormat("json")}
                    type="button"
                  >
                    JSON
                  </button>
                </div>

                <p className="text-[11px] text-muted-foreground">
                  {importFormat === "csv"
                    ? t("manage.csvHeadersHint")
                    : t("manage.jsonArrayHint")}
                </p>

                <textarea
                  className="h-48 w-full rounded-xl border border-slate-200 bg-white p-3 font-mono text-xs text-foreground outline-none focus:border-violet-500 dark:border-white/10 dark:bg-black/30"
                  onChange={(e) => setImportText(e.target.value)}
                  placeholder={
                    importFormat === "csv"
                      ? "JP,桜,さくら,cherry blossom,桜が咲いた,The cherry blossoms bloomed\nJP,雨,あめ,rain,雨が降っている,It is raining"
                      : '[{"language":"JP","word":"桜","reading":"さくら","meaning":"cherry blossom"}]'
                  }
                  value={importText}
                />

                <div className="flex justify-end gap-2 pt-2">
                  <button
                    className="rounded-xl border border-slate-200 bg-slate-100 px-4 py-2 text-xs font-semibold text-foreground hover:bg-slate-200 dark:border-white/10 dark:bg-transparent dark:text-muted-foreground dark:hover:bg-white/5"
                    onClick={() => setIsImportOpen(false)}
                    type="button"
                  >
                    {t("manage.cancel")}
                  </button>
                  <button
                    className="flex items-center gap-1 rounded-xl bg-violet-600 px-4 py-2 text-xs font-semibold text-white shadow-lg shadow-violet-600/20 hover:bg-violet-500"
                    disabled={importMutation.isPending || !importText.trim()}
                    onClick={handleImportSubmit}
                    type="button"
                  >
                    <Check className="h-3.5 w-3.5" />
                    {importMutation.isPending ? t("manage.importing") : t("manage.executeImport")}
                  </button>
                </div>
              </div>
            </DoubleBezelCard>
          </div>
        </div>
      )}
    </div>
  );
}
