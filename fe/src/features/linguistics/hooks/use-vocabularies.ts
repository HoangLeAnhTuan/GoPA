import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { createVocabulary, listVocabularies, reviewQueue, reviewVocabulary } from "../api/vocabulary-api";
import type { VocabularyInput } from "../types";
const keys = { all:["vocabularies"] as const, list:()=>[...keys.all,"list"] as const, due:()=>[...keys.all,"due"] as const };
export function useVocabularies() { return useQuery({queryKey:keys.list(),queryFn:listVocabularies}); }
export function useReviewQueue() { return useQuery({queryKey:keys.due(),queryFn:reviewQueue}); }
export function useCreateVocabulary() { const client=useQueryClient(); return useMutation({mutationFn:(input:VocabularyInput)=>createVocabulary(input),onSuccess:()=>client.invalidateQueries({queryKey:keys.all})}); }
export function useReviewVocabulary() { const client=useQueryClient(); return useMutation({mutationFn:({id,quality}:{id:string;quality:number})=>reviewVocabulary(id,quality),onSuccess:()=>client.invalidateQueries({queryKey:keys.all})}); }
