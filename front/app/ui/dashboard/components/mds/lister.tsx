'use client';

import React from 'react';

import MachineDeploymentTable from '@/app/ui/dashboard/components/mds/table'
import MachineDeploymentDetails from "@/app/ui/dashboard/components/mds/details";

import {MachineDeploymentType} from "@/app/ui/dashboard/components/mds/types";
import BaseLister from "@/app/ui/dashboard/base/lister";
/**
 * Thin composition of `BaseLister` with the machine-deployment-specific table/details renderers.
 * `BaseLister` owns the live WebSocket stream, loading/empty/error states, and selection.
 */
export default function MachineDeploymentLister() {
  return <BaseLister
    objectType="machine-deployment"
    items={[]}
    renderDetails={(item: MachineDeploymentType) => <MachineDeploymentDetails md={item}/>}
    renderTable={(items : MachineDeploymentType[], handleSelect) =>  (
      <MachineDeploymentTable select={handleSelect} mds={items}/>
    )}
    title="Machine Deployments / cluster.x-k8s.io"
  />
}
