'use client';

import React, {useEffect, useState} from 'react';
import {Group, RenderTreeNodePayload, Text, Tree, useTree} from '@mantine/core';
import {IconChevronRight} from '@tabler/icons-react';
import {getRawObject, ResourceGVR} from '@/app/lib/data';
import {toTreeData} from '@/app/ui/dashboard/shared/to-tree-data';
import {CenteredLoader} from '@/app/ui/dashboard/utils/loader';
import {ErrorState} from '@/app/ui/dashboard/shared/error-state';

/** Renders each node with a rotating chevron when it has children, so it's visually clear which
 * items can be expanded — Mantine's Tree has no such indicator by default. */
function renderNode({node, expanded, hasChildren, elementProps}: RenderTreeNodePayload) {
  return (
    <Group gap={4} {...elementProps}>
      {hasChildren ? (
        <IconChevronRight
          size={14}
          style={{
            transform: expanded ? 'rotate(90deg)' : 'rotate(0deg)',
            transition: 'transform 150ms ease',
          }}
        />
      ) : (
        <span style={{width: 14, display: 'inline-block'}}/>
      )}
      <Text size="sm">{node.label}</Text>
    </Group>
  );
}

/**
 * Renders the complete underlying object (every field the backend returns, not the curated
 * subset the rest of the screen shows) as an expandable/collapsible tree, fetched on demand via
 * GET /api/raw. Re-fetches when `resourceVersion` changes — driven by the same live WebSocket
 * stream already feeding the Specification tab, not independent polling (research.md R2).
 */
export function ObjectTree({
  gvr, namespace, name, resourceVersion,
}: {
  gvr: ResourceGVR;
  namespace: string;
  name: string;
  resourceVersion?: string;
}) {
  const [data, setData] = useState<unknown>(null);
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const tree = useTree();

  useEffect(() => {
    let cancelled = false;
    setIsLoading(true);
    setError(null);
    getRawObject({...gvr, namespace, name})
      .then((response) => {
        if (!cancelled) setData(response);
      })
      .catch((err) => {
        console.error('Failed to load raw object:', err);
        if (!cancelled) setError('Failed to load the complete object.');
      })
      .finally(() => {
        if (!cancelled) setIsLoading(false);
      });
    return () => {
      cancelled = true;
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [gvr.group, gvr.version, gvr.resource, namespace, name, resourceVersion]);

  if (isLoading) {
    return <CenteredLoader/>;
  }
  if (error) {
    return <ErrorState message={error}/>;
  }

  return (
    <div style={{maxHeight: 600, overflowY: 'auto'}}>
      <Tree data={toTreeData(data)} tree={tree} renderNode={renderNode}/>
    </div>
  );
}
