import { apiClient } from "../../../lib/api-client";
import { unwrap, unwrapCollection, type ApiEnvelope } from "../../../lib/api-envelope";
import type { Journal, JournalInput } from "../types";
export async function listJournals(query:string){return unwrapCollection(await apiClient.get<ApiEnvelope<Journal[] | null>>("/journals",{params:query.trim()===""?{}:{q:query}}));}
export async function createJournal(input:JournalInput){return unwrap(await apiClient.post<ApiEnvelope<Journal>>("/journals",input));}
export async function updateJournal(id:string,input:JournalInput){return unwrap(await apiClient.patch<ApiEnvelope<Journal>>(`/journals/${id}`,input));}
