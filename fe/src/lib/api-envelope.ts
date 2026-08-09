import type { AxiosResponse } from "axios";

export interface ApiEnvelope<T> {
  data: T;
  meta: { request_id: string };
}

export function unwrap<T>(response: AxiosResponse<ApiEnvelope<T>>): T {
  return response.data.data;
}

/**
 * API collection endpoints are rendered as empty lists when an older API
 * instance returns `null` instead of the JSON array contract (`[]`).
 */
export function unwrapCollection<T>(
  response: AxiosResponse<ApiEnvelope<T[] | null>>,
): T[] {
  const data = response.data.data;
  return Array.isArray(data) ? data : [];
}
