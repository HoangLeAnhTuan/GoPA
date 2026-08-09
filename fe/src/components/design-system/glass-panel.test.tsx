import { render, screen } from "@testing-library/react";
import { describe, expect, it } from "vitest";
import { GlassPanel } from "./glass-panel";

describe("GlassPanel", () => {
  it("renders its content in the shared acrylic surface", () => {
    render(<GlassPanel>Workspace</GlassPanel>);
    expect(screen.getByText("Workspace")).toHaveClass("backdrop-blur-xl");
  });
});
