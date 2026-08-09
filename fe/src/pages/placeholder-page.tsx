import { motion, useReducedMotion } from "framer-motion";
import { useTranslation } from "react-i18next";
import { GlassPanel } from "../components/design-system/glass-panel";
import { PageHeader } from "../components/design-system/page-header";

interface PlaceholderPageProps {
  titleKey: string;
  descriptionKey: string;
}

export function PlaceholderPage({ titleKey, descriptionKey }: PlaceholderPageProps) {
  const { t } = useTranslation();
  const reducedMotion = useReducedMotion();
  return (
    <motion.div animate={{ opacity: 1, y: 0 }} initial={{ opacity: 0, y: reducedMotion ? 0 : 6 }} transition={{ duration: reducedMotion ? 0 : 0.18 }}>
      <PageHeader description={t(descriptionKey)} title={t(titleKey)} />
      <GlassPanel className="mt-7 max-w-3xl p-6 sm:p-8">
        <p className="text-sm font-medium text-primary">{t("bootstrap.eyebrow")}</p>
        <p className="mt-2 text-lg font-medium tracking-tight">{t("bootstrap.title")}</p>
        <p className="mt-2 leading-relaxed text-muted-foreground">{t("bootstrap.description")}</p>
      </GlassPanel>
    </motion.div>
  );
}
