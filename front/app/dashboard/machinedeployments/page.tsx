import Link from 'next/link';

import { Suspense } from 'react';
import { getMachinesDeployments} from "@/app/lib/data";
import MachineDeploymentLister from "@/app/ui/dashboard/components/MachineDeploymentLister";

import Search from '@/app/ui/dashboard/search'
import { Title, Grid, GridCol } from '@mantine/core';

export default async function MachineDeployments() {
  const mds = await getMachinesDeployments()

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
            <Suspense fallback={<div>Loading...</div>}>
              <Search placeholder="Machine deployment name"/>
            </Suspense>
          </GridCol>
          <GridCol span={12}>
            <MachineDeploymentLister mds={mds} />
          </GridCol>
        </Grid>
      </main>
    </div>
  )
}