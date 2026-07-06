import "@testing-library/jest-dom";
import {fireEvent, screen} from "@testing-library/react";
import {render} from "@/app/ui/dashboard/utils/test-render";
import {FilterItemsByName} from "@/app/dashboard/utils";
import Search from "@/app/ui/dashboard/search";

describe("Search", () => {
  it("filters the visible list as the user types (client-side, by name)", () => {
    const items = [{metadata: {name: "alpha"}}, {metadata: {name: "beta"}}];
    expect(FilterItemsByName("al", items)).toEqual([{metadata: {name: "alpha"}}]);
    expect(FilterItemsByName("", items)).toEqual(items);
    expect(FilterItemsByName("zzz", items)).toEqual([]);
  });

  it("calls onChange with the typed value, not a static pick-one selection", () => {
    const onChange = jest.fn();
    render(<Search value="" onChange={onChange}/>);
    const input = screen.getByRole("textbox", {name: /search by name/i});
    fireEvent.change(input, {target: {value: "clu"}});
    expect(onChange).toHaveBeenCalledWith("clu");
  });
});
