'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';

import Search from "@/app/ui/dashboard/search";
import {FilterItems} from "@/app/dashboard/utils";
import { Grid, GridCol, Title, Badge } from '@mantine/core';

import { getMachines } from "@/app/lib/data";
import MachinesTable from "@/app/ui/dashboard/components/machines/table";
import ClusterDetails from "@/app/ui/dashboard/components/clusters/details";
import ClusterTable from "@/app/ui/dashboard/components/clusters/table";

type Status = {
  failed: number;
  total: number;
}
// MachineLister: List machines existent in the cluster.
export default function MachineLister () {
  const [status, setStatus] = useState<Status>({failed: 0, total: 0})
  const [machines,setMachines] = useState<[]>([])
  const [selected, setSelected] = useState('')
  const [loading, setLoading] = useState(true)

  let filteredMachines = undefined;
  if (selected) {
    filteredMachines = FilterItems(selected, machines);
  }

  useEffect(() => {
    const fetchData = async () => {
      const response = await getMachines();
      setMachines(response.machines)
      setLoading(false)
      setStatus({
        "failed": response.failed,
        "total": response.total,
      })
    }
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
      <GridCol className="text-right" h={60} span={2}>
        <Badge className="m-1" radius="sm" variant="dot" color="blue" size="lg">{status.total}</Badge>
        { status.failed > 0 ? <Badge radius="sm" variant="dot" color="red" size="lg">{status.failed}</Badge> : <div></div> }
      </GridCol>
      <Search
        options={machines}
        onChange={setSelected}/>
      {
        filteredMachines
          ? <ClusterDetails cluster={filteredMachines} />
          : <MachinesTable loading={loading} machines={machines}/>
      }
    </Grid>
  )
}
