import type { AxiosResponse } from "axios";
import { describe, expect, it } from "vitest";
import { unwrapCollection, type ApiEnvelope } from "./api-envelope";

function collectionResponse<T>(data: T[] | null): AxiosResponse<ApiEnvelope<T[] | null>> {
  return { data: { data, meta: { request_id: "test-request" } } } as AxiosResponse<ApiEnvelope<T[] | null>>;
}

describe("unwrapCollection", () => {
  it("turns a legacy null collection into an empty array", () => {
    expect(unwrapCollection(collectionResponse<string>(null))).toEqual([]);
  });

  it("keeps API collection values intact", () => {
    expect(unwrapCollection(collectionResponse(["first"]))).toEqual(["first"]);
  });
});
