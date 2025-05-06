/* eslint-disable @typescript-eslint/no-explicit-any */
'use client';

import Link from 'next/link';
import React, { useState, useEffect } from 'react';
import { sourceCodePro400 } from "@/fonts";

import Search from "@/app/ui/dashboard/search";
import {FilterItems} from "@/app/dashboard/utils";
import { Grid, GridCol, Title, Loader } from '@mantine/core';

import ClusterInfraTable from '@/app/ui/dashboard/components/clusters/infra-table'
import ClusterInfraDetails from "@/app/ui/dashboard/components/clusters/infra-details";

import {ClusterInfraType, ClusterType} from "@/app/ui/dashboard/components/clusters/types";
import {receiveAndPopulate, WebSocket} from "@/app/lib/websocket";
import {sendInitialRequest} from "@/app/lib/websocket";

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

  const {sendJsonMessage, lastJsonMessage, readyState} = WebSocket()

  useEffect(() => {
    sendInitialRequest(readyState, "cluster-infra", sendJsonMessage)
  }, [readyState, sendJsonMessage])

  useEffect(() => {
    setVsphereClusters(receiveAndPopulate(lastJsonMessage, [...vsphereClusters]))
    setLoading(false)
  }, [lastJsonMessage])

  const filteredCluster: ClusterInfraType | undefined = selected
    ? FilterItems(selected, vsphereClusters)
    : undefined;

  if (loading) {
    return(<div className="text-center"><Loader color="teal" size="xl"/></div>)
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