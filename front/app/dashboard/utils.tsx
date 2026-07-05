/* eslint-disable @typescript-eslint/no-explicit-any */

/**
 * Returns the first item whose `metadata.name` contains `query` (case-insensitive),
 * or `undefined`. Null-safe: items missing `metadata`/`name` are skipped rather than
 * throwing.
 */
export function FilterItems(query: string, items: any[]) {
  const q = query.toLowerCase();
  return items
    .filter((i: { metadata?: {name?: string} }) =>
      (i?.metadata?.name ?? '').toLowerCase().includes(q))
    .at(0);
}
