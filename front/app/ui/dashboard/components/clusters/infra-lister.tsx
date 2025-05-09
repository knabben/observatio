'use client';

import Link from 'next/link';
import React, { useState, useEffect } from 'react';
import { sourceCodePro400 } from "@/fonts";

import Search from "@/app/ui/dashboard/search";
import {FilterItems} from "@/app/dashboard/utils";
import { Grid, GridCol, Title } from '@mantine/core';

import ClusterInfraTable from '@/app/ui/dashboard/components/clusters/infra-table'
import ClusterInfraDetails from "@/app/ui/dashboard/components/clusters/infra-details";

import {ClusterInfraType} from "@/app/ui/dashboard/components/clusters/types";
import {receiveAndPopulate, sendInitialRequest, WebSocket} from "@/app/lib/websocket";
import {CenteredLoader} from "@/app/ui/dashboard/utils/loader";

/**
 * The `ClusterInfraLister` function is a React functional component responsible for rendering
 * a user interface to manage and display vSphere cluster infrastructure details.
 * It handles WebSocket communication, loading states, search functionality, and conditionally
 * displays a detailed view or a table of vSphere clusters based on the data provided.
 */
export default function ClusterInfraLister() {
  const [vsphereClusters,setVsphereClusters] = useState<ClusterInfraType[]>([])
  const [selected, setSelected] = useState('')
  const [loading, setLoading] = useState(true)

  const filteredCluster: ClusterInfraType | undefined = selected
    ? FilterItems(selected, vsphereClusters)
    : undefined;

  const {sendJsonMessage, lastJsonMessage, readyState} = WebSocket()

  useEffect(() => {
    sendInitialRequest(readyState, "cluster-infra", sendJsonMessage)
  }, [readyState, sendJsonMessage])

  useEffect(() => {
    const newVsphereClusters: ClusterInfraType[] = receiveAndPopulate(lastJsonMessage, [...vsphereClusters])
    setVsphereClusters(newVsphereClusters.sort((a: ClusterInfraType, b: ClusterInfraType) => a.name.localeCompare(b.name)))
    setLoading(false)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [lastJsonMessage, setVsphereClusters])

  if (loading) {
    return <CenteredLoader/>;
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
        filteredCluster
          ? <ClusterInfraDetails cluster={filteredCluster} />
          : <ClusterInfraTable clusters={vsphereClusters}/>
      }
    </Grid>
  );
}