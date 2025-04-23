'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';

import Search from "@/app/ui/dashboard/search";
import {FilterItems} from "@/app/dashboard/utils";
import { Grid, GridCol, Title } from '@mantine/core';

import { getMachines } from "@/app/lib/data";
import MachinesTable from "@/app/ui/dashboard/components/machines/table";

// MachineLister: List machines existent in the cluster.
export default function MachineLister () {
  const [machines,setMachines] = useState<[]>([])
  const [query,setQuery] = useState('')
  const filteredMachines = FilterItems(query, machines);

  useEffect(() => {
    const fetchData = async () => { setMachines(await getMachines()) }
    fetchData().catch((e) => { console.error('error', e) })
  }, [])

  return (
    <Grid justify="flex-end" align="flex-start">
      <GridCol h={60} span={8}>
        <Link href="/dashboard/clusters">
          <Title className="hidden md:block" order={2}>
            Machines / cluster.x-k8s.io
          </Title>
        </Link>
      </GridCol>
      <Search
        value={query}
        onChange={(e: { currentTarget: { value: React.SetStateAction<string>; }; }) => setQuery(e.currentTarget.value)}
        placeholder="Cluster name" />
      <MachinesTable machines={filteredMachines} />
    </Grid>
  )
}
