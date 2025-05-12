'use client';

import Link from 'next/link';
import React, { useState, useEffect } from 'react';
import { sourceCodePro400 } from "@/fonts";

import Search from "@/app/ui/dashboard/search";
import {FilterItems} from "@/app/dashboard/utils";
import { Grid, GridCol, Title } from '@mantine/core';

import MachineInfraDetails from '@/app/ui/dashboard/components/machines/infra-details'
import MachineInfraTable from '@/app/ui/dashboard/components/machines/infra-table'

import {MachineInfraType} from "@/app/ui/dashboard/components/machines/types";
import {receiveAndPopulate, sendInitialRequest, WebSocket} from "@/app/lib/websocket";
import {CenteredLoader} from "@/app/ui/dashboard/utils/loader";

export default function MachineInfraLister() {
  const [vsphereMachine ,setVsphereMachine] = useState<MachineInfraType[]>([])
  const [selected, setSelected] = useState('')
  const [loading, setLoading] = useState(true)

  const filteredCluster: MachineInfraType | undefined = selected
    ? FilterItems(selected, vsphereMachine)
    : undefined;

  const {sendJsonMessage, lastJsonMessage, readyState} = WebSocket()

  useEffect(() => {
    sendInitialRequest(readyState, "machine-infra", sendJsonMessage)
  }, [readyState, sendJsonMessage])

  useEffect(() => {
    const newVsphereMachine: MachineInfraType[] = receiveAndPopulate(lastJsonMessage, [...vsphereMachine])
    setVsphereMachine(newVsphereMachine.sort((a: MachineInfraType, b: MachineInfraType) => a.name.localeCompare(b.name)))
    setLoading(false)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [lastJsonMessage, setVsphereMachine])

  if (loading) {
    return <CenteredLoader/>;
  }

  return (
    <Grid justify="flex-end" align="flex-start">
      <GridCol h={60} span={8}>
        <Link href="/dashboard/machines">
          <Title className={sourceCodePro400.className} order={2}>
            VSphereMachine / infrastructure.cluster.x-k8s.io/v1beta1
          </Title>
        </Link>
      </GridCol>
      <Search
        options={vsphereMachine}
        onChange={setSelected}/>
      {
        filteredCluster
          ? <MachineInfraDetails machine={filteredCluster} />
          : <MachineInfraTable machines={vsphereMachine}/>
      }
    </Grid>
  );
}