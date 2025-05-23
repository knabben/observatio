'use client';

import React from 'react';
import MachinesTable from "@/app/ui/dashboard/components/machines/table";
import MachineDetails from "@/app/ui/dashboard/components/machines/details";

import {MachineType} from "@/app/ui/dashboard/components/machines/types";
import BaseLister from "@/app/ui/dashboard/base/lister";

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