import { apiClient } from "../../../lib/api-client";
import { unwrap, unwrapCollection, type ApiEnvelope } from "../../../lib/api-envelope";
import type { LearningSession, Vocabulary, VocabularyInput, VocabularyStats } from "../types";

export async function listVocabularies(limit: number = 100): Promise<Vocabulary[]> {
  return unwrapCollection(
    await apiClient.get<ApiEnvelope<Vocabulary[] | null>>("/vocabularies", {
      params: { limit },
    })
  );
}

export async function reviewQueue(limit: number = 50): Promise<Vocabulary[]> {
  return unwrapCollection(
    await apiClient.get<ApiEnvelope<Vocabulary[] | null>>("/vocabularies/review-queue", {
      params: { limit },
    })
  );
}

export async function getStats(): Promise<VocabularyStats> {
  return unwrap(await apiClient.get<ApiEnvelope<VocabularyStats>>("/vocabularies/stats"));
}

export async function createVocabulary(input: VocabularyInput): Promise<Vocabulary> {
  return unwrap(
    await apiClient.post<ApiEnvelope<Vocabulary>>("/vocabularies", {
      language: input.language,
      word: input.word,
      reading: input.reading,
      meaning: input.meaning,
      example_sentence: input.exampleSentence,
      example_translation: input.exampleTranslation,
      tags: input.tags ?? [],
      difficulty_level: input.difficultyLevel ?? "MEDIUM",
    })
  );
}

export async function updateVocabulary(
  id: string,
  input: Partial<VocabularyInput>
): Promise<Vocabulary> {
  return unwrap(
    await apiClient.patch<ApiEnvelope<Vocabulary>>(`/vocabularies/${id}`, {
      language: input.language,
      word: input.word,
      reading: input.reading,
      meaning: input.meaning,
      example_sentence: input.exampleSentence,
      example_translation: input.exampleTranslation,
      tags: input.tags,
      difficulty_level: input.difficultyLevel,
    })
  );
}

export async function deleteVocabulary(id: string): Promise<void> {
  await apiClient.delete(`/vocabularies/${id}`);
}

export async function reviewVocabulary(id: string, quality: number): Promise<void> {
  await apiClient.post(`/vocabularies/${id}/reviews`, { quality });
}

export async function generateAudio(
  id: string
): Promise<{ vocabulary_id: string; audio_url: string | null }> {
  return unwrap(
    await apiClient.post<ApiEnvelope<{ vocabulary_id: string; audio_url: string | null }>>(
      `/vocabularies/${id}/audio`
    )
  );
}

export async function importVocabularies(
  data: string,
  format: "csv" | "json" = "json"
): Promise<{ imported_count: number }> {
  return unwrap(
    await apiClient.post<ApiEnvelope<{ imported_count: number }>>(
      `/vocabularies/import?format=${format}`,
      data,
      {
        headers: {
          "Content-Type": format === "csv" ? "text/csv" : "application/json",
        },
      }
    )
  );
}

export async function startLearningSession(
  language: string,
  sessionType: string = "review"
): Promise<LearningSession> {
  return unwrap(
    await apiClient.post<ApiEnvelope<LearningSession>>("/learning-sessions", {
      language,
      session_type: sessionType,
    })
  );
}

export async function endLearningSession(
  id: string,
  itemsReviewed: number,
  itemsCorrect: number,
  durationSeconds: number
): Promise<LearningSession> {
  return unwrap(
    await apiClient.patch<ApiEnvelope<LearningSession>>(`/learning-sessions/${id}`, {
      items_reviewed: itemsReviewed,
      items_correct: itemsCorrect,
      duration_seconds: durationSeconds,
    })
  );
}
