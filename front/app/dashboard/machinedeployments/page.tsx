import Link from 'next/link';

import { getMachinesDeployments} from "@/app/lib/data";
import MachineDeploymentLister from "@/app/ui/dashboard/components/MachineDeploymentLister";

import Search from '@/app/ui/dashboard/search'
import { Title, Grid, GridCol, Space } from '@mantine/core';
import ClusterLister from "@/app/ui/dashboard/components/ClusterLister";

export default async function MachineDeployments(props: {
  searchParams?: Promise<{
    query?: string;
  }>;
}) {
  const searchParams = await props.searchParams;
  const query = searchParams?.query || '';
  let mds = await getMachinesDeployments(query)

  if (query != "") {
    mds = mds.filter((i: { name: string; }) =>
      i.name.toLowerCase().includes(query.toLowerCase()));

  }
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
            <MachineDeploymentLister mds={mds} />
          </GridCol>
        </Grid>
      </main>
    </div>
  )
}