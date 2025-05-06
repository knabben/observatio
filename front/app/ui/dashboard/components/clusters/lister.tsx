'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';
import { sourceCodePro400 } from "@/fonts";

import Search from "@/app/ui/dashboard/search";
import {FilterItems} from "@/app/dashboard/utils";
import { Grid, GridCol, Title, Badge, Alert, Loader } from '@mantine/core';

import {ClusterType} from "@/app/ui/dashboard/components/clusters/types";
import ClusterTable from '@/app/ui/dashboard/components/clusters/table'
import ClusterDetails from "@/app/ui/dashboard/components/clusters/details";
import {receiveAndPopulate, sendInitialRequest, WebSocket} from "@/app/lib/websocket";

/**
 * A functional component that fetches, filters, and displays a list of clusters.
 * The component integrates WebSocket for real-time communication and enables
 * cluster search functionality. Displays a loader while data is being fetched.
 */
export default function ClusterLister() {
  const [clusters,setClusters] = useState<ClusterType[]>([])
  const [selected,setSelected] = useState('')
  const [loading, setLoading] = useState(true)

  const {sendJsonMessage, lastJsonMessage, readyState} = WebSocket()

  useEffect(() => {
    sendInitialRequest(readyState, "cluster", sendJsonMessage)
  }, [readyState, sendJsonMessage])

  useEffect(() => {
    setClusters(receiveAndPopulate(lastJsonMessage, [...clusters]))
    setLoading(false)
  }, [lastJsonMessage])

  const filteredCluster: ClusterType | undefined = selected
    ? FilterItems(selected, clusters)
    : undefined;

  if (loading) {
    return (<div className="text-center"><Loader color="teal" size="xl"/></div>)
  }

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
      {
        filteredCluster
          ? <ClusterDetails cluster={filteredCluster} />
          : <ClusterTable clusters={clusters}/>
      }
    </Grid>
  );
}