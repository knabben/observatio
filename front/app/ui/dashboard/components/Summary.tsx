'use client';

import { Paper, Space, Text } from '@mantine/core';

type Summary = {
  failed: bigint
  provisioned: bigint,
}

// Summary : resume clusters response
export default function Summary({
  clusterSummary,
}: {
  clusterSummary: Summary
}) {
  return (
    <div>
      <Paper ta="center" shadow="xs" withBorder p="xl">
        Provisioned
        <Space h="xs" />
        <Text fw={700} size="xl"  c="teal.4">{clusterSummary.provisioned}</Text>
      </Paper>
      <Paper ta="center" shadow="xs" withBorder p="xl">
        Failing
        <Space h="xs" />
        <Text fw={700} size="xl" c="red">{clusterSummary.failed}</Text>
      </Paper>
    </div>
  );
}