'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';
import { sourceCodePro400 } from "@/fonts";

import Search from "@/app/ui/dashboard/search";
import {FilterItems} from "@/app/dashboard/utils";
import { Grid, GridCol, Title } from '@mantine/core';

import { getClusterList } from "@/app/lib/data";
import ClusterTable from '@/app/ui/dashboard/components/clusters/table'
import ClusterDetails from "@/app/ui/dashboard/components/clusters/details";

// ClusterLister: Cluster list and details component.
export default function ClusterLister() {
  const [clusters,setClusters] = useState<[]>([])
  const [selected, setSelected] = useState('')
  const filteredClusters = FilterItems(selected, clusters);

  useEffect(() => {
    const fetchData = async () => {
      const response = await getClusterList()
      setClusters(response.clusters)
    }
    fetchData().catch((e) => { console.error('error', e) })
  }, [])

  return (
    <Grid justify="flex-end" align="flex-start">
      <GridCol h={60} span={8}>
        <Link href="/dashboard/clusters">
          <Title className={sourceCodePro400.className} order={2}>
            Clusters / cluster.x-k8s.io
          </Title>
        </Link>
      </GridCol>
      <Search
        options={clusters}
        onChange={setSelected}/>
      <ClusterTable clusters={filteredClusters}/>
      <ClusterDetails selected={selected} />
    </Grid>
  );
}