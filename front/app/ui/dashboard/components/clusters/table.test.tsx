import "@testing-library/jest-dom";
import {screen} from "@testing-library/react";
import {Grid} from "@mantine/core";
import {render} from "@/app/ui/dashboard/utils/test-render";
import ClusterTable from "./table";
import {ClusterType} from "./types";

const wrap = (ui: React.ReactNode) => render(<Grid>{ui}</Grid>);

describe("ClusterTable", () => {
  it("renders a partial cluster (missing status/metadata fields) without throwing", () => {
    const partial = [{metadata: {name: "c1"}}] as ClusterType[];
    wrap(<ClusterTable clusters={partial} select={() => {}}/>);
    expect(screen.getByRole("button", {name: /select c1/i})).toBeInTheDocument();
    // absent version/phase render a placeholder, not a crash
    expect(screen.getAllByText("—").length).toBeGreaterThan(0);
    // a cluster with no/unrecognized provider still appears, with an Unknown badge
    expect(screen.getByText("Unknown")).toBeInTheDocument();
  });

  it("renders a labeled empty state for an empty collection (no rows, no crash)", () => {
    wrap(<ClusterTable clusters={[]} select={() => {}}/>);
    expect(screen.getByText(/no clusters found/i)).toBeInTheDocument();
    expect(screen.queryByRole("button", {name: /select/i})).not.toBeInTheDocument();
  });
});
