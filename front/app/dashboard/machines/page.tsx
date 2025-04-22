import Link from 'next/link';

import { Suspense } from 'react';
import { getMachines } from "@/app/lib/data";
import MachineLister from "@/app/ui/dashboard/components/MachineLister";

import Search from "@/app/ui/dashboard/search";
import { Title, Grid, GridCol } from '@mantine/core';

export default async function Machines() {
  const machines = await getMachines()

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
            <Suspense fallback={<div>Loading...</div>}>
              <Search placeholder="Machine name"/>
            </Suspense>
          </GridCol>
          <GridCol span={12}>
            <MachineLister machines={machines} />
          </GridCol>
        </Grid>
      </main>
    </div>
  )
}