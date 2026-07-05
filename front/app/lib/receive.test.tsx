import {describe, it, expect} from "@jest/globals";
import {receiveAndPopulate, byMetadataName, WSOperationType} from "./websocket";

const item = (name: string) => ({metadata: {name}});

describe("receiveAndPopulate", () => {
  it("returns the current list UNCHANGED on an empty/malformed frame (no wipe)", () => {
    const items = [item("a"), item("b")];
    expect(receiveAndPopulate(null, items)).toBe(items);
    expect(receiveAndPopulate({type: "ADDED"}, items)).toBe(items);
    expect(receiveAndPopulate({}, items)).toBe(items);
  });

  it("adds a new item on ADDED/MODIFIED", () => {
    const items = [item("a")];
    const next = receiveAndPopulate({type: WSOperationType.ADDED, data: item("b")}, items);
    expect(next.map((i: {metadata: {name: string}}) => i.metadata.name)).toEqual(["a", "b"]);
  });

  it("replaces an existing item by metadata.name on MODIFIED", () => {
    const items = [{metadata: {name: "a"}, v: 1}];
    const next = receiveAndPopulate({type: WSOperationType.MODIFIED, data: {metadata: {name: "a"}, v: 2}}, items);
    expect(next).toHaveLength(1);
    expect(next[0].v).toBe(2);
  });

  it("removes an item on DELETED", () => {
    const items = [item("a"), item("b")];
    const next = receiveAndPopulate({type: WSOperationType.DELETED, data: item("a")}, items);
    expect(next.map((i: {metadata: {name: string}}) => i.metadata.name)).toEqual(["b"]);
  });
});

describe("byMetadataName", () => {
  it("sorts by name and tolerates missing metadata/name", () => {
    const list = [{metadata: {name: "b"}}, {}, {metadata: {name: "a"}}];
    const sorted = [...list].sort(byMetadataName);
    // the item with no name sorts first (empty string), then a, then b
    expect(sorted.map((i) => (i as {metadata?: {name?: string}}).metadata?.name ?? "")).toEqual(["", "a", "b"]);
  });
});
