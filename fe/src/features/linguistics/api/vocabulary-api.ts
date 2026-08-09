import { apiClient } from "../../../lib/api-client";
import { unwrap, unwrapCollection, type ApiEnvelope } from "../../../lib/api-envelope";
import type { Vocabulary, VocabularyInput } from "../types";
export async function listVocabularies() { return unwrapCollection(await apiClient.get<ApiEnvelope<Vocabulary[] | null>>("/vocabularies")); }
export async function reviewQueue() { return unwrapCollection(await apiClient.get<ApiEnvelope<Vocabulary[] | null>>("/vocabularies/review-queue")); }
export async function createVocabulary(input: VocabularyInput) { return unwrap(await apiClient.post<ApiEnvelope<Vocabulary>>("/vocabularies", { language:input.language, word:input.word, reading:input.reading, meaning:input.meaning, example_sentence:input.exampleSentence })); }
export async function reviewVocabulary(id:string, quality:number) { await apiClient.post(`/vocabularies/${id}/reviews`, { quality }); }
