import Link from 'next/link';

import { Suspense } from 'react';

import { getClusterList } from "@/app/lib/data";
import ClusterLister from '@/app/ui/dashboard/components/ClusterLister'

import Search from '@/app/ui/dashboard/search'
import { Grid, GridCol, Title } from '@mantine/core';

export default async function Clusters() {
  const clusters = await getClusterList()

  return (
    <div>
      <main>
        <Grid justify="flex-end" align="flex-start">
          <GridCol h={60} span={8}>
            <Link href="/dashboard/clusters">
              <Title className="hidden md:block" order={2}>
                Clusters / cluster.x-k8s.io
              </Title>
            </Link>
          </GridCol>
          <GridCol span={4}>
            <Suspense fallback={<div>Loading...</div>}>
              <Search placeholder="Cluster name"/>
            </Suspense>
          </GridCol>
          <GridCol span={12}>
            <ClusterLister clusterList={clusters} />
          </GridCol>
        </Grid>
      </main>
    </div>
  )
}