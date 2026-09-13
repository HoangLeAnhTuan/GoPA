import {
  ArrowRight,
  BookOpen,
  CheckSquare,
  Coins,
  Compass,
  FileText,
  Languages,
  LayoutDashboard,
  ListTodo,
  NotebookPen,
  Search,
  Settings,
  Sparkles,
  Timer,
  WalletCards,
  type LucideIcon,
} from "lucide-react";
import { useEffect, useMemo, useRef, useState } from "react";
import { useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { useJournals } from "../../features/journal/hooks/use-journals";
import { useVocabularies } from "../../features/linguistics/hooks/use-vocabularies";
import { useTasks } from "../../features/tasks/hooks/use-tasks";
import { DoubleBezelCard } from "./double-bezel-card";

interface CommandItem {
  id: string;
  categoryKey: "navigation" | "quickAction" | "content";
  title: string;
  subtitle?: string;
  icon: LucideIcon;
  onSelect: () => void;
}

interface CommandPaletteProps {
  isOpen: boolean;
  onClose: () => void;
}

export function CommandPalette({ isOpen, onClose }: CommandPaletteProps) {
  const { t } = useTranslation();
  const navigate = useNavigate();
  const [query, setQuery] = useState("");
  const [selectedIndex, setSelectedIndex] = useState(0);

  const inputRef = useRef<HTMLInputElement>(null);

  // Live queries for tasks, vocabularies, journals
  const tasksQ = useTasks();
  const vocabsQ = useVocabularies();
  const journalsQ = useJournals();

  useEffect(() => {
    if (isOpen) {
      setQuery("");
      setSelectedIndex(0);
      setTimeout(() => inputRef.current?.focus(), 50);
    }
  }, [isOpen]);

  const navigationCommands: CommandItem[] = useMemo(
    () => [
      {
        id: "nav-today",
        categoryKey: "navigation",
        title: t("commandPalette.navToday"),
        subtitle: t("commandPalette.navTodaySub"),
        icon: LayoutDashboard,
        onSelect: () => {
          navigate("/app/today");
          onClose();
        },
      },
      {
        id: "nav-tasks",
        categoryKey: "navigation",
        title: t("commandPalette.navTasks"),
        subtitle: t("commandPalette.navTasksSub"),
        icon: ListTodo,
        onSelect: () => {
          navigate("/app/tasks");
          onClose();
        },
      },
      {
        id: "nav-pomodoro",
        categoryKey: "navigation",
        title: t("commandPalette.navPomodoro"),
        subtitle: t("commandPalette.navPomodoroSub"),
        icon: Timer,
        onSelect: () => {
          navigate("/app/pomodoro");
          onClose();
        },
      },
      {
        id: "nav-learn",
        categoryKey: "navigation",
        title: t("commandPalette.navLearn"),
        subtitle: t("commandPalette.navLearnSub"),
        icon: Languages,
        onSelect: () => {
          navigate("/app/learn");
          onClose();
        },
      },
      {
        id: "nav-review",
        categoryKey: "navigation",
        title: t("commandPalette.navReview"),
        subtitle: t("commandPalette.navReviewSub"),
        icon: Sparkles,
        onSelect: () => {
          navigate("/app/learn/review");
          onClose();
        },
      },
      {
        id: "nav-vocabulary",
        categoryKey: "navigation",
        title: t("commandPalette.navVocabulary"),
        subtitle: t("commandPalette.navVocabularySub"),
        icon: BookOpen,
        onSelect: () => {
          navigate("/app/learn/vocabulary");
          onClose();
        },
      },
      {
        id: "nav-finance",
        categoryKey: "navigation",
        title: t("commandPalette.navFinance"),
        subtitle: t("commandPalette.navFinanceSub"),
        icon: WalletCards,
        onSelect: () => {
          navigate("/app/finance");
          onClose();
        },
      },
      {
        id: "nav-budgets",
        categoryKey: "navigation",
        title: t("commandPalette.navBudgets"),
        subtitle: t("commandPalette.navBudgetsSub"),
        icon: Coins,
        onSelect: () => {
          navigate("/app/finance/budgets");
          onClose();
        },
      },
      {
        id: "nav-goals",
        categoryKey: "navigation",
        title: t("commandPalette.navGoals"),
        subtitle: t("commandPalette.navGoalsSub"),
        icon: Compass,
        onSelect: () => {
          navigate("/app/finance/goals");
          onClose();
        },
      },
      {
        id: "nav-journal",
        categoryKey: "navigation",
        title: t("commandPalette.navJournal"),
        subtitle: t("commandPalette.navJournalSub"),
        icon: NotebookPen,
        onSelect: () => {
          navigate("/app/journal");
          onClose();
        },
      },
      {
        id: "nav-settings",
        categoryKey: "navigation",
        title: t("commandPalette.navSettings"),
        subtitle: t("commandPalette.navSettingsSub"),
        icon: Settings,
        onSelect: () => {
          navigate("/app/settings");
          onClose();
        },
      },
    ],
    [navigate, onClose, t]
  );

  // Dynamic live search items from backend data
  const contentCommands: CommandItem[] = useMemo(() => {
    if (!query.trim()) return [];
    const q = query.toLowerCase();
    const results: CommandItem[] = [];

    // Tasks matches
    for (const task of tasksQ.data ?? []) {
      if (task.title.toLowerCase().includes(q)) {
        results.push({
          id: `task-${task.id}`,
          categoryKey: "content",
          title: task.title,
          subtitle: t("commandPalette.contentTask", { priority: task.priority }),
          icon: CheckSquare,
          onSelect: () => {
            navigate("/app/tasks");
            onClose();
          },
        });
      }
    }

    // Vocabulary matches
    for (const v of vocabsQ.data ?? []) {
      if (v.word.toLowerCase().includes(q) || v.meaning.toLowerCase().includes(q)) {
        results.push({
          id: `vocab-${v.id}`,
          categoryKey: "content",
          title: `${v.word} (${v.meaning})`,
          subtitle: t("commandPalette.contentVocab", { box: v.current_box, language: v.language }),
          icon: BookOpen,
          onSelect: () => {
            navigate("/app/learn/vocabulary");
            onClose();
          },
        });
      }
    }

    // Journal matches
    for (const j of journalsQ.data ?? []) {
      if (j.title.toLowerCase().includes(q) || j.content.toLowerCase().includes(q)) {
        results.push({
          id: `journal-${j.id}`,
          categoryKey: "content",
          title: j.title,
          subtitle: t("commandPalette.contentJournal", {
            date: new Date(j.published_date).toLocaleDateString(),
          }),
          icon: FileText,
          onSelect: () => {
            navigate(`/app/journal/${j.id}`);
            onClose();
          },
        });
      }
    }

    return results.slice(0, 8);
  }, [query, tasksQ.data, vocabsQ.data, journalsQ.data, navigate, onClose, t]);

  // Combine and filter commands
  const filteredCommands = useMemo(() => {
    const q = query.trim().toLowerCase();
    const navMatches = navigationCommands.filter(
      (cmd) =>
        q === "" ||
        cmd.title.toLowerCase().includes(q) ||
        (cmd.subtitle && cmd.subtitle.toLowerCase().includes(q))
    );
    return [...contentCommands, ...navMatches];
  }, [navigationCommands, contentCommands, query]);

  // Keyboard navigation
  useEffect(() => {
    if (!isOpen) return;

    const onKeyDown = (e: KeyboardEvent) => {
      if (e.key === "Escape") {
        e.preventDefault();
        onClose();
      } else if (e.key === "ArrowDown") {
        e.preventDefault();
        setSelectedIndex((prev) => (prev + 1) % Math.max(1, filteredCommands.length));
      } else if (e.key === "ArrowUp") {
        e.preventDefault();
        setSelectedIndex((prev) => (prev - 1 + filteredCommands.length) % Math.max(1, filteredCommands.length));
      } else if (e.key === "Enter") {
        e.preventDefault();
        if (filteredCommands[selectedIndex]) {
          filteredCommands[selectedIndex].onSelect();
        }
      }
    };

    window.addEventListener("keydown", onKeyDown);
    return () => window.removeEventListener("keydown", onKeyDown);
  }, [isOpen, filteredCommands, selectedIndex, onClose]);

  if (!isOpen) return null;

  return (
    <div
      className="fixed inset-0 z-50 flex items-start justify-center bg-black/65 p-4 pt-[12vh] backdrop-blur-md"
      onClick={onClose}
    >
      <div className="w-full max-w-xl" onClick={(e) => e.stopPropagation()}>
        <DoubleBezelCard className="overflow-hidden p-0 shadow-2xl">
          {/* Search Input Bar */}
          <div className="flex items-center gap-3 border-b border-white/10 px-4 py-3">
            <Search className="h-5 w-5 text-muted-foreground" />
            <input
              className="flex-1 bg-transparent text-sm text-foreground outline-none placeholder:text-muted-foreground"
              onChange={(e) => {
                setQuery(e.target.value);
                setSelectedIndex(0);
              }}
              placeholder={t("commandPalette.placeholder")}
              ref={inputRef}
              value={query}
            />
            <kbd className="rounded border border-white/10 bg-white/5 px-1.5 py-0.5 font-mono text-[10px] text-muted-foreground">
              Esc
            </kbd>
          </div>

          {/* Results List */}
          <div className="max-h-[380px] overflow-y-auto p-2 scrollbar-thin">
            {filteredCommands.length === 0 ? (
              <div className="py-12 text-center text-xs text-muted-foreground">
                {t("commandPalette.noResults", { query })}
              </div>
            ) : (
              filteredCommands.map((item, idx) => {
                const Icon = item.icon;
                const isSelected = idx === selectedIndex;
                const categoryLabel = t(`commandPalette.categories.${item.categoryKey}`);

                return (
                  <div
                    className={`flex cursor-pointer items-center justify-between rounded-xl px-3 py-2.5 transition ${
                      isSelected ? "bg-violet-600 text-white shadow" : "text-foreground hover:bg-white/5"
                    }`}
                    key={item.id}
                    onClick={item.onSelect}
                    onMouseEnter={() => setSelectedIndex(idx)}
                  >
                    <div className="flex items-center gap-3 min-w-0">
                      <div
                        className={`flex h-8 w-8 shrink-0 items-center justify-center rounded-lg ${
                          isSelected ? "bg-white/20 text-white" : "bg-white/5 text-muted-foreground"
                        }`}
                      >
                        <Icon className="h-4 w-4" />
                      </div>
                      <div className="truncate">
                        <p className="text-xs font-semibold truncate">{item.title}</p>
                        {item.subtitle && (
                          <p
                            className={`text-[11px] truncate ${
                              isSelected ? "text-violet-200" : "text-muted-foreground"
                            }`}
                          >
                            {item.subtitle}
                          </p>
                        )}
                      </div>
                    </div>

                    <div className="flex items-center gap-2 pl-2 shrink-0">
                      <span
                        className={`rounded px-1.5 py-0.5 font-mono text-[9px] uppercase tracking-wider ${
                          isSelected ? "bg-white/20 text-white" : "bg-white/5 text-muted-foreground"
                        }`}
                      >
                        {categoryLabel}
                      </span>
                      {isSelected && <ArrowRight className="h-3.5 w-3.5" />}
                    </div>
                  </div>
                );
              })
            )}
          </div>

          {/* Footer Shortcuts */}
          <div className="flex items-center justify-between border-t border-white/10 bg-white/5 px-4 py-2 text-[10px] text-muted-foreground">
            <div className="flex items-center gap-3">
              <span>&uarr;&darr; {t("commandPalette.hints.navigate")}</span>
              <span>&crarr; {t("commandPalette.hints.select")}</span>
              <span>esc {t("commandPalette.hints.close")}</span>
            </div>
            <span>{t("commandPalette.hints.spotlight")}</span>
          </div>
        </DoubleBezelCard>
      </div>
    </div>
  );
}
