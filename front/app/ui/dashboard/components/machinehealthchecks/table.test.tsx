import "@testing-library/jest-dom";
import {screen} from "@testing-library/react";
import {Grid} from "@mantine/core";
import {render} from "@/app/ui/dashboard/utils/test-render";
import MachineHealthCheckTable from "./table";
import {MachineHealthCheckType} from "./types";

const wrap = (ui: React.ReactNode) => render(<Grid>{ui}</Grid>);

describe("MachineHealthCheckTable", () => {
  it("renders a partial MachineHealthCheck (missing status) without throwing", () => {
    const partial = [{metadata: {name: "mhc1"}}] as MachineHealthCheckType[];
    wrap(<MachineHealthCheckTable mhcs={partial} select={() => {}}/>);
    expect(screen.getByRole("button", {name: /select mhc1/i})).toBeInTheDocument();
  });

  it("shows unknown health as 'Unknown', not failed", () => {
    const partial = [{metadata: {name: "mhc1"}}] as MachineHealthCheckType[];
    wrap(<MachineHealthCheckTable mhcs={partial} select={() => {}}/>);
    expect(screen.getByRole("img", {name: "Unknown"})).toBeInTheDocument();
    expect(screen.queryByRole("img", {name: "Not ready"})).not.toBeInTheDocument();
  });

  it("renders a labeled empty state for an empty collection without crashing", () => {
    wrap(<MachineHealthCheckTable mhcs={[]} select={() => {}}/>);
    expect(screen.getByText(/no machine health checks found/i)).toBeInTheDocument();
  });
});
