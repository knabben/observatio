'use client';

import React from "react";
import {Badge} from '@mantine/core';
import {GridCol} from '@mantine/core';

import {MachineSetType} from '@/app/ui/dashboard/components/machinesets/types';
import {StatusIndicator} from '@/app/ui/dashboard/shared/status-indicator';
import {toStatusState} from '@/app/ui/dashboard/shared/status';
import {ObjectTable} from '@/app/ui/dashboard/shared/object-table';
import {ColumnDef} from '@/app/ui/dashboard/base/types';

/** MachineSet is ready when every desired replica is currently available; unknown when absent. */
function msReady(desired: number | undefined, available: number | undefined): boolean | undefined {
  if (desired == null || available == null) return undefined;
  return available === desired;
}

const columns: ColumnDef<MachineSetType>[] = [
  {header: 'Name', render: (ms) => ms.metadata?.name ?? '—'},
  {header: 'Namespace', render: (ms) => <Badge variant="light" color="gray">{ms.metadata?.namespace ?? '—'}</Badge>},
  {header: 'Cluster', render: (ms) => ms.cluster ?? '—'},
  {header: 'Machine Deployment', render: (ms) => ms.machineDeployment ?? '—'},
  {header: 'Available / Replicas', align: 'center', render: (ms) => `${ms.status?.availableReplicas ?? '—'} / ${ms.replicas ?? '—'}`},
  {header: 'Age', align: 'center', render: (ms) => ms.age ?? '—'},
  {
    header: 'Status',
    align: 'center',
    render: (ms) => <StatusIndicator state={toStatusState(msReady(ms.replicas, ms.status?.availableReplicas))} dotOnly/>,
  },
];

export default function MachineSetTable({
  mss, select
}: {
  mss: MachineSetType[]
  select: (ms: MachineSetType) => void
}) {
  return (
    <GridCol span={12}>
      <ObjectTable
        items={mss}
        columns={columns}
        getRowKey={(ms, i) => ms.metadata?.name ?? `row-${i}`}
        onSelect={select}
        emptyLabel="No machine sets found"
      />
    </GridCol>
  );
}
