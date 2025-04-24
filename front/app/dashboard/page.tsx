import ClusterInfo from '@/app/ui/dashboard/components/dashboard/clusterinfo'
import ClusterVersions from '@/app/ui/dashboard/components/dashboard/versions'
import ClusterSummary from '@/app/ui/dashboard/components/dashboard/summary'
import ClusterClass from '@/app/ui/dashboard/components/dashboard/clusterclass'

import Link from 'next/link';
import { Grid, GridCol, Title, Space } from '@mantine/core';

export default async function Dashboard() {
  return (
    <Grid grow justify="center" align="top">
      <GridCol span={12}>
        <Link href="/dashboard">
          <Title className="hidden md:block" order={2}>
            Dashboard
          </Title>
        </Link>
      </GridCol>
      <GridCol span={5}>
        <ClusterSummary />
        <Space h="md"/>
        <ClusterInfo />
      </GridCol>
      <GridCol span={7}>
        <ClusterClass />
        <Space h="md"/>
        <ClusterVersions />
      </GridCol>
    </Grid>
  );
}

export const dynamic = 'error'
