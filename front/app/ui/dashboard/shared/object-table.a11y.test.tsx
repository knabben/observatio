import "@testing-library/jest-dom";
import {fireEvent, screen} from "@testing-library/react";
import {render} from "@/app/ui/dashboard/utils/test-render";
import {ObjectTable} from "./object-table";
import {ColumnDef} from "@/app/ui/dashboard/base/types";

type Row = {id: string; name: string};

const columns: ColumnDef<Row>[] = [
  {header: "Name", render: (row) => row.name},
];

describe("ObjectTable accessibility", () => {
  it("makes a selectable row a keyboard-focusable, labeled control", () => {
    const onSelect = jest.fn();
    render(
      <ObjectTable
        items={[{id: "row-1", name: "Alpha"}]}
        columns={columns}
        getRowKey={(r) => r.id}
        onSelect={onSelect}
        emptyLabel="No rows found"
      />,
    );

    const control = screen.getByRole("button", {name: /select row-1/i});
    expect(control.tagName).toBe("BUTTON");

    control.focus();
    expect(control).toHaveFocus();

    fireEvent.click(control);
    expect(onSelect).toHaveBeenCalledWith({id: "row-1", name: "Alpha"});
  });

  it("does not render a selectable control when onSelect is omitted", () => {
    render(
      <ObjectTable
        items={[{id: "row-1", name: "Alpha"}]}
        columns={columns}
        getRowKey={(r) => r.id}
        emptyLabel="No rows found"
      />,
    );
    expect(screen.queryByRole("button")).not.toBeInTheDocument();
  });
});
