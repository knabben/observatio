import ClusterInfo from '@/app/ui/dashboard/components/ClusterInfo'
import Versions from '@/app/ui/dashboard/components/Versions'
import Summary from '@/app/ui/dashboard/components/Summary'

import { getClusterInformation } from "@/app/lib/data";
import { getClusterSummary } from "@/app/lib/data";
import { getComponentsVersion } from "@/app/lib/data";

import { Suspense } from 'react';
import { Card, Grid, GridCol, Text, Divider } from '@mantine/core';
import Loading from "@/app/dashboard/loading";

export default async function Dahsboard() {
  const clusterInfo = await getClusterInformation()
  const componentsVersion = await getComponentsVersion()
  const clusterSummary = await getClusterSummary()
  return (
    <main>
      <Grid grow>
        <GridCol span={8}>
          <Card shadow="md"  radius="md" withBorder>
            <Text tt="uppercase"  fw={600} c="teal.8" ta="center">Management Cluster</Text>
            <Divider my="sm" variant="dashed" />
          </Card>
        </GridCol>
        <GridCol span={4}>
          <Card shadow="md"  radius="md" withBorder>
            <Text tt="uppercase"  fw={600} c="teal.8" ta="center">Clusters Health</Text>
            <Divider my="sm" variant="dashed" />
            <Summary clusterSummary={clusterSummary} />
          </Card>
        </GridCol>
      </Grid>
      <Grid grow>
        <GridCol span={6}>
          <Card shadow="md"  radius="md" withBorder>
            <Suspense fallback={<Loading />}>
              <Text tt="uppercase"  fw={600} c="teal.8" ta="center">Components version</Text>
              <Divider my="sm" variant="dashed" />
              <Versions components={componentsVersion} />
            </Suspense>
          </Card>
        </GridCol>
        <GridCol span={6}>
          <Card shadow="md" padding="sm" radius="md" withBorder>
            <Suspense fallback={<Loading />}>
              <Text tt="uppercase"  fw={600} c="teal.8" ta="center">Cluster Information</Text>
              <Divider my="sm" variant="dashed" />
              <ClusterInfo clusterInfo={clusterInfo}/>
            </Suspense>
          </Card>
        </GridCol>
      </Grid>
    </main>
  );
}
