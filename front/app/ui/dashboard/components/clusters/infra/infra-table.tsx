'use client';

import React from "react";

import {Badge} from '@mantine/core';
import { GridCol } from '@mantine/core';

import {ClusterInfraType} from '@/app/ui/dashboard/components/clusters/types';
import {StatusIndicator} from '@/app/ui/dashboard/shared/status-indicator';
import {toStatusState} from '@/app/ui/dashboard/shared/status';
import {ObjectTable} from '@/app/ui/dashboard/shared/object-table';
import {ColumnDef} from '@/app/ui/dashboard/base/types';

const columns: ColumnDef<ClusterInfraType>[] = [
  {header: 'Name', render: (c) => c.metadata?.name ?? '—'},
  {header: 'Namespace', render: (c) => <Badge variant="light" color="gray">{c.metadata?.namespace ?? '—'}</Badge>},
  {header: 'Cluster', render: (c) => c.cluster ?? '—'},
  {header: 'Server', render: (c) => c.server ?? '—'},
  {header: 'Age', render: (c) => c.age ?? '—'},
  {header: 'Status', align: 'center', render: (c) => <StatusIndicator state={toStatusState(c.status?.ready)} dotOnly/>},
];

export default function ClusterInfraTable({
  clusters, select
}: {
  clusters: ClusterInfraType[]
  select: (cluster: ClusterInfraType) => void
}) {
  return (
    <GridCol span={12}>
      <ObjectTable
        items={clusters}
        columns={columns}
        getRowKey={(c, i) => c.metadata?.name ?? `row-${i}`}
        onSelect={select}
        emptyLabel="No vSphere clusters found"
      />
    </GridCol>
  )
}
