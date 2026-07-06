'use client';

import React from 'react';
import MachineInfraDockerDetails from '@/app/ui/dashboard/components/machines/infra/docker-details'
import MachineInfraDockerTable from '@/app/ui/dashboard/components/machines/infra/docker-table'

import {MachineInfraDockerType} from "@/app/ui/dashboard/components/machines/types";
import BaseLister from "@/app/ui/dashboard/base/lister";

/**
 * Renders the Docker (CAPD) equivalent of `MachineInfraLister`: live machine infrastructure
 * details, loading/empty/error states, and search, for Docker-backed machines.
 */
export default function MachineInfraDockerLister() {
  return <BaseLister
      objectType="machine-infra-docker"
      items={[]}
      renderDetails={(item: MachineInfraDockerType) => <MachineInfraDockerDetails machine={item}/>}
      renderTable={(items : MachineInfraDockerType[], handleSelect) =>  (
        <MachineInfraDockerTable select={handleSelect} machines={items}/>
      )}
      title="DockerMachine / infrastructure.cluster.x-k8s.io/v1beta1"
    />
}
