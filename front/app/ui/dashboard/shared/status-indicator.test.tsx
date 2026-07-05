import {describe, it, expect} from "@jest/globals";
import {screen} from "@testing-library/react";
import {render} from "@/app/ui/dashboard/utils/test-render";
import {toStatusState, allReady} from "./status";
import {StatusIndicator} from "./status-indicator";
import {EmptyState} from "./empty-state";
import {ErrorState} from "./error-state";

describe("toStatusState", () => {
  it("maps true → healthy, false → notready, absent → unknown", () => {
    expect(toStatusState(true)).toBe("healthy");
    expect(toStatusState(false)).toBe("notready");
    expect(toStatusState(undefined)).toBe("unknown");
    expect(toStatusState(null)).toBe("unknown");
  });

  it("does not treat a zero-ish/absent value as failed", () => {
    // absent readiness must be unknown, never notready
    expect(toStatusState(undefined)).not.toBe("notready");
  });
});

describe("allReady", () => {
  it("is healthy only when every flag is true", () => {
    expect(allReady(true, true)).toBe("healthy");
    expect(allReady(true, false)).toBe("notready");
    expect(allReady(true, undefined)).toBe("unknown");
  });
});

describe("StatusIndicator", () => {
  it("renders an accessible label per state (not color alone)", () => {
    render(<StatusIndicator state="notready"/>);
    expect(screen.getByText("Not ready")).toBeInTheDocument();
    expect(screen.getByRole("img", {name: "Not ready"})).toBeInTheDocument();
  });

  it("renders unknown distinctly from healthy and failed", () => {
    render(<StatusIndicator state="unknown"/>);
    expect(screen.getByText("Unknown")).toBeInTheDocument();
  });
});

describe("EmptyState / ErrorState", () => {
  it("EmptyState shows its label", () => {
    render(<EmptyState label="No clusters found"/>);
    expect(screen.getByText("No clusters found")).toBeInTheDocument();
  });

  it("ErrorState shows the message and a retry button when onRetry is provided", () => {
    render(<ErrorState message="Unable to load clusters" onRetry={() => {}}/>);
    expect(screen.getByText("Unable to load clusters")).toBeInTheDocument();
    expect(screen.getByRole("button", {name: /retry/i})).toBeInTheDocument();
  });

  it("ErrorState omits the retry button when no onRetry", () => {
    render(<ErrorState message="boom"/>);
    expect(screen.queryByRole("button", {name: /retry/i})).not.toBeInTheDocument();
  });
});
