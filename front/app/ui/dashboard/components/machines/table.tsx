'use client';

import React from "react";
import {Badge} from '@mantine/core';
import { GridCol } from '@mantine/core';
import {MachineType} from "@/app/ui/dashboard/components/machines/types";
import {StatusIndicator} from "@/app/ui/dashboard/shared/status-indicator";
import {allReady} from "@/app/ui/dashboard/shared/status";
import {ObjectTable} from '@/app/ui/dashboard/shared/object-table';
import {ColumnDef} from '@/app/ui/dashboard/base/types';

const columns: ColumnDef<MachineType>[] = [
  {header: 'Name', render: (m) => m.metadata?.name ?? '—'},
  {header: 'Namespace', render: (m) => <Badge variant="light" color="gray">{m.metadata?.namespace ?? '—'}</Badge>},
  {header: 'Version', render: (m) => m.version ?? '—'},
  {header: 'Cluster', render: (m) => m.cluster ?? '—'},
  {header: 'Age', render: (m) => m.age ?? '—', align: 'center'},
  {
    header: 'Status',
    align: 'center',
    render: (m) => (
      <StatusIndicator state={allReady(m.status?.bootstrapReady, m.status?.infrastructureReady)} dotOnly/>
    ),
  },
];

export default function MachinesTable({
  machines, select
}: {
  machines: MachineType[],
  select: (machine: MachineType) => void
}) {
  return (
    <GridCol span={12}>
      <ObjectTable
        items={machines}
        columns={columns}
        getRowKey={(m, i) => m.metadata?.name ?? `row-${i}`}
        onSelect={select}
        emptyLabel="No machines found"
      />
    </GridCol>
  )
}
