import axios, { AxiosError } from "axios";
import { DEFAULT_API_TIMEOUT } from "../constants/constants";

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
  baseURL: import.meta.env.VITE_API_BASE_URL ?? "http://localhost:8080/api/v1",
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
