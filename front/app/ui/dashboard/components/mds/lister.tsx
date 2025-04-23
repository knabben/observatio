'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';

import Search from '@/app/ui/dashboard/search'
import {FilterItems} from "@/app/dashboard/utils";
import { Title, Grid, GridCol } from '@mantine/core';

import {getMachinesDeployments} from "@/app/lib/data";
import MDTable from '@/app/ui/dashboard/components/mds/table'

// MDLister: List the MDs existent in the cluster.
export default function MDLister() {
  const [mds, setMD] = useState<[]>([])
  const [query, setQuery] = useState('')
  const filteredMDs = FilterItems(query, mds);

  useEffect(() => {
    const fetchData = async () => { setMD(await getMachinesDeployments()) }
    fetchData().catch((e) => { console.error('error', e) })
  }, [])

  return (
    <Grid justify="flex-end" align="flex-start">
      <GridCol h={60} span={8}>
        <Link href="/dashboard/machinedeployments">
          <Title className="hidden md:block" order={2}>
            Machine Deployments / cluster.x-k8s.io
          </Title>
        </Link>
      </GridCol>
      <Search
      value={query}
      onChange={(e: { currentTarget: { value: React.SetStateAction<string>; }; }) => setQuery(e.currentTarget.value)}
      placeholder="Cluster name" />
      <MDTable mds={filteredMDs} />
    </Grid>
  )
}
