import { apiClient } from "../../../lib/api-client";
import { unwrap, unwrapCollection, type ApiEnvelope } from "../../../lib/api-envelope";
import type { Journal, JournalFilter, JournalInput, JournalStats } from "../types";

export async function listJournals(filter?: JournalFilter | string) {
  const params: Record<string, string | number> = {};
  if (typeof filter === "string") {
    if (filter.trim() !== "") params.q = filter.trim();
  } else if (filter) {
    if (filter.q?.trim()) params.q = filter.q.trim();
    if (filter.tag?.trim()) params.tag = filter.tag.trim();
    if (filter.mood) params.mood = filter.mood;
    if (filter.from) params.from = filter.from;
    if (filter.to) params.to = filter.to;
    if (filter.limit) params.limit = filter.limit;
  }
  return unwrapCollection(
    await apiClient.get<ApiEnvelope<Journal[] | null>>("/journals", { params })
  );
}

export async function getJournal(id: string) {
  return unwrap(await apiClient.get<ApiEnvelope<Journal>>(`/journals/${id}`));
}

export async function getJournalStats() {
  return unwrap(await apiClient.get<ApiEnvelope<JournalStats>>("/journals/stats"));
}

export async function createJournal(input: JournalInput) {
  return unwrap(await apiClient.post<ApiEnvelope<Journal>>("/journals", input));
}

export async function updateJournal(id: string, input: JournalInput) {
  return unwrap(await apiClient.patch<ApiEnvelope<Journal>>(`/journals/${id}`, input));
}

export async function deleteJournal(id: string) {
  return unwrap(await apiClient.delete<ApiEnvelope<Record<string, never>>>(`/journals/${id}`));
}

export async function linkJournal(journalId: string, linkedId: string) {
  return unwrap(
    await apiClient.post<ApiEnvelope<Record<string, never>>>(`/journals/${journalId}/link`, {
      linked_id: linkedId,
    })
  );
}

export async function unlinkJournal(journalId: string, linkedId: string) {
  return unwrap(
    await apiClient.delete<ApiEnvelope<Record<string, never>>>(
      `/journals/${journalId}/link/${linkedId}`
    )
  );
}

