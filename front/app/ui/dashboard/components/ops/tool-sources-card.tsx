'use client';

import React from 'react';
import {Badge, Card, Group, Stack, Text} from '@mantine/core';
import {sourceSans400} from '@/fonts';
import Header from '@/app/ui/dashboard/utils/header';
import {CenteredLoader} from '@/app/ui/dashboard/utils/loader';
import {StatusIndicator} from '@/app/ui/dashboard/shared/status-indicator';
import {StatusState} from '@/app/ui/dashboard/shared/status';
import {useFetchState} from '@/app/ui/dashboard/shared/use-fetch-state';
import {emptyMCPSourcesResponse, getMCPSources, MCPHealthState, MCPSourceStatus} from '@/app/lib/data';

function healthToStatusState(state: MCPHealthState): StatusState {
  if (state === 'healthy') return 'healthy';
  if (state === 'unhealthy') return 'notready';
  return 'unknown';
}

/** Every capability name a conflict rejected, so a source's list can flag which of its own
 * would-be capabilities didn't survive registration (FR-007). */
function rejectedCapabilityNames(source: MCPSourceStatus, conflicts: {capabilityName: string; rejectedSource: string}[]): Set<string> {
  return new Set(conflicts.filter((c) => c.rejectedSource === source.name).map((c) => c.capabilityName));
}

/**
 * AI assistant tool source status for the Day-2 Ops landing screen (FR-003, SC-005,
 * specs/009-mcp-server-client-aggregator). Fetched independently of the WS-driven day2ops
 * payload via GET /api/mcp/sources (research.md R8) — tool-source health is assistant-tooling
 * state, not CAPI cluster state, so it isn't pushed over the existing WebSocket.
 */
export function ToolSourcesCard() {
  const {data, isLoading, error} = useFetchState(
    getMCPSources,
    emptyMCPSourcesResponse,
    'Failed to load AI assistant tool sources',
  );

  return (
    <Card shadow="md" className={sourceSans400.className} padding="lg" radius="md" withBorder>
      <Header title="AI Tool Sources"/>
      {isLoading ? (
        <CenteredLoader/>
      ) : error ? (
        <Text c="red.6" ta="center" fw={600}>{error}</Text>
      ) : (
        <Stack gap="sm">
          {data.sources.map((source) => {
            const rejected = rejectedCapabilityNames(source, data.conflicts);
            return (
              <Stack key={source.name} gap={4}>
                <Group justify="space-between" gap="xs" wrap="nowrap">
                  <Group gap="xs" wrap="nowrap">
                    <Text size="sm" fw={600}>{source.name}</Text>
                    <Badge size="xs" variant="light" color="gray">{source.kind}</Badge>
                  </Group>
                  <StatusIndicator state={healthToStatusState(source.health.state)} dotOnly/>
                </Group>
                {source.capabilities.length > 0 && (
                  <Text size="xs" c="dimmed">
                    {source.capabilities.join(', ')}
                  </Text>
                )}
                {rejected.size > 0 && (
                  <Text size="xs" c="orange.6">
                    Naming conflict: {Array.from(rejected).join(', ')}
                  </Text>
                )}
              </Stack>
            );
          })}
        </Stack>
      )}
    </Card>
  );
}
