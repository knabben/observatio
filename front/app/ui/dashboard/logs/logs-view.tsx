'use client';

import React, {useCallback, useEffect, useState} from 'react';
import {useSearchParams} from 'next/navigation';
import {Button, Card, Group, ScrollArea, Text} from '@mantine/core';
import {IconAlertTriangle} from '@tabler/icons-react';
import {getControllerLogs} from '@/app/lib/data';

export interface ControllerRef {
  namespace: string;
  deployment: string;
}

/** Known CAPI core + provider controllers (research.md R7's ControllerNamespaces). */
const KNOWN_CONTROLLERS: ControllerRef[] = [
  {namespace: 'capi-system', deployment: 'capi-controller-manager'},
  {namespace: 'capd-system', deployment: 'capd-controller-manager'},
  {namespace: 'capv-system', deployment: 'capv-controller-manager'},
];

interface LogsViewProps {
  /** Preselects a controller, e.g. when arriving via a debugging-path deep-dive (FR-019). */
  initialController?: ControllerRef;
}

/**
 * Controller selector + log pane (User Story 5): streams a controller's Pod log output via
 * GET /api/logs/controller — the same data `kubectl logs` provides (FR-020). A fetch failure
 * (no Pod, no retained log history) shows an explicit "logs unavailable" state (FR-023) rather
 * than a blank pane.
 */
export function LogsView({initialController}: LogsViewProps) {
  // Query params take precedence when arriving via a debugging-path deep-dive link
  // (front/app/ui/dashboard/components/ops/debugging-path.tsx); a page-level Suspense boundary
  // is required around this component because output:"export" (static export, see
  // front/next.config.ts) has no server to resolve search params at request time.
  const searchParams = useSearchParams();
  const fromQuery: ControllerRef | undefined = (() => {
    const namespace = searchParams.get('namespace');
    const deployment = searchParams.get('deployment');
    return namespace && deployment ? {namespace, deployment} : undefined;
  })();

  const [selected, setSelected] = useState<ControllerRef | null>(fromQuery ?? initialController ?? null);
  const [logs, setLogs] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  const load = useCallback(async (controller: ControllerRef) => {
    setLoading(true);
    setError(null);
    setLogs(null);
    try {
      const text = await getControllerLogs(controller.namespace, controller.deployment);
      setLogs(text);
    } catch {
      setError('Logs unavailable — the controller may not be running or has no retained log history.');
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    if (selected) void load(selected);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  function select(controller: ControllerRef) {
    setSelected(controller);
    void load(controller);
  }

  return (
    <Card withBorder padding="md" radius="md">
      <Group gap="xs" mb="sm">
        {KNOWN_CONTROLLERS.map((c) => (
          <Button
            key={`${c.namespace}/${c.deployment}`}
            size="xs"
            variant={selected?.deployment === c.deployment && selected?.namespace === c.namespace ? 'filled' : 'default'}
            onClick={() => select(c)}
          >
            {c.deployment}
          </Button>
        ))}
      </Group>
      {loading && <Text size="sm" c="dimmed">Loading…</Text>}
      {error && (
        <Group gap="xs">
          <IconAlertTriangle size={16} color="var(--mantine-color-red-6)"/>
          <Text size="sm" c="red.6">{error}</Text>
        </Group>
      )}
      {logs != null && (
        <ScrollArea h={400}>
          <Text
            component="pre"
            size="xs"
            style={{whiteSpace: 'pre-wrap', fontFamily: 'monospace'}}
          >
            {logs}
          </Text>
        </ScrollArea>
      )}
    </Card>
  );
}
