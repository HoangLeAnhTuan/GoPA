import axios, { AxiosError, type InternalAxiosRequestConfig } from "axios";
import { apiClient, toApiError } from "../../../lib/api-client";
import type { AuthDto, AuthResult, LoginInput, RegisterInput, User, UserDto } from "../types";

interface Envelope<T> {
  data: T;
  meta: { request_id: string };
}

let accessToken: string | null = null;
let refreshPromise: Promise<string> | null = null;

const refreshClient = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL ?? "http://localhost:8080/api/v1",
  timeout: 10_000,
  withCredentials: true,
});

apiClient.interceptors.request.use((config) => {
  if (accessToken !== null) {
    config.headers.Authorization = `Bearer ${accessToken}`;
  }
  return config;
});

apiClient.interceptors.response.use(
  (response) => response,
  async (error: unknown) => {
    if (!(error instanceof AxiosError) || error.response?.status !== 401 || isAuthRequest(error.config)) {
      return Promise.reject(error);
    }
    const request = error.config as (InternalAxiosRequestConfig & { _retry?: boolean }) | undefined;
    if (request === undefined || request._retry === true) {
      return Promise.reject(error);
    }
    request._retry = true;
    try {
      request.headers.Authorization = `Bearer ${await refreshAccessToken()}`;
      return await apiClient.request(request);
    } catch (refreshError: unknown) {
      clearAccessToken();
      return Promise.reject(refreshError);
    }
  },
);

export async function register(input: RegisterInput): Promise<User> {
  const response = await apiClient.post<Envelope<UserDto>>("/auth/register", {
    email: input.email,
    password: input.password,
    confirm_password: input.confirmPassword,
  });
  return toUser(response.data.data);
}

export async function login(input: LoginInput): Promise<AuthResult> {
  const response = await apiClient.post<Envelope<AuthDto>>("/auth/login", input);
  return applyAuthResult(response.data.data);
}

export async function restoreSession(): Promise<User> {
  const token = await refreshAccessToken();
  if (token === "") {
    throw new Error("Refresh did not return an access token.");
  }
  const response = await apiClient.get<Envelope<UserDto>>("/me");
  return toUser(response.data.data);
}

export async function logout(): Promise<void> {
  try {
    await refreshClient.post("/auth/logout");
  } catch (error: unknown) {
    const apiError = toApiError(error);
    if (apiError.status !== 401) {
      throw error;
    }
  } finally {
    clearAccessToken();
  }
}

export function clearAccessToken() {
  accessToken = null;
}

async function refreshAccessToken(): Promise<string> {
  if (refreshPromise === null) {
    refreshPromise = refreshClient
      .post<Envelope<AuthDto>>("/auth/refresh")
      .then((response) => applyAuthResult(response.data.data).accessToken)
      .finally(() => {
        refreshPromise = null;
      });
  }
  return refreshPromise;
}

function applyAuthResult(dto: AuthDto): AuthResult {
  accessToken = dto.access_token;
  return { accessToken: dto.access_token, user: toUser(dto.user) };
}

function toUser(dto: UserDto): User {
  return { id: dto.id, email: dto.email, role: dto.role, createdAt: new Date(dto.created_at) };
}

function isAuthRequest(config: InternalAxiosRequestConfig | undefined): boolean {
  return config?.url?.includes("/auth/") ?? false;
}
