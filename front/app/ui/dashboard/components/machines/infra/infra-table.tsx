'use client';

import {MachineInfraType} from "@/app/ui/dashboard/components/machines/types";
import React from "react";
import {GridCol, Badge} from "@mantine/core";
import {StatusIndicator} from "@/app/ui/dashboard/shared/status-indicator";
import {toStatusState} from "@/app/ui/dashboard/shared/status";
import {ObjectTable} from "@/app/ui/dashboard/shared/object-table";
import {ColumnDef} from "@/app/ui/dashboard/base/types";

const columns: ColumnDef<MachineInfraType>[] = [
  {header: 'Name', render: (m) => m.metadata?.name ?? '—'},
  {header: 'Namespace', render: (m) => <Badge variant="light" color="gray">{m.metadata?.namespace ?? '—'}</Badge>},
  {header: 'ProviderID', render: (m) => m.providerID ?? '—'},
  {header: 'Template', render: (m) => m.template ?? '—'},
  {header: 'Age', render: (m) => m.age ?? '—'},
  {header: 'Status', align: 'center', render: (m) => <StatusIndicator state={toStatusState(m.status?.ready)} dotOnly/>},
];

export default function MachineInfraTable({
  machines, select
}: {
  machines: MachineInfraType[],
  select: (machine: MachineInfraType) => void
}) {
  return (
    <GridCol span={12}>
      <ObjectTable
        items={machines}
        columns={columns}
        getRowKey={(m, i) => m.metadata?.name ?? `row-${i}`}
        onSelect={select}
        emptyLabel="No vSphere machines found"
      />
    </GridCol>
  )
}
