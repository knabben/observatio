'use client';

import React from 'react';

import MachineHealthCheckTable from '@/app/ui/dashboard/components/machinehealthchecks/table'
import MachineHealthCheckDetails from "@/app/ui/dashboard/components/machinehealthchecks/details";

import {MachineHealthCheckType} from "@/app/ui/dashboard/components/machinehealthchecks/types";
import BaseLister from "@/app/ui/dashboard/base/lister";
/**
 * Thin composition of `BaseLister` with the MachineHealthCheck-specific table/details renderers.
 * `BaseLister` owns the live WebSocket stream, loading/empty/error states, and selection.
 */
export default function MachineHealthCheckLister() {
  return <BaseLister
    objectType="machinehealthcheck"
    items={[]}
    renderDetails={(item: MachineHealthCheckType) => <MachineHealthCheckDetails mhc={item}/>}
    renderTable={(items : MachineHealthCheckType[], handleSelect) =>  (
      <MachineHealthCheckTable select={handleSelect} mhcs={items}/>
    )}
    title="Machine Health Checks / cluster.x-k8s.io"
  />
}
