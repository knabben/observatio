import "@testing-library/jest-dom";
import {screen} from "@testing-library/react";
import {Grid} from "@mantine/core";
import {render} from "@/app/ui/dashboard/utils/test-render";
import MachineSetTable from "./table";
import {MachineSetType} from "./types";

const wrap = (ui: React.ReactNode) => render(<Grid>{ui}</Grid>);

describe("MachineSetTable", () => {
  it("renders a partial MachineSet (missing status) without throwing", () => {
    const partial = [{metadata: {name: "ms1"}}] as MachineSetType[];
    wrap(<MachineSetTable mss={partial} select={() => {}}/>);
    expect(screen.getByRole("button", {name: /select ms1/i})).toBeInTheDocument();
  });

  it("shows unknown availability as 'Unknown', not failed", () => {
    const partial = [{metadata: {name: "ms1"}}] as MachineSetType[];
    wrap(<MachineSetTable mss={partial} select={() => {}}/>);
    expect(screen.getByRole("img", {name: "Unknown"})).toBeInTheDocument();
    expect(screen.queryByRole("img", {name: "Not ready"})).not.toBeInTheDocument();
  });

  it("renders a labeled empty state for an empty collection without crashing", () => {
    wrap(<MachineSetTable mss={[]} select={() => {}}/>);
    expect(screen.getByText(/no machine sets found/i)).toBeInTheDocument();
  });
});
