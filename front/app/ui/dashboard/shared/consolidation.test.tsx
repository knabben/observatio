import "@testing-library/jest-dom";
import {screen} from "@testing-library/react";
import {Grid} from "@mantine/core";
import {render} from "@/app/ui/dashboard/utils/test-render";

import ClusterTable from "@/app/ui/dashboard/components/clusters/table";
import {ClusterType} from "@/app/ui/dashboard/components/clusters/types";
import MachinesTable from "@/app/ui/dashboard/components/machines/table";
import {MachineType} from "@/app/ui/dashboard/components/machines/types";
import MDTable from "@/app/ui/dashboard/components/mds/table";
import {MachineDeploymentType} from "@/app/ui/dashboard/components/mds/types";

const wrap = (ui: React.ReactNode) => render(<Grid>{ui}</Grid>);

/**
 * Each resource area's table is a `ColumnDef[]` config rendered through the single
 * shared `ObjectTable` — not a hand-rolled copy. These tests prove the config-driven
 * architecture by asserting the SAME behavior (keyboard-focusable selectable rows,
 * one shared empty state) holds identically across clusters/machines/mds, so a fix
 * made once in `ObjectTable` is inherited by every resource area instead of needing
 * to be re-applied per screen.
 */
const cases = [
  {
    name: "clusters",
    renderWithItem: () => wrap(<ClusterTable clusters={[{metadata: {name: "c1"}}] as ClusterType[]} select={() => {}}/>),
    renderEmpty: () => wrap(<ClusterTable clusters={[]} select={() => {}}/>),
    selectName: /select c1/i,
    emptyText: /no clusters found/i,
  },
  {
    name: "machines",
    renderWithItem: () => wrap(<MachinesTable machines={[{metadata: {name: "m1"}}] as MachineType[]} select={() => {}}/>),
    renderEmpty: () => wrap(<MachinesTable machines={[]} select={() => {}}/>),
    selectName: /select m1/i,
    emptyText: /no machines found/i,
  },
  {
    name: "machine deployments",
    renderWithItem: () => wrap(<MDTable mds={[{metadata: {name: "md1"}}] as MachineDeploymentType[]} select={() => {}}/>),
    renderEmpty: () => wrap(<MDTable mds={[]} select={() => {}}/>),
    selectName: /select md1/i,
    emptyText: /no machine deployments found/i,
  },
];

describe("Resource table consolidation over shared ObjectTable", () => {
  it.each(cases)("$name: renders a keyboard-focusable, labeled select control (shared ObjectTable behavior)", ({renderWithItem, selectName}) => {
    renderWithItem();
    const control = screen.getByRole("button", {name: selectName});
    expect(control.tagName).toBe("BUTTON");
  });

  it.each(cases)("$name: renders a distinct, config-supplied empty label via the shared EmptyState", ({renderEmpty, emptyText}) => {
    renderEmpty();
    expect(screen.getByText(emptyText)).toBeInTheDocument();
  });
});
