'use client';

import React, {useState} from 'react';
import {Badge, Button, Card, Group, Stack, Text} from '@mantine/core';
import Link from 'next/link';
import {getDay2OpsDetail} from '@/app/lib/data';
import {DebugLayerName, DebugLayerStatus, DebugPath} from '@/app/ui/dashboard/shared/use-day2-ops';
import {NodeAccessPanel} from '@/app/ui/dashboard/logs/node-access-panel';

const LAYER_LABELS: Record<DebugLayerName, string> = {
  conditions: 'Object conditions',
  phase: 'Machine phase',
  provider_resource: 'Provider resource',
  controller_activity: 'Controller activity',
};

const LAYER_STATUS_COLORS: Record<DebugLayerStatus, string> = {
  ok: 'green',
  implicated: 'red',
  inconclusive: 'gray',
};

/** Best-effort provider-Kind -> controller mapping, derived from the provider_resource layer's
 * source (e.g. "DockerMachine/worker-0"). Falls back to CAPI core when not determinable. */
function controllerForProviderKind(source: string): {namespace: string; deployment: string} {
  if (source.startsWith('DockerMachine')) return {namespace: 'capd-system', deployment: 'capd-controller-manager'};
  if (source.startsWith('VSphereMachine')) return {namespace: 'capv-system', deployment: 'capv-controller-manager'};
  return {namespace: 'capi-system', deployment: 'capi-controller-manager'};
}

interface DebuggingPathProps {
  path: DebugPath;
}

/**
 * Renders the ordered, labeled debugging path for one unhealthy object (FR-004, FR-005), inline on
 * the landing screen using the WS-pushed capped evidence. "Show full evidence" fetches the
 * uncapped path from GET /api/day2ops/detail on demand (research.md R9) — a deep dive, not a
 * requirement to see which layer is at fault in the first place. When the controller_activity
 * layer is implicated, a "View controller logs" deep-dive links into the Logs destination
 * (User Story 5, FR-019); "Node access" toggles the SSH connection-instructions deep-dive.
 */
export function DebuggingPath({path}: DebuggingPathProps) {
  const [detail, setDetail] = useState<DebugPath | null>(null);
  const [expanding, setExpanding] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [showNodeAccess, setShowNodeAccess] = useState(false);

  const layers = detail?.layers ?? path.layers;
  const summary = detail?.summary ?? path.summary;

  const controllerActivityLayer = layers.find((l) => l.layer === 'controller_activity');
  const providerLayer = layers.find((l) => l.layer === 'provider_resource');
  const controller = controllerForProviderKind(providerLayer?.source ?? '');

  async function expand() {
    setExpanding(true);
    setError(null);
    try {
      const response = await getDay2OpsDetail({
        group: path.objectRef.group,
        version: path.objectRef.version,
        resource: path.objectRef.resource,
        namespace: path.objectRef.namespace,
        name: path.objectRef.name,
      });
      setDetail(response.path);
    } catch {
      setError('Failed to load the full evidence list.');
    } finally {
      setExpanding(false);
    }
  }

  return (
    <Card withBorder padding="sm" radius="md">
      <Text size="sm" fw={600} mb="xs">{summary}</Text>
      <Stack gap={4}>
        {layers.map((layer, index) => (
          <Group key={layer.layer} gap={6} wrap="nowrap" align="flex-start">
            <Badge size="xs" circle color={LAYER_STATUS_COLORS[layer.status]}>{index + 1}</Badge>
            <Text size="xs">
              <Text span fw={600}>{LAYER_LABELS[layer.layer]}: </Text>
              {layer.evidence.length > 0 ? layer.evidence.join('; ') : '—'}
            </Text>
          </Group>
        ))}
      </Stack>
      <Group gap="xs" mt="xs">
        {!detail && (
          <Button size="xs" variant="subtle" onClick={expand} loading={expanding}>
            Show full evidence
          </Button>
        )}
        {controllerActivityLayer?.status === 'implicated' && (
          <Button
            component={Link}
            href={`/dashboard/logs?namespace=${controller.namespace}&deployment=${controller.deployment}`}
            size="xs"
            variant="subtle"
          >
            View controller logs
          </Button>
        )}
        <Button size="xs" variant="subtle" onClick={() => setShowNodeAccess((v) => !v)}>
          {showNodeAccess ? 'Hide node access' : 'Node access'}
        </Button>
      </Group>
      {error && <Text size="xs" c="red.6" mt="xs">{error}</Text>}
      {showNodeAccess && (
        <Card withBorder padding="xs" radius="sm" mt="xs">
          <NodeAccessPanel objectRef={path.objectRef}/>
        </Card>
      )}
    </Card>
  );
}
