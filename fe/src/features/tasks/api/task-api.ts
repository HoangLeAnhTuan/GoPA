import { apiClient } from "../../../lib/api-client";
import { unwrap, unwrapCollection, type ApiEnvelope } from "../../../lib/api-envelope";
import type { Task, TaskDto, TaskInput } from "../types";

export async function listTasks(): Promise<Task[]> { return unwrapCollection(await apiClient.get<ApiEnvelope<TaskDto[] | null>>("/tasks")).map(toTask); }
export async function createTask(input: TaskInput): Promise<Task> { return toTask(unwrap(await apiClient.post<ApiEnvelope<TaskDto>>("/tasks", toDto(input)))); }
export async function updateTask(id: string, input: TaskInput): Promise<Task> { return toTask(unwrap(await apiClient.patch<ApiEnvelope<TaskDto>>(`/tasks/${id}`, toDto(input)))); }
export async function deleteTask(id: string): Promise<void> { await apiClient.delete(`/tasks/${id}`); }
function toTask(dto: TaskDto): Task { return { id:dto.id, title:dto.title, description:dto.description, status:dto.status, priority:dto.priority, category:dto.category, dueDate:dto.due_date === null ? null : new Date(dto.due_date), position:dto.position, createdAt:new Date(dto.created_at), updatedAt:new Date(dto.updated_at) }; }
function toDto(input: TaskInput) { return { title:input.title, description:input.description, status:input.status, priority:input.priority, category:input.category, due_date:input.dueDate === null ? null : `${input.dueDate}T00:00:00Z`, ...(input.position === undefined ? {} : { position:input.position }) }; }
