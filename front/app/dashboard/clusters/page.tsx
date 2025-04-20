import Link from 'next/link';

import { getClusterList } from "@/app/lib/data";
import ClusterLister from '@/app/ui/dashboard/components/ClusterLister'
import { Suspense } from 'react';
import { Loader } from '@mantine/core';
import Search from '@/app/ui/dashboard/search'
import { Grid, GridCol, Title } from '@mantine/core';

export default async function Clusters(props: {
  searchParams?: Promise<{
    query?: string;
    page?: string;
  }>;
}) {
  const searchParams = await props.searchParams;
  const query = searchParams?.query || '';
  let clusters = await getClusterList()
  const currentPage = Number(searchParams?.page) || 1;

    if (query != "") {
      clusters = clusters.filter((i: { name: string; }) =>
        i.name.toLowerCase().includes(query.toLowerCase()));
  }

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
            <Suspense key={query + currentPage} fallback={<Loader />}>
              <ClusterLister clusterList={clusters} />
            </Suspense>
          </GridCol>
        </Grid>
      </main>
    </div>
  )
}