import { Outlet } from "react-router-dom";
import { LanguageSwitcher } from "../components/design-system/language-switcher";
import { ThemeSwitcher } from "../components/design-system/theme-switcher";

export function AuthLayout() {
  return (
    <main className="grid min-h-[100dvh] place-items-center bg-background p-5 text-foreground">
      <div className="pointer-events-none fixed inset-0 -z-10 bg-[radial-gradient(circle_at_20%_20%,rgba(96,165,250,0.18),transparent_28%),radial-gradient(circle_at_80%_80%,rgba(196,181,253,0.2),transparent_30%)]" />
      <div className="absolute right-5 top-5 flex items-center gap-1"><LanguageSwitcher /><ThemeSwitcher /></div>
      <section className="w-full max-w-md rounded-3xl border border-white/50 bg-white/75 p-8 shadow-xl shadow-slate-950/5 backdrop-blur-xl dark:border-white/10 dark:bg-slate-950/75"><Outlet /></section>
    </main>
  );
}
