import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  createVocabulary,
  deleteVocabulary,
  generateAudio,
  getStats,
  importVocabularies,
  listVocabularies,
  reviewQueue,
  reviewVocabulary,
  updateVocabulary,
} from "../api/vocabulary-api";
import type { VocabularyInput } from "../types";

export const vocabularyKeys = {
  all: ["vocabularies"] as const,
  list: () => [...vocabularyKeys.all, "list"] as const,
  due: () => [...vocabularyKeys.all, "due"] as const,
  stats: () => [...vocabularyKeys.all, "stats"] as const,
};

export function useVocabularies() {
  return useQuery({
    queryKey: vocabularyKeys.list(),
    queryFn: () => listVocabularies(150),
  });
}

export function useReviewQueue() {
  return useQuery({
    queryKey: vocabularyKeys.due(),
    queryFn: () => reviewQueue(50),
  });
}

export function useVocabularyStats() {
  return useQuery({
    queryKey: vocabularyKeys.stats(),
    queryFn: getStats,
  });
}

export function useCreateVocabulary() {
  const client = useQueryClient();
  return useMutation({
    mutationFn: (input: VocabularyInput) => createVocabulary(input),
    onSuccess: () => client.invalidateQueries({ queryKey: vocabularyKeys.all }),
  });
}

export function useUpdateVocabulary() {
  const client = useQueryClient();
  return useMutation({
    mutationFn: ({ id, input }: { id: string; input: Partial<VocabularyInput> }) =>
      updateVocabulary(id, input),
    onSuccess: () => client.invalidateQueries({ queryKey: vocabularyKeys.all }),
  });
}

export function useDeleteVocabulary() {
  const client = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => deleteVocabulary(id),
    onSuccess: () => client.invalidateQueries({ queryKey: vocabularyKeys.all }),
  });
}

export function useReviewVocabulary() {
  const client = useQueryClient();
  return useMutation({
    mutationFn: ({ id, quality }: { id: string; quality: number }) =>
      reviewVocabulary(id, quality),
    onSuccess: () => client.invalidateQueries({ queryKey: vocabularyKeys.all }),
  });
}

export function useGenerateAudio() {
  return useMutation({
    mutationFn: (id: string) => generateAudio(id),
  });
}

export function useImportVocabularies() {
  const client = useQueryClient();
  return useMutation({
    mutationFn: ({ data, format }: { data: string; format?: "csv" | "json" }) =>
      importVocabularies(data, format),
    onSuccess: () => client.invalidateQueries({ queryKey: vocabularyKeys.all }),
  });
}
