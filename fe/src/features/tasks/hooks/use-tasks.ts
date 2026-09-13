import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { createTask, deleteTask, listTasks, updateTask, updateTaskStatus } from "../api/task-api";
import { taskKeys } from "../api/task-keys";
import type { TaskInput, TaskStatus } from "../types";

export function useTasks() {
  return useQuery({ queryKey: taskKeys.list(), queryFn: listTasks });
}

export function useCreateTask() {
  const client = useQueryClient();
  return useMutation({
    mutationFn: (input: TaskInput) => createTask(input),
    onSuccess: () => client.invalidateQueries({ queryKey: taskKeys.all }),
  });
}

export function useUpdateTask() {
  const client = useQueryClient();
  return useMutation({
    mutationFn: ({ id, input }: { id: string; input: TaskInput }) => updateTask(id, input),
    onSuccess: () => client.invalidateQueries({ queryKey: taskKeys.all }),
  });
}

export function useUpdateTaskStatus() {
  const client = useQueryClient();
  return useMutation({
    mutationFn: ({ id, status, targetPosition }: { id: string; status: TaskStatus; targetPosition?: number }) =>
      updateTaskStatus(id, status, targetPosition),
    onSuccess: () => client.invalidateQueries({ queryKey: taskKeys.all }),
  });
}

export function useDeleteTask() {
  const client = useQueryClient();
  return useMutation({
    mutationFn: deleteTask,
    onSuccess: () => client.invalidateQueries({ queryKey: taskKeys.all }),
  });
}
