import { apiClient } from "../../../lib/api-client";
import { unwrap, unwrapCollection, type ApiEnvelope } from "../../../lib/api-envelope";

export interface PomodoroState {
  status: "running" | "paused";
  started_at: string;
  duration_seconds: number;
  task_id: string | null;
  updated_at: string;
}

export interface PomodoroHistoryItem {
  id: string;
  event_id: string;
  user_id: string;
  task_id?: string | null;
  started_at: string;
  ended_at: string;
  duration_seconds: number;
}

export async function getPomodoro(): Promise<PomodoroState | null> {
  return unwrap(await apiClient.get<ApiEnvelope<PomodoroState | null>>("/pomodoro"));
}

export async function setPomodoro(state: PomodoroState): Promise<void> {
  await apiClient.put("/pomodoro", state);
}

export async function stopPomodoro(): Promise<void> {
  await apiClient.post("/pomodoro/stop");
}

export async function listPomodoroHistory(limit = 20, cursor?: string): Promise<PomodoroHistoryItem[]> {
  const params: Record<string, string | number> = { limit };
  if (cursor) params.cursor = cursor;
  return unwrapCollection(
    await apiClient.get<ApiEnvelope<PomodoroHistoryItem[] | null>>("/pomodoro/history", { params })
  );
}
