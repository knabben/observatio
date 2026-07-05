import "@testing-library/jest-dom";
import {screen} from "@testing-library/react";
import {render} from "@/app/ui/dashboard/utils/test-render";
import NavLinks from "./nav-links";

const mockUsePathname = jest.fn();
jest.mock("next/navigation", () => ({
  usePathname: () => mockUsePathname(),
}));

describe("NavLinks", () => {
  afterEach(() => {
    jest.resetAllMocks();
  });

  it("marks the exact-match link active with aria-current", () => {
    mockUsePathname.mockReturnValue("/dashboard/clusters");
    render(<NavLinks/>);
    const clustersLink = screen.getByRole("link", {name: "Clusters"});
    expect(clustersLink).toHaveAttribute("aria-current", "page");
  });

  it("highlights the parent link on a nested route", () => {
    mockUsePathname.mockReturnValue("/dashboard/clusters/some-cluster");
    render(<NavLinks/>);
    const clustersLink = screen.getByRole("link", {name: "Clusters"});
    expect(clustersLink).toHaveAttribute("aria-current", "page");
    const dashboardLink = screen.getByRole("link", {name: "Dashboard"});
    expect(dashboardLink).not.toHaveAttribute("aria-current");
  });

  it("only matches the dashboard root link exactly (root guard)", () => {
    mockUsePathname.mockReturnValue("/dashboard/machines");
    render(<NavLinks/>);
    const dashboardLink = screen.getByRole("link", {name: "Dashboard"});
    expect(dashboardLink).not.toHaveAttribute("aria-current");
  });

  it("gives every link an accessible name via aria-label, independent of visible text", () => {
    mockUsePathname.mockReturnValue("/dashboard");
    render(<NavLinks/>);
    expect(screen.getByRole("link", {name: "Dashboard"})).toBeInTheDocument();
    expect(screen.getByRole("link", {name: "Machines"})).toBeInTheDocument();
  });
});
