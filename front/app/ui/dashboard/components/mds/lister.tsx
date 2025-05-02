'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';

import Search from '@/app/ui/dashboard/search'
import {FilterItems} from "@/app/dashboard/utils";
import { Grid, GridCol, Title, Badge } from '@mantine/core';

import {getMachinesDeployments} from "@/app/lib/data";
import MachineDeploymentTable from '@/app/ui/dashboard/components/mds/table'
import MachineDeploymentDetails from "@/app/ui/dashboard/components/mds/details";
import {sourceCodePro400} from "@/fonts";

type Status = {
  failed: number;
  total: number;
}
// MDLister: List the MDs existent in the cluster.
export default function MDLister() {
  const [status, setStatus] = useState<Status>({failed: 0, total: 0})
  const [mds, setMD] = useState<[]>([])
  const [selected, setSelected] = useState('')

  let filteredMD = undefined;
  if (selected) {
    filteredMD = FilterItems(selected, mds);
  }

  useEffect(() => {
    const fetchData = async () => {
      const response = await getMachinesDeployments()
      setMD(response.machineDeployments)
      setStatus({
        "failed": response.failed,
        "total": response.total,
      })
    }
    fetchData().catch((e) => { console.error('error', e) })
  }, [])

  return (
    <Grid justify="flex-end" align="flex-start">
      <GridCol h={60} span={6}>
        <Link href="/dashboard/machinedeployments">
          <Title className={sourceCodePro400.className} order={2}>
            Machine Deployments / cluster.x-k8s.io
          </Title>
        </Link>
      </GridCol>
      <GridCol className="text-right" h={60} span={2}>
        <Badge className="m-1" radius="sm" variant="dot" color="blue" size="lg">{status.total}</Badge>
        { status.failed > 0 ? <Badge radius="sm" variant="dot" color="red" size="lg">{status.failed}</Badge> : <div></div> }
      </GridCol>
      <Search
        options={mds}
        onChange={setSelected}/>
      {
        filteredMD
          ? <MachineDeploymentDetails md={filteredMD} />
          : <MachineDeploymentTable mds={mds} />
      }
    </Grid>
  )
}
