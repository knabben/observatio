import ClusterInfo from '@/app/ui/dashboard/components/ClusterInfo'
import Versions from '@/app/ui/dashboard/components/Versions'
import Summary from '@/app/ui/dashboard/components/Summary'
import ClusterClass from '@/app/ui/dashboard/components/ClusterClass'

import { getClusterInformation } from "@/app/lib/data";
import { getClusterSummary } from "@/app/lib/data";
import { getComponentsVersion } from "@/app/lib/data";
import { getClusterClasses } from "@/app/lib/data";

import Link from 'next/link';
import { Card, Grid, GridCol, Text, Divider, Title, Space } from '@mantine/core';

export default async function Dashboard() {
  const clusterInfo = await getClusterInformation()
  const componentsVersion = await getComponentsVersion()
  const clusterSummary = await getClusterSummary()
  const clusterClasses = await getClusterClasses()

  return (
    <main>
      <Grid grow justify="center" align="top">
        <GridCol span={12}>
          <Link href="/dashboard">
            <Title className="hidden md:block" order={2}>
              Dashboard
            </Title>
          </Link>
        </GridCol>

        <GridCol span={5}>
          <Card shadow="md" padding="lg" radius="md" withBorder>
            <Text tt="uppercase"  fw={600} c="teal.8" ta="center">Clusters Health</Text>
            <Divider my="sm" variant="dashed" />
            <Summary clusterSummary={clusterSummary} />
          </Card>
          <Space h="md"/>
          <Card shadow="md"  radius="md" withBorder>
            <Text tt="uppercase"  fw={600} c="teal.8" ta="center">Cluster Information</Text>
            <Divider my="sm" variant="dashed" />
            <ClusterInfo clusterInfo={clusterInfo}/>
          </Card>
        </GridCol>

        <GridCol span={7}>
          <Card shadow="md"  radius="md" withBorder>
            <Text tt="uppercase"  fw={600} c="teal.8" ta="center">Cluster Class</Text>
            <Divider my="sm" variant="dashed" />
            <ClusterClass clusterClass={clusterClasses} />
          </Card>
          <Space h="md"/>
          <Card shadow="md"  radius="md" withBorder>
            <Text tt="uppercase"  fw={600} c="teal.8" ta="center">Components version</Text>
            <Divider my="sm" variant="dashed" />
            <Versions components={componentsVersion} />
          </Card>
        </GridCol>
      </Grid>
    </main>
  );
}

export const dynamic = 'error'
