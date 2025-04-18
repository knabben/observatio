import ClusterInfo from '@/app/ui/dashboard/components/ClusterInfo'
import Versions from '@/app/ui/dashboard/components/Versions'
import Summary from '@/app/ui/dashboard/components/Summary'
import ClusterClass from '@/app/ui/dashboard/components/ClusterClass'

import { getClusterInformation } from "@/app/lib/data";
import { getClusterSummary } from "@/app/lib/data";
import { getComponentsVersion } from "@/app/lib/data";
import { getClusterClasses } from "@/app/lib/data";

import { Suspense } from 'react';
import { Card, Grid, GridCol, Text, Divider } from '@mantine/core';
import Loading from "@/app/dashboard/loading";

export default async function Dashboard() {
  const clusterInfo = await getClusterInformation()
  const componentsVersion = await getComponentsVersion()
  const clusterSummary = await getClusterSummary()
  const clusterClasses = await getClusterClasses()

  return (
    <main>
      <Grid grow>
        <GridCol span={7}>
          <Card shadow="md"  radius="md" withBorder>
            <Text tt="uppercase"  fw={600} c="teal.8" ta="center">Cluster Class</Text>
            <Divider my="sm" variant="dashed" />
            <ClusterClass clusterClass={clusterClasses} />
          </Card>
        </GridCol>
        <GridCol span={5}>
          <Card shadow="md" padding="sm" radius="md" withBorder>
            <Suspense fallback={<Loading />}>
              <Text tt="uppercase"  fw={600} c="teal.8" ta="center">Clusters Health</Text>
              <Divider my="sm" variant="dashed" />
              <Summary clusterSummary={clusterSummary} />
            </Suspense>
          </Card>

        </GridCol>
      </Grid>
      <Grid grow>
        <GridCol span={7}>
          <Card shadow="md"  radius="md" withBorder>
            <Suspense fallback={<Loading />}>
              <Text tt="uppercase"  fw={600} c="teal.8" ta="center">Components version</Text>
              <Divider my="sm" variant="dashed" />
              <Versions components={componentsVersion} />
            </Suspense>
          </Card>
        </GridCol>
        <GridCol span={5}>
          <Card shadow="md"  radius="md" withBorder>
            <Text tt="uppercase"  fw={600} c="teal.8" ta="center">Cluster Information</Text>
            <Divider my="sm" variant="dashed" />
            <ClusterInfo clusterInfo={clusterInfo}/>
          </Card>
        </GridCol>
      </Grid>
    </main>
  );
}
