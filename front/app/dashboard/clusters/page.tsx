import Link from 'next/link';

import { getClusterList } from "@/app/lib/data";
import ClusterLister from '@/app/ui/dashboard/components/ClusterLister'
import { FilterItems } from "@/app/dashboard/utils";

import Search from '@/app/ui/dashboard/search'
import { Grid, GridCol, Title } from '@mantine/core';

export default async function Clusters(props: {
  searchParams?: Promise<{
    query?: string;
  }>;
}) {
  const searchParams = await props.searchParams;
  const query = searchParams?.query || '';
  const clusters = await getClusterList(query)

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
            <Search placeholder="Cluster name"/>
          </GridCol>
          <GridCol span={12}>
            <ClusterLister clusterList={FilterItems(query, clusters)} />
          </GridCol>
        </Grid>
      </main>
    </div>
  )
}