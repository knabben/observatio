'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';
import { sourceCodePro400 } from "@/fonts";

import Search from "@/app/ui/dashboard/search";
import {FilterItems} from "@/app/dashboard/utils";
import { Grid, GridCol, Title, Badge, Loader, Alert } from '@mantine/core';

import ClusterInfraTable from '@/app/ui/dashboard/components/clusters/infra-table'
import ClusterInfraDetails from "@/app/ui/dashboard/components/clusters/infra-details";

import { getClusterInfraList } from "@/app/lib/data";

type Status = {
  failed: number;
  total: number;
}
// ClusterLister: Cluster list and details component.
export default function ClusterInfraLister() {
  const [vsphereClusters,setVsphereClusters] = useState<[]>([])
  const [selected, setSelected] = useState('')
  const [status, setStatus] = useState<Status>({failed: 0, total: 0})
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState("")

  let filteredClusterInfra= undefined;
  if (selected) {
    filteredClusterInfra = FilterItems(selected, vsphereClusters);
  }

  useEffect(() => {
    const fetchData = async () => {
      const response = await getClusterInfraList()
      setVsphereClusters(response.clusters)
      setLoading(false)
      setStatus({
        "failed": response.failed,
        "total": response.total,
      })
    }
    fetchData().catch((e) => {
      setLoading(false)
      setError(e.toString())
    })
  }, [])

  if (loading) {
    return(<div className="text-center"><Loader color="teal" size="xl"/></div>)
  }
  if (error) {
    return (<Alert variant="light" color="red" title="Endpoint Error"> {error} </Alert>)
  }
  
  return (
    <Grid justify="flex-end" align="flex-start">
      <GridCol h={60} span={7}>
        <Link href="/dashboard/clusters">
          <Title className={sourceCodePro400.className} order={2}>
            VSphereCluster / infrastructure.cluster.x-k8s.io/v1beta1
          </Title>
        </Link>
      </GridCol>
      <GridCol className="text-right" h={60} span={1}>
        <Badge className="m-1" radius="sm" variant="dot" color="blue" size="lg">{status.total}</Badge>
        { status.failed > 0 ? <Badge radius="sm" variant="dot" color="red" size="lg">{status.failed}</Badge> : <div></div> }
      </GridCol>
      <Search
        options={vsphereClusters}
        onChange={setSelected}/>
      {
        filteredClusterInfra
          ? <ClusterInfraDetails cluster={filteredClusterInfra} />
          : <ClusterInfraTable clusters={vsphereClusters}/>
      }
    </Grid>
  );
}