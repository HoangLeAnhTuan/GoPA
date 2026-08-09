import { useTranslation } from "react-i18next";

interface ErrorStateProps { onRetry?: () => void; }
export function ErrorState({ onRetry }: ErrorStateProps) { const { t } = useTranslation("feedback"); return <div role="alert" className="rounded-2xl border border-red-200 bg-red-50 p-5 text-red-800 dark:border-red-900 dark:bg-red-950/30 dark:text-red-200"><p className="font-semibold">{t("error.title")}</p><p className="mt-1 text-sm">{t("error.description")}</p>{onRetry === undefined ? null : <button className="mt-3 min-h-11 rounded-xl border border-current px-3 text-sm font-semibold" onClick={onRetry} type="button">{t("error.retry")}</button>}</div>; }
