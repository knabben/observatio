'use client';

import Link from 'next/link';
import React, { useState, useEffect } from 'react';
import {sourceCodePro400} from "@/fonts";

import Search from '@/app/ui/dashboard/search'
import {FilterItems} from "@/app/dashboard/utils";
import { Grid, GridCol, Title } from '@mantine/core';

import MachineDeploymentTable from '@/app/ui/dashboard/components/mds/table'
import MachineDeploymentDetails from "@/app/ui/dashboard/components/mds/details";

import {MachineDeploymentType} from "@/app/ui/dashboard/components/mds/types";
import {receiveAndPopulate, sendInitialRequest, WebSocket} from "@/app/lib/websocket";
import {CenteredLoader} from "@/app/ui/dashboard/utils/loader";

/**
 * MachineDeploymentLister is a functional component responsible for listing machine deployments
 * and optionally showing details for a selected deployment. It manages state for the deployments,
 * filters based on user selection, and handles WebSocket communication for data fetching.
 */
export default function MachineDeploymentLister() {
  const [machineDeployments, setMachineDeployment] = useState<MachineDeploymentType[]>([])
  const [selected, setSelected] = useState('')
  const [loading, setLoading] = useState(true)

  const filteredMD: MachineDeploymentType | undefined = selected
    ? FilterItems(selected, machineDeployments)
    : undefined;

  const {sendJsonMessage, lastJsonMessage, readyState} = WebSocket()

  useEffect(() => {
    sendInitialRequest(readyState, "machine-deployment", sendJsonMessage)
  }, [readyState, sendJsonMessage])

  useEffect(() => {
    setMachineDeployment(receiveAndPopulate(lastJsonMessage, [...machineDeployments]))
    setLoading(false)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [lastJsonMessage, setMachineDeployment])

  if (loading) {
    return <CenteredLoader/>;
  }

  return (
    <Grid justify="flex-end" align="flex-start">
      <GridCol h={60} span={8}>
        <Link href="/dashboard/machinedeployments">
          <Title className={sourceCodePro400.className} order={2}>
            Machine Deployments / cluster.x-k8s.io
          </Title>
        </Link>
      </GridCol>
      <Search
        options={machineDeployments}
        onChange={setSelected}/>
      {
        filteredMD
          ? <MachineDeploymentDetails md={filteredMD} />
          : <MachineDeploymentTable mds={machineDeployments} />
      }
    </Grid>
  )
}
