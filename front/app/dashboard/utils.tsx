/* eslint-disable @typescript-eslint/no-explicit-any */

/**
 * Returns the first item whose `metadata.name` contains `query` (case-insensitive),
 * or `undefined`. Null-safe: items missing `metadata`/`name` are skipped rather than
 * throwing.
 */
export function FilterItems(query: string, items: any[]) {
  return FilterItemsByName(query, items).at(0);
}

/**
 * Returns every item whose `metadata.name` contains `query` (case-insensitive).
 * Null-safe: items missing `metadata`/`name` are skipped rather than throwing.
 * An empty query matches everything (no filter applied).
 */
export function FilterItemsByName<T extends {metadata?: {name?: string}}>(query: string, items: T[]): T[] {
  const q = query.trim().toLowerCase();
  if (!q) return items;
  return items.filter((i) => (i?.metadata?.name ?? '').toLowerCase().includes(q));
}
