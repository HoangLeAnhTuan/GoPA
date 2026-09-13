export type TaskStatus = "TODO" | "IN_PROGRESS" | "DONE";
export type TaskPriority = "LOW" | "MEDIUM" | "HIGH" | "URGENT";
export type TaskCategory = "WORK" | "STUDY" | "LIFE" | "LEETCODE";

export interface TaskDto {
  id: string;
  title: string;
  description: string;
  status: TaskStatus;
  priority: TaskPriority;
  category: TaskCategory;
  due_date: string | null;
  position: number;
  estimated_pomodoros?: number;
  completed_pomodoros?: number;
  created_at: string;
  updated_at: string;
}

export interface Task {
  id: string;
  title: string;
  description: string;
  status: TaskStatus;
  priority: TaskPriority;
  category: TaskCategory;
  dueDate: Date | null;
  position: number;
  estimatedPomodoros: number;
  completedPomodoros: number;
  createdAt: Date;
  updatedAt: Date;
}

export interface TaskInput {
  title: string;
  description: string;
  status: TaskStatus;
  priority: TaskPriority;
  category: TaskCategory;
  dueDate: string | null;
  position?: number;
  estimated_pomodoros?: number;
}
