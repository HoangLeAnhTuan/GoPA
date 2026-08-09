import { describe, expect, it, vi } from "vitest";
import { apiClient } from "../../../lib/api-client";
import {
  createJournal,
  getJournalStats,
  linkJournal,
  listJournals,
  unlinkJournal,
} from "./journal-api";

vi.mock("../../../lib/api-client", () => ({
  apiClient: {
    get: vi.fn(),
    post: vi.fn(),
    patch: vi.fn(),
    delete: vi.fn(),
  },
}));

describe("journal-api", () => {
  it("listJournals passes query filter params", async () => {
    vi.mocked(apiClient.get).mockResolvedValueOnce({
      data: { status: "success", data: [] },
    });

    await listJournals({ q: "markdown", mood: "POSITIVE", tag: "ideas" });

    expect(apiClient.get).toHaveBeenCalledWith("/journals", {
      params: { q: "markdown", mood: "POSITIVE", tag: "ideas" },
    });
  });

  it("getJournalStats calls stats endpoint", async () => {
    vi.mocked(apiClient.get).mockResolvedValueOnce({
      data: {
        status: "success",
        data: {
          current_streak: 5,
          longest_streak: 10,
          word_count: 500,
          entries: 3,
          mood_counts: { POSITIVE: 3 },
        },
      },
    });

    const res = await getJournalStats();
    expect(res.current_streak).toBe(5);
    expect(apiClient.get).toHaveBeenCalledWith("/journals/stats");
  });

  it("createJournal posts input with mood & published_date", async () => {
    vi.mocked(apiClient.post).mockResolvedValueOnce({
      data: {
        status: "success",
        data: {
          id: "j-1",
          user_id: "u-1",
          title: "My Entry",
          content: "Content",
          tags: ["go"],
          mood: "VERY_POSITIVE",
          published_date: "2026-08-09",
          created_at: "2026-08-09T00:00:00Z",
          updated_at: "2026-08-09T00:00:00Z",
        },
      },
    });

    const res = await createJournal({
      title: "My Entry",
      content: "Content",
      tags: ["go"],
      mood: "VERY_POSITIVE",
      published_date: "2026-08-09",
    });

    expect(res.id).toBe("j-1");
    expect(apiClient.post).toHaveBeenCalledWith("/journals", {
      title: "My Entry",
      content: "Content",
      tags: ["go"],
      mood: "VERY_POSITIVE",
      published_date: "2026-08-09",
    });
  });

  it("linkJournal calls post to link endpoint", async () => {
    vi.mocked(apiClient.post).mockResolvedValueOnce({
      data: { status: "success", data: {} },
    });

    await linkJournal("j-1", "j-2");

    expect(apiClient.post).toHaveBeenCalledWith("/journals/j-1/link", {
      linked_id: "j-2",
    });
  });

  it("unlinkJournal calls delete on link endpoint", async () => {
    vi.mocked(apiClient.delete).mockResolvedValueOnce({
      data: { status: "success", data: {} },
    });

    await unlinkJournal("j-1", "j-2");

    expect(apiClient.delete).toHaveBeenCalledWith("/journals/j-1/link/j-2");
  });
});
