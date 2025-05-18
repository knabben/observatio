'use client';

import Link from 'next/link';
import React, { useState, useEffect } from 'react';
import { sourceCodePro400 } from "@/fonts";

import {FilterItems} from "@/app/dashboard/utils";
import {Grid, GridCol, Title} from '@mantine/core';

import MachineInfraDetails from '@/app/ui/dashboard/components/machines/infra-details'
import MachineInfraTable from '@/app/ui/dashboard/components/machines/infra-table'

import {MachineInfraType} from "@/app/ui/dashboard/components/machines/types";
import {receiveAndPopulate, sendInitialRequest, WebSocket} from "@/app/lib/websocket";
import {CenteredLoader} from "@/app/ui/dashboard/utils/loader";
import {IconArrowBigLeft} from "@tabler/icons-react";

export default function MachineInfraLister() {
  const [vsphereMachines,setVsphereMachine] = useState<MachineInfraType[]>([])
  const [selected, setSelected] = useState<string>('')
  const [loading, setLoading] = useState(true)

  const handleSelect = (machine: MachineInfraType | null) => {
    if (machine === null) {
      setSelected('')
    }
    // @ts-expect-error machine
    setSelected(machine?.name)
  }

  const filteredMachine: MachineInfraType | undefined = selected
    ? FilterItems(selected, vsphereMachines)
    : undefined;

  const {sendJsonMessage, lastJsonMessage, readyState} = WebSocket()

  useEffect(() => {
    sendInitialRequest(readyState, "machine-infra", sendJsonMessage)
  }, [readyState, sendJsonMessage])

  useEffect(() => {
    const newVsphereMachine: MachineInfraType[] = receiveAndPopulate(lastJsonMessage, [...vsphereMachines])
    setVsphereMachine(newVsphereMachine.sort((a: MachineInfraType, b: MachineInfraType) =>
      a.metadata?.name.localeCompare(b.metadata?.name))
    )
    setLoading(false)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [lastJsonMessage, setVsphereMachine])

  if (loading) {
    return <CenteredLoader/>;
  }

  return (
    <Grid justify="flex-end" align="flex-start">
      <GridCol span={9}>
        <Link href="/dashboard/machines">
          <Title className={sourceCodePro400.className} order={2}>
            VSphereMachine / infrastructure.cluster.x-k8s.io/v1beta1
          </Title>
        </Link>
      </GridCol>
      <GridCol span={3} className="flex justify-end items-center">
        { selected &&
          <div>
            <span className="text-sm text-gray-500">{selected}</span>
            <IconArrowBigLeft onClick={() => handleSelect(null)} size={32} className="cursor-pointer hover:opacity-70"/>
          </div>
        }
      </GridCol>
      {
        filteredMachine
          ? <MachineInfraDetails machine={filteredMachine} />
          : <MachineInfraTable select={handleSelect} machines={vsphereMachines}/>
      }
    </Grid>
  );
}