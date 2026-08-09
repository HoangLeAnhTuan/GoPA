import { motion } from "framer-motion";
import { Moon, Sun } from "lucide-react";
import { useTranslation } from "react-i18next";
import { useUIPreferencesStore } from "../../stores/ui-preferences-store";

export function ThemeSwitcher() {
  const { t } = useTranslation();
  const theme = useUIPreferencesStore((state) => state.theme);
  const setTheme = useUIPreferencesStore((state) => state.setTheme);
  const isDark = theme === "dark";

  return (
    <button
      aria-label={t("controls.theme")}
      className="relative flex h-9 w-16 items-center rounded-full border border-white/40 bg-slate-200/80 p-1 shadow-inner backdrop-blur-xl transition dark:border-white/10 dark:bg-slate-900/80"
      onClick={() => setTheme(isDark ? "light" : "dark")}
      type="button"
    >
      <motion.div
        animate={{ x: isDark ? 28 : 0 }}
        className="flex size-7 items-center justify-center rounded-full bg-white shadow-md text-amber-500 dark:bg-slate-800 dark:text-cyan-400"
        transition={{ type: "spring", stiffness: 500, damping: 30 }}
      >
        {isDark ? <Moon aria-hidden="true" size={14} /> : <Sun aria-hidden="true" size={14} />}
      </motion.div>
    </button>
  );
}

