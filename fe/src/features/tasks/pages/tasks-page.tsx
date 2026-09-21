import React, { useState } from "react";
import {
  AlertCircle,
  ChevronRight,
  ChevronLeft,
  Plus,
  Trash2,
  Code2,
  Briefcase,
  GraduationCap,
  Heart,
  Timer,
} from "lucide-react";
import { useTranslation } from "react-i18next";
import { ErrorState } from "../../../components/feedback/error-state";
import { LoadingState } from "../../../components/feedback/loading-state";
import { DoubleBezelCard } from "../../../components/design-system/double-bezel-card";
import { ButtonInButton } from "../../../components/design-system/button-in-button";
import { PageHeader } from "../../../components/design-system/page-header";
import { MacSelect } from "../../../components/design-system/mac-select";
import {
  useCreateTask,
  useDeleteTask,
  useTasks,
  useUpdateTaskStatus,
} from "../hooks/use-tasks";
import type { Task, TaskCategory, TaskInput, TaskPriority, TaskStatus } from "../types";
import { cn } from "../../../lib/cn";
import { toApiError } from "../../../lib/api-client";

const columns: TaskStatus[] = ["TODO", "IN_PROGRESS", "DONE"];

const columnBadgeStyles: Record<TaskStatus, string> = {
  TODO: "bg-blue-500/10 text-blue-400 border-blue-500/20",
  IN_PROGRESS: "bg-amber-500/10 text-amber-400 border-amber-500/20",
  DONE: "bg-emerald-500/10 text-emerald-400 border-emerald-500/20",
};

const defaults: TaskInput = {
  title: "",
  description: "",
  status: "TODO",
  priority: "MEDIUM",
  category: "WORK",
  dueDate: null,
  estimated_pomodoros: 2,
};

export function TasksPage() {
  const { t } = useTranslation("tasks");
  const tasks = useTasks();
  const create = useCreateTask();
  const updateStatus = useUpdateTaskStatus();
  const remove = useDeleteTask();
  const [input, setInput] = useState<TaskInput>(defaults);
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  if (tasks.isPending) return <LoadingState />;
  if (tasks.isError) return <ErrorState onRetry={() => void tasks.refetch()} />;

  const submit = async (event: React.FormEvent) => {
    event.preventDefault();
    if (input.title.trim() === "") {
      setErrorMessage(t("form.titleRequired"));
      return;
    }
    setErrorMessage(null);
    try {
      await create.mutateAsync(input);
      setInput(defaults);
    } catch (error: unknown) {
      const apiMsg = toApiError(error).message;
      if (!apiMsg.toLowerCase().includes("status code")) {
        setErrorMessage(apiMsg);
      } else {
        setErrorMessage(t("form.createFailed"));
      }
    }
  };

  const moveNext = (task: Task) => {
    const currentIndex = columns.indexOf(task.status);
    if (currentIndex < columns.length - 1) {
      const nextStatus = columns[currentIndex + 1];
      if (nextStatus) {
        updateStatus.mutate({ id: task.id, status: nextStatus });
      }
    }
  };

  const movePrev = (task: Task) => {
    const currentIndex = columns.indexOf(task.status);
    if (currentIndex > 0) {
      const prevStatus = columns[currentIndex - 1];
      if (prevStatus) {
        updateStatus.mutate({ id: task.id, status: prevStatus });
      }
    }
  };

  const priorityOptions = [
    { value: "LOW", label: t("priorityOptions.LOW") },
    { value: "MEDIUM", label: t("priorityOptions.MEDIUM") },
    { value: "HIGH", label: t("priorityOptions.HIGH") },
    { value: "URGENT", label: t("priorityOptions.URGENT") },
  ];

  const categoryOptions = [
    { value: "WORK", label: t("categoryOptions.WORK") },
    { value: "STUDY", label: t("categoryOptions.STUDY") },
    { value: "LEETCODE", label: t("categoryOptions.LEETCODE") },
    { value: "LIFE", label: t("categoryOptions.LIFE") },
  ];

  const getCategoryIcon = (category: TaskCategory) => {
    switch (category) {
      case "WORK":
        return <Briefcase size={12} />;
      case "STUDY":
        return <GraduationCap size={12} />;
      case "LEETCODE":
        return <Code2 size={12} className="text-amber-400" />;
      case "LIFE":
        return <Heart size={12} />;
    }
  };

  return (
    <div className="flex flex-col gap-6 pb-12 animate-in fade-in duration-500">
      <PageHeader
        title={t("headerTitle")}
        description={t("headerDescription")}
      />

      {/* Quick Add Bar with Double-Bezel */}
      <DoubleBezelCard glowColor="rgba(59, 130, 246, 0.1)">
        <form className="flex flex-col gap-2" onSubmit={submit}>
          <div className="grid gap-3 lg:grid-cols-[1fr_auto_auto_auto_auto]">
            <input
              aria-describedby={errorMessage ? "task-title-error" : undefined}
              aria-invalid={errorMessage !== null}
              aria-label={t("form.title")}
              className={cn(
                "h-11 rounded-xl border px-4 text-sm text-foreground placeholder:text-muted-foreground outline-none transition",
                errorMessage
                  ? "border-rose-500 bg-rose-50/30 focus:ring-2 focus:ring-rose-500/30 dark:border-rose-500 dark:bg-rose-950/20"
                  : "border-slate-200 bg-white focus:border-blue-500 focus:ring-1 focus:ring-blue-500 dark:border-white/10 dark:bg-white/5"
              )}
              id="task-title"
              onChange={(event) => {
                setInput({ ...input, title: event.target.value });
                if (errorMessage) setErrorMessage(null);
              }}
              placeholder={t("form.quickAddPlaceholder")}
              value={input.title}
            />
            <MacSelect
              aria-label={t("form.priorityLabel")}
              onChange={(val) => setInput({ ...input, priority: val as TaskPriority })}
              options={priorityOptions}
              value={input.priority}
            />
            <MacSelect
              aria-label={t("form.categoryLabel")}
              onChange={(val) => setInput({ ...input, category: val as TaskCategory })}
              options={categoryOptions}
              value={input.category}
            />
            <div className="flex items-center gap-2 px-3 rounded-xl border border-slate-200 bg-white dark:border-white/10 dark:bg-white/5">
              <Timer size={15} className="text-rose-400 shrink-0" />
              <span className="text-xs text-muted-foreground">{t("form.pomodoro")}</span>
              <input
                type="number"
                min={1}
                max={20}
                value={input.estimated_pomodoros ?? 2}
                onChange={(e) => setInput({ ...input, estimated_pomodoros: Number(e.target.value) })}
                className="w-12 bg-transparent text-xs font-semibold text-foreground text-center outline-none"
              />
            </div>
            <ButtonInButton
              icon={Plus}
              variant="primary"
              size="md"
              type="submit"
              disabled={create.isPending}
            >
              {t("form.add")}
            </ButtonInButton>
          </div>
          {errorMessage && (
            <div className="flex items-center gap-1.5 text-xs font-semibold text-rose-500 dark:text-rose-400 px-1 pt-1 animate-in fade-in slide-in-from-top-1" id="task-title-error" role="alert">
              <AlertCircle size={14} className="shrink-0" />
              <span>{errorMessage}</span>
            </div>
          )}
        </form>
      </DoubleBezelCard>

      {/* Kanban Board Columns */}
      <div className="grid gap-5 lg:grid-cols-3">
        {columns.map((colStatus) => {
          const colTasks = (tasks.data ?? []).filter((task) => task.status === colStatus);
          const label = t(`columns.${colStatus}`);
          const badgeColor = columnBadgeStyles[colStatus];

          return (
            <div
              key={colStatus}
              className="flex flex-col gap-3 rounded-[1.5rem] border border-slate-200/80 bg-slate-50/50 p-4 backdrop-blur-md dark:border-white/10 dark:bg-white/[0.02]"
            >
              {/* Column Header */}
              <div className="flex items-center justify-between px-2 pb-2 border-b border-slate-200/60 dark:border-white/5">
                <span className="text-sm font-bold tracking-tight text-foreground">
                  {label}
                </span>
                <span
                  className={cn(
                    "text-xs font-mono font-semibold px-2 py-0.5 rounded-full border",
                    badgeColor
                  )}
                >
                  {colTasks.length}
                </span>
              </div>

              {/* Cards list */}
              <div className="flex flex-col gap-3 min-h-[360px]">
                {colTasks.length === 0 ? (
                  <div className="flex flex-1 flex-col items-center justify-center p-6 text-center text-xs text-muted-foreground border border-dashed border-slate-200 dark:border-white/10 rounded-2xl">
                    {t("columnEmpty")}
                  </div>
                ) : (
                  colTasks.map((task) => {
                    const isUrgent = task.priority === "URGENT";
                    const isHigh = task.priority === "HIGH";

                    return (
                      <div
                        key={task.id}
                        className={cn(
                          "group relative rounded-2xl p-4 transition-all duration-300",
                          "border border-slate-200/80 bg-white/90 shadow-sm shadow-slate-900/5 hover:border-slate-300 hover:bg-white dark:border-white/10 dark:bg-slate-900/60 dark:shadow-lg dark:shadow-black/20 dark:hover:border-white/20 dark:hover:bg-slate-900/80",
                          isUrgent && "border-rose-500/40 shadow-rose-950/20"
                        )}
                      >
                        {/* Header tag row */}
                        <div className="flex items-center justify-between gap-2 mb-2">
                          <div className="flex items-center gap-1.5">
                            <span
                              className={cn(
                                "size-2 rounded-full",
                                isUrgent
                                  ? "bg-rose-500 animate-ping"
                                  : isHigh
                                  ? "bg-orange-500"
                                  : task.priority === "MEDIUM"
                                  ? "bg-amber-500"
                                  : "bg-blue-500"
                              )}
                            />
                            <span
                              className={cn(
                                "text-[10px] font-bold uppercase tracking-wider",
                                isUrgent
                                  ? "text-rose-400"
                                  : isHigh
                                  ? "text-orange-400"
                                  : task.priority === "MEDIUM"
                                  ? "text-amber-400"
                                  : "text-blue-400"
                              )}
                            >
                              {t(`priorityOptions.${task.priority}`)}
                            </span>
                          </div>

                          <div className="flex items-center gap-2">
                            <span className="inline-flex items-center gap-1 px-2 py-0.5 rounded-md bg-slate-100 text-slate-600 border border-slate-200 dark:bg-white/5 dark:border-white/5 dark:text-muted-foreground text-[11px] font-medium">
                              {getCategoryIcon(task.category)}
                              <span>{t(`categoryOptions.${task.category}`)}</span>
                            </span>
                          </div>
                        </div>

                        {/* Title */}
                        <h3 className="text-sm font-semibold text-foreground leading-snug">
                          {task.title}
                        </h3>

                        {/* Description if present */}
                        {task.description && (
                          <p className="mt-1 text-xs text-muted-foreground line-clamp-2">
                            {task.description}
                          </p>
                        )}

                        {/* Footer info: Pomodoro counter + Actions */}
                        <div className="mt-4 pt-3 border-t border-slate-200/60 dark:border-white/5 flex items-center justify-between text-xs text-muted-foreground">
                          <div className="flex items-center gap-1 font-mono text-[11px] text-rose-500 dark:text-rose-400/90 font-medium">
                            <Timer size={13} />
                            <span>
                              {task.completedPomodoros} / {task.estimatedPomodoros || 2} pmd
                            </span>
                          </div>

                          <div className="flex items-center gap-1">
                            {colStatus !== "TODO" && (
                              <button
                                type="button"
                                onClick={() => movePrev(task)}
                                className="size-7 rounded-lg bg-slate-100 hover:bg-slate-200 text-muted-foreground hover:text-foreground dark:bg-white/5 dark:hover:bg-white/10 flex items-center justify-center transition-colors cursor-pointer"
                                title={t("actions.prevStatus")}
                              >
                                <ChevronLeft size={14} />
                              </button>
                            )}

                            {colStatus !== "DONE" && (
                              <button
                                type="button"
                                onClick={() => moveNext(task)}
                                className="size-7 rounded-lg bg-blue-500/10 hover:bg-blue-500/20 text-blue-400 flex items-center justify-center transition-colors cursor-pointer"
                                title={t("actions.nextStatus")}
                              >
                                <ChevronRight size={14} />
                              </button>
                            )}

                            <button
                              type="button"
                              onClick={() => remove.mutate(task.id)}
                              className="size-7 rounded-lg hover:bg-rose-500/10 hover:text-rose-400 flex items-center justify-center text-muted-foreground transition-colors cursor-pointer ml-1"
                              title={t("actions.delete")}
                            >
                              <Trash2 size={13} />
                            </button>
                          </div>
                        </div>
                      </div>
                    );
                  })
                )}
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
