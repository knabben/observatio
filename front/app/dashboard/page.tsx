import ClusterVersions from '@/app/ui/dashboard/components/dashboard/clusterversions'
import ClusterClassLister from '@/app/ui/dashboard/components/dashboard/clusterclass'
import ClusterHierarchy from '@/app/ui/dashboard/components/dashboard/clusterhierarchy'
import {OpsDashboard} from '@/app/ui/dashboard/components/ops/ops-dashboard'

import { sourceSans400 } from "@/fonts";
import Link from 'next/link';
import { Grid, GridCol, Title, Space } from '@mantine/core';

/**
 * Day-2 Operations landing view (spec 006): consolidated, live per-category health rollups
 * replace the old standalone ClusterSummary widget as the entry point. ClusterHierarchy and
 * ClusterVersions are retained below it as supplementary detail.
 */
export default async function Dashboard() {
  return (
    <Grid grow justify="center" align="top">
      <GridCol span={12}>
        <Link href="/dashboard">
          <Title className={sourceSans400.className} order={2}>
            Day-2 Operations
          </Title>
        </Link>
      </GridCol>
      <GridCol span={12}>
        <OpsDashboard/>
      </GridCol>
      <GridCol span={{base: 12, md: 5}}>
        <ClusterClassLister />
      </GridCol>
      <GridCol span={{base: 12, md: 7}}>
        <ClusterHierarchy />
        <Space h="md"/>
        <ClusterVersions />
      </GridCol>
    </Grid>
  );
}

export const dynamic = 'error'
