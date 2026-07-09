'use client';

import React from 'react';

import MachineSetTable from '@/app/ui/dashboard/components/machinesets/table'
import MachineSetDetails from "@/app/ui/dashboard/components/machinesets/details";

import {MachineSetType} from "@/app/ui/dashboard/components/machinesets/types";
import BaseLister from "@/app/ui/dashboard/base/lister";
/**
 * Thin composition of `BaseLister` with the MachineSet-specific table/details renderers.
 * `BaseLister` owns the live WebSocket stream, loading/empty/error states, and selection.
 */
export default function MachineSetLister() {
  return <BaseLister
    objectType="machineset"
    items={[]}
    renderDetails={(item: MachineSetType) => <MachineSetDetails ms={item}/>}
    renderTable={(items : MachineSetType[], handleSelect) =>  (
      <MachineSetTable select={handleSelect} mss={items}/>
    )}
    title="Machine Sets / cluster.x-k8s.io"
  />
}
