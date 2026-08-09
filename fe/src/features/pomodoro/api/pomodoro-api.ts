import { apiClient } from "../../../lib/api-client";
import { unwrap, type ApiEnvelope } from "../../../lib/api-envelope";
export interface PomodoroState { status:"running"|"paused"; started_at:string; duration_seconds:number; task_id:string|null; updated_at:string; }
export async function getPomodoro(){return unwrap(await apiClient.get<ApiEnvelope<PomodoroState|null>>("/pomodoro"));} export async function setPomodoro(state:PomodoroState){await apiClient.put("/pomodoro",state);} export async function stopPomodoro(){await apiClient.post("/pomodoro/stop");}
