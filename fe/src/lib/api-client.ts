import axios, { AxiosError } from "axios";
import { DEFAULT_API_TIMEOUT } from "../constants/constants";

const configuredApiBaseUrl = import.meta.env.VITE_API_BASE_URL?.trim();

if (import.meta.env.PROD && !configuredApiBaseUrl) {
  throw new Error("VITE_API_BASE_URL is required for production builds.");
}

export const API_BASE_URL = (configuredApiBaseUrl || "http://localhost:8081/api/v1").replace(/\/+$/, "");

export interface ApiError {
  code: string;
  message: string;
  status?: number;
  requestId?: string;
  details?: Array<{ field: string; message: string }>;
}

interface ErrorEnvelope {
  error: { code: string; message: string; details?: Array<{ field: string; message: string }> };
  meta?: { request_id?: string };
}

export const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: DEFAULT_API_TIMEOUT,
  withCredentials: true,
});

export function toApiError(error: unknown): ApiError {
  if (error instanceof AxiosError) {
    const payload = error.response?.data as ErrorEnvelope | undefined;
    return {
      code: payload?.error.code ?? "NETWORK_ERROR",
      message: payload?.error.message ?? "Unable to reach GoPA. Please try again.",
      status: error.response?.status,
      requestId: payload?.meta?.request_id,
      details: payload?.error.details,
    };
  }
  return { code: "UNKNOWN_ERROR", message: "An unexpected error occurred." };
}
