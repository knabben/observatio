import {TreeNodeData} from '@mantine/core';

function isExpandable(value: unknown): value is Record<string, unknown> | unknown[] {
  return value !== null && typeof value === 'object';
}

function isEmpty(value: Record<string, unknown> | unknown[]): boolean {
  return Array.isArray(value) ? value.length === 0 : Object.keys(value).length === 0;
}

function formatScalar(value: unknown): string {
  if (value === null || value === undefined) return '—';
  if (Array.isArray(value)) return '[]';
  if (typeof value === 'object') return '{}';
  return String(value);
}

/**
 * Converts an arbitrary JSON value (the raw Kubernetes object returned by GET /api/raw) into
 * Mantine `Tree`'s `TreeNodeData[]` shape. Object/array keys with content become expandable
 * nodes; scalars (and empty objects/arrays) render as a single "key: value" leaf so an empty
 * `{}` doesn't show a misleading expand affordance with nothing inside it.
 */
export function toTreeData(value: unknown, path = ''): TreeNodeData[] {
  if (!isExpandable(value)) return [];

  const entries: [string, unknown][] = Array.isArray(value)
    ? value.map((item, i) => [`[${i}]`, item])
    : Object.entries(value);

  return entries.map(([key, val]) => {
    const nodePath = path ? `${path}.${key}` : key;
    if (isExpandable(val) && !isEmpty(val)) {
      return {label: key, value: nodePath, children: toTreeData(val, nodePath)};
    }
    return {label: `${key}: ${formatScalar(val)}`, value: nodePath};
  });
}
