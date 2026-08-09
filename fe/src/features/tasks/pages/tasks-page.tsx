import { Check, ChevronRight, Plus, Trash2 } from "lucide-react";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { ErrorState } from "../../../components/feedback/error-state";
import { LoadingState } from "../../../components/feedback/loading-state";
import { GlassPanel } from "../../../components/design-system/glass-panel";
import { PageHeader } from "../../../components/design-system/page-header";
import { MacSelect } from "../../../components/design-system/mac-select";
import { useCreateTask, useDeleteTask, useTasks, useUpdateTask } from "../hooks/use-tasks";
import type { Task, TaskInput, TaskStatus } from "../types";

const columns: TaskStatus[] = ["TODO", "IN_PROGRESS", "DONE"];
const defaults: TaskInput = { title: "", description: "", status: "TODO", priority: "MEDIUM", category: "WORK", dueDate: null };

export function TasksPage() {
  const { t } = useTranslation("tasks");
  const tasks = useTasks();
  const create = useCreateTask();
  const update = useUpdateTask();
  const remove = useDeleteTask();
  const [input, setInput] = useState(defaults);

  if (tasks.isPending) return <LoadingState />;
  if (tasks.isError) return <ErrorState />;

  const submit = async (event: React.FormEvent) => {
    event.preventDefault();
    if (input.title.trim() === "") return;
    await create.mutateAsync(input);
    setInput(defaults);
  };

  const move = (task: Task) => {
    const index = columns.indexOf(task.status);
    if (index === columns.length - 1) return;
    const status = columns[index + 1];
    if (status === undefined) return;
    update.mutate({
      id: task.id,
      input: {
        title: task.title,
        description: task.description,
        status,
        priority: task.priority,
        category: task.category,
        dueDate: task.dueDate?.toISOString().slice(0, 10) ?? null,
      },
    });
  };

  const priorityOptions = [
    { value: "LOW", label: t("priority.low") },
    { value: "MEDIUM", label: t("priority.medium") },
    { value: "HIGH", label: t("priority.high") },
  ];

  const categoryOptions = [
    { value: "WORK", label: t("category.work") },
    { value: "STUDY", label: t("category.study") },
    { value: "LIFE", label: t("category.life") },
  ];

  return (
    <div>
      <PageHeader description={t("description")} title={t("title")} />
      <GlassPanel className="mt-6 p-4">
        <form className="grid gap-3 md:grid-cols-[1fr_auto_auto_auto]" onSubmit={submit}>
          <label className="sr-only" htmlFor="task-title">
            {t("form.title")}
          </label>
          <input
            className="h-11 rounded-xl border border-white/40 bg-white/70 px-3.5 text-sm outline-none transition focus:border-primary dark:border-white/10 dark:bg-slate-900/70"
            id="task-title"
            onChange={(event) => setInput({ ...input, title: event.target.value })}
            placeholder={t("form.placeholder")}
            value={input.title}
          />
          <MacSelect
            onChange={(val) => setInput({ ...input, priority: val as TaskInput["priority"] })}
            options={priorityOptions}
            value={input.priority}
          />
          <MacSelect
            onChange={(val) => setInput({ ...input, category: val as TaskInput["category"] })}
            options={categoryOptions}
            value={input.category}
          />
          <button
            className="inline-flex h-11 items-center justify-center gap-2 rounded-xl bg-primary px-5 text-sm font-semibold text-primary-foreground shadow-lg shadow-primary/25 transition hover:opacity-90 disabled:opacity-60"
            disabled={create.isPending}
            type="submit"
          >
            <Plus size={17} />
            {t("form.add")}
          </button>
        </form>
      </GlassPanel>
      <section className="mt-6 grid gap-4 xl:grid-cols-3">
        {columns.map((status) => (
          <GlassPanel className="min-h-72 p-4" key={status}>
            <h2 className="mb-4 text-sm font-semibold">{t(`status.${status}`)}</h2>
            <div className="space-y-3">
              {tasks.data
                .filter((task) => task.status === status)
                .map((task) => (
                  <article className="rounded-xl border border-white/40 bg-white/80 p-3.5 shadow-sm dark:border-white/10 dark:bg-slate-900/80" key={task.id}>
                    <div className="flex gap-2">
                      <p className="flex-1 font-medium">{task.title}</p>
                      <button
                        aria-label={t("actions.delete")}
                        className="text-muted-foreground transition hover:text-red-500"
                        onClick={() => remove.mutate(task.id)}
                        type="button"
                      >
                        <Trash2 size={16} />
                      </button>
                    </div>
                    <p className="mt-2 text-xs text-muted-foreground">
                      {t(`priority.${task.priority.toLowerCase()}`)} · {t(`category.${task.category.toLowerCase()}`)}
                    </p>
                    {status === "DONE" ? (
                      <p className="mt-3 inline-flex items-center gap-1 text-xs font-semibold text-emerald-500">
                        <Check size={14} />
                        {t("complete")}
                      </p>
                    ) : (
                      <button
                        className="mt-3 inline-flex items-center gap-1 text-xs font-medium text-primary transition hover:underline"
                        onClick={() => move(task)}
                        type="button"
                      >
                        {t("actions.move")}
                        <ChevronRight size={14} />
                      </button>
                    )}
                  </article>
                ))}
              {tasks.data.filter((task) => task.status === status).length === 0 ? (
                <p className="text-sm text-muted-foreground">{t("empty")}</p>
              ) : null}
            </div>
          </GlassPanel>
        ))}
      </section>
    </div>
  );
}

