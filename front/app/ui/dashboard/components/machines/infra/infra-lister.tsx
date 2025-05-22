'use client';

import React from 'react';
import MachineInfraDetails from '@/app/ui/dashboard/components/machines/infra/infra-details'
import MachineInfraTable from '@/app/ui/dashboard/components/machines/infra/infra-table'

import {MachineInfraType} from "@/app/ui/dashboard/components/machines/types";
import BaseLister from "@/app/ui/dashboard/base/lister";

/**
 * A functional component that fetches, manages, and displays a list of vSphere machines along with their infrastructure details.
 * It uses WebSocket messaging to receive and update machine details in real time.
 * The component allows for selecting a specific machine to display its details or displaying the entire list in a table format.
 */
export default function MachineInfraLister() {
  return <BaseLister
      objectType="machine-infra"
      items={[]}
      renderDetails={(item: MachineInfraType) => <MachineInfraDetails machine={item}/>}
      renderTable={(items : MachineInfraType[], handleSelect) =>  (
        <MachineInfraTable select={handleSelect} machines={items}/>
      )}
      title="Machine Infra / infrastructure.cluster.x-k8s.io/v1beta1"
    />
}