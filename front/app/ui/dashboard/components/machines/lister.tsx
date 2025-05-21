'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';
import {sourceCodePro400} from "@/fonts";

import {FilterItems} from "@/app/dashboard/utils";
import {Grid, GridCol, Title } from '@mantine/core';

import MachinesTable from "@/app/ui/dashboard/components/machines/table";
import MachineDetails from "@/app/ui/dashboard/components/machines/details";

import {MachineType} from "@/app/ui/dashboard/components/machines/types";
import {receiveAndPopulate, sendInitialRequest, WebSocket} from "@/app/lib/websocket";
import {CenteredLoader} from "@/app/ui/dashboard/utils/loader";
import {IconArrowBigLeft} from "@tabler/icons-react";

/**
 * MachineLister component handles listing and displaying details of machines.
 * It manages state for machines, selected machine, and loading status.
 * Uses WebSocket to fetch and populate machine data and renders a loader while data is being loaded.
 * Once the data is available, it renders a search interface and displays filtered results or a table of machines.
 */
export default function MachineLister() {
  const [machines, setMachines] = useState<MachineType[]>([])
  const [selected, setSelected] = useState('')
  const [loading, setLoading] = useState(true)

  const handleSelect = (machine: MachineType | null) => {
    if (machine === null) {
      setSelected('')
    }
    // @ts-expect-error machine
    setSelected(machine.metadata?.name)
  }

  const filteredMachines: MachineType | undefined = selected
    ? FilterItems(selected, machines)
    : undefined;

  const {sendJsonMessage, lastJsonMessage, readyState} = WebSocket()

  useEffect(() => {
    sendInitialRequest(readyState, "machine", sendJsonMessage)
  }, [readyState, sendJsonMessage])

  useEffect(() => {
    const newMachines: MachineType[] = receiveAndPopulate(lastJsonMessage, [...machines])
    setMachines(newMachines.sort((a: MachineType, b: MachineType) =>
      a.metadata?.name.localeCompare(b.metadata?.name)))
    setLoading(false)
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [lastJsonMessage, setMachines])

  if (loading) {
    return <CenteredLoader/>;
  }

  return (
    <Grid justify="flex-end" align="flex-start">
      <GridCol span={9}>
        <Link href="/dashboard/clusters">
          <Title className={sourceCodePro400.className} order={2}>
            Machines / cluster.x-k8s.io
          </Title>
        </Link>
      </GridCol>
      <GridCol span={3} className="flex justify-end items-center">
        { selected &&
          <div>
            <IconArrowBigLeft onClick={() => handleSelect(null)} size={32} className="cursor-pointer hover:opacity-70"/>
          </div>
        }
      </GridCol>
      {
        filteredMachines
          ? <MachineDetails machine={filteredMachines} />
          : <MachinesTable select={handleSelect} machines={machines}/>
      }
    </Grid>
  )
}
