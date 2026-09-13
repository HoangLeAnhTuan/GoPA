export type JournalMood =
  | "VERY_NEGATIVE"
  | "NEGATIVE"
  | "NEUTRAL"
  | "POSITIVE"
  | "VERY_POSITIVE";

export interface Journal {
  id: string;
  user_id: string;
  title: string;
  content: string;
  tags: string[];
  mood?: JournalMood | null;
  energy_level?: number | null;
  pinned?: boolean;
  word_count?: number;
  published_date: string;
  created_at: string;
  updated_at: string;
}

export interface JournalInput {
  title: string;
  content: string;
  tags: string[];
  mood?: JournalMood | null;
  energy_level?: number | null;
  pinned?: boolean;
  published_date?: string;
}

export interface JournalFilter {
  q?: string;
  tag?: string;
  mood?: JournalMood | "";
  from?: string;
  to?: string;
  limit?: number;
}

export interface JournalStats {
  current_streak: number;
  longest_streak: number;
  word_count: number;
  entries: number;
  mood_counts: Record<JournalMood, number>;
}
