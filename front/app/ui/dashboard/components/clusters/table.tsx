'use client';

import React from "react";

import {Badge} from '@mantine/core';
import { GridCol } from '@mantine/core';

import {ClusterType} from '@/app/ui/dashboard/components/clusters/types';
import {StatusIndicator} from '@/app/ui/dashboard/shared/status-indicator';
import {allReady} from '@/app/ui/dashboard/shared/status';
import {ObjectTable} from '@/app/ui/dashboard/shared/object-table';
import {ColumnDef} from '@/app/ui/dashboard/base/types';
import {ProviderBadge} from '@/app/ui/dashboard/shared/provider-badge';

const columns: ColumnDef<ClusterType>[] = [
  {header: 'Name', render: (c) => c.metadata?.name ?? '—'},
  {header: 'Namespace', render: (c) => <Badge variant="light" color="gray">{c.metadata?.namespace ?? '—'}</Badge>},
  {header: 'Provider', render: (c) => <ProviderBadge provider={c.provider}/>},
  {header: 'Version', render: (c) => c.topology?.kubernetesVersion ?? '—'},
  {header: 'Phase', render: (c) => c.status?.phase ?? '—'},
  {header: 'Age', render: (c) => c.age ?? '—', align: 'center'},
  {
    header: 'Status',
    align: 'center',
    render: (c) => (
      <StatusIndicator state={allReady(c.status?.controlPlaneReady, c.status?.infrastructureReady)} dotOnly/>
    ),
  },
];

export default function ClusterTable({
  clusters, select
}: {
  clusters: ClusterType[],
  select: (cluster: ClusterType) => void
}) {
  return (
    <GridCol span={12}>
      <ObjectTable
        items={clusters}
        columns={columns}
        getRowKey={(c, i) => c.metadata?.name ?? `row-${i}`}
        onSelect={select}
        emptyLabel="No clusters found"
      />
    </GridCol>
  )
}
