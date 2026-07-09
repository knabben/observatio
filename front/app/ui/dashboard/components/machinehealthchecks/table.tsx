'use client';

import React from "react";
import {Badge} from '@mantine/core';
import {GridCol} from '@mantine/core';

import {MachineHealthCheckType} from '@/app/ui/dashboard/components/machinehealthchecks/types';
import {StatusIndicator} from '@/app/ui/dashboard/shared/status-indicator';
import {toStatusState} from '@/app/ui/dashboard/shared/status';
import {ObjectTable} from '@/app/ui/dashboard/shared/object-table';
import {ColumnDef} from '@/app/ui/dashboard/base/types';

/** MachineHealthCheck is ready when every expected Machine is currently healthy; unknown when absent. */
function mhcReady(expected: number | undefined, healthy: number | undefined): boolean | undefined {
  if (expected == null || healthy == null) return undefined;
  return healthy === expected;
}

const columns: ColumnDef<MachineHealthCheckType>[] = [
  {header: 'Name', render: (mhc) => mhc.metadata?.name ?? '—'},
  {header: 'Namespace', render: (mhc) => <Badge variant="light" color="gray">{mhc.metadata?.namespace ?? '—'}</Badge>},
  {header: 'Cluster', render: (mhc) => mhc.cluster ?? '—'},
  {header: 'Max Unhealthy', align: 'center', render: (mhc) => mhc.maxUnhealthy ?? '—'},
  {header: 'Healthy / Expected', align: 'center', render: (mhc) => `${mhc.status?.currentHealthy ?? '—'} / ${mhc.status?.expectedMachines ?? '—'}`},
  {header: 'Age', align: 'center', render: (mhc) => mhc.age ?? '—'},
  {
    header: 'Status',
    align: 'center',
    render: (mhc) => <StatusIndicator state={toStatusState(mhcReady(mhc.status?.expectedMachines, mhc.status?.currentHealthy))} dotOnly/>,
  },
];

export default function MachineHealthCheckTable({
  mhcs, select
}: {
  mhcs: MachineHealthCheckType[]
  select: (mhc: MachineHealthCheckType) => void
}) {
  return (
    <GridCol span={12}>
      <ObjectTable
        items={mhcs}
        columns={columns}
        getRowKey={(mhc, i) => mhc.metadata?.name ?? `row-${i}`}
        onSelect={select}
        emptyLabel="No machine health checks found"
      />
    </GridCol>
  );
}
