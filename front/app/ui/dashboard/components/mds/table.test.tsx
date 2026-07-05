import "@testing-library/jest-dom";
import {screen} from "@testing-library/react";
import {Grid} from "@mantine/core";
import {render} from "@/app/ui/dashboard/utils/test-render";
import MDTable from "./table";
import {MachineDeploymentType} from "./types";

const wrap = (ui: React.ReactNode) => render(<Grid>{ui}</Grid>);

describe("MDTable", () => {
  it("renders a partial machine deployment (missing status) without throwing", () => {
    const partial = [{metadata: {name: "md1"}}] as MachineDeploymentType[];
    wrap(<MDTable mds={partial} select={() => {}}/>);
    expect(screen.getByRole("button", {name: /select md1/i})).toBeInTheDocument();
  });

  it("shows unknown availability as 'Unknown', not failed", () => {
    // no status ⇒ unavailableReplicas absent ⇒ unknown, never notready
    const partial = [{metadata: {name: "md1"}}] as MachineDeploymentType[];
    wrap(<MDTable mds={partial} select={() => {}}/>);
    expect(screen.getByRole("img", {name: "Unknown"})).toBeInTheDocument();
    expect(screen.queryByRole("img", {name: "Not ready"})).not.toBeInTheDocument();
  });

  it("renders a labeled empty state for an empty collection without crashing", () => {
    wrap(<MDTable mds={[]} select={() => {}}/>);
    expect(screen.getByText(/no machine deployments found/i)).toBeInTheDocument();
  });
});
