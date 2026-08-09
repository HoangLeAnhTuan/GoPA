import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { createJournal, listJournals, updateJournal } from "../api/journal-api";
import type { JournalInput } from "../types";
const keys={all:["journals"] as const,list:(query:string)=>[...keys.all,"list",query] as const};
export function useJournals(query:string){return useQuery({queryKey:keys.list(query),queryFn:()=>listJournals(query)});} export function useCreateJournal(){const c=useQueryClient();return useMutation({mutationFn:(input:JournalInput)=>createJournal(input),onSuccess:()=>c.invalidateQueries({queryKey:keys.all})});} export function useUpdateJournal(){const c=useQueryClient();return useMutation({mutationFn:({id,input}:{id:string;input:JournalInput})=>updateJournal(id,input),onSuccess:()=>c.invalidateQueries({queryKey:keys.all})});}
