import { ChevronLeft, Languages, LayoutDashboard, ListTodo, Menu, NotebookPen, Search, Settings, Timer, WalletCards, X } from "lucide-react";
import { useEffect, useRef, useState } from "react";
import { useTranslation } from "react-i18next";
import { NavLink, Outlet, useLocation } from "react-router-dom";
import { CommandPalette } from "../components/design-system/command-palette";
import { LanguageSwitcher } from "../components/design-system/language-switcher";
import { PomodoroWidget } from "../components/design-system/pomodoro-widget";
import { ThemeSwitcher } from "../components/design-system/theme-switcher";
import { cn } from "../lib/cn";
import { useUIPreferencesStore } from "../stores/ui-preferences-store";

const navigation = [
  { to: "/app/today", key: "navigation.today", icon: LayoutDashboard },
  { to: "/app/tasks", key: "navigation.tasks", icon: ListTodo },
  { to: "/app/pomodoro", key: "pomodoro.title", icon: Timer },
  { to: "/app/learn", key: "navigation.learn", icon: Languages },
  { to: "/app/finance", key: "finance:title", icon: WalletCards },
  { to: "/app/journal", key: "navigation.journal", icon: NotebookPen },
  { to: "/app/settings", key: "navigation.settings", icon: Settings },
];

function ambientGradient(path: string) {
  if (path.includes("/finance")) return "bg-[radial-gradient(ellipse_60%_50%_at_10%_0%,rgba(16,185,129,0.2),transparent),radial-gradient(ellipse_50%_40%_at_90%_100%,rgba(20,184,166,0.15),transparent)]";
  if (path.includes("/learn")) return "bg-[radial-gradient(ellipse_60%_50%_at_10%_0%,rgba(139,92,246,0.2),transparent),radial-gradient(ellipse_50%_40%_at_90%_100%,rgba(168,85,247,0.15),transparent)]";
  if (path.includes("/journal")) return "bg-[radial-gradient(ellipse_60%_50%_at_10%_0%,rgba(245,158,11,0.18),transparent),radial-gradient(ellipse_50%_40%_at_90%_100%,rgba(217,119,6,0.13),transparent)]";
  if (path.includes("/tasks")) return "bg-[radial-gradient(ellipse_60%_50%_at_10%_0%,rgba(59,130,246,0.2),transparent),radial-gradient(ellipse_50%_40%_at_90%_100%,rgba(14,165,233,0.15),transparent)]";
  if (path.includes("/pomodoro")) return "bg-[radial-gradient(ellipse_60%_50%_at_10%_0%,rgba(99,102,241,0.2),transparent),radial-gradient(ellipse_50%_40%_at_90%_100%,rgba(139,92,246,0.15),transparent)]";
  return "bg-[radial-gradient(ellipse_60%_50%_at_10%_0%,rgba(96,165,250,0.18),transparent),radial-gradient(ellipse_50%_40%_at_90%_100%,rgba(196,181,253,0.15),transparent)]";
}

interface SidebarProps {
  collapsed?: boolean;
  onClose?: () => void;
  onToggleCollapse?: () => void;
}

function Sidebar({ collapsed = false, onClose, onToggleCollapse }: SidebarProps) {
  const { t } = useTranslation();
  const compact = collapsed && onToggleCollapse !== undefined;

  return (
    <>
      <div className="mb-6 flex min-h-10 items-center justify-between px-2 pt-1">
        <span className={cn("text-base font-bold tracking-tight", compact && "sr-only")}>GoPA</span>
        {onToggleCollapse ? (
          <button aria-label={t("controls.toggleSidebar")} className="icon-button" onClick={onToggleCollapse} type="button">
            <ChevronLeft className={cn("transition-transform duration-200", compact && "rotate-180")} size={18} />
          </button>
        ) : (
          <button aria-label={t("controls.closeNavigation")} className="icon-button sm:hidden" onClick={onClose} type="button">
            <X aria-hidden="true" size={18} />
          </button>
        )}
      </div>
      <nav aria-label={t("navigation.label")} className="space-y-1">
        {navigation.map(({ to, key, icon: Icon }) => (
          <NavLink
            className={({ isActive }) => cn("flex min-h-11 items-center gap-3 rounded-2xl px-3.5 text-sm font-medium text-muted-foreground transition-colors hover:bg-slate-200/60 hover:text-foreground focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary dark:hover:bg-slate-800/80", isActive && "bg-primary font-semibold text-primary-foreground shadow-lg shadow-primary/25 hover:bg-primary hover:text-primary-foreground", compact && "justify-center px-2")}
            key={to}
            onClick={onClose}
            to={to}
          >
            <Icon aria-hidden="true" size={18} />
            <span className={cn(compact && "sr-only")}>{t(key)}</span>
          </NavLink>
        ))}
      </nav>
    </>
  );
}

export function AppLayout() {
  const { t } = useTranslation();
  const location = useLocation();
  const collapsed = useUIPreferencesStore((state) => state.sidebarCollapsed);
  const toggleSidebar = useUIPreferencesStore((state) => state.toggleSidebar);
  const [mobileOpen, setMobileOpen] = useState(false);
  const [commandPaletteOpen, setCommandPaletteOpen] = useState(false);
  const mainRef = useRef<HTMLElement>(null);
  const mobileNavigationRef = useRef<HTMLElement>(null);

  useEffect(() => {
    setMobileOpen(false);
    mainRef.current?.focus();
  }, [location.pathname]);
  useEffect(() => {
    if (mobileNavigationRef.current) {
      mobileNavigationRef.current.inert = !mobileOpen;
    }
  }, [mobileOpen]);
  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") setMobileOpen(false);
      if ((event.ctrlKey || event.metaKey) && event.key.toLowerCase() === "k") {
        event.preventDefault();
        setCommandPaletteOpen((prev) => !prev);
      }
    };
    window.addEventListener("keydown", handleKeyDown);
    return () => window.removeEventListener("keydown", handleKeyDown);
  }, []);

  return (
    <div className="relative h-dvh overflow-hidden bg-background text-foreground transition-colors duration-300">
      <a className="sr-only fixed left-4 top-4 z-50 rounded-xl bg-primary px-4 py-3 font-semibold text-primary-foreground focus:not-sr-only" href="#main-content">
        {t("controls.skipToContent")}
      </a>
      <div aria-hidden="true" className="pointer-events-none fixed inset-0 -z-20" style={{ backgroundImage: "radial-gradient(circle, rgba(148,163,184,0.13) 1px, transparent 1px)", backgroundSize: "28px 28px" }} />
      <div className={cn("pointer-events-none fixed inset-0 -z-10 transition-all duration-700", ambientGradient(location.pathname))} />

      <aside className={cn("fixed inset-y-3 left-3 z-20 hidden flex-col rounded-3xl border border-white/60 bg-white/70 p-3 shadow-2xl shadow-slate-950/10 backdrop-blur-2xl transition-[width] duration-300 dark:border-white/10 dark:bg-slate-950/75 sm:flex", collapsed ? "w-20" : "w-60")}>
        <Sidebar collapsed={collapsed} onToggleCollapse={toggleSidebar} />
      </aside>

      <button aria-label={t("controls.openNavigation")} className={cn("fixed inset-0 z-30 bg-slate-950/40 transition-opacity sm:hidden", mobileOpen ? "opacity-100" : "pointer-events-none opacity-0")} onClick={() => setMobileOpen(false)} type="button" />
      <aside aria-hidden={!mobileOpen} aria-label={t("navigation.label")} className={cn("fixed inset-y-3 left-3 z-40 flex w-[min(20rem,calc(100vw-1.5rem))] flex-col rounded-3xl border border-white/60 bg-white/95 p-3 shadow-2xl shadow-slate-950/25 backdrop-blur-2xl transition-transform duration-300 dark:border-white/10 dark:bg-slate-950/95 sm:hidden", mobileOpen ? "translate-x-0" : "-translate-x-[calc(100%+1.5rem)]")} ref={mobileNavigationRef}>
        <Sidebar onClose={() => setMobileOpen(false)} />
      </aside>

      <div className={cn("flex h-dvh min-w-0 flex-col transition-[padding] duration-300 sm:pl-[16.5rem]", collapsed && "sm:pl-[6.5rem]")}>
        <header className="relative z-10 mx-3 mt-1.5 flex h-11 shrink-0 items-center justify-between gap-3 rounded-xl border border-white/60 bg-white/75 px-3.5 shadow-sm backdrop-blur-xl dark:border-white/10 dark:bg-slate-950/75 sm:px-5">
          <div className="flex min-w-0 items-center gap-3">
            <button aria-expanded={mobileOpen} aria-label={t("controls.openNavigation")} className="icon-button sm:hidden" onClick={() => setMobileOpen(true)} type="button"><Menu aria-hidden="true" size={19} /></button>
            <p className="truncate text-sm font-bold tracking-tight">GoPA</p>
            <button
              aria-label={t("controls.searchAria")}
              className="ml-2 flex min-h-11 items-center gap-2 rounded-xl border border-white/40 bg-white/50 px-2.5 py-1 text-xs text-muted-foreground transition hover:bg-white/80 hover:text-foreground dark:border-white/10 dark:bg-slate-900/60"
              onClick={() => setCommandPaletteOpen(true)}
              type="button"
            >
              <Search aria-hidden="true" size={13} />
              <span className="hidden sm:inline">{t("controls.search")}</span>
              <kbd className="rounded border border-white/20 bg-white/40 px-1 py-0.2 font-mono text-[10px] dark:bg-white/10">
                Ctrl K
              </kbd>
            </button>
          </div>
          <div className="flex shrink-0 items-center gap-1 sm:gap-3"><LanguageSwitcher /><ThemeSwitcher /></div>
        </header>
        <main className="flex-1 min-h-0 flex flex-col overflow-hidden px-3 pb-3 pt-3 focus:outline-none sm:px-4 sm:pb-4 sm:pt-3" id="main-content" ref={mainRef} tabIndex={-1}>
          <div className="h-full w-full mx-auto max-w-[1600px] flex flex-col min-h-0 overflow-hidden">
            <Outlet />
          </div>
        </main>
      </div>
      <PomodoroWidget />
      <CommandPalette isOpen={commandPaletteOpen} onClose={() => setCommandPaletteOpen(false)} />
    </div>
  );
}
