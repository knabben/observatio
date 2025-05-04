'use client';

import Link from 'next/link';
import React, { useState, useEffect } from 'react';
import useWebSocket, { ReadyState } from 'react-use-websocket';
import { sourceCodePro400 } from "@/fonts";

import Search from "@/app/ui/dashboard/search";
import {FilterItems} from "@/app/dashboard/utils";
import { Grid, GridCol, Title, Badge, Loader, Alert } from '@mantine/core';

import ClusterInfraTable from '@/app/ui/dashboard/components/clusters/infra-table'
import ClusterInfraDetails from "@/app/ui/dashboard/components/clusters/infra-details";

import {WS_URL} from '@/app/ui/dashboard/utils/consts'
import {ClusterInfraType} from "@/app/ui/dashboard/components/clusters/types";

type WSResponse = {
  type: string;
  data: any;
}

// ClusterLister: Cluster list and details component.
export default function ClusterInfraLister() {
  const [vsphereClusters,setVsphereClusters] = useState<ClusterInfraType[]>([])
  const [selected, setSelected] = useState('')
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState("")

  const { sendJsonMessage, lastJsonMessage, readyState } = useWebSocket(
    WS_URL, {share: false, shouldReconnect: () => true},
  )

  let filteredClusterInfra= undefined;
  if (selected) {
    filteredClusterInfra = FilterItems(selected, vsphereClusters);
  }

  const updateOrAddItem = (newItem: ClusterInfraType) => {
    const index = vsphereClusters.findIndex(item => item.name === newItem.name);
    if (index !== -1) {
      const updatedItems = [...vsphereClusters];
      updatedItems[index] = newItem;
      return updatedItems;
    } else {
      return [...vsphereClusters, newItem];
    }
  };

  useEffect(() => {
    if (readyState === ReadyState.OPEN) {
      sendJsonMessage({types: ['cluster-infra']});
    }
  }, [readyState])

  useEffect(() => {
    let response: WSResponse = (lastJsonMessage as WSResponse)
    if (response?.type == "ADDED" || response?.type == "MODIFIED") {
      const data: ClusterInfraType = (response?.data as ClusterInfraType)
      setVsphereClusters(updateOrAddItem(data));
    } else {
      const data = response?.data
      setVsphereClusters(vsphereClusters.filter(
        item => item.name !== data.name
      ));
    }
    setLoading(false)
  }, [lastJsonMessage])

  if (loading) {
    return(<div className="text-center"><Loader color="teal" size="xl"/></div>)
  }
  if (error) {
    return (<Alert variant="light" color="red" title="Endpoint Error">{error}</Alert>)
  }

  return (
    <Grid justify="flex-end" align="flex-start">
      <GridCol h={60} span={8}>
        <Link href="/dashboard/clusters">
          <Title className={sourceCodePro400.className} order={2}>
            VSphereCluster / infrastructure.cluster.x-k8s.io/v1beta1
          </Title>
        </Link>
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