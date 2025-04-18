'use client';

import { Grid, Paper, Space, Text } from '@mantine/core';

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
    <div  style={{ resize: 'vertical', overflow: 'hidden', maxHeight: '100%' }}>
        <Grid ta="center">
          <Grid.Col span={6}>
            Provisioned
            <Space h="lg" />
            <Text fw={900} size="xl"  c="teal.4">{clusterSummary.provisioned}</Text>
          </Grid.Col>
          <Grid.Col span={6}>
            Failing
            <Space h="lg" />
            <Text fw={900} size="xl" c="red">{clusterSummary.failed}</Text>
          </Grid.Col>
        </Grid>
    </div>
  );
}