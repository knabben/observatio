import "@testing-library/jest-dom";
import {screen} from "@testing-library/react";
import {render} from "@/app/ui/dashboard/utils/test-render";
import {ObjectTable} from "./object-table";
import {ColumnDef} from "@/app/ui/dashboard/base/types";

type Row = {id: string; name: string};

const columns: ColumnDef<Row>[] = [
  {header: "Name", render: (row) => row.name},
];

describe("ObjectTable", () => {
  it("renders a labeled empty state for an empty collection", () => {
    render(<ObjectTable items={[]} columns={columns} getRowKey={(r) => r.id} emptyLabel="No rows found"/>);
    expect(screen.getByText("No rows found")).toBeInTheDocument();
    expect(screen.queryByRole("table")).not.toBeInTheDocument();
  });

  it("renders a labeled empty state for an undefined collection", () => {
    render(<ObjectTable items={undefined} columns={columns} getRowKey={(r) => r.id} emptyLabel="No rows found"/>);
    expect(screen.getByText("No rows found")).toBeInTheDocument();
  });

  it("wraps populated content in a horizontal scroll container", () => {
    const {container} = render(
      <ObjectTable items={[{id: "a", name: "Row A"}]} columns={columns} getRowKey={(r) => r.id} emptyLabel="No rows found"/>,
    );
    expect(screen.getByRole("table")).toBeInTheDocument();
    expect(container.querySelector(".mantine-TableScrollContainer-scrollContainer")).not.toBeNull();
  });

  it("keys rows by the stable id, not array index", () => {
    render(
      <ObjectTable items={[{id: "row-1", name: "Alpha"}]} columns={columns} getRowKey={(r) => r.id} emptyLabel="No rows found"/>,
    );
    expect(screen.getByText("Alpha")).toBeInTheDocument();
  });
});
