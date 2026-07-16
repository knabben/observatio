'use client';

import React from 'react';
import {Badge, Card, Group, Stack, Text} from '@mantine/core';
import {IconAlertTriangle} from '@tabler/icons-react';
import {sourceSans400} from '@/fonts';
import Header from '@/app/ui/dashboard/utils/header';
import {StatusIndicator} from '@/app/ui/dashboard/shared/status-indicator';
import {StatusState, toStatusState} from '@/app/ui/dashboard/shared/status';
import {BackupHealth, ClusterBackupCoverage} from '@/app/ui/dashboard/shared/use-day2-ops';

/** A cluster with no backup at all is worse than one with a merely stale backup (spec.md Edge Cases). */
function coverageStatus(c: ClusterBackupCoverage): StatusState {
  if (!c.covered) return 'notready';
  if (c.stale) return 'degraded';
  return 'healthy';
}

function coverageLabel(c: ClusterBackupCoverage): string {
  if (!c.covered) return 'No backup coverage';
  if (c.stale) return 'Stale';
  return 'On-time';
}

interface BackupHealthCardProps {
  health: BackupHealth;
}

/** Backup Health summary for the Day-2 Ops landing screen (008/US1) — storage-location
 * reachability, per-cluster backup staleness against the configured RPO, and restore activity.
 * Rendered alongside the existing HealthRollupCards, but with its own tailored shape rather than
 * forced into the healthy/degraded/failed 3-bucket model (research.md R3). */
export function BackupHealthCard({health}: BackupHealthCardProps) {
  return (
    <Card shadow="md" className={sourceSans400.className} padding="lg" radius="md" withBorder>
      <Header title="Backup Health"/>
      {!health.available ? (
        <Group gap="xs" justify="center">
          <IconAlertTriangle size={18} color="var(--mantine-color-gray-6)"/>
          <Text c="dimmed" fw={600}>Velero not available</Text>
        </Group>
      ) : (
        <Stack gap="sm">
          {health.restoresInProgress > 0 && (
            <Badge color="blue" variant="light" fullWidth>
              {health.restoresInProgress} restore{health.restoresInProgress === 1 ? '' : 's'} in progress
            </Badge>
          )}
          {health.storageLocations.length > 0 && (
            <Stack gap={4}>
              <Text size="xs" c="dimmed" fw={600}>Storage locations</Text>
              {health.storageLocations.map((loc) => (
                <Group key={`${loc.namespace}/${loc.name}`} justify="space-between" gap="xs" wrap="nowrap">
                  <Text size="sm">{loc.name}</Text>
                  <StatusIndicator state={toStatusState(loc.reachable)} dotOnly/>
                </Group>
              ))}
            </Stack>
          )}
          {health.clusterCoverage.length > 0 && (
            <Stack gap={4}>
              <Text size="xs" c="dimmed" fw={600}>Cluster coverage</Text>
              {health.clusterCoverage.map((c) => (
                <Group key={`${c.clusterRef.namespace}/${c.clusterRef.name}`} justify="space-between" gap="xs" wrap="nowrap">
                  <Stack gap={0}>
                    <Text size="sm">{c.clusterRef.name}</Text>
                    <Text size="xs" c="dimmed">
                      {coverageLabel(c)}{c.mostRecentBackupAge ? ` · ${c.mostRecentBackupAge}` : ''}
                    </Text>
                  </Stack>
                  <StatusIndicator state={coverageStatus(c)} dotOnly/>
                </Group>
              ))}
            </Stack>
          )}
        </Stack>
      )}
    </Card>
  );
}
