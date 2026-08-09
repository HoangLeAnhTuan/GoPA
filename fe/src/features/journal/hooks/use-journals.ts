import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  createJournal,
  deleteJournal,
  getJournal,
  getJournalStats,
  linkJournal,
  listJournals,
  unlinkJournal,
  updateJournal,
} from "../api/journal-api";
import type { JournalFilter, JournalInput } from "../types";

const keys = {
  all: ["journals"] as const,
  list: (filter?: JournalFilter | string) => [...keys.all, "list", filter] as const,
  detail: (id: string) => [...keys.all, "detail", id] as const,
  stats: () => [...keys.all, "stats"] as const,
};

export function useJournals(filter?: JournalFilter | string) {
  return useQuery({
    queryKey: keys.list(filter),
    queryFn: () => listJournals(filter),
  });
}

export function useJournal(id?: string) {
  return useQuery({
    queryKey: keys.detail(id ?? ""),
    queryFn: () => getJournal(id!),
    enabled: Boolean(id),
  });
}

export function useJournalStats() {
  return useQuery({
    queryKey: keys.stats(),
    queryFn: () => getJournalStats(),
  });
}

export function useCreateJournal() {
  const c = useQueryClient();
  return useMutation({
    mutationFn: (input: JournalInput) => createJournal(input),
    onSuccess: () => c.invalidateQueries({ queryKey: keys.all }),
  });
}

export function useUpdateJournal() {
  const c = useQueryClient();
  return useMutation({
    mutationFn: ({ id, input }: { id: string; input: JournalInput }) => updateJournal(id, input),
    onSuccess: () => c.invalidateQueries({ queryKey: keys.all }),
  });
}

export function useDeleteJournal() {
  const c = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => deleteJournal(id),
    onSuccess: () => c.invalidateQueries({ queryKey: keys.all }),
  });
}

export function useLinkJournal() {
  const c = useQueryClient();
  return useMutation({
    mutationFn: ({ journalId, linkedId }: { journalId: string; linkedId: string }) =>
      linkJournal(journalId, linkedId),
    onSuccess: () => c.invalidateQueries({ queryKey: keys.all }),
  });
}

export function useUnlinkJournal() {
  const c = useQueryClient();
  return useMutation({
    mutationFn: ({ journalId, linkedId }: { journalId: string; linkedId: string }) =>
      unlinkJournal(journalId, linkedId),
    onSuccess: () => c.invalidateQueries({ queryKey: keys.all }),
  });
}

