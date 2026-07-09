import "@testing-library/jest-dom";
import {screen} from "@testing-library/react";
import {Grid} from "@mantine/core";
import {render} from "@/app/ui/dashboard/utils/test-render";
import ClusterClassTable from "./table";
import {ClusterClassType} from "./types";

const wrap = (ui: React.ReactNode) => render(<Grid>{ui}</Grid>);

describe("ClusterClassTable", () => {
  it("renders a partial ClusterClass (missing conditions) without throwing", () => {
    const partial = [{metadata: {name: "cc1"}, name: "cc1"}] as ClusterClassType[];
    wrap(<ClusterClassTable ccs={partial} select={() => {}}/>);
    expect(screen.getByRole("button", {name: /select cc1/i})).toBeInTheDocument();
  });

  it("shows unknown status as 'Unknown' when no conditions are reported", () => {
    const partial = [{metadata: {name: "cc1"}, name: "cc1"}] as ClusterClassType[];
    wrap(<ClusterClassTable ccs={partial} select={() => {}}/>);
    expect(screen.getByRole("img", {name: "Unknown"})).toBeInTheDocument();
    expect(screen.queryByRole("img", {name: "Not ready"})).not.toBeInTheDocument();
  });

  it("renders a labeled empty state for an empty collection without crashing", () => {
    wrap(<ClusterClassTable ccs={[]} select={() => {}}/>);
    expect(screen.getByText(/no cluster classes found/i)).toBeInTheDocument();
  });
});
