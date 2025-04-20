import Link from 'next/link';

import { getClusterList } from "@/app/lib/data";
import ClusterLister from '@/app/ui/dashboard/components/ClusterLister'

import { Grid, GridCol, Title, Space, Input } from '@mantine/core';

export default async function Clusters() {
  const clusters = await getClusterList()
  return (
    <div>
      <main>
        <Grid grow>
          <GridCol h={60} span={8}>
            <Link href="/dashboard/clusters">
              <Title className="hidden md:block" order={2}>
                Clusters / cluster.x-k8s.io
              </Title>
            </Link>
          </GridCol>
          <GridCol span={4}>
            <Input placeholder="Search Clusters" />
          </GridCol>

          <GridCol span={12}>
            <ClusterLister clusterList={clusters} />
          </GridCol>
        </Grid>
      </main>
    </div>
  )
}