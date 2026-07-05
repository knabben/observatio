import "@testing-library/jest-dom";
import {screen} from "@testing-library/react";
import {Grid} from "@mantine/core";
import {render} from "@/app/ui/dashboard/utils/test-render";
import MachinesTable from "./table";
import InfraSpecification from "./infra/specification";
import {MachineType, MachineInfraType} from "./types";

const wrap = (ui: React.ReactNode) => render(<Grid>{ui}</Grid>);

describe("MachinesTable", () => {
  it("renders a partial machine (missing status) without throwing", () => {
    const partial = [{metadata: {name: "m1"}}] as MachineType[];
    wrap(<MachinesTable machines={partial} select={() => {}}/>);
    expect(screen.getByRole("button", {name: /select m1/i})).toBeInTheDocument();
  });

  it("renders a labeled empty state for an empty collection without crashing", () => {
    wrap(<MachinesTable machines={[]} select={() => {}}/>);
    expect(screen.getByText(/no machines found/i)).toBeInTheDocument();
  });
});

describe("Machine infra specification", () => {
  it("renders numCoresPerSocket === 0 as data (no stray 0 / no crash)", () => {
    const machine = {numCoresPerSocket: 0, numCPUs: 2} as MachineInfraType;
    render(<InfraSpecification machine={machine}/>);
    // the CPU-per-socket row is present and shows the zero value
    expect(screen.getByText("CPU Per Socket")).toBeInTheDocument();
    expect(screen.getByText("0")).toBeInTheDocument();
  });
});
