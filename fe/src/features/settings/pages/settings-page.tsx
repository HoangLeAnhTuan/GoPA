import React, { useState } from "react";
import { useTranslation } from "react-i18next";
import {
  User as UserIcon,
  Palette,
  Keyboard,
  Save,
  Check,
  Languages,
  Shield,
  Moon,
  Sun,
} from "lucide-react";
import { DoubleBezelCard } from "../../../components/design-system/double-bezel-card";
import { ButtonInButton } from "../../../components/design-system/button-in-button";
import { useAuth } from "../../auth/hooks/use-auth";
import { useUpdateProfile } from "../../auth/hooks/use-update-profile";
import type { Locale } from "../../../stores/ui-preferences-store";
import { useUIPreferencesStore } from "../../../stores/ui-preferences-store";
import { cn } from "../../../lib/cn";

export function SettingsPage() {
  const { t, i18n } = useTranslation();
  const { user } = useAuth();
  const updateProfileMutation = useUpdateProfile();

  const currentTheme = useUIPreferencesStore((state) => state.theme);
  const setTheme = useUIPreferencesStore((state) => state.setTheme);
  const setLocale = useUIPreferencesStore((state) => state.setLocale);

  const [displayName, setDisplayName] = useState(user?.displayName ?? "");
  const [avatarUrl, setAvatarUrl] = useState(user?.avatarUrl ?? "");
  const [savedSuccess, setSavedSuccess] = useState(false);

  const handleSaveProfile = (e: React.FormEvent) => {
    e.preventDefault();
    updateProfileMutation.mutate(
      {
        display_name: displayName.trim(),
        avatar_url: avatarUrl.trim() ? avatarUrl.trim() : undefined,
      },
      {
        onSuccess: () => {
          setSavedSuccess(true);
          setTimeout(() => setSavedSuccess(false), 3000);
        },
      }
    );
  };

  const handleLanguageChange = (lang: string) => {
    setLocale(lang as Locale);
    void i18n.changeLanguage(lang);
  };

  const shortcuts = [
    { key: "Space", desc: t("settings.shortcuts.pomodoroDesc") },
    { key: "1 / 2 / 3 / 4", desc: t("settings.shortcuts.sm2Desc") },
    { key: "Ctrl + K", desc: t("settings.shortcuts.commandPaletteDesc") },
    { key: "Esc", desc: t("settings.shortcuts.escDesc") },
  ];

  return (
    <div className="flex flex-col gap-6 pb-12 animate-in fade-in duration-500 max-w-5xl mx-auto w-full">
      {/* Page Header */}
      <div>
        <h1 className="text-2xl sm:text-3xl font-bold tracking-tight text-foreground">
          {t("navigation.settings")}
        </h1>
        <p className="mt-1 text-sm text-muted-foreground">
          {t("settings.subtitle")}
        </p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-12 gap-6">
        {/* Left Column: Profile Card (col-span-7) */}
        <div className="md:col-span-7 space-y-6">
          <DoubleBezelCard glowColor="rgba(99, 102, 241, 0.1)">
            <div className="flex items-center justify-between border-b border-white/10 pb-4">
              <div className="flex items-center gap-3">
                <div className="size-10 rounded-xl bg-blue-500/10 border border-blue-500/20 flex items-center justify-center text-blue-400">
                  <UserIcon size={20} />
                </div>
                <div>
                  <h2 className="text-base font-semibold text-foreground">{t("settings.profileTitle")}</h2>
                  <p className="text-xs text-muted-foreground">{t("settings.profileSubtitle")}</p>
                </div>
              </div>
              <span className="inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium bg-white/5 border border-white/10 text-muted-foreground">
                <Shield size={12} className="text-emerald-400" />
                {user?.role ?? "USER"}
              </span>
            </div>

            <form onSubmit={handleSaveProfile} className="mt-5 space-y-4">
              <div>
                <label className="block text-xs font-semibold uppercase tracking-wider text-muted-foreground mb-1.5">
                  {t("settings.emailLabel")}
                </label>
                <input
                  type="email"
                  disabled
                  value={user?.email ?? ""}
                  className="w-full rounded-xl border border-white/10 bg-white/5 px-3.5 py-2 text-sm text-muted-foreground cursor-not-allowed"
                />
                <span className="text-[11px] text-muted-foreground mt-1 block">
                  {t("settings.emailHint")}
                </span>
              </div>

              <div>
                <label className="block text-xs font-semibold uppercase tracking-wider text-muted-foreground mb-1.5">
                  {t("settings.displayNameLabel")}
                </label>
                <input
                  type="text"
                  required
                  value={displayName}
                  onChange={(e) => setDisplayName(e.target.value)}
                  placeholder={t("settings.displayNamePlaceholder")}
                  className="w-full rounded-xl border border-white/10 bg-white/5 px-3.5 py-2 text-sm text-foreground focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition-colors"
                />
              </div>

              <div>
                <label className="block text-xs font-semibold uppercase tracking-wider text-muted-foreground mb-1.5">
                  {t("settings.avatarLabel")}
                </label>
                <input
                  type="url"
                  value={avatarUrl}
                  onChange={(e) => setAvatarUrl(e.target.value)}
                  placeholder="https://example.com/avatar.jpg"
                  className="w-full rounded-xl border border-white/10 bg-white/5 px-3.5 py-2 text-sm text-foreground focus:border-blue-500 focus:ring-1 focus:ring-blue-500 outline-none transition-colors"
                />
              </div>

              <div className="pt-2 flex items-center justify-between">
                {savedSuccess && (
                  <span className="inline-flex items-center gap-1.5 text-xs font-medium text-emerald-400">
                    <Check size={14} />
                    {t("settings.saveSuccess")}
                  </span>
                )}
                {!savedSuccess && <div />}

                <ButtonInButton
                  type="submit"
                  icon={Save}
                  variant="primary"
                  size="md"
                  disabled={updateProfileMutation.isPending}
                >
                  {updateProfileMutation.isPending ? t("settings.saving") : t("settings.save")}
                </ButtonInButton>
              </div>
            </form>
          </DoubleBezelCard>

          {/* Keyboard Shortcuts Card */}
          <DoubleBezelCard glowColor="rgba(245, 158, 11, 0.08)">
            <div className="flex items-center gap-3 border-b border-white/10 pb-4">
              <div className="size-10 rounded-xl bg-amber-500/10 border border-amber-500/20 flex items-center justify-center text-amber-400">
                <Keyboard size={20} />
              </div>
              <div>
                <h2 className="text-base font-semibold text-foreground">{t("settings.shortcutsTitle")}</h2>
                <p className="text-xs text-muted-foreground">{t("settings.shortcutsSubtitle")}</p>
              </div>
            </div>

            <div className="mt-4 divide-y divide-white/5">
              {shortcuts.map((sc) => (
                <div key={sc.key} className="py-2.5 flex items-center justify-between gap-4">
                  <span className="text-xs text-muted-foreground">{sc.desc}</span>
                  <kbd className="px-2.5 py-1 rounded-md bg-white/10 border border-white/10 text-xs font-mono font-semibold text-foreground shrink-0">
                    {sc.key}
                  </kbd>
                </div>
              ))}
            </div>
          </DoubleBezelCard>
        </div>

        {/* Right Column: Preferences (col-span-5) */}
        <div className="md:col-span-5 space-y-6">
          {/* Appearance Card */}
          <DoubleBezelCard glowColor="rgba(168, 85, 247, 0.1)">
            <div className="flex items-center gap-3 border-b border-white/10 pb-4">
              <div className="size-10 rounded-xl bg-violet-500/10 border border-violet-500/20 flex items-center justify-center text-violet-400">
                <Palette size={20} />
              </div>
              <div>
                <h2 className="text-base font-semibold text-foreground">{t("settings.appearanceTitle")}</h2>
                <p className="text-xs text-muted-foreground">{t("settings.appearanceSubtitle")}</p>
              </div>
            </div>

            <div className="mt-5 grid grid-cols-2 gap-3">
              <button
                type="button"
                onClick={() => setTheme("dark")}
                className={cn(
                  "p-3 rounded-xl border flex flex-col items-center gap-2 transition-all cursor-pointer",
                  currentTheme === "dark"
                    ? "border-primary bg-primary/10 text-foreground ring-2 ring-primary/20"
                    : "border-white/10 bg-white/5 text-muted-foreground hover:bg-white/10"
                )}
              >
                <Moon size={20} className={currentTheme === "dark" ? "text-primary" : ""} />
                <span className="text-xs font-semibold">{t("settings.themeDark")}</span>
              </button>

              <button
                type="button"
                onClick={() => setTheme("light")}
                className={cn(
                  "p-3 rounded-xl border flex flex-col items-center gap-2 transition-all cursor-pointer",
                  currentTheme === "light"
                    ? "border-primary bg-primary/10 text-foreground ring-2 ring-primary/20"
                    : "border-white/10 bg-white/5 text-muted-foreground hover:bg-white/10"
                )}
              >
                <Sun size={20} className={currentTheme === "light" ? "text-primary" : ""} />
                <span className="text-xs font-semibold">{t("settings.themeLight")}</span>
              </button>
            </div>
          </DoubleBezelCard>

          {/* Language Selection Card */}
          <DoubleBezelCard glowColor="rgba(16, 185, 129, 0.1)">
            <div className="flex items-center gap-3 border-b border-white/10 pb-4">
              <div className="size-10 rounded-xl bg-emerald-500/10 border border-emerald-500/20 flex items-center justify-center text-emerald-400">
                <Languages size={20} />
              </div>
              <div>
                <h2 className="text-base font-semibold text-foreground">{t("settings.languageTitle")}</h2>
                <p className="text-xs text-muted-foreground">{t("settings.languageSubtitle")}</p>
              </div>
            </div>

            <div className="mt-5 space-y-2">
              {[
                { code: "vi", label: "Tiếng Việt", flag: "🇻🇳", desc: t("settings.langViDesc") },
                { code: "en", label: "English", flag: "🇬🇧", desc: t("settings.langEnDesc") },
                { code: "ja", label: "日本語", flag: "🇯🇵", desc: t("settings.langJaDesc") },
              ].map((lang) => (
                <button
                  key={lang.code}
                  type="button"
                  onClick={() => handleLanguageChange(lang.code)}
                  className={cn(
                    "w-full p-3 rounded-xl border flex items-center justify-between text-left transition-all cursor-pointer",
                    i18n.language === lang.code
                      ? "border-emerald-500/50 bg-emerald-500/10 text-foreground"
                      : "border-white/10 bg-white/5 text-muted-foreground hover:bg-white/10"
                  )}
                >
                  <div className="flex items-center gap-2.5">
                    <span className="text-base">{lang.flag}</span>
                    <div>
                      <p className="text-xs font-semibold text-foreground">{lang.label}</p>
                      <p className="text-[11px] text-muted-foreground">{lang.desc}</p>
                    </div>
                  </div>
                  {i18n.language === lang.code && (
                    <Check size={16} className="text-emerald-400 shrink-0" />
                  )}
                </button>
              ))}
            </div>
          </DoubleBezelCard>
        </div>
      </div>
    </div>
  );
}
