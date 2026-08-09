import type { JournalMood } from "./types";

export interface MoodConfig {
  value: JournalMood;
  translationKey: string;
  badgeStyle: string;
  buttonActiveStyle: string;
  dotColor: string;
}

export const PRESET_TAGS = [
  "Nhật ký",
  "Lịch trình",
  "Ý tưởng",
  "Học tập",
  "Công việc",
  "Mục tiêu",
];

export const MOOD_CONFIGS: MoodConfig[] = [
  {
    value: "VERY_POSITIVE",
    translationKey: "mood.VERY_POSITIVE",
    // FIX: use text-emerald-700 for light mode (was text-emerald-800 — too dark/near-black on emerald bg)
    badgeStyle: "bg-emerald-500/15 text-emerald-700 dark:text-emerald-300 border-emerald-500/30",
    buttonActiveStyle:
      "bg-emerald-500 text-white dark:bg-emerald-500 dark:text-slate-950 border-emerald-500 shadow-md shadow-emerald-500/25 font-bold",
    dotColor: "bg-emerald-500",
  },
  {
    value: "POSITIVE",
    translationKey: "mood.POSITIVE",
    badgeStyle: "bg-teal-500/15 text-teal-700 dark:text-teal-300 border-teal-500/30",
    buttonActiveStyle:
      "bg-teal-500 text-white dark:bg-teal-500 dark:text-slate-950 border-teal-500 shadow-md shadow-teal-500/25 font-bold",
    dotColor: "bg-teal-500",
  },
  {
    value: "NEUTRAL",
    translationKey: "mood.NEUTRAL",
    badgeStyle: "bg-sky-500/15 text-sky-700 dark:text-sky-300 border-sky-500/30",
    buttonActiveStyle:
      "bg-sky-500 text-white dark:bg-sky-500 dark:text-slate-950 border-sky-500 shadow-md shadow-sky-500/25 font-bold",
    dotColor: "bg-sky-500",
  },
  {
    value: "NEGATIVE",
    translationKey: "mood.NEGATIVE",
    badgeStyle: "bg-amber-500/15 text-amber-700 dark:text-amber-300 border-amber-500/30",
    buttonActiveStyle:
      "bg-amber-500 text-white dark:bg-amber-500 dark:text-slate-950 border-amber-500 shadow-md shadow-amber-500/25 font-bold",
    dotColor: "bg-amber-500",
  },
  {
    value: "VERY_NEGATIVE",
    translationKey: "mood.VERY_NEGATIVE",
    badgeStyle: "bg-rose-500/15 text-rose-700 dark:text-rose-300 border-rose-500/30",
    buttonActiveStyle:
      "bg-rose-500 text-white dark:bg-rose-500 dark:text-slate-950 border-rose-500 shadow-md shadow-rose-500/25 font-bold",
    dotColor: "bg-rose-500",
  },
];

export function getMoodConfig(mood?: JournalMood | null): MoodConfig | undefined {
  if (!mood) return undefined;
  return MOOD_CONFIGS.find((m) => m.value === mood);
}
