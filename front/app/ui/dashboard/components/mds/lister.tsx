'use client';

import React from 'react';

import MachineDeploymentTable from '@/app/ui/dashboard/components/mds/table'
import MachineDeploymentDetails from "@/app/ui/dashboard/components/mds/details";

import {MachineDeploymentType} from "@/app/ui/dashboard/components/mds/types";
import BaseLister from "@/app/ui/dashboard/base/lister";
/**
 * MachineDeploymentLister is a functional component responsible for listing machine deployments
 * and optionally showing details for a selected deployment. It manages state for the deployments,
 * filters based on user selection, and handles WebSocket communication for data fetching.
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
