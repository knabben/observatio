'use client';

import ClusterInfo from '@/app/ui/dashboard/components/ClusterInfo'
import { Card, Grid, Text, Title} from '@mantine/core';

export default async function Home() {
  return (
    <main>
      <Grid grow>
        <Grid.Col span={12}>
          <Title order={2} tt="capitalize">
            Dashboard
          </Title>
        </Grid.Col>
      </Grid>
      <Grid grow>
        <Grid.Col span={4}>
          <Card shadow="sm" padding="lg" radius="md" withBorder>
            <Text>
              CAPI Versions
            </Text>
          </Card>
        </Grid.Col>
        <Grid.Col span={8}>
          <Card shadow="sm" padding="lg" radius="md" withBorder>
            <ClusterInfo />
          </Card>
        </Grid.Col>
      </Grid>
    </main>
  );
}
