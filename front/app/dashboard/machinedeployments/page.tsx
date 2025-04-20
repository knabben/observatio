import Link from 'next/link';

import { getMachinesDeployments} from "@/app/lib/data";
import MachineDeploymentLister from "@/app/ui/dashboard/components/MachineDeploymentLister";
import { FilterItems } from "@/app/dashboard/utils";

import Search from '@/app/ui/dashboard/search'
import { Title, Grid, GridCol } from '@mantine/core';
import {Suspense} from "react";

export default async function MachineDeployments(props: {
  searchParams?: Promise<{
    query?: string;
  }>;
}) {
  const searchParams = await props.searchParams;
  const query = searchParams?.query || '';
  const mds = await getMachinesDeployments(query)

  return (
    <div>
      <main>
        <Grid justify="flex-end" align="flex-start">
          <GridCol h={60} span={8}>
            <Link href="/dashboard/machinedeployments">
              <Title className="hidden md:block" order={2}>
                Machine Deployments / cluster.x-k8s.io
              </Title>
            </Link>
          </GridCol>
          <GridCol span={4}>
            <Search placeholder="Machine deployment name"/>
          </GridCol>
          <GridCol span={12}>
            <Suspense fallback={<p>Loading feed...</p>}>
              <MachineDeploymentLister mds={FilterItems(query, mds)} />
            </Suspense>
          </GridCol>
        </Grid>
      </main>
    </div>
  )
}