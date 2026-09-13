export type VocabularyLanguage = "JP" | "EN";

export interface Vocabulary {
  id: string;
  user_id?: string;
  language: VocabularyLanguage;
  word: string;
  reading: string | null;
  meaning: string;
  example_sentence: string | null;
  example_translation?: string | null;
  tags?: string[];
  difficulty_level?: string;
  audio_url?: string | null;
  image_url?: string | null;
  ease_factor?: number;
  interval_days?: number;
  repetition_count?: number;
  current_box: number;
  next_review_at: string;
  mastery_score?: number;
  created_at?: string;
  updated_at?: string;
}

export interface VocabularyInput {
  language: VocabularyLanguage;
  word: string;
  reading: string | null;
  meaning: string;
  exampleSentence: string | null;
  exampleTranslation?: string | null;
  tags?: string[];
  difficultyLevel?: string;
}

export interface VocabularyStats {
  total_words: number;
  mastered_words: number;
  due_count: number;
  box_distribution: Record<number, number>;
  review_heatmap: Record<string, number>;
}

export interface LearningSession {
  id: string;
  user_id: string;
  language: VocabularyLanguage;
  session_type: string;
  items_reviewed: number;
  items_correct: number;
  duration_seconds: number;
  started_at: string;
  ended_at?: string | null;
}

export type StudyMode = "flashcard" | "multiple_choice" | "type_in" | "audio_quiz";
