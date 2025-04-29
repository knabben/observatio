'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';
import { sourceCodePro400 } from "@/fonts";

import Search from "@/app/ui/dashboard/search";
import {FilterItems} from "@/app/dashboard/utils";
import { Grid, GridCol, Title, Badge } from '@mantine/core';

import { getClusterList } from "@/app/lib/data";
import ClusterTable from '@/app/ui/dashboard/components/clusters/table'
import ClusterDetails from "@/app/ui/dashboard/components/clusters/details";

type Status = {
  failed: number;
  total: number;
}
// ClusterLister: Cluster list and details component.
export default function ClusterLister() {
  const [status, setStatus] = useState<Status>({failed: 0, total: 0})
  const [clusters,setClusters] = useState<[]>([])
  const [selected, setSelected] = useState('')

  let filteredCluster = undefined;
  if (selected) {
    filteredCluster = FilterItems(selected, clusters);
  }

  useEffect(() => {
    const fetchData = async () => {
      const response = await getClusterList()
      setClusters(response.clusters)
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
        <Link href="/dashboard/clusters">
          <Title className={sourceCodePro400.className} order={2}>
            Clusters / cluster.x-k8s.io
          </Title>
        </Link>
      </GridCol>
      <GridCol className="text-right" h={60} span={2}>
        <Badge className="m-1" radius="sm" variant="dot" color="blue" size="lg">{status.total}</Badge>
        { status.failed > 0 ? <Badge radius="sm" variant="dot" color="red" size="lg">{status.failed}</Badge> : <div></div> }
      </GridCol>
      <Search
        options={clusters}
        onChange={setSelected}/>
      {
        filteredCluster
          ? <ClusterDetails cluster={filteredCluster} />
          : <ClusterTable clusters={clusters}/>
      }
    </Grid>
  );
}