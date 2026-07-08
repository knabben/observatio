'use client';

import React from 'react';
import {Card, Group, Stack, Text, Title} from '@mantine/core';
import {IconAlertTriangle, IconCircleCheck} from '@tabler/icons-react';
import {sourceSans400} from '@/fonts';
import Header from '@/app/ui/dashboard/utils/header';
import {STATUS_COLORS} from '@/app/styles/theme';
import {Category, HealthRollup} from '@/app/ui/dashboard/shared/use-day2-ops';

const CATEGORY_TITLES: Record<Category, string> = {
  cluster: 'Clusters',
  machine_deployment: 'Machine Deployments',
  machine: 'Machines',
};

interface HealthRollupCardProps {
  rollup: HealthRollup;
}

/** Per-category healthy/degraded/failed rollup card for the Day-2 Ops landing screen (FR-002). */
export function HealthRollupCard({rollup}: HealthRollupCardProps) {
  const allClear = !rollup.unavailable && rollup.degraded === 0 && rollup.failed === 0;

  return (
    <Card shadow="md" className={sourceSans400.className} padding="lg" radius="md" withBorder>
      <Header title={CATEGORY_TITLES[rollup.category]}/>
      {rollup.unavailable ? (
        <Group gap="xs" justify="center">
          <IconAlertTriangle size={18} color="var(--mantine-color-red-6)"/>
          <Text c="red.6" fw={600}>Data unavailable</Text>
        </Group>
      ) : allClear ? (
        <Group gap="xs" justify="center">
          <IconCircleCheck size={18} color={`var(--mantine-color-${STATUS_COLORS.healthy}-6)`}/>
          <Text c={`${STATUS_COLORS.healthy}.6`} fw={600}>All clear — {rollup.healthy} healthy</Text>
        </Group>
      ) : (
        <Group justify="center" gap="xl">
          <Stack gap={0} align="center">
            <Title order={3} c={`${STATUS_COLORS.healthy}.6`}>{rollup.healthy}</Title>
            <Text size="sm">Healthy</Text>
          </Stack>
          <Stack gap={0} align="center">
            <Title order={3} c={`${STATUS_COLORS.degraded}.6`}>{rollup.degraded}</Title>
            <Text size="sm">Degraded</Text>
          </Stack>
          <Stack gap={0} align="center">
            <Title order={3} c={`${STATUS_COLORS.notready}.6`}>{rollup.failed}</Title>
            <Text size="sm">Failed</Text>
          </Stack>
        </Group>
      )}
    </Card>
  );
}
