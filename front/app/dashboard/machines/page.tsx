import Link from 'next/link';

import { getMachines } from "@/app/lib/data";
import MachineLister from "@/app/ui/dashboard/components/MachineLister";
import { FilterItems } from "@/app/dashboard/utils";

import Search from "@/app/ui/dashboard/search";
import { Title, Grid, GridCol } from '@mantine/core';

export default async function Machines(props: {
  searchParams?: Promise<{
    query?: string;
  }>;
}) {
  const searchParams = await props.searchParams;
  const query = searchParams?.query || '';
  const machines = await getMachines(query)

  return (
    <div>
      <main>
        <Grid justify="flex-end" align="flex-start">
          <GridCol h={60} span={8}>
            <Link href="/dashboard/machines">
              <Title className="hidden md:block" order={2}>
                Machines / cluster.x-k8s.io
              </Title>
            </Link>
          </GridCol>
          <GridCol span={4}>
            <Search placeholder="Machine name"/>
          </GridCol>
          <GridCol span={12}>
            <MachineLister machines={FilterItems(query, machines)} />
          </GridCol>
        </Grid>
      </main>
    </div>
  )
}