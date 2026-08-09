import { useEffect, type PropsWithChildren } from "react";
import { useTranslation } from "react-i18next";
import { useUIPreferencesStore } from "../../stores/ui-preferences-store";

export function ThemeController({ children }: PropsWithChildren) {
  const { i18n } = useTranslation();
  const theme = useUIPreferencesStore((state) => state.theme);
  const locale = useUIPreferencesStore((state) => state.locale);

  useEffect(() => {
    const prefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
    document.documentElement.classList.toggle("dark", theme === "dark" || (theme === "system" && prefersDark));
  }, [theme]);

  useEffect(() => {
    void i18n.changeLanguage(locale);
    document.documentElement.lang = locale;
  }, [i18n, locale]);

  return <>{children}</>;
}
