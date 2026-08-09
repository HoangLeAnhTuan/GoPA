import { create } from "zustand";
import { persist } from "zustand/middleware";

export type ThemePreference = "system" | "light" | "dark";
export type Locale = "vi" | "en" | "ja";

interface UIPreferencesState {
  theme: ThemePreference;
  locale: Locale;
  sidebarCollapsed: boolean;
  showFloatingPomodoro: boolean;
  setTheme: (theme: ThemePreference) => void;
  setLocale: (locale: Locale) => void;
  toggleSidebar: () => void;
  toggleFloatingPomodoro: () => void;
}

export const useUIPreferencesStore = create<UIPreferencesState>()(
  persist(
    (set) => ({
      theme: "system",
      locale: "vi",
      sidebarCollapsed: false,
      showFloatingPomodoro: false, // Default minimized / hidden as requested by user
      setTheme: (theme) => set({ theme }),
      setLocale: (locale) => set({ locale }),
      toggleSidebar: () => set((state) => ({ sidebarCollapsed: !state.sidebarCollapsed })),
      toggleFloatingPomodoro: () => set((state) => ({ showFloatingPomodoro: !state.showFloatingPomodoro })),
    }),
    { name: "gopa-ui-preferences" },
  ),
);
