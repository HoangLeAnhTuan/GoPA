/**
 * Domain constants for GoPA frontend.
 * Maps directly to backend domain models in be/internal/core/domain.
 */

export const USER_ROLE = {
  USER: "USER",
  ADMIN: "ADMIN",
} as const;

export const TASK_STATUS = {
  TODO: "TODO",
  IN_PROGRESS: "IN_PROGRESS",
  DONE: "DONE",
} as const;

export const TASK_PRIORITY = {
  LOW: "LOW",
  MEDIUM: "MEDIUM",
  HIGH: "HIGH",
} as const;

export const TASK_CATEGORY = {
  WORK: "WORK",
  STUDY: "STUDY",
  LIFE: "LIFE",
} as const;

export const VOCABULARY_LANGUAGE = {
  JP: "JP",
  EN: "EN",
} as const;

export const JOURNAL_MOOD = {
  VERY_NEGATIVE: "VERY_NEGATIVE",
  NEGATIVE: "NEGATIVE",
  NEUTRAL: "NEUTRAL",
  POSITIVE: "POSITIVE",
  VERY_POSITIVE: "VERY_POSITIVE",
} as const;

export const POMODORO_STATUS = {
  RUNNING: "running",
  PAUSED: "paused",
} as const;

export const ACCOUNT_TYPE = {
  CASH: "CASH",
  BANK: "BANK",
  CREDIT_CARD: "CREDIT_CARD",
  SAVINGS: "SAVINGS",
  INVESTMENT: "INVESTMENT",
  CRYPTO: "CRYPTO",
} as const;

export const CATEGORY_TYPE = {
  INCOME: "INCOME",
  EXPENSE: "EXPENSE",
} as const;

export const TRANSACTION_TYPE = {
  INCOME: "INCOME",
  EXPENSE: "EXPENSE",
  TRANSFER: "TRANSFER",
} as const;
