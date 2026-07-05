'use client';

import React from "react";
import {Badge} from '@mantine/core';
import { GridCol } from '@mantine/core';

import {MachineDeploymentType} from '@/app/ui/dashboard/components/mds/types';
import {StatusIndicator} from '@/app/ui/dashboard/shared/status-indicator';
import {toStatusState} from '@/app/ui/dashboard/shared/status';
import {ObjectTable} from '@/app/ui/dashboard/shared/object-table';
import {ColumnDef} from '@/app/ui/dashboard/base/types';

/** MachineDeployment is ready when it has zero unavailable replicas; unknown when absent. */
function mdReady(unavailable: number | undefined): boolean | undefined {
  if (unavailable == null) return undefined;
  return unavailable === 0;
}

const columns: ColumnDef<MachineDeploymentType>[] = [
  {header: 'Name', render: (md) => md.metadata?.name ?? '—'},
  {header: 'Namespace', render: (md) => <Badge variant="light" color="gray">{md.metadata?.namespace ?? '—'}</Badge>},
  {header: 'Replicas', render: (md) => md.replicas ?? '—'},
  {header: 'Cluster', render: (md) => md.cluster ?? '—'},
  {header: 'Phase', align: 'center', render: (md) => md.status?.phase ?? '—'},
  {header: 'Age', align: 'center', render: (md) => md.age ?? '—'},
  {
    header: 'Status',
    align: 'center',
    render: (md) => <StatusIndicator state={toStatusState(mdReady(md.status?.unavailableReplicas))} dotOnly/>,
  },
];

export default function MDTable({
  mds, select
}: {
  mds: MachineDeploymentType[]
  select: (machine: MachineDeploymentType) => void
}) {
  return (
    <GridCol span={12}>
      <ObjectTable
        items={mds}
        columns={columns}
        getRowKey={(md, i) => md.metadata?.name ?? `row-${i}`}
        onSelect={select}
        emptyLabel="No machine deployments found"
      />
    </GridCol>
  );
}
