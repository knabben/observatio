'use client';

import React from 'react';
import MachinesTable from "@/app/ui/dashboard/components/machines/table";
import MachineDetails from "@/app/ui/dashboard/components/machines/details";

import {MachineType} from "@/app/ui/dashboard/components/machines/types";
import BaseLister from "@/app/ui/dashboard/base/lister";

/**
 * Thin composition of `BaseLister` with the machine-specific table/details renderers.
 * `BaseLister` owns the live WebSocket stream, loading/empty/error states, and selection.
 */
export default function MachineLister() {
  return <BaseLister
    objectType="machine"
    items={[]}
    renderDetails={(item: MachineType) => <MachineDetails machine={item}/>}
    renderTable={(items : MachineType[], handleSelect) =>  (
      <MachinesTable select={handleSelect} machines={items}/>
    )}
    title="Machines / cluster.x-k8s.io"
  />
}