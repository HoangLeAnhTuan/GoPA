import { Link } from "react-router-dom";
import { useTranslation } from "react-i18next";

export function NotFoundPage() {
  const { t } = useTranslation();
  return <main className="grid min-h-[100dvh] place-items-center bg-background p-6 text-center"><div><h1 className="text-3xl font-semibold">{t("notFound.title")}</h1><p className="mt-2 text-muted-foreground">{t("notFound.description")}</p><Link className="mt-5 inline-flex min-h-11 items-center rounded-xl bg-primary px-4 py-2 text-sm font-medium text-primary-foreground" to="/app/today">{t("notFound.action")}</Link></div></main>;
}
