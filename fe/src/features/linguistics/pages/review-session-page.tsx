import { motion, AnimatePresence } from "framer-motion";
import {
  ArrowLeft,
  CheckCircle2,
  HelpCircle,
  Keyboard,
  RotateCw,
  Trophy,
  Volume2,
} from "lucide-react";
import { useEffect, useMemo, useRef, useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { ButtonInButton } from "../../../components/design-system/button-in-button";
import { DoubleBezelCard } from "../../../components/design-system/double-bezel-card";
import { ErrorState } from "../../../components/feedback/error-state";
import { LoadingState } from "../../../components/feedback/loading-state";
import {
  endLearningSession,
  generateAudio,
  startLearningSession,
} from "../api/vocabulary-api";
import { useReviewQueue, useReviewVocabulary, useVocabularies } from "../hooks/use-vocabularies";
import type { StudyMode, Vocabulary } from "../types";

export function ReviewSessionPage() {
  const { t } = useTranslation("learning");
  const navigate = useNavigate();
  const queueQuery = useReviewQueue();
  const allVocabsQuery = useVocabularies();
  const reviewMutation = useReviewVocabulary();

  const [currentIndex, setCurrentIndex] = useState(0);
  const [isFlipped, setIsFlipped] = useState(false);
  const [studyMode, setStudyMode] = useState<StudyMode>("flashcard");
  const [selectedChoice, setSelectedChoice] = useState<string | null>(null);
  const [typeInput, setTypeInput] = useState("");
  const [isTypeSubmitted, setIsTypeSubmitted] = useState(false);
  const [sessionStartTime] = useState<number>(Date.now());
  const [sessionId, setSessionId] = useState<string | null>(null);
  const [correctCount, setCorrectCount] = useState(0);
  const [reviewedCount, setReviewedCount] = useState(0);
  const [isFinished, setIsFinished] = useState(false);
  const [isPlayingAudio, setIsPlayingAudio] = useState(false);

  const audioRef = useRef<HTMLAudioElement | null>(null);

  // Initialize Learning Session on backend
  useEffect(() => {
    let active = true;
    startLearningSession("JP", "review")
      .then((session) => {
        if (active) setSessionId(session.id);
      })
      .catch(() => {
        // Fallback gracefully if session endpoint has transient issue
      });
    return () => {
      active = false;
    };
  }, []);

  const queue = queueQuery.data ?? [];
  const currentCard: Vocabulary | undefined = queue[currentIndex];

  // Distractors for Multiple Choice mode
  const multipleChoiceOptions = useMemo(() => {
    if (!currentCard) return [];
    const all = allVocabsQuery.data ?? [];
    const distractors = all
      .filter((v) => v.id !== currentCard.id)
      .sort(() => 0.5 - Math.random())
      .slice(0, 3)
      .map((v) => v.meaning);
    const choices = [...distractors, currentCard.meaning];
    return choices.sort(() => 0.5 - Math.random());
  }, [currentCard, allVocabsQuery.data]);

  // Audio speech synthesis helper
  const handlePlayAudio = async () => {
    if (!currentCard) return;
    setIsPlayingAudio(true);
    try {
      if (currentCard.audio_url) {
        if (!audioRef.current) {
          audioRef.current = new Audio(currentCard.audio_url);
        } else {
          audioRef.current.src = currentCard.audio_url;
        }
        await audioRef.current.play();
      } else {
        // Request backend audio generation or use Web Speech API fallback
        const res = await generateAudio(currentCard.id).catch(() => null);
        if (res?.audio_url) {
          if (!audioRef.current) audioRef.current = new Audio(res.audio_url);
          else audioRef.current.src = res.audio_url;
          await audioRef.current.play();
        } else if ("speechSynthesis" in window) {
          const utterance = new SpeechSynthesisUtterance(currentCard.word);
          utterance.lang = currentCard.language === "JP" ? "ja-JP" : "en-US";
          window.speechSynthesis.speak(utterance);
        }
      }
    } catch {
      // Audio playback fallback ignored
    } finally {
      setIsPlayingAudio(false);
    }
  };

  const handleGrade = async (quality: number) => {
    if (!currentCard) return;
    const isCorrect = quality >= 3;
    if (isCorrect) setCorrectCount((prev) => prev + 1);
    const newReviewedCount = reviewedCount + 1;
    setReviewedCount(newReviewedCount);

    await reviewMutation.mutateAsync({ id: currentCard.id, quality });

    // Transition to next card or finish session
    if (currentIndex + 1 < queue.length) {
      setIsFlipped(false);
      setSelectedChoice(null);
      setTypeInput("");
      setIsTypeSubmitted(false);
      setCurrentIndex((prev) => prev + 1);
    } else {
      setIsFinished(true);
      const durationSeconds = Math.round((Date.now() - sessionStartTime) / 1000);
      if (sessionId) {
        endLearningSession(sessionId, newReviewedCount, correctCount + (isCorrect ? 1 : 0), durationSeconds).catch(() => {});
      }
    }
  };

  // Keyboard navigation
  useEffect(() => {
    const onKeyDown = (e: KeyboardEvent) => {
      if (isFinished) return;
      if (e.target instanceof HTMLInputElement || e.target instanceof HTMLTextAreaElement) {
        if (e.key === "Enter" && studyMode === "type_in" && !isTypeSubmitted) {
          e.preventDefault();
          setIsTypeSubmitted(true);
          setIsFlipped(true);
        }
        return;
      }

      if (e.code === "Space") {
        e.preventDefault();
        setIsFlipped((prev) => !prev);
      } else if (e.key === "1") {
        e.preventDefault();
        handleGrade(1);
      } else if (e.key === "2") {
        e.preventDefault();
        handleGrade(2);
      } else if (e.key === "3") {
        e.preventDefault();
        handleGrade(4);
      } else if (e.key === "4") {
        e.preventDefault();
        handleGrade(5);
      } else if (e.key.toLowerCase() === "a") {
        e.preventDefault();
        handlePlayAudio();
      } else if (e.key === "Escape") {
        navigate("/app/learn");
      }
    };

    window.addEventListener("keydown", onKeyDown);
    return () => window.removeEventListener("keydown", onKeyDown);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [currentCard, isFlipped, isFinished, studyMode, isTypeSubmitted]);

  if (queueQuery.isPending) return <LoadingState />;
  if (queueQuery.isError) return <ErrorState />;

  if (queue.length === 0 || isFinished || !currentCard) {
    const accuracy = reviewedCount > 0 ? Math.round((correctCount / reviewedCount) * 100) : 100;
    const xpEarned = reviewedCount * 15;

    return (
      <div className="flex h-full min-h-[500px] w-full items-center justify-center p-4">
        <motion.div
          animate={{ opacity: 1, scale: 1 }}
          className="w-full max-w-lg"
          initial={{ opacity: 0, scale: 0.95 }}
          transition={{ duration: 0.4, ease: [0.32, 0.72, 0, 1] }}
        >
          <DoubleBezelCard className="p-8 text-center">
            <div className="mx-auto flex h-16 w-16 items-center justify-center rounded-2xl bg-violet-500/10 text-violet-500 shadow-inner">
              <Trophy className="h-8 w-8" />
            </div>
            <h1 className="mt-4 text-2xl font-bold tracking-tight">
              {queue.length === 0 ? t("reviewSession.queueClearTitle") : t("reviewSession.sessionCompletedTitle")}
            </h1>
            <p className="mt-1 text-sm text-muted-foreground">
              {queue.length === 0
                ? t("reviewSession.queueClearSubtitle")
                : t("reviewSession.sessionCompletedSubtitle")}
            </p>

            <div className="mt-6 grid grid-cols-3 gap-3">
              <div className="rounded-xl border border-white/10 bg-white/5 p-3">
                <p className="text-xs text-muted-foreground">{t("reviewSession.reviewed")}</p>
                <p className="mt-1 font-mono text-xl font-bold">{reviewedCount}</p>
              </div>
              <div className="rounded-xl border border-white/10 bg-white/5 p-3">
                <p className="text-xs text-muted-foreground">{t("reviewSession.accuracy")}</p>
                <p className="mt-1 font-mono text-xl font-bold text-emerald-400">{accuracy}%</p>
              </div>
              <div className="rounded-xl border border-white/10 bg-white/5 p-3">
                <p className="text-xs text-muted-foreground">{t("reviewSession.xpGained")}</p>
                <p className="mt-1 font-mono text-xl font-bold text-violet-400">+{xpEarned}</p>
              </div>
            </div>

            <div className="mt-8 flex justify-center">
              <Link to="/app/learn">
                <ButtonInButton icon={ArrowLeft} variant="primary">
                  {t("reviewSession.backToHub")}
                </ButtonInButton>
              </Link>
            </div>
          </DoubleBezelCard>
        </motion.div>
      </div>
    );
  }

  const progressPercent = Math.round(((currentIndex + 1) / queue.length) * 100);

  const modeLabels: Record<StudyMode, string> = {
    flashcard: t("reviewSession.modes.flashcard"),
    multiple_choice: t("reviewSession.modes.choice"),
    type_in: t("reviewSession.modes.typeIn"),
    audio_quiz: t("reviewSession.modes.audio"),
  };

  return (
    <div className="flex h-full w-full flex-col overflow-y-auto">
      {/* Top Session Bar */}
      <div className="flex shrink-0 items-center justify-between border-b border-white/10 px-4 py-3">
        <div className="flex items-center gap-3">
          <Link
            aria-label="Exit session"
            className="flex h-9 w-9 items-center justify-center rounded-xl border border-white/10 bg-white/5 text-muted-foreground transition hover:bg-white/10 hover:text-foreground"
            to="/app/learn"
          >
            <ArrowLeft className="h-4 w-4" />
          </Link>
          <div>
            <p className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
              {t("reviewSession.cardProgress", { current: currentIndex + 1, total: queue.length })}
            </p>
            <div className="mt-1 h-1.5 w-32 overflow-hidden rounded-full bg-slate-800">
              <motion.div
                animate={{ width: `${progressPercent}%` }}
                className="h-full bg-gradient-to-r from-violet-500 to-indigo-500"
                transition={{ duration: 0.3 }}
              />
            </div>
          </div>
        </div>

        {/* Study Mode Selector */}
        <div className="flex items-center gap-1 rounded-xl border border-white/10 bg-white/5 p-1">
          {(["flashcard", "multiple_choice", "type_in", "audio_quiz"] as StudyMode[]).map((mode) => (
            <button
              className={`rounded-lg px-2.5 py-1 text-xs font-medium transition ${
                studyMode === mode
                  ? "bg-violet-600 text-white shadow-sm"
                  : "text-muted-foreground hover:text-foreground"
              }`}
              key={mode}
              onClick={() => {
                setStudyMode(mode);
                setIsFlipped(false);
                setSelectedChoice(null);
                setTypeInput("");
                setIsTypeSubmitted(false);
              }}
              type="button"
            >
              {modeLabels[mode]}
            </button>
          ))}
        </div>
      </div>

      {/* Main Review Canvas */}
      <div className="flex flex-1 flex-col items-center justify-center p-4">
        <div className="w-full max-w-xl perspective-[1200px]">
          <AnimatePresence mode="wait">
            <motion.div
              animate={{ rotateY: isFlipped ? 180 : 0, scale: 1, opacity: 1 }}
              className="relative min-h-[360px] w-full cursor-pointer select-none rounded-[1.75rem]"
              exit={{ opacity: 0, scale: 0.95 }}
              initial={{ scale: 0.95, opacity: 0 }}
              key={currentCard.id}
              onClick={() => {
                if (studyMode === "flashcard" || studyMode === "audio_quiz") {
                  setIsFlipped((prev) => !prev);
                }
              }}
              style={{ transformStyle: "preserve-3d" }}
              transition={{ duration: 0.5, ease: [0.32, 0.72, 0, 1] }}
            >
              {/* FRONT OF CARD */}
              <div
                className="absolute inset-0 flex flex-col justify-between rounded-[1.75rem] border border-white/15 bg-gradient-to-b from-slate-900/90 to-slate-950/95 p-8 shadow-2xl backdrop-blur-xl"
                style={{ backfaceVisibility: "hidden" }}
              >
                <div className="flex items-center justify-between">
                  <span className="rounded-full border border-violet-500/20 bg-violet-500/10 px-3 py-0.5 text-xs font-semibold text-violet-400">
                    Box {currentCard.current_box} · {currentCard.language}
                  </span>
                  <button
                    aria-label="Play audio"
                    className="flex h-9 w-9 items-center justify-center rounded-full border border-white/10 bg-white/5 text-muted-foreground transition hover:bg-white/15 hover:text-foreground"
                    onClick={(e) => {
                      e.stopPropagation();
                      handlePlayAudio();
                    }}
                    type="button"
                  >
                    <Volume2 className={`h-4 w-4 ${isPlayingAudio ? "animate-pulse text-violet-400" : ""}`} />
                  </button>
                </div>

                {/* Central Question Display */}
                <div className="my-auto text-center">
                  {studyMode === "audio_quiz" ? (
                    <div className="flex flex-col items-center">
                      <div className="flex h-20 w-20 items-center justify-center rounded-full bg-violet-500/10 text-violet-400 shadow-inner">
                        <Volume2 className="h-10 w-10 animate-pulse" />
                      </div>
                      <p className="mt-4 text-sm text-muted-foreground">{t("reviewSession.audioPrompt")}</p>
                    </div>
                  ) : (
                    <>
                      <p className="text-4xl font-extrabold tracking-tight sm:text-5xl" lang={currentCard.language === "JP" ? "ja" : undefined}>
                        {currentCard.word}
                      </p>
                      {currentCard.reading && !isFlipped && studyMode === "flashcard" && (
                        <p className="mt-2 text-sm text-muted-foreground/60">{currentCard.reading}</p>
                      )}
                    </>
                  )}

                  {/* Multiple Choice interactive buttons */}
                  {studyMode === "multiple_choice" && (
                    <div className="mt-6 grid grid-cols-2 gap-2" onClick={(e) => e.stopPropagation()}>
                      {multipleChoiceOptions.map((choice, idx) => {
                        const isChosen = selectedChoice === choice;
                        const isCorrect = choice === currentCard.meaning;
                        let btnStyle = "border-white/10 bg-white/5 hover:bg-white/10 text-foreground";
                        if (selectedChoice !== null) {
                          if (isCorrect) btnStyle = "border-emerald-500/50 bg-emerald-500/20 text-emerald-300 ring-2 ring-emerald-500";
                          else if (isChosen) btnStyle = "border-rose-500/50 bg-rose-500/20 text-rose-300 ring-2 ring-rose-500";
                        }

                        return (
                          <button
                            className={`rounded-xl border p-3 text-left text-xs font-medium transition ${btnStyle}`}
                            key={idx}
                            onClick={() => {
                              setSelectedChoice(choice);
                              setIsFlipped(true);
                            }}
                            type="button"
                          >
                            <span className="mr-2 font-mono text-muted-foreground">{idx + 1}.</span>
                            {choice}
                          </button>
                        );
                      })}
                    </div>
                  )}

                  {/* Type-in Answer Input */}
                  {studyMode === "type_in" && (
                    <div className="mt-6 flex flex-col items-center gap-2" onClick={(e) => e.stopPropagation()}>
                      <input
                        autoFocus
                        className="h-11 w-full max-w-sm rounded-xl border border-white/20 bg-black/40 px-4 text-center text-sm outline-none transition focus:border-violet-500 focus:ring-2 focus:ring-violet-500/20"
                        onChange={(e) => setTypeInput(e.target.value)}
                        onKeyDown={(e) => {
                          if (e.key === "Enter" && !isTypeSubmitted) {
                            setIsTypeSubmitted(true);
                            setIsFlipped(true);
                          }
                        }}
                        placeholder={t("reviewSession.typePlaceholder")}
                        value={typeInput}
                      />
                      <button
                        className="rounded-lg bg-violet-600 px-4 py-1.5 text-xs font-semibold text-white shadow transition hover:bg-violet-500"
                        onClick={() => {
                          setIsTypeSubmitted(true);
                          setIsFlipped(true);
                        }}
                        type="button"
                      >
                        {t("reviewSession.checkAnswer")}
                      </button>
                    </div>
                  )}
                </div>

                <div className="flex items-center justify-between text-xs text-muted-foreground">
                  <span className="flex items-center gap-1">
                    <Keyboard className="h-3 w-3" /> {t("reviewSession.spaceToFlip")}
                  </span>
                  <span className="flex items-center gap-1">
                    <RotateCw className="h-3 w-3" /> {t("reviewSession.flipCard")}
                  </span>
                </div>
              </div>

              {/* BACK OF CARD */}
              <div
                className="absolute inset-0 flex flex-col justify-between rounded-[1.75rem] border border-violet-500/30 bg-gradient-to-b from-slate-900/95 to-slate-950/98 p-8 shadow-2xl backdrop-blur-xl [transform:rotateY(180deg)]"
                style={{ backfaceVisibility: "hidden" }}
              >
                <div className="flex items-center justify-between">
                  <span className="rounded-full border border-emerald-500/30 bg-emerald-500/10 px-3 py-0.5 text-xs font-semibold text-emerald-400">
                    {t("reviewSession.answerRevealed")}
                  </span>
                  <button
                    aria-label="Play audio again"
                    className="flex h-9 w-9 items-center justify-center rounded-full border border-white/10 bg-white/5 text-muted-foreground transition hover:bg-white/15 hover:text-foreground"
                    onClick={(e) => {
                      e.stopPropagation();
                      handlePlayAudio();
                    }}
                    type="button"
                  >
                    <Volume2 className="h-4 w-4" />
                  </button>
                </div>

                <div className="my-auto text-center">
                  <p className="text-3xl font-extrabold tracking-tight" lang={currentCard.language === "JP" ? "ja" : undefined}>
                    {currentCard.word}
                  </p>
                  {currentCard.reading && (
                    <p className="mt-1 font-mono text-base text-violet-400">{currentCard.reading}</p>
                  )}
                  <p className="mt-4 text-xl font-semibold text-emerald-400">{currentCard.meaning}</p>

                  {currentCard.example_sentence && (
                    <div className="mt-4 rounded-xl border border-white/10 bg-white/5 p-3 text-left">
                      <p className="text-sm font-medium text-slate-200">{currentCard.example_sentence}</p>
                      {currentCard.example_translation && (
                        <p className="mt-1 text-xs text-muted-foreground">{currentCard.example_translation}</p>
                      )}
                    </div>
                  )}

                  {/* Type-in feedback */}
                  {studyMode === "type_in" && isTypeSubmitted && (
                    <div className="mt-3 flex items-center justify-center gap-2 text-xs">
                      {typeInput.trim().toLowerCase() === currentCard.meaning.trim().toLowerCase() ? (
                        <span className="flex items-center gap-1 font-semibold text-emerald-400">
                          <CheckCircle2 className="h-4 w-4" /> {t("reviewSession.exactMatch")}
                        </span>
                      ) : (
                        <span className="flex items-center gap-1 text-amber-400">
                          <HelpCircle className="h-4 w-4" /> {t("reviewSession.yourAnswer", { answer: typeInput || t("reviewSession.emptyAnswer") })}
                        </span>
                      )}
                    </div>
                  )}
                </div>

                <p className="text-center text-xs text-muted-foreground">{t("reviewSession.gradePrompt")}</p>
              </div>
            </motion.div>
          </AnimatePresence>

          {/* SM-2 Recall Rating Bar */}
          <div className="mt-6 flex flex-col gap-2">
            <div className="grid grid-cols-4 gap-2">
              <button
                className="flex flex-col items-center justify-center rounded-xl border border-rose-500/20 bg-rose-500/10 py-2.5 transition active:scale-95 hover:bg-rose-500/20 hover:border-rose-500/40"
                onClick={() => handleGrade(1)}
                type="button"
              >
                <span className="font-mono text-xs font-semibold text-rose-400">{t("reviewSession.rating1")}</span>
                <span className="text-[10px] text-muted-foreground">{t("reviewSession.rating1Interval")}</span>
              </button>

              <button
                className="flex flex-col items-center justify-center rounded-xl border border-amber-500/20 bg-amber-500/10 py-2.5 transition active:scale-95 hover:bg-amber-500/20 hover:border-amber-500/40"
                onClick={() => handleGrade(2)}
                type="button"
              >
                <span className="font-mono text-xs font-semibold text-amber-400">{t("reviewSession.rating2")}</span>
                <span className="text-[10px] text-muted-foreground">{t("reviewSession.rating2Interval")}</span>
              </button>

              <button
                className="flex flex-col items-center justify-center rounded-xl border border-emerald-500/20 bg-emerald-500/10 py-2.5 transition active:scale-95 hover:bg-emerald-500/20 hover:border-emerald-500/40"
                onClick={() => handleGrade(4)}
                type="button"
              >
                <span className="font-mono text-xs font-semibold text-emerald-400">{t("reviewSession.rating3")}</span>
                <span className="text-[10px] text-muted-foreground">{t("reviewSession.rating3Interval")}</span>
              </button>

              <button
                className="flex flex-col items-center justify-center rounded-xl border border-blue-500/20 bg-blue-500/10 py-2.5 transition active:scale-95 hover:bg-blue-500/20 hover:border-blue-500/40"
                onClick={() => handleGrade(5)}
                type="button"
              >
                <span className="font-mono text-xs font-semibold text-blue-400">{t("reviewSession.rating4")}</span>
                <span className="text-[10px] text-muted-foreground">{t("reviewSession.rating4Interval")}</span>
              </button>
            </div>

            <div className="flex items-center justify-between px-1 text-[11px] text-muted-foreground">
              <span>{t("reviewSession.shortcuts")}</span>
              <span>{t("reviewSession.shortcutsExtra")}</span>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
