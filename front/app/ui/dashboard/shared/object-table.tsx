'use client';

import React from "react";
import {Table, UnstyledButton} from "@mantine/core";
import {ColumnDef} from "@/app/ui/dashboard/base/types";
import {EmptyState} from "@/app/ui/dashboard/shared/empty-state";

interface ObjectTableProps<T> {
  items: T[] | undefined;
  columns: ColumnDef<T>[];
  /** Stable unique row identity — NEVER the array index. */
  getRowKey: (item: T) => string;
  onSelect?: (item: T) => void;
  emptyLabel: string;
}

/**
 * Generic, config-driven table shared by every resource area. Renders a labeled empty
 * state for an empty/undefined collection, wraps content in a horizontal scroll
 * container, keys rows by a stable id, and makes the primary cell a keyboard-focusable
 * button when a row is selectable.
 */
export function ObjectTable<T>({
  items,
  columns,
  getRowKey,
  onSelect,
  emptyLabel,
}: ObjectTableProps<T>) {
  if (!items || items.length === 0) {
    return <EmptyState label={emptyLabel}/>;
  }

  return (
    <Table.ScrollContainer minWidth={640} type="native">
      <Table highlightOnHover striped>
        <Table.Thead>
          <Table.Tr>
            {columns.map((col) => (
              <Table.Th key={col.header} ta={col.align} w={col.width}>
                {col.header}
              </Table.Th>
            ))}
          </Table.Tr>
        </Table.Thead>
        <Table.Tbody>
          {items.map((item) => (
            <Table.Tr key={getRowKey(item)}>
              {columns.map((col, ci) => {
                const content = col.render(item);
                return (
                  <Table.Td key={col.header} ta={col.align}>
                    {ci === 0 && onSelect ? (
                      <UnstyledButton
                        component="button"
                        type="button"
                        onClick={() => onSelect(item)}
                        aria-label={`Select ${getRowKey(item)}`}
                        className="cursor-pointer hover:opacity-70"
                      >
                        {content}
                      </UnstyledButton>
                    ) : (
                      content
                    )}
                  </Table.Td>
                );
              })}
            </Table.Tr>
          ))}
        </Table.Tbody>
      </Table>
    </Table.ScrollContainer>
  );
}
