import { Globe } from "lucide-react";
import { useTranslation } from "react-i18next";
import type { Locale } from "../../stores/ui-preferences-store";
import { useUIPreferencesStore } from "../../stores/ui-preferences-store";
import { MacSelect } from "./mac-select";

export function LanguageSwitcher() {
  const { t } = useTranslation();
  const locale = useUIPreferencesStore((state) => state.locale);
  const setLocale = useUIPreferencesStore((state) => state.setLocale);

  const languageOptions = [
    { value: "vi", label: "Tiếng Việt", icon: Globe },
    { value: "en", label: "English", icon: Globe },
    { value: "ja", label: "日本語", icon: Globe },
  ];

  return (
    <MacSelect
      aria-label={t("controls.language")}
      options={languageOptions}
      triggerClassName="min-h-10 border-none bg-transparent hover:bg-slate-200/50 dark:hover:bg-slate-800/60 shadow-none px-2.5"
      onChange={(val) => setLocale(val as Locale)}
      value={locale}
    />
  );
}

