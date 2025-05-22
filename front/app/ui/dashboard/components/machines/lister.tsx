'use client';

import React, { useState, useEffect } from 'react';
import Link from 'next/link';
import {sourceCodePro400} from "@/fonts";

import {FilterItems} from "@/app/dashboard/utils";
import {Grid, GridCol, Title } from '@mantine/core';

import MachinesTable from "@/app/ui/dashboard/components/machines/table";
import MachineDetails from "@/app/ui/dashboard/components/machines/details";

import {MachineInfraType, MachineType} from "@/app/ui/dashboard/components/machines/types";
import {receiveAndPopulate, sendInitialRequest, WebSocket} from "@/app/lib/websocket";
import {CenteredLoader} from "@/app/ui/dashboard/utils/loader";
import {IconArrowBigLeft} from "@tabler/icons-react";
import BaseLister from "@/app/ui/dashboard/base/lister";
import MachineInfraDetails from "@/app/ui/dashboard/components/machines/infra/infra-details";
import MachineInfraTable from "@/app/ui/dashboard/components/machines/infra/infra-table";

/**
 * MachineLister component handles listing and displaying details of machines.
 * It manages state for machines, selected machine, and loading status.
 * Uses WebSocket to fetch and populate machine data and renders a loader while data is being loaded.
 * Once the data is available, it renders a search interface and displays filtered results or a table of machines.
 */
export default function MachineLister() {
  return <BaseLister
    objectType="machineinfra"
    items={[]}
    renderDetails={(item: MachineType) => <MachineDetails machine={item}/>}
    renderTable={(items : MachineType[], handleSelect) =>  (
      <MachinesTable select={handleSelect} machines={items}/>
    )}
    title="Machines / cluster.x-k8s.io"
  />
}