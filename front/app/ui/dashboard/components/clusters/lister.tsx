'use client';

import React, { useState, useEffect } from 'react';
import { getClusterList } from "@/app/lib/data";

import Link from 'next/link';
import { Grid, GridCol, Title } from '@mantine/core';

import Search from "@/app/ui/dashboard/search";
import ClusterTable from '@/app/ui/dashboard/components/clusters/table'

import {FilterItems} from "@/app/dashboard/utils";

// ClusterLister: Cluster list and details component.
export default function ClusterLister() {
  const [clusters,setClusters] = useState<any[]>([])
  const [query,setQuery] = useState('')
  const filteredClusters = FilterItems(query, clusters);

  useEffect(() => {
    const fetchData = async () => {
      setClusters(await getClusterList())
    }
    fetchData().catch((e) => {
      console.error('error', e)
    })
  }, [])

  return (
    <Grid justify="flex-end" align="flex-start">
      <GridCol h={60} span={8}>
        <Link href="/dashboard/clusters">
          <Title className="hidden md:block" order={2}>
            Clusters / cluster.x-k8s.io
          </Title>
        </Link>
      </GridCol>
      <Search
        value={query}
        onChange={(e: { currentTarget: { value: React.SetStateAction<string>; }; }) => setQuery(e.currentTarget.value)}
        placeholder="Cluster name" />
      <ClusterTable clusters={filteredClusters} />
    </Grid>
  );
}